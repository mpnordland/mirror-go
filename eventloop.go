package main

import(
        "fmt"
        "reflect"
        "github.com/BurntSushi/xgb"
        "github.com/BurntSushi/xgb/xproto"
        "github.com/BurntSushi/xgbutil"
)

func getEventName(ev xgb.Event) string{
    return reflect.TypeOf(ev).Name()
}

func winFromValue(v reflect.Value) xproto.Window{
    return xproto.Window(v.Uint())
}

func getEventWin(ev xgb.Event) []xproto.Window{
    evVal := reflect.ValueOf(ev)
    evType := evVal.Type()
    var w xproto.Window
    ret := make([]xproto.Window, 2)

    if t := evType.Name(); t == "ConfigureRequestEvent" || t ==  "MapRequestEvent"{
        ret = append(ret, winFromValue(evVal.FieldByName("Parent")))
    }
    if _, ok := evType.FieldByName("Window"); ok{
        w = winFromValue(evVal.FieldByName("Window"))
    }else if _, ok := evType.FieldByName("Event"); ok{
        w =  winFromValue(evVal.FieldByName("Event"))
    }else if _, ok := evType.FieldByName("Owner"); ok{
        w = winFromValue(evVal.FieldByName("Owner"))
    }else if _, ok := evType.FieldByName("Requestor"); ok{
        w =  winFromValue(evVal.FieldByName("Requestor"))
    }
    return append(ret, w)
}

type EventCallback func(ev xgb.Event)

type EventLoop struct{
    theEnd bool
    X *xgbutil.XUtil
    callbacks map[string] map[xproto.Window] []EventCallback
}

func NewEventLoop(X *xgbutil.XUtil) *EventLoop{
    return &EventLoop{theEnd:false, X:X, callbacks:make(map[string] map[xproto.Window] []EventCallback)}
}

func (e *EventLoop) runCallbacks(event xgb.Event){
    for _, w := range getEventWin(event){
        callbacks, ok:= e.callbacks[getEventName(event)][w]
        if ok{
            for _, c := range callbacks{
                c(event)
            }
        }
    }
}

func (e *EventLoop) Run(){
    for !e.theEnd{
        ev, err := e.X.Conn().WaitForEvent()
        if ev != nil{
            e.runCallbacks(ev)
        }
        if err != nil{
            fmt.Println("got error:", err)
        }
    }
}

func (e *EventLoop) Stop(){
    e.theEnd = true
}

func (e *EventLoop) AddCallback(eventName string, win xproto.Window, callback EventCallback){
    if _, ok:= e.callbacks[eventName]; !ok{
        e.callbacks[eventName] = make(map[xproto.Window] []EventCallback)
    }

    if _, ok:=e.callbacks[eventName][win]; !ok{
        e.callbacks[eventName][win] = make([]EventCallback, 0)
    }
    e.callbacks[eventName][win] = append(e.callbacks[eventName][win], callback)
}

func (e *EventLoop) DetachCallbacks(win xproto.Window){
    for evType := range e.callbacks{
        delete(e.callbacks[evType], win)
    }
}
