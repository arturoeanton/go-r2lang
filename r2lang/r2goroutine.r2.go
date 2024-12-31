package r2lang

import (
	"sync"
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
func builtinSemaphore(env *Environment, args ...interface{}) interface{} {
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
func builtinAcquireSemaphore(env *Environment, args ...interface{}) interface{} {
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

func builtinReleaseSemaphore(env *Environment, args ...interface{}) interface{} {
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
func builtinMonitor(env *Environment, args ...interface{}) interface{} {
	if len(args) != 0 {
		panic("monitor no necesita argumentos")
	}
	mon := NewMonitor()
	return mon
}

// Builtin functions para operar sobre monitores
func builtinLock(env *Environment, args ...interface{}) interface{} {
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

func builtinUnlock(env *Environment, args ...interface{}) interface{} {
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

func builtinWait(env *Environment, args ...interface{}) interface{} {
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

func builtinSignal(env *Environment, args ...interface{}) interface{} {
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

func builtinBroadcast(env *Environment, args ...interface{}) interface{} {
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

func builtinWaitAll(env *Environment, args ...interface{}) interface{} {
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

func RegisterConcurrency(env *Environment) {
	// Semáforos
	env.Set("semaphore", BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinSemaphore(env, args...)
	}))

	env.Set("acquire", BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinAcquireSemaphore(env, args...)
	}))

	env.Set("release", BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinReleaseSemaphore(env, args...)
	}))

	// Monitores
	env.Set("monitor", BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinMonitor(env, args...)
	}))
	env.Set("lock", BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinLock(env, args...)
	}))
	env.Set("unlock", BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinUnlock(env, args...)
	}))
	env.Set("wait", BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinWait(env, args...)
	}))
	env.Set("signal", BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinSignal(env, args...)
	}))
	env.Set("broadcast", BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinBroadcast(env, args...)
	}))
	env.Set("waitAll", BuiltinFunction(func(args ...interface{}) interface{} {
		return builtinWaitAll(env, args...)
	}))
}
