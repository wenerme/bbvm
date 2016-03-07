package vm

type BackgroundMode int
const (
//透明显示，即字体的背景颜色无效。
	TRANSPARENT BackgroundMode = 1
//不透明显示，即字体的背景颜色有效
	OPAQUE BackgroundMode = 2
)

type EnvType int
const (
	ENV_SIM EnvType = 0
	ENV_9288 EnvType = 9288
	ENV_9188 EnvType = 9188
	ENV_9288T EnvType = 9287
	ENV_9288S EnvType = 9286
	ENV_9388 EnvType = 9388
)

type FontType int
const (
	FONT_12SONG FontType = 0
	FONT_12KAI FontType = 1
	FONT_12HEI FontType = 2
	FONT_16SONG FontType = 3
	FONT_16KAI FontType = 4
	FONT_16HEI FontType = 5
	FONT_24SONG FontType = 6
	FONT_24KAI FontType = 7
	FONT_24HEI FontType = 8
)

type KeyCode int
const (
	KEY_UP KeyCode = 38
	KEY_DOWN KeyCode = 40
	KEY_LEFT KeyCode = 37
	KEY_RIGHT KeyCode = 39
	KEY_SPACE KeyCode = 32
	KEY_ESCAPE KeyCode = 27
	KEY_ENTER KeyCode = 13
)

type DrawMode int
const (
	KEY_COLOR DrawMode = 1
)
type BrushStyle int
const (
	BRUSH_SOLID BrushStyle = 0
)