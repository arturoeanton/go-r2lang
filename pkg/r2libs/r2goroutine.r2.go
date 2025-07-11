package r2libs

import (
	"sync"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
)

// Semaphore estructura
type Semaphore struct {
	ch chan struct{}
}

// NewSemaphore crea un nuevo semáforo con el número dado de permisos
func NewSemaphore(permits int) *Semaphore {
	return &Semaphore{
		ch: make(chan struct{}, permits),
	}
}

// Acquire obtiene un permiso del semáforo
func (s *Semaphore) Acquire() {
	s.ch <- struct{}{}
}

// Release libera un permiso del semáforo
func (s *Semaphore) Release() {
	<-s.ch
}

// Builtin function para crear un semáforo
func builtinSemaphore(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("semaphore necesita exactamente un argumento: número de permisos")
	}
	permitCount, ok := args[0].(float64)
	if !ok {
		panic("semaphore: el argumento debe ser un número")
	}
	if permitCount < 1 {
		panic("semaphore: el número de permisos debe ser al menos 1")
	}
	sem := NewSemaphore(int(permitCount))
	return sem
}

// Builtin functions para adquirir y liberar semáforos
func builtinAcquireSemaphore(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("acquire necesita exactamente un argumento: semáforo")
	}
	sem, ok := args[0].(*Semaphore)
	if !ok {
		panic("acquire: el argumento debe ser un semáforo")
	}
	sem.Acquire()
	return nil
}

func builtinReleaseSemaphore(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("release necesita exactamente un argumento: semáforo")
	}
	sem, ok := args[0].(*Semaphore)
	if !ok {
		panic("release: el argumento debe ser un semáforo")
	}
	sem.Release()
	return nil
}

// r2goroutine.go (continuación)

// Monitor estructura
type Monitor struct {
	mutex sync.Mutex
	cond  *sync.Cond
}

// NewMonitor crea un nuevo monitor
func NewMonitor() *Monitor {
	m := &Monitor{}
	m.cond = sync.NewCond(&m.mutex)
	return m
}

// Lock adquiere el mutex del monitor
func (m *Monitor) Lock() {
	m.mutex.Lock()
}

// Unlock libera el mutex del monitor
func (m *Monitor) Unlock() {
	m.mutex.Unlock()
}

// Wait espera en la condición del monitor
func (m *Monitor) Wait() {
	m.cond.Wait()
}

// Signal despierta una goroutine esperando en la condición
func (m *Monitor) Signal() {
	m.cond.Signal()
}

// Broadcast despierta todas las goroutines esperando en la condición
func (m *Monitor) Broadcast() {
	m.cond.Broadcast()
}

// Builtin function para crear un monitor
func builtinMonitor(args ...interface{}) interface{} {
	if len(args) != 0 {
		panic("monitor no necesita argumentos")
	}
	mon := NewMonitor()
	return mon
}

// Builtin functions para operar sobre monitores
func builtinLock(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("lock necesita exactamente un argumento: monitor")
	}
	mon, ok := args[0].(*Monitor)
	if !ok {
		panic("lock: el argumento debe ser un monitor")
	}
	mon.Lock()
	return nil
}

func builtinUnlock(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("unlock necesita exactamente un argumento: monitor")
	}
	mon, ok := args[0].(*Monitor)
	if !ok {
		panic("unlock: el argumento debe ser un monitor")
	}
	mon.Unlock()
	return nil
}

func builtinWait(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("wait necesita exactamente un argumento: monitor")
	}
	mon, ok := args[0].(*Monitor)
	if !ok {
		panic("wait: el argumento debe ser un monitor")
	}
	mon.Wait()
	return nil
}

func builtinSignal(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("signal necesita exactamente un argumento: monitor")
	}
	mon, ok := args[0].(*Monitor)
	if !ok {
		panic("signal: el argumento debe ser un monitor")
	}
	mon.Signal()
	return nil
}

func builtinBroadcast(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("broadcast necesita exactamente un argumento: monitor")
	}
	mon, ok := args[0].(*Monitor)
	if !ok {
		panic("broadcast: el argumento debe ser un monitor")
	}
	mon.Broadcast()
	return nil
}

func builtinWaitAll(args ...interface{}) interface{} {
	if len(args) != 1 {
		panic("waitAll necesita exactamente un argumento: monitor o semáforo")
	}
	mon, ok := args[0].(*Monitor)
	if !ok {
		sem, ok := args[0].(*Semaphore)
		if !ok {
			panic("waitAll: el argumento debe ser un monitor o semáforo")
		}
		sem.Acquire()
		return nil
	}
	mon.Broadcast()
	return nil
}

// r2goroutine.go (continuación)

func RegisterConcurrency(env *r2core.Environment) {
	// Semáforos
	env.Set("semaphore", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinSemaphore(args...)
	}))

	env.Set("acquire", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinAcquireSemaphore(args...)
	}))

	env.Set("release", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinReleaseSemaphore(args...)
	}))

	// Monitores
	env.Set("monitor", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinMonitor(args...)
	}))
	env.Set("lock", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinLock(args...)
	}))
	env.Set("unlock", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinUnlock(args...)
	}))
	env.Set("wait", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinWait(args...)
	}))
	env.Set("signal", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinSignal(args...)
	}))
	env.Set("broadcast", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinBroadcast(args...)
	}))
	env.Set("waitAll", r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinWaitAll(args...)
	}))
}
