package parser

import (
	"encoding/binary"
	"fmt"
	"github.com/juju/errors"
	"github.com/wenerme/bbvm/libbbvm/asm"
	"strconv"
)

type pseudoDataInt int32

func (v pseudoDataInt) MarshalBinary() (data []byte, err error) {
	data = make([]byte, 4)
	binary.LittleEndian.PutUint32(data, uint32(v))
	return
}
func (v pseudoDataInt) Len() int {
	return 4
}
func (v pseudoDataInt) Assembly() string {
	return strconv.Itoa(int(v))
}

type pseudoDataStr []byte

func (v pseudoDataStr) MarshalBinary() (data []byte, err error) {
	data = []byte(v)
	return
}
func (v pseudoDataStr) Len() int {
	return len(v)
}
func (v pseudoDataStr) Assembly() string {
	return fmt.Sprintf(`"%s"`, string(v))
}

func createPseudoDataValue(v interface{}) (d asm.PseudoDataValue, e error) {
	switch v.(type) {
	case int:
		d = pseudoDataInt(v.(int))
	case string:
		b := []byte(v.(string))
		if b[0] == '"' {
			d = pseudoDataStr([]byte(b[1 : len(b)-1]))
		}
	// TODO Symbol reference
	default:
		e = errors.Errorf("Can not create pseudo data value by %#v", v)
	}
	return
}
