package main

import (
	"github.com/codegangsta/cli"
	"os"
	"fmt"
	"path"
	"github.com/op/go-logging"
	"bytes"
	"os/exec"
	"golang.org/x/text/encoding/simplifiedchinese"
	"regexp"
	"strings"
	"path/filepath"
	"io"
)
var log = logging.MustGetLogger("bb")
var env = make(map[string]string)
var paths = make([]string, 0)
func main() {
	app := cli.NewApp()
	app.Name = "bb"
	app.Usage = "BB command tool"
	app.Version = "0.0.1"
	app.Author = "wener"
	app.Email = "wener@wener.me"
	app.Action = func(c *cli.Context) {
		println("bb help for useage")
	}
	app.Before = check
	app.Flags = []cli.Flag{
		cli.StringFlag{
			Name: "bb, b",

			Usage: "BB home directory",
			EnvVar: "BB_HOME",
		},
		cli.StringFlag{
			Name: "gamdev",
			Usage: "Gamdev.exe directory",
		},
		cli.BoolFlag{
			Name: "no-autopath",
			Usage: "Disable auto detect path",
		},
		cli.BoolFlag{
			Name: "verbose",
			Usage: "Verbose log output",
		},
	}
	app.Commands = []cli.Command{
		{
			Name:      "compile",
			Aliases:     []string{"c"},
			Usage:     "Compile basm,bas",
			Action: func(c *cli.Context) {
				compile(c.Args().First())
			},

		},
		{
			Name:      "run",
			Aliases:     []string{"r"},
			Usage:     "Run bin,basm,bas",
			Action: func(c *cli.Context) {
				run(c.Args().First(), true)
			},
		},
		{
			Name:      "prepare",
			Aliases:     []string{"r"},
			Usage:     "Copy compile bin to BBasic, but not run",
			Action: func(c *cli.Context) {
				run(c.Args().First(), false)
			},
		},
		{
			Name:      "env",
			Usage:     "Show env setting",
			Action: func(c *cli.Context) {
				print(reportEnv())
			},
		},
		{
			Name:      "tool",
			Usage:     "Run a tool",
			Action: func(c *cli.Context) {
				if c.Args().First() == "" {
					log.Warning("No tool name")
					os.Exit(1)
				}
				fmt.Print(tool(c.Args().First(), c.Args().Tail()...))
			},
		},
	}

	//	fmt.Println(os.Args)
	//	app.Run(os.Args)
	os.Chdir("doc/testsuit")
//	os.Args = []string{"bb", "--bb", "doc/testsuit/BB", "prepare", "test.bas"}
	os.Args = []string{"bb", "--bb", "doc/testsuit/BB", "prepare", "../../tests/case/in/38.basm"}
	app.RunAndExitOnError()
	//	app.Run([]string{"bb", "--bb", "BB/Tool", "help"})
	//	fmt.Println("BBDIR ", bbDir.Value)
}

func init() {

}

var GBKDecoder = simplifiedchinese.GBK.NewDecoder()
func tool(exe string, args ...string) string {
	if p, ok := env[exe]; ok {

		//		cmd := exec.Command(p, args...)
		args = append([]string{p}, args...)

		p, args = adapterCommand(p, args)
		log.Info("Run command %s ARGS: %s", p, args)
		cmd := exec.Command(p, args...)

		out := &bytes.Buffer{}
		cmd.Stdout = out
		cmd.Stderr = os.Stderr

		_ = cmd.Run()

		decoded := make([]byte, out.Len()*2)
		GBKDecoder.Reset()
		n, _, decErr := GBKDecoder.Transform(decoded, out.Bytes(), true)
		if decErr != nil {panic(decErr)}
		decoded = decoded[:n]
		return string(decoded)
	}
	panic(fmt.Sprintf("tool %s not found", exe))
}

func adapterCommand(p string, args []string) (string, []string) {
	//	t := "open"
	args = append([]string{p}, args...)
	return "wine", args
}

func check(c *cli.Context) error {
	logging.SetLevel(logging.WARNING, "bb")
	if c.Bool("verbose") {
		logging.SetLevel(logging.INFO, "bb")
	}

	cwd, err := os.Getwd()
	if err!=nil {panic(err)}
	paths = append(paths, cwd)
	log.Info("CWD %s", cwd)

	if v := c.String("bb"); v != "" {
		if !path.IsAbs(v) {
			v = path.Join(cwd, v)
		}
		tryPath(v+"/Tool")
		tryPath(v+"/Sim/Debug")
	}

	if !c.Bool("no-autopath") {
		tryPath(path.Join(cwd, "BB/Tool"))
		tryPath(path.Join(cwd, "BB/Sim/Debug"))
	}

	uniqueStr(paths)

	trySet("BBasic.exe", "bbasic")
	trySet("BBTool.exe", "bbtool")
	trySet("Blink.exe", "blink")
	trySet("GamDev.exe", "gamdev")


	log.Info(reportEnv())
	return nil
}
func tryPath(p string) {
	log.Info("Try Path %s", p)
	if _, err := os.Stat(p); !os.IsNotExist(err) {
		paths = append(paths, p)
	}
}

func uniqueStr(col []string) []string {
	m := map[string]struct {}{}
	for _, v := range col {
		if _, ok := m[v]; !ok {
			m[v] = struct {}{}
		}
	}
	list := make([]string, len(m))

	i := 0
	for v := range m {
		list[i] = v
		i++
	}
	return list
}

func reportEnv() string {
	sb := bytes.NewBufferString("# ENV\n")
	fmt.Fprintf(sb, "PATH: %s\n", paths)
	for k, v := range env {
		fmt.Fprintf(sb, "%s: %s\n", k, v)
	}
	return sb.String()
}
func findUnderPath(fn string) (string, bool) {
	for _, p := range paths {
		fp := path.Join(p, fn)
		if _, err := os.Stat(fp); !os.IsNotExist(err) {
			return fp, true
		}
	}
	return "", false
}
func trySet(fn string, key string) bool {
	if p, ok := findUnderPath(fn); ok {
		env[key]=p
		return true
	}
	log.Warning("%s not found in path", fn)
	return false
}

func compile(fn string) (string, bool) {
	if checked, ok := checkFile(fn); ok {
		fn = checked
	}else {
		return "", false
	}

	if m, _ := regexp.MatchString(".bas$", fn); m {
		log.Info("Compile %s to basm", fn)
		r := tool("bbasic", fn)
		if !strings.Contains(r, "编译成功") {
			fmt.Println(r)
			return "", false
		}
		fn = strings.TrimSuffix(fn, filepath.Ext(fn))+".obj"
	}

	if m, _ := regexp.MatchString(".(obj|basm)$", fn); m {
		log.Info("Link %s to bin", fn)
		bin := strings.TrimSuffix(fn, filepath.Ext(fn))+".bin"

		r := tool("blink", fn, bin)
		//		fmt.Println(r)
		if !strings.Contains(r, "连接成功") {
			fmt.Println(r)
			return "", false
		}
		return bin, true
	}
	log.Error("Can not compile %s", fn)

	return "", false
}

func run(fn string, doRun bool) (bool) {
	if checked, ok := checkFile(fn); ok {
		fn = checked
	}else {
		return false
	}

	var bin string = fn
	switch filepath.Ext(fn){
		case ".bas", ".basm", ".obj":

		var ok bool
		if bin, ok = compile(fn); !ok {
			return false
		}
		fallthrough
		case ".bin":

		dest := path.Join(path.Dir(env["gamdev"]), "../BBasic", "test.bin")
		//		fmt.Println("Will copy to ", bin, dest)
		log.Info("Copy bin %s to %s", bin, dest)
		copyFileContents(bin, dest)
		cwd, _ := os.Getwd()
		os.Chdir(path.Dir(env["gamdev"]))
		defer func() {
			os.Chdir(cwd)
		}()
		if doRun {
			tool("gamdev")
		}

		default:
		fmt.Println("Can not run %s, need bas|basm|obj|bin", fn)
		return false
	}
	return false
}



// copyFileContents copies the contents of the file named src to the file named
// by dst. The file will be created if it does not already exist. If the
// destination file exists, all it's contents will be replaced by the contents
// of the source file.
func copyFileContents(src, dst string) (err error) {
	in, err := os.Open(src)
	if err != nil {
		return
	}
	defer in.Close()
	out, err := os.Create(dst)
	if err != nil {
		return
	}
	defer func() {
		cerr := out.Close()
		if err == nil {
			err = cerr
		}
	}()
	if _, err = io.Copy(out, in); err != nil {
		return
	}
	err = out.Sync()
	return
}
func checkFile(fn string) (string, bool) {
	if fn == "" {
		fmt.Println("No input file")
		return fn, false
	}


	if !path.IsAbs(fn) {
		cwd, err := os.Getwd()
		if err !=nil {fmt.Print(err)}
		fn = path.Join(cwd, fn)
	}

	if _, err := os.Stat(fn); os.IsNotExist(err) {
		fmt.Printf("%s not found\n", fn)
		return fn, false
	}
	return fn, true
}