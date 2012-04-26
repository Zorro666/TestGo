package jake_graphics

import "os"
import "fmt"
import "code.google.com/p/x-go-binding/xgb"
import "image"
import "image/draw"

const (
	keyCodeStart = 8
	keyCodeEnd = 255
)

type WindowEvent interface{}

type MouseMoveEvent struct {
	X int
	Y int
}

type MouseButtonEvent struct {
	X int
	Y int
	ButtonMask int
	Buttons int
	EventType int
}

type KeyEvent struct {
	X int
	Y int
	KeyMask int
	Key int
}

type Jake_Graphics struct {
	m_c* xgb.Conn
	m_win xgb.Id
	m_gc xgb.Id
	m_drawable xgb.Id
	m_backbuffer* image.RGBA
	m_windowWidth int
	m_windowHeight int
	m_keyboardMapping* xgb.GetKeyboardMappingReply
}

func NewInstance() *Jake_Graphics {
	jg := Jake_Graphics{m_c:nil}
	c, err := xgb.Dial(os.Getenv("DISPLAY"))
		if err != nil {
			fmt.Printf("Jake_Graphics: cannot connect to X server: '%v'\n", err)
			return nil
		}
	jg.m_c = c
	return &jg
}

func (jg* Jake_Graphics) CloseWindow() {
	jg.m_c.Close()
}

func (jg* Jake_Graphics) GetBackBuffer() draw.Image {
	return jg.m_backbuffer
}

func (jg* Jake_Graphics) CreateWindow(width int, height int, x0 int, y0 int) bool {
	if (jg.m_c == nil) {
			fmt.Printf("Jake_Graphics: connection is null\n")
			return false
	}

	jg.m_win = jg.m_c.NewId()
	jg.m_gc = jg.m_c.NewId()
	jg.m_drawable = jg.m_c.NewId()
	jg.m_windowWidth = width
	jg.m_windowHeight = height

	var depth byte = 0

	jg.m_c.CreateWindow(depth, jg.m_win, jg.m_c.DefaultScreen().Root,
										  int16(x0), int16(y0), uint16(width), uint16(height), 0, 0, 0, 0, nil)
	jg.m_c.ChangeWindowAttributes(jg.m_win, xgb.CWEventMask,
			[]uint32{xgb.EventMaskExposure |
			xgb.EventMaskKeyRelease | xgb.EventMaskKeyPress |
			xgb.EventMaskButtonPress | xgb.EventMaskButtonRelease |
			xgb.EventMaskPointerMotion})
	jg.m_c.CreateGC(jg.m_gc, jg.m_win, 0, nil)
	jg.m_c.MapWindow(jg.m_win)

	r := image.Rect(0, 0, width, height)
	jg.m_backbuffer = image.NewRGBA(r)
	jg.FlipBackBuffer()

	firstKeyCode := keyCodeStart
	count := keyCodeEnd-keyCodeStart+1
	keyboardMapping, _ := jg.m_c.GetKeyboardMapping(byte(firstKeyCode), byte(count))
	jg.m_keyboardMapping = keyboardMapping
/*
	for i := 0; i < int(jg.m_keyboardMapping.Length); i++ {
		keysym := jg.m_keyboardMapping.Keysyms[i]
		keycode := (i/int(jg.m_keyboardMapping.KeysymsPerKeycode))+8
		fmt.Printf("i:%v keycode:0x%X keySym:0x%X '%c'\n", i, keycode, keysym, keysym)
	}
*/
	fmt.Printf("Length:%d keysymsPerKeycode:%d\n", jg.m_keyboardMapping.Length, jg.m_keyboardMapping.KeysymsPerKeycode)

	return true
}

func (jg* Jake_Graphics) FlipBackBuffer() {
	var format byte = xgb.ImageFormatZPixmap
	widthLeft := jg.m_windowWidth
	heightLeft := jg.m_windowHeight
	backbuffer := jg.m_backbuffer

	var leftPad byte = 0
	var depth byte = 24

	storage := make([]byte, 0, 256*256*4)
	width := widthLeft
	height := heightLeft

	dstX := 0
	dstY := 0
  for {
	  maxBytes:= 256*256
		if (width*height >= maxBytes) {
	    height = maxBytes/width
    }
	  data := storage[:width*height*4]

		for y:= 0; y < height; y++ {
			for x:= 0; x < width; x++ {
		    pixel := backbuffer.At(x+dstX,y+dstY)
				red, green, blue, _ := pixel.RGBA()
				i := 4*(y*width+x)
				data[i+0] = byte(blue)
				data[i+1] = byte(green)
				data[i+2] = byte(red)
			}
		}
		jg.m_c.PutImage(format, jg.m_win, jg.m_gc, uint16(width), uint16(height), int16(dstX), int16(dstY), leftPad, depth, data)
    dstY += height
		heightLeft -= height
		height = heightLeft
		if heightLeft <= 0 {
			break
    }
	}
}

func (jg* Jake_Graphics) WaitForEvent() (event WindowEvent) {
  reply, err := jg.m_c.WaitForEvent()
	if err != nil {
		fmt.Printf("error: %v\n", err)
		os.Exit(1)
	}
	switch x11Event := reply.(type) {
		case xgb.KeyPressEvent:
			event := KeyEvent{}
		  event.X = int(x11Event.EventX)
		  event.Y = int(x11Event.EventY)
			keycode := int(x11Event.Detail)
			modifier := int(x11Event.State)
			keycodeMapping := (keycode-8)*int(jg.m_keyboardMapping.KeysymsPerKeycode) + modifier
		  event.Key = int(jg.m_keyboardMapping.Keysyms[keycodeMapping])
			fmt.Printf("keycode:0x%X keysym:0x%X\n", keycode, event.Key)
			return event
		case xgb.KeyReleaseEvent:
			event := KeyEvent{}
		  event.X = int(x11Event.EventX)
		  event.Y = int(x11Event.EventY)
			keycode := int(x11Event.Detail)
			modifier := int(x11Event.State)
			keycodeMapping := (keycode-8)*int(jg.m_keyboardMapping.KeysymsPerKeycode) + modifier
		  event.Key = -int(jg.m_keyboardMapping.Keysyms[keycodeMapping])
			return event
		case xgb.ButtonPressEvent:
			event := MouseButtonEvent{}
		  event.X = int(x11Event.EventX)
		  event.Y = int(x11Event.EventY)
		  event.Buttons = int(x11Event.Detail)
			return event
		case xgb.ButtonReleaseEvent:
			event := MouseButtonEvent{}
		  event.X = int(x11Event.EventX)
		  event.Y = int(x11Event.EventY)
		  event.Buttons = -int(x11Event.Detail)
			return event
		case xgb.MotionNotifyEvent:
			event := MouseMoveEvent{}
		  event.X = int(x11Event.EventX)
		  event.Y = int(x11Event.EventY)
			return event
  }

	return event
}

/*
		reply, err := c.WaitForEvent()
			if err != nil {
				fmt.Printf("error: %v\n", err)
					os.Exit(1)
			}
	c.Close()
}
*/
