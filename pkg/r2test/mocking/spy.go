package mocking

import (
	"fmt"
	"reflect"
	"sync"
)

// Spy represents a spy that can intercept and record function calls
// while optionally calling through to the original function
type Spy struct {
	name         string
	originalFunc interface{}
	calls        []Call
	stubs        map[string]*Stub
	mu           sync.RWMutex
	callThrough  bool
}

// Stub represents a stubbed function that replaces the original
type Stub struct {
	FunctionName   string
	Implementation func(args ...interface{}) []interface{}
	ReturnValues   []interface{}
	Error          error
	CallCount      int
}

// SpyManager manages spies and stubs
type SpyManager struct {
	spies map[string]*Spy
	stubs map[string]*Stub
	mu    sync.RWMutex
}

var GlobalSpyManager = &SpyManager{
	spies: make(map[string]*Spy),
	stubs: make(map[string]*Stub),
}

// NewSpy creates a new spy for a function
func NewSpy(name string, originalFunc interface{}) *Spy {
	return &Spy{
		name:         name,
		originalFunc: originalFunc,
		calls:        make([]Call, 0),
		stubs:        make(map[string]*Stub),
		callThrough:  false,
	}
}

// SpyOn creates a spy on an existing function
func SpyOn(functionName string, originalFunc interface{}) *Spy {
	spy := NewSpy(functionName, originalFunc)
	GlobalSpyManager.RegisterSpy(functionName, spy)
	return spy
}

// CallThrough enables calling through to the original function
func (s *Spy) CallThrough() *Spy {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.callThrough = true
	return s
}

// DontCallThrough disables calling through to the original function
func (s *Spy) DontCallThrough() *Spy {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.callThrough = false
	return s
}

// And creates a stub for this spy
func (s *Spy) And() *SpyStubBuilder {
	return &SpyStubBuilder{spy: s}
}

// SpyStubBuilder helps build stubs for spies
type SpyStubBuilder struct {
	spy *Spy
}

// ReturnValue sets the spy to return specific values
func (ssb *SpyStubBuilder) ReturnValue(values ...interface{}) *Spy {
	ssb.spy.mu.Lock()
	defer ssb.spy.mu.Unlock()

	stub := &Stub{
		FunctionName: ssb.spy.name,
		ReturnValues: values,
	}
	ssb.spy.stubs[ssb.spy.name] = stub
	return ssb.spy
}

// ReturnError sets the spy to return an error
func (ssb *SpyStubBuilder) ReturnError(err error) *Spy {
	ssb.spy.mu.Lock()
	defer ssb.spy.mu.Unlock()

	stub := &Stub{
		FunctionName: ssb.spy.name,
		Error:        err,
	}
	ssb.spy.stubs[ssb.spy.name] = stub
	return ssb.spy
}

// CallFake sets the spy to call a fake implementation
func (ssb *SpyStubBuilder) CallFake(implementation func(args ...interface{}) []interface{}) *Spy {
	ssb.spy.mu.Lock()
	defer ssb.spy.mu.Unlock()

	stub := &Stub{
		FunctionName:   ssb.spy.name,
		Implementation: implementation,
	}
	ssb.spy.stubs[ssb.spy.name] = stub
	return ssb.spy
}

// Call executes the spy, recording the call and potentially calling through
func (s *Spy) Call(args ...interface{}) ([]interface{}, error) {
	s.mu.Lock()
	defer s.mu.Unlock()

	// Check if there's a stub for this function
	if stub, exists := s.stubs[s.name]; exists {
		stub.CallCount++

		call := Call{
			FunctionName: s.name,
			Args:         args,
			ReturnValues: stub.ReturnValues,
			Error:        stub.Error,
		}
		s.calls = append(s.calls, call)

		if stub.Implementation != nil {
			returnValues := stub.Implementation(args...)
			call.ReturnValues = returnValues
			s.calls[len(s.calls)-1] = call
			return returnValues, stub.Error
		}

		return stub.ReturnValues, stub.Error
	}

	// If no stub and call through is enabled, call original function
	if s.callThrough && s.originalFunc != nil {
		returnValues, err := s.callOriginal(args...)

		call := Call{
			FunctionName: s.name,
			Args:         args,
			ReturnValues: returnValues,
			Error:        err,
		}
		s.calls = append(s.calls, call)

		return returnValues, err
	}

	// Default behavior: just record the call
	call := Call{
		FunctionName: s.name,
		Args:         args,
		ReturnValues: nil,
		Error:        nil,
	}
	s.calls = append(s.calls, call)

	return nil, nil
}

// callOriginal calls the original function using reflection
func (s *Spy) callOriginal(args ...interface{}) ([]interface{}, error) {
	if s.originalFunc == nil {
		return nil, fmt.Errorf("no original function to call")
	}

	funcValue := reflect.ValueOf(s.originalFunc)
	funcType := funcValue.Type()

	if funcType.Kind() != reflect.Func {
		return nil, fmt.Errorf("original is not a function")
	}

	// Convert args to reflect.Value
	in := make([]reflect.Value, len(args))
	for i, arg := range args {
		in[i] = reflect.ValueOf(arg)
	}

	// Call the function
	results := funcValue.Call(in)

	// Convert results to interface{}
	returnValues := make([]interface{}, len(results))
	for i, result := range results {
		returnValues[i] = result.Interface()
	}

	// Check if last return value is an error
	var err error
	if len(results) > 0 {
		if lastResult := results[len(results)-1]; lastResult.Type().Implements(reflect.TypeOf((*error)(nil)).Elem()) {
			if !lastResult.IsNil() {
				err = lastResult.Interface().(error)
				// Remove error from return values
				returnValues = returnValues[:len(returnValues)-1]
			}
		}
	}

	return returnValues, err
}

// GetCalls returns all recorded calls for this spy
func (s *Spy) GetCalls() []Call {
	s.mu.RLock()
	defer s.mu.RUnlock()

	calls := make([]Call, len(s.calls))
	copy(calls, s.calls)
	return calls
}

// WasCalled checks if the spy was called
func (s *Spy) WasCalled() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.calls) > 0
}

// WasCalledWith checks if the spy was called with specific arguments
func (s *Spy) WasCalledWith(args ...interface{}) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()

	for _, call := range s.calls {
		if reflect.DeepEqual(call.Args, args) {
			return true
		}
	}
	return false
}

// WasCalledTimes checks if the spy was called exactly n times
func (s *Spy) WasCalledTimes(times int) bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.calls) == times
}

// Reset clears all recorded calls
func (s *Spy) Reset() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.calls = make([]Call, 0)
	for _, stub := range s.stubs {
		stub.CallCount = 0
	}
}

// Restore removes the spy and restores original function
func (s *Spy) Restore() {
	GlobalSpyManager.RemoveSpy(s.name)
}

// SpyManager methods

// RegisterSpy registers a spy with the manager
func (sm *SpyManager) RegisterSpy(name string, spy *Spy) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.spies[name] = spy
}

// GetSpy retrieves a spy by name
func (sm *SpyManager) GetSpy(name string) (*Spy, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	spy, exists := sm.spies[name]
	return spy, exists
}

// RemoveSpy removes a spy from the manager
func (sm *SpyManager) RemoveSpy(name string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	delete(sm.spies, name)
}

// CreateStub creates a standalone stub function
func (sm *SpyManager) CreateStub(name string) *Stub {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	stub := &Stub{
		FunctionName: name,
		CallCount:    0,
	}
	sm.stubs[name] = stub
	return stub
}

// GetStub retrieves a stub by name
func (sm *SpyManager) GetStub(name string) (*Stub, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	stub, exists := sm.stubs[name]
	return stub, exists
}

// ResetAll resets all spies and stubs
func (sm *SpyManager) ResetAll() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	for _, spy := range sm.spies {
		spy.Reset()
	}

	for _, stub := range sm.stubs {
		stub.CallCount = 0
	}
}

// RestoreAll restores all spies
func (sm *SpyManager) RestoreAll() {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Clear all spies (they should restore their original functions)
	sm.spies = make(map[string]*Spy)
	sm.stubs = make(map[string]*Stub)
}

// Stub methods

// Returns sets the return values for a stub
func (s *Stub) Returns(values ...interface{}) *Stub {
	s.ReturnValues = values
	return s
}

// ReturnsError sets an error for a stub
func (s *Stub) ReturnsError(err error) *Stub {
	s.Error = err
	return s
}

// CallsThrough sets a custom implementation for a stub
func (s *Stub) CallsThrough(implementation func(args ...interface{}) []interface{}) *Stub {
	s.Implementation = implementation
	return s
}

// Call executes the stub
func (s *Stub) Call(args ...interface{}) ([]interface{}, error) {
	s.CallCount++

	if s.Implementation != nil {
		return s.Implementation(args...), s.Error
	}

	return s.ReturnValues, s.Error
}

// Global convenience functions

// CreateStub creates a new stub using the global manager
func CreateStub(name string) *Stub {
	return GlobalSpyManager.CreateStub(name)
}

// GetSpy retrieves a spy by name using the global manager
func GetSpy(name string) (*Spy, bool) {
	return GlobalSpyManager.GetSpy(name)
}

// GetStub retrieves a stub by name using the global manager
func GetStub(name string) (*Stub, bool) {
	return GlobalSpyManager.GetStub(name)
}

// ResetAllSpies resets all spies and stubs using the global manager
func ResetAllSpies() {
	GlobalSpyManager.ResetAll()
}

// RestoreAllSpies restores all spies using the global manager
func RestoreAllSpies() {
	GlobalSpyManager.RestoreAll()
}
