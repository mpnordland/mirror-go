package main

import( "github.com/BurntSushi/xgbutil/xrect"
        "github.com/BurntSushi/xgb"
        "github.com/BurntSushi/xgb/xproto"
        "github.com/BurntSushi/xgbutil"
)

type Window struct{
    *Dispatcher
    Id xproto.Window
    X *xgbutil.XUtil
    Geom *xrect.XRect
   /* Name string
    Class string
    States []xprop.Atom
    Types []xprop.Atom
    Protocols []xprop.Atom*/
}

func NewWindow(X *xgbutil.XUtil, Id xproto.Window) *Window{
    win := new(Window)
    win.Id = Id
    win.X = X
    win.Dispatcher = NewDispatcher()
    win.Geometry()
    return win
}

func (w *Window) Geometry() (*xrect.XRect, error){
    cookie := xproto.GetGeometry(w.X.Conn(), xproto.Drawable(w.Id))
    reply, err := cookie.Reply()
    if err!= nil{
        return nil, err
    }
    w.Geom = xrect.New(int(reply.X), int(reply.Y),
       int(reply.Width), int(reply.Height))
    return w.Geom, nil
}

func (w* Window) GeometrySet(rect *xrect.XRect){
    if rect.Width() == 0 || rect.Height() == 0{
        w.SendConfigNotify()
        return
    }
    if rect.Width() != w.Geom.Width() ||
        rect.Height() != w.Geom.Height() ||
        rect.X() != w.Geom.X()||
        rect.Y() != w.Geom.Y(){
        w.Configure(xproto.ConfigWindowX | xproto.ConfigWindowY |
			xproto.ConfigWindowWidth | xproto.ConfigWindowHeight, rect.X(), rect.Y(), rect.Width(), rect.Height(), 0, 0, 0)
        if rect.Width() == w.Geom.Width() && rect.Height() == w.Geom.Height(){
            w.SendConfigNotify()
        }
    }
    w.Geom = rect
}

func (win* Window) Configure(flags, x, y, w, h, borderWidth int, sibling xproto.Window, stackMode byte){
	if win == nil {
		return
	}
	vals := []uint32{}
	if xproto.ConfigWindowX&flags > 0 {
		vals = append(vals, uint32(x))
	}
    if xproto.ConfigWindowY&flags > 0 {
		vals = append(vals, uint32(y))
	}
	if xproto.ConfigWindowWidth&flags > 0 {
		if int16(w) <= 0 {
				w = 1
		}
		vals = append(vals, uint32(w))
	}
	if xproto.ConfigWindowHeight&flags > 0 {
		if int16(h) <= 0 {
				h = 1
		}
		vals = append(vals, uint32(h))
	}
    if xproto.ConfigWindowBorderWidth&flags > 0{
        vals = append(vals, uint32(sibling))
    }
	if xproto.ConfigWindowSibling&flags > 0 {
		vals = append(vals, uint32(sibling))
	}
	if xproto.ConfigWindowStackMode&flags > 0 {
		vals = append(vals, uint32(stackMode))
	}

	xproto.ConfigureWindow(win.X.Conn(), win.Id, uint16(flags), vals)
}

func (w *Window) Map(){
    xproto.MapWindow(w.X.Conn(), w.Id)
}
func (w *Window) Unmap(){
    xproto.UnmapWindow(w.X.Conn(), w.Id)
}

func (w *Window) SendConfigNotify(){
	ev := xproto.ConfigureNotifyEvent{
		Event:            w.Id,
		Window:           w.Id,
		AboveSibling:     0,
		X:                int16(w.Geom.X()),
		Y:                int16(w.Geom.Y()),
		Width:            uint16(w.Geom.Width()),
		Height:           uint16(w.Geom.Height()),
		BorderWidth:      0,
		OverrideRedirect: false,
	}
	xproto.SendEvent(w.X.Conn(), false, w.Id,
		xproto.EventMaskStructureNotify, string(ev.Bytes()))
}

func (w *Window) ConfigureRequest(event xgb.Event){
    ev := event.(xproto.ConfigureRequestEvent)
    w.GeometrySet(xrect.New(int(ev.X), int(ev.Y), int(ev.Width), int(ev.Height)))
    w.PullHook("window::configured", w.Id)
}

func (w *Window) UnmapNotify(event xgb.Event){
    w.PullHook("window::unmapped", w.Id)
}

func (w *Window) DestroyNotify(event xgb.Event){
    w.PullHook("window::destroyed", w.Id)
}


