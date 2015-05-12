package bc
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"image/color"
	"log"
)


func TestColor(t *testing.T) {
	assert := assert.New(t)

	c := BGR888{255}
	assert.EqualValues(c, BGR888Model.Convert(color.RGBAModel.Convert(c)))
	assert.EqualValues(31, BGR565Model.Convert(c))

	cb := BGR565{46584}
	log.Print("Color %+v", color.RGBAModel.Convert(cb))
}