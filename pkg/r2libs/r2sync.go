package r2libs

import (
	"sync"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// MutexObject wraps a *sync.Mutex. The field must be a pointer (not sync.Mutex
// by value) because R2Lang values get copied around by the interpreter (map
// assignment, closures, etc.); a copied sync.Mutex loses its shared lock state.
type MutexObject struct {
	mu *sync.Mutex
}

func (m *MutexObject) Eval(env *r2core.Environment) interface{} {
	return m
}

func (m *MutexObject) Getattr(name string) (r2core.Node, bool) {
	switch name {
	case "lock":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			m.mu.Lock()
			return nil
		}}, true
	case "unlock":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			m.mu.Unlock()
			return nil
		}}, true
	case "tryLock":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			return m.mu.TryLock()
		}}, true
	}
	return nil, false
}

// WaitGroupObject wraps a *sync.WaitGroup for the same pointer-sharing reason
// as MutexObject above.
type WaitGroupObject struct {
	wg *sync.WaitGroup
}

func (w *WaitGroupObject) Eval(env *r2core.Environment) interface{} {
	return w
}

func (w *WaitGroupObject) Getattr(name string) (r2core.Node, bool) {
	switch name {
	case "add":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			if len(args) != 1 {
				panic("WaitGroup.add needs exactly one argument: delta")
			}
			delta, ok := args[0].(float64)
			if !ok {
				panic("WaitGroup.add: argument must be a number")
			}
			w.wg.Add(int(delta))
			return nil
		}}, true
	case "done":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			w.wg.Done()
			return nil
		}}, true
	case "wait":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			w.wg.Wait()
			return nil
		}}, true
	}
	return nil, false
}

// SemaphoreObject implements a counting semaphore on top of a buffered
// channel: acquiring is a send, releasing is a receive, and channel capacity
// is the permit count.
type SemaphoreObject struct {
	ch chan struct{}
}

func (s *SemaphoreObject) Eval(env *r2core.Environment) interface{} {
	return s
}

func (s *SemaphoreObject) Getattr(name string) (r2core.Node, bool) {
	switch name {
	case "acquire":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			s.ch <- struct{}{}
			return nil
		}}, true
	case "release":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			select {
			case <-s.ch:
			default:
				panic("Semaphore.release: no permit currently held")
			}
			return nil
		}}, true
	case "tryAcquire":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			select {
			case s.ch <- struct{}{}:
				return true
			default:
				return false
			}
		}}, true
	}
	return nil, false
}

// OnceObject wraps a *sync.Once so `.do(fn)` runs fn at most once across
// however many goroutines race to call it.
type OnceObject struct {
	once *sync.Once
}

func (o *OnceObject) Eval(env *r2core.Environment) interface{} {
	return o
}

func (o *OnceObject) Getattr(name string) (r2core.Node, bool) {
	switch name {
	case "do":
		return &NativeFunction{Fn: func(args ...interface{}) interface{} {
			if len(args) != 1 {
				panic("Once.do needs exactly one argument: a function")
			}
			fn, ok := args[0].(*r2core.UserFunction)
			if !ok {
				panic("Once.do: argument must be a function")
			}
			o.once.Do(func() {
				fn.Call()
			})
			return nil
		}}, true
	}
	return nil, false
}

func RegisterSync(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"Mutex": func(args ...interface{}) interface{} {
			return &MutexObject{mu: &sync.Mutex{}}
		},
		"WaitGroup": func(args ...interface{}) interface{} {
			return &WaitGroupObject{wg: &sync.WaitGroup{}}
		},
		"Semaphore": func(args ...interface{}) interface{} {
			if len(args) != 1 {
				panic("sync.Semaphore needs exactly one argument: permit count")
			}
			permits, ok := args[0].(float64)
			if !ok {
				panic("sync.Semaphore: argument must be a number")
			}
			if permits < 1 {
				panic("sync.Semaphore: permit count must be at least 1")
			}
			return &SemaphoreObject{ch: make(chan struct{}, int(permits))}
		},
		"Once": func(args ...interface{}) interface{} {
			return &OnceObject{once: &sync.Once{}}
		},
	}

	RegisterModule(env, "sync", functions)
}
