package jake_graphics

import (
		"os"
		"fmt"
		"code.google.com/p/x-go-binding/xgb"
		)

type Jake_Graphics struct {
	m_c* xgb.Conn
	m_win xgb.Id
	m_gc xgb.Id
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

func (jg* Jake_Graphics) CreateWindow(width uint16, height uint16) bool {
	if (jg.m_c == nil) {
			fmt.Printf("Jake_Graphics: connection is null\n")
			return false
	}

	jg.m_win = jg.m_c.NewId()
	jg.m_gc = jg.m_c.NewId()

	jg.m_c.CreateWindow(0, jg.m_win, jg.m_c.DefaultScreen().Root, 100, 100, width, height, 0, 0, 0, 0, nil)
	jg.m_c.CreateGC(jg.m_gc, jg.m_win, 0, nil)
	jg.m_c.MapWindow(jg.m_win)
	return true
}

/*
func main() {
	c, err := xgb.Dial(os.Getenv("DISPLAY"))
		if err != nil {
			fmt.Printf("cannot connect: %v\n", err)
				os.Exit(1)
		}

	fmt.Printf("vendor = '%s'\n", string(c.Setup.Vendor))

		win := c.NewId()
		gc := c.NewId()

		c.CreateWindow(0, win, c.DefaultScreen().Root, 150, 150, 200, 200, 0, 0, 0, 0, nil)
		c.ChangeWindowAttributes(win, xgb.CWEventMask,
				[]uint32{xgb.EventMaskExposure | xgb.EventMaskKeyRelease})
		c.CreateGC(gc, win, 0, nil)
		c.MapWindow(win)

		atom, _ := c.InternAtom(false, "HELLO")
		fmt.Printf("atom = %d\n", atom.Atom)

		points := make([]xgb.Point, 2)
		points[0] = xgb.Point{5, 5}
	points[1] = xgb.Point{100, 120}

	hosts, _ := c.ListHosts()
		fmt.Printf("hosts = %+v\n", hosts)

		ecookie := c.ListExtensionsRequest()
		exts, _ := c.ListExtensionsReply(ecookie)
		for _, name := range exts.Names {
			fmt.Printf("exts = '%s'\n", name.Name)
		}

	for {
		reply, err := c.WaitForEvent()
			if err != nil {
				fmt.Printf("error: %v\n", err)
					os.Exit(1)
			}
		fmt.Printf("event %T\n", reply)
			switch event := reply.(type) {
				case xgb.ExposeEvent:
					c.PolyLine(xgb.CoordModeOrigin, win, gc, points)
				case xgb.KeyReleaseEvent:
						fmt.Printf("key release!\n")
							points[0].X = event.EventX
							points[0].Y = event.EventY
							c.PolyLine(xgb.CoordModeOrigin, win, gc, points)
							c.Bell(75)
			}
	}

	c.Close()
}
*/
