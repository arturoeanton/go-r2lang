package mocking

import (
	"fmt"
	"sync"
)

// TestIsolation provides isolated environments for test execution
type TestIsolation struct {
	contexts map[string]*IsolationContext
	mu       sync.RWMutex
}

// IsolationContext represents an isolated test execution context
type IsolationContext struct {
	ID           string
	TestName     string
	Mocks        map[string]*Mock
	Spies        map[string]*Spy
	Stubs        map[string]*Stub
	Environment  interface{} // Will be *r2core.Environment when integrated
	GlobalState  map[string]interface{}
	Interceptors map[string]*FunctionInterceptor
	IsActive     bool
}

// FunctionInterceptor intercepts function calls within an isolation context
type FunctionInterceptor struct {
	OriginalFunction interface{}
	InterceptorFunc  func(args ...interface{}) ([]interface{}, error)
	IsEnabled        bool
}

// GlobalTestIsolation is the global isolation manager
var GlobalTestIsolation = &TestIsolation{
	contexts: make(map[string]*IsolationContext),
}

// NewTestIsolation creates a new test isolation manager
func NewTestIsolation() *TestIsolation {
	return &TestIsolation{
		contexts: make(map[string]*IsolationContext),
	}
}

// CreateContext creates a new isolation context for a test
func (ti *TestIsolation) CreateContext(testName string) *IsolationContext {
	ti.mu.Lock()
	defer ti.mu.Unlock()

	contextID := fmt.Sprintf("test_%s_%d", testName, len(ti.contexts))

	context := &IsolationContext{
		ID:           contextID,
		TestName:     testName,
		Mocks:        make(map[string]*Mock),
		Spies:        make(map[string]*Spy),
		Stubs:        make(map[string]*Stub),
		GlobalState:  make(map[string]interface{}),
		Interceptors: make(map[string]*FunctionInterceptor),
		IsActive:     true,
	}

	ti.contexts[contextID] = context
	return context
}

// GetContext retrieves an isolation context by ID
func (ti *TestIsolation) GetContext(contextID string) (*IsolationContext, bool) {
	ti.mu.RLock()
	defer ti.mu.RUnlock()

	context, exists := ti.contexts[contextID]
	return context, exists
}

// ActivateContext activates an isolation context
func (ti *TestIsolation) ActivateContext(contextID string) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()

	context, exists := ti.contexts[contextID]
	if !exists {
		return fmt.Errorf("context not found: %s", contextID)
	}

	// Deactivate all other contexts
	for _, ctx := range ti.contexts {
		ctx.IsActive = false
	}

	// Activate this context
	context.IsActive = true

	// Apply all interceptors in this context
	for _, interceptor := range context.Interceptors {
		interceptor.IsEnabled = true
	}

	return nil
}

// DeactivateContext deactivates an isolation context
func (ti *TestIsolation) DeactivateContext(contextID string) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()

	context, exists := ti.contexts[contextID]
	if !exists {
		return fmt.Errorf("context not found: %s", contextID)
	}

	context.IsActive = false

	// Disable all interceptors in this context
	for _, interceptor := range context.Interceptors {
		interceptor.IsEnabled = false
	}

	return nil
}

// CleanupContext cleans up and removes an isolation context
func (ti *TestIsolation) CleanupContext(contextID string) error {
	ti.mu.Lock()
	defer ti.mu.Unlock()

	context, exists := ti.contexts[contextID]
	if !exists {
		return fmt.Errorf("context not found: %s", contextID)
	}

	// Restore all mocks and spies
	for _, mock := range context.Mocks {
		mock.Restore()
	}

	for _, spy := range context.Spies {
		spy.Restore()
	}

	// Disable all interceptors
	for _, interceptor := range context.Interceptors {
		interceptor.IsEnabled = false
	}

	// Remove context
	delete(ti.contexts, contextID)

	return nil
}

// GetActiveContext returns the currently active context
func (ti *TestIsolation) GetActiveContext() *IsolationContext {
	ti.mu.RLock()
	defer ti.mu.RUnlock()

	for _, context := range ti.contexts {
		if context.IsActive {
			return context
		}
	}

	return nil
}

// IsolationContext methods

// CreateMock creates a mock within this isolation context
func (ic *IsolationContext) CreateMock(name string) *Mock {
	mock := NewMock(name)
	ic.Mocks[name] = mock
	return mock
}

// CreateSpy creates a spy within this isolation context
func (ic *IsolationContext) CreateSpy(name string, originalFunc interface{}) *Spy {
	spy := NewSpy(name, originalFunc)
	ic.Spies[name] = spy
	return spy
}

// CreateStub creates a stub within this isolation context
func (ic *IsolationContext) CreateStub(name string) *Stub {
	stub := &Stub{
		FunctionName: name,
		CallCount:    0,
	}
	ic.Stubs[name] = stub
	return stub
}

// SetGlobalVariable sets a global variable within this context
func (ic *IsolationContext) SetGlobalVariable(name string, value interface{}) {
	ic.GlobalState[name] = value
}

// GetGlobalVariable gets a global variable from this context
func (ic *IsolationContext) GetGlobalVariable(name string) (interface{}, bool) {
	value, exists := ic.GlobalState[name]
	return value, exists
}

// InterceptFunction intercepts calls to a function within this context
func (ic *IsolationContext) InterceptFunction(name string, originalFunc interface{}, interceptor func(args ...interface{}) ([]interface{}, error)) {
	ic.Interceptors[name] = &FunctionInterceptor{
		OriginalFunction: originalFunc,
		InterceptorFunc:  interceptor,
		IsEnabled:        ic.IsActive,
	}
}

// RemoveInterceptor removes a function interceptor
func (ic *IsolationContext) RemoveInterceptor(name string) {
	delete(ic.Interceptors, name)
}

// CallIntercepted calls a function through its interceptor if available
func (ic *IsolationContext) CallIntercepted(name string, args ...interface{}) ([]interface{}, error, bool) {
	interceptor, exists := ic.Interceptors[name]
	if !exists || !interceptor.IsEnabled {
		return nil, nil, false
	}

	result, err := interceptor.InterceptorFunc(args...)
	return result, err, true
}

// Reset resets all mocks, spies, and stubs in this context
func (ic *IsolationContext) Reset() {
	for _, mock := range ic.Mocks {
		mock.Reset()
	}

	for _, spy := range ic.Spies {
		spy.Reset()
	}

	for _, stub := range ic.Stubs {
		stub.CallCount = 0
	}

	ic.GlobalState = make(map[string]interface{})
}

// Verify verifies all mocks in this context
func (ic *IsolationContext) Verify() error {
	for name, mock := range ic.Mocks {
		if err := mock.Verify(); err != nil {
			return fmt.Errorf("mock '%s' verification failed in context '%s': %w", name, ic.ID, err)
		}
	}
	return nil
}

// IsolatedTestRunner provides isolated test execution
type IsolatedTestRunner struct {
	isolation *TestIsolation
}

// NewIsolatedTestRunner creates a new isolated test runner
func NewIsolatedTestRunner() *IsolatedTestRunner {
	return &IsolatedTestRunner{
		isolation: NewTestIsolation(),
	}
}

// RunIsolated runs a test function in an isolated context
func (itr *IsolatedTestRunner) RunIsolated(testName string, testFunc func(*IsolationContext) error) error {
	// Create isolation context
	context := itr.isolation.CreateContext(testName)

	// Activate context
	if err := itr.isolation.ActivateContext(context.ID); err != nil {
		return fmt.Errorf("failed to activate context: %w", err)
	}

	// Defer cleanup
	defer func() {
		itr.isolation.CleanupContext(context.ID)
	}()

	// Run test function
	if err := testFunc(context); err != nil {
		return err
	}

	// Verify mocks
	if err := context.Verify(); err != nil {
		return err
	}

	return nil
}

// Global convenience functions

// CreateIsolationContext creates a new isolation context using the global manager
func CreateIsolationContext(testName string) *IsolationContext {
	return GlobalTestIsolation.CreateContext(testName)
}

// ActivateIsolationContext activates an isolation context using the global manager
func ActivateIsolationContext(contextID string) error {
	return GlobalTestIsolation.ActivateContext(contextID)
}

// DeactivateIsolationContext deactivates an isolation context using the global manager
func DeactivateIsolationContext(contextID string) error {
	return GlobalTestIsolation.DeactivateContext(contextID)
}

// CleanupIsolationContext cleans up an isolation context using the global manager
func CleanupIsolationContext(contextID string) error {
	return GlobalTestIsolation.CleanupContext(contextID)
}

// GetActiveIsolationContext returns the currently active context using the global manager
func GetActiveIsolationContext() *IsolationContext {
	return GlobalTestIsolation.GetActiveContext()
}

// RunInIsolation runs a test function in an isolated context
func RunInIsolation(testName string, testFunc func(*IsolationContext) error) error {
	runner := NewIsolatedTestRunner()
	return runner.RunIsolated(testName, testFunc)
}
