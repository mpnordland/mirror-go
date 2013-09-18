package main

type HookFunc func(args ...interface{})

type Signalable interface {
	AttachToHook(name string, hookfunc HookFunc)
	PullHook(name string, args ...interface{})
}

type Dispatcher struct {
	hooks map[string][]HookFunc
}

func NewDispatcher() *Dispatcher {
	return &Dispatcher{make(map[string][]HookFunc)}
}

func (d *Dispatcher) AttachToHook(name string, hookfunc HookFunc) {
	if _, ok := d.hooks[name]; !ok {
		d.hooks[name] = make([]HookFunc, 0)
	}
	d.hooks[name] = append(d.hooks[name], hookfunc)
}

func (d *Dispatcher) PullHook(name string, args ...interface{}) {
	if hookSlice, ok := d.hooks[name]; ok {
		for _, h := range hookSlice {
			h(args...)
		}
	}
}
