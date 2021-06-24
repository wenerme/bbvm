package bbvm

import (
	"bytes"
	"github.com/op/go-logging"
	"golang.org/x/image/bmp"
	"image"
	"image/png"
	"io/ioutil"
	"os"
	"regexp"
	"strings"
	"testing"
)

func testByBAsm(file string, t *testing.T) bool {
	t.Logf("\nTest basm file %s\n", file)

	input := &bytes.Buffer{}
	expected := &bytes.Buffer{}
	output := &bytes.Buffer{}

	b, err := ioutil.ReadFile(file)
	if err != nil {
		panic(err)
	}

	extractIO(string(b), input, expected)

	t.Logf("%10s: %#v\n", "expected", string(expected.Bytes()))
	t.Logf("%10s: %#v\n", "input", string(input.Bytes()))

	v := NewVM()

	rom, err := ioutil.ReadFile(strings.Replace(file, ".basm", ".bin", -1))
	if err != nil {
		t.Error(err)
		t.Fail()
	}
	v.Load(rom[16:])
	IN.All(v)
	OUT.OutputToWriter(v, output)
	OUT.InputByReader(v, input)
	OUT.File(v)
	OUT.Graphic(v)
	Misc.All(v)

	logging.SetLevel(logging.DEBUG, "bbvm")
	for !v.IsExited() {
		v.Loop() // call
		//		t.Log(v.Report())
		//		t.Logf("%10s: %#v\n", "output", string(output.Bytes()))
	}
	t.Logf("%10s: %#v\n", "output", string(output.Bytes()))

	// Debug page 1
	saveImage(v.Attr()["graph-dev"].(GraphDev).Screen(), "screen.bmp")
	//	saveImage(v.Attr()["graph-dev"].(GraphDev).PicPool().Get(0).Get().(Picture), "pic.bmp")

	for {
		o, oe := output.ReadString('\n')
		e, ee := expected.ReadString('\n')
		if len(o) > 1 {
			o = o[:len(o)-1]
		}
		if len(e) > 1 {
			e = e[:len(e)-1]
		}

		// support skip syntax
		if e == "skip" {
			continue
		}

		if o != e {
			t.Errorf("Output is not expected: %s != %s", o, e)
			return false
		}
		if oe != nil || ee != nil {
			if len(expected.Bytes()) != len(output.Bytes()) {
				t.Error("Final output is not expected")
				return false
			}
			break
		}
	}
	return true
}

func extractIO(basm string, input *bytes.Buffer, expected *bytes.Buffer) {
	regIO := regexp.MustCompile(`(?:;|')\s*(\|?(\<|\>))([^\r\n]*)`)

	matches := regIO.FindAllStringSubmatch(basm, -1)
	//	t.Log(matches)
	for _, v := range matches {
		val := string(v[3])
		val = strings.TrimRight(val, " ")
		val = strings.Replace(val, `\n`, "\n", -1)
		//		t.Logf("%#v", v)
		switch string(v[1]) {
		case "|>":
			expected.WriteString(val)
		case ">":
			expected.WriteString(val)
			expected.WriteString("\n")
		//			t.Logf("APPEND LINE %s =  %s\n", val, expected.String())
		case "|<":
			input.WriteString(val)
		case "<":
			input.WriteString(val)
			input.WriteString("\n")
		}
	}
}
func TestIn9(t *testing.T) {
	//	testByBAsm("case/in/38.basm", t)
	//	testByBAsm("case/out/read-restore.basm", t)
	testByBAsm("case/draw-line.basm", t)
}

func saveTemp(i image.Image) {
	p, err := os.Create("temp.png")
	if err != nil {
		panic(err)
	}
	err = png.Encode(p, i)
	if err != nil {
		panic(err)
	}
}

func saveImage(i image.Image, fn string) {
	p, err := os.Create(fn)
	if err != nil {
		panic(err)
	}
	err = bmp.Encode(p, i)
	if err != nil {
		panic(err)
	}
}
