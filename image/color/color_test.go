package color
import (
	"testing"
	"github.com/stretchr/testify/assert"
	"image/color"
)


func TestColor(t *testing.T) {
	assert := assert.New(t)

	c := BGR888{255}
	assert.EqualValues(c, BGR888Model.Convert(color.RGBAModel.Convert(c)))
	assert.EqualValues(31, BGR565Model.Convert(c).(BGR565).V)

	cb := BGR565{46584}
	r, g, b, a := color.RGBAModel.Convert(cb).RGBA()
	assert.EqualValues([]uint32{192, 188, 176, 255}, []uint32{r>>8, g>>8, b>>8, a>>8})
}