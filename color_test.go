package bbvm
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"image/color"
)


func TestColor(t *testing.T) {
	assert := assert.New(t)

	c := BGR888Color(255)
	assert.EqualValues(c, BGR888Model.Convert(color.RGBAModel.Convert(c)))
	assert.EqualValues(31, BGR565Model.Convert(c))
}