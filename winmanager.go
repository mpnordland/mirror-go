package main

import (
	"github.com/BurntSushi/xgb"
	"github.com/BurntSushi/xgb/xproto"
	"github.com/BurntSushi/xgbutil"
	"github.com/BurntSushi/xgbutil/xrect"
	"github.com/BurntSushi/xgbutil/xwindow"
)

type WindowManager struct {
	X       *xgbutil.XUtil
	WinTree *Tree
	*EventLoop
}

func NewWindowManager(X *xgbutil.XUtil) *WindowManager {
	w := &WindowManager{X: X, Windows: &Tree{}, EventLoop: NewEventLoop(X)}
	w.AddCallback("MapRequestEvent", w.X.RootWin(), w.MapRequest)
	w.AddCallback("ConfigureRequestEvent", w.X.RootWin(), w.ConfigureRequest)
	return w
}

func (w *WindowManager) LayoutWindows(args ...interface{}) {
	w.Windows.Layout(xrect.New(0, 0, 800, 600))
}

func (w *WindowManager) Manage(window xproto.Window) {
	win := NewWindow(w.X, window)
    //attach window hooks to allow it to handle unmanaging and configuring itself
	w.AddCallback("UnmapNotifyEvent", win.Id, win.UnmapNotify)
	w.AddCallback("ConfigureRequestEvent", win.Id, win.ConfigureRequest)
    //this order is important because there are some hooks that are
    //attached inside the tree that need to be pulled before
    //we call the layout again.
	w.WinTree.Insert(win)
	win.AttachToHook("window::unmapped", w.LayoutWindows)
	win.AttachToHook("window::configured", w.LayoutWindows)
	win.Map()
	w.LayoutWindows()
}

//Event Handlers
func (w *WindowManager) MapRequest(event xgb.Event) {
	ev := event.(xproto.MapRequestEvent)
	w.Manage(ev.Window)
}

func (w *WindowManager) ConfigureRequest(event xgb.Event) {
	ev := event.(xproto.ConfigureRequestEvent)
	xwindow.New(w.X, ev.Window).Configure(int(ev.ValueMask),
		int(ev.X), int(ev.Y), int(ev.Width), int(ev.Height),
		ev.Sibling, ev.StackMode)
}
