package parser

import (
	"encoding/binary"
	"encoding/hex"
	"fmt"
	"github.com/juju/errors"
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

func (v pseudoDataStr) MarshalBinary() ([]byte, error) {
	return v, nil
}
func (v pseudoDataStr) Len() int {
	return len(v)
}
func (v pseudoDataStr) Assembly() string {
	return fmt.Sprintf(`"%s"`, string(v))
}

type pseudoDataBytes []byte

func (v pseudoDataBytes) MarshalBinary() ([]byte, error) {
	return v, nil
}
func (v pseudoDataBytes) Len() int {
	return len(v)
}
func (v pseudoDataBytes) Assembly() string {
	return fmt.Sprintf(`%%%s%%`, hex.EncodeToString(v))
}

func createPseudoDataValue(v interface{}) (d PseudoDataValue, e error) {
	switch v.(type) {
	case int:
		d = pseudoDataInt(v.(int))
	case string:
		b := []byte(v.(string))
		switch b[0] {
		case '"':
			d = pseudoDataStr(b[1 : len(b)-1])
		case '%':
			b, e := hex.DecodeString(string(b[1 : len(b)-1]))
			if e != nil {
				return nil, e
			}
			d = pseudoDataBytes(b)
		default:
			e = errors.Errorf("Currently can not create pseudo data value by %#v", v)
		}
	// TODO Symbol reference
	default:
		e = errors.Errorf("Can not create pseudo data value by %#v", v)
	}
	return
}
