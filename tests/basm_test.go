package bbvm
import (
	"testing"
	"bytes"
	"io/ioutil"
	"regexp"
	"strings"
	. "../."
	"github.com/op/go-logging"
)


func testByBAsm(file string, t *testing.T) bool {
	t.Logf("\nTest basm file %s\n", file)

	input := &bytes.Buffer{}
	expected := &bytes.Buffer{}
	output := &bytes.Buffer{}

	b, err := ioutil.ReadFile(file)
	if err!= nil { panic(err) }
	basm := string(b)

	regIO := regexp.MustCompile(`(?m)^\s*;\s*(\|?\<|\>)([^\r\n]*)`)

	matches := regIO.FindAllStringSubmatch(basm, -1)
	//	t.Log(matches)
	for _, v := range matches {
		val := string(v[2])
		val = strings.TrimRight(val, " ")
		val = strings.Replace(val, `\n`, "\n", -1)
		//		t.Logf("%s %#v",string(v[1]),string(v[2]))
		switch string(v[1]){
			case "|>":
			expected.WriteString(val)
			case ">":
			expected.WriteString(val)
			expected.WriteString("\n")
			case "|<":
			input.WriteString(val)
			case "<":
			input.WriteString(val)
			input.WriteString("\n")
		}
	}

	t.Logf("%10s: %#v\n","expected", string(expected.Bytes()))
	t.Logf("%10s: %#v\n","input", string(input.Bytes()))

	v := NewVM()

	rom, err := ioutil.ReadFile(strings.Replace(file, ".basm", ".bin", -1))
	if err!= nil {
		t.Error(err)
		t.Fail()
	}
	v.Load(rom[16:])
	HandInStr(v)
	OUT.OutputToWriter(v,output)
	OUT.InputByReader(v,input)
	logging.SetLevel(logging.INFO, "bbvm")
	for !v.IsExited() {
		v.Loop()// call
		//		t.Log(v.Report())
	}
	t.Logf("%10s: %#v\n","output", string(output.Bytes()))

	if bytes.Compare(output.Bytes(), expected.Bytes())!= 0 {
		t.Error("Output is not expected")
		return false
	}
	return true
}

func TestIn9(t *testing.T) {
	testByBAsm("case/out/10.basm", t)
}