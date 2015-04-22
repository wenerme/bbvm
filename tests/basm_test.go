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

	regIO := regexp.MustCompile(`;\s*(\|?(\<|\>))([^\r\n]*)`)

	matches := regIO.FindAllStringSubmatch(basm, -1)
	//	t.Log(matches)
	for _, v := range matches {
		val := string(v[3])
		val = strings.TrimRight(val, " ")
		val = strings.Replace(val, `\n`, "\n", -1)
		//		t.Logf("%#v", v)
		switch string(v[1]){
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

	t.Logf("%10s: %#v\n", "expected", string(expected.Bytes()))
	t.Logf("%10s: %#v\n", "input", string(input.Bytes()))

	v := NewVM()

	rom, err := ioutil.ReadFile(strings.Replace(file, ".basm", ".bin", -1))
	if err!= nil {
		t.Error(err)
		t.Fail()
	}
	v.Load(rom[16:])
	IN.StrFunc(v)
	IN.ConvFunc(v)
	IN.Misc(v)
	OUT.OutputToWriter(v, output)
	OUT.InputByReader(v, input)
	logging.SetLevel(logging.INFO, "bbvm")
	for !v.IsExited() {
		v.Loop()// call
//		t.Log(v.Report())
//		t.Logf("%10s: %#v\n", "output", string(output.Bytes()))
	}
	t.Logf("%10s: %#v\n", "output", string(output.Bytes()))

	for {
		o, oe := output.ReadString('\n')
		e, ee := expected.ReadString('\n')
		if len(o) > 1 {o = o[:len(o)-1]}
		if len(e) > 1 {e = e[:len(e)-1]}

		// support skip syntax
		if e == "skip" { continue }

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

func TestIn9(t *testing.T) {
	testByBAsm("case/in/13.basm", t)
}
