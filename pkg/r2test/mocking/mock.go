package mocking

import (
	"fmt"
	"reflect"
	"sync"
)

// Mock represents a mock object that can intercept function calls
type Mock struct {
	name         string
	expectations map[string]*Expectation
	calls        []Call
	mu           sync.RWMutex
	autoRestore  bool
	originalRefs map[string]interface{}
}

// Expectation defines what should happen when a mocked function is called
type Expectation struct {
	FunctionName string
	Args         []interface{}
	ReturnValues []interface{}
	CallCount    int
	MaxCalls     int
	MinCalls     int
	Error        error
	Callback     func(args ...interface{}) []interface{}
}

// Call represents a recorded function call
type Call struct {
	FunctionName string
	Args         []interface{}
	ReturnValues []interface{}
	Error        error
}

// MockVerificationError represents an error during mock verification
type MockVerificationError struct {
	Message string
}

func (e *MockVerificationError) Error() string {
	return e.Message
}

// NewMock creates a new mock object
func NewMock(name string) *Mock {
	return &Mock{
		name:         name,
		expectations: make(map[string]*Expectation),
		calls:        make([]Call, 0),
		autoRestore:  true,
		originalRefs: make(map[string]interface{}),
	}
}

// When sets up an expectation for a function call
func (m *Mock) When(functionName string, args ...interface{}) *Expectation {
	m.mu.Lock()
	defer m.mu.Unlock()

	expectation := &Expectation{
		FunctionName: functionName,
		Args:         args,
		CallCount:    0,
		MaxCalls:     -1, // unlimited by default
		MinCalls:     1,  // at least one call expected
	}

	m.expectations[functionName] = expectation
	return expectation
}

// Returns sets the return values for an expectation
func (e *Expectation) Returns(values ...interface{}) *Expectation {
	e.ReturnValues = values
	return e
}

// ReturnsError sets an error to be returned by the expectation
func (e *Expectation) ReturnsError(err error) *Expectation {
	e.Error = err
	return e
}

// Times sets the expected number of calls (exact)
func (e *Expectation) Times(count int) *Expectation {
	e.MinCalls = count
	e.MaxCalls = count
	return e
}

// AtLeast sets the minimum number of expected calls
func (e *Expectation) AtLeast(count int) *Expectation {
	e.MinCalls = count
	return e
}

// AtMost sets the maximum number of expected calls
func (e *Expectation) AtMost(count int) *Expectation {
	e.MaxCalls = count
	return e
}

// WithCallback sets a callback function to execute when the mock is called
func (e *Expectation) WithCallback(callback func(args ...interface{}) []interface{}) *Expectation {
	e.Callback = callback
	return e
}

// Call records a function call and returns the expected values
func (m *Mock) Call(functionName string, args ...interface{}) ([]interface{}, error) {
	m.mu.Lock()
	defer m.mu.Unlock()

	expectation, exists := m.expectations[functionName]
	if !exists {
		return nil, &MockVerificationError{
			Message: fmt.Sprintf("unexpected call to %s with args %v", functionName, args),
		}
	}

	// Check if we've exceeded max calls
	if expectation.MaxCalls > 0 && expectation.CallCount >= expectation.MaxCalls {
		return nil, &MockVerificationError{
			Message: fmt.Sprintf("function %s called too many times (max: %d)", functionName, expectation.MaxCalls),
		}
	}

	expectation.CallCount++

	var returnValues []interface{}
	var err error

	// Execute callback if provided
	if expectation.Callback != nil {
		returnValues = expectation.Callback(args...)
	} else {
		returnValues = expectation.ReturnValues
	}

	if expectation.Error != nil {
		err = expectation.Error
	}

	// Record the call
	call := Call{
		FunctionName: functionName,
		Args:         args,
		ReturnValues: returnValues,
		Error:        err,
	}
	m.calls = append(m.calls, call)

	return returnValues, err
}

// Verify checks that all expectations have been met
func (m *Mock) Verify() error {
	m.mu.RLock()
	defer m.mu.RUnlock()

	for funcName, expectation := range m.expectations {
		if expectation.CallCount < expectation.MinCalls {
			return &MockVerificationError{
				Message: fmt.Sprintf("function %s was called %d times, expected at least %d",
					funcName, expectation.CallCount, expectation.MinCalls),
			}
		}

		if expectation.MaxCalls > 0 && expectation.CallCount > expectation.MaxCalls {
			return &MockVerificationError{
				Message: fmt.Sprintf("function %s was called %d times, expected at most %d",
					funcName, expectation.CallCount, expectation.MaxCalls),
			}
		}
	}

	return nil
}

// Reset clears all expectations and recorded calls
func (m *Mock) Reset() {
	m.mu.Lock()
	defer m.mu.Unlock()

	m.expectations = make(map[string]*Expectation)
	m.calls = make([]Call, 0)
}

// GetCalls returns all recorded calls
func (m *Mock) GetCalls() []Call {
	m.mu.RLock()
	defer m.mu.RUnlock()

	calls := make([]Call, len(m.calls))
	copy(calls, m.calls)
	return calls
}

// GetCallsFor returns all calls for a specific function
func (m *Mock) GetCallsFor(functionName string) []Call {
	m.mu.RLock()
	defer m.mu.RUnlock()

	var functionCalls []Call
	for _, call := range m.calls {
		if call.FunctionName == functionName {
			functionCalls = append(functionCalls, call)
		}
	}
	return functionCalls
}

// WasCalled checks if a function was called with specific arguments
func (m *Mock) WasCalled(functionName string, args ...interface{}) bool {
	calls := m.GetCallsFor(functionName)

	if len(args) == 0 {
		return len(calls) > 0
	}

	for _, call := range calls {
		if reflect.DeepEqual(call.Args, args) {
			return true
		}
	}
	return false
}

// WasCalledTimes checks if a function was called exactly n times
func (m *Mock) WasCalledTimes(functionName string, times int) bool {
	calls := m.GetCallsFor(functionName)
	return len(calls) == times
}

// Restore restores original function references if auto-restore is enabled
func (m *Mock) Restore() {
	if !m.autoRestore {
		return
	}

	m.mu.Lock()
	defer m.mu.Unlock()

	for name, originalRef := range m.originalRefs {
		// This would restore the original function reference
		// Implementation depends on how functions are stored in R2Lang
		_ = name
		_ = originalRef
	}

	m.originalRefs = make(map[string]interface{})
}

// MockManager manages multiple mocks and provides global operations
type MockManager struct {
	mocks map[string]*Mock
	mu    sync.RWMutex
}

// GlobalMockManager is the global instance
var GlobalMockManager = &MockManager{
	mocks: make(map[string]*Mock),
}

// CreateMock creates a new mock and registers it with the manager
func (mm *MockManager) CreateMock(name string) *Mock {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	mock := NewMock(name)
	mm.mocks[name] = mock
	return mock
}

// GetMock retrieves a mock by name
func (mm *MockManager) GetMock(name string) (*Mock, bool) {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	mock, exists := mm.mocks[name]
	return mock, exists
}

// VerifyAll verifies all registered mocks
func (mm *MockManager) VerifyAll() error {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	for name, mock := range mm.mocks {
		if err := mock.Verify(); err != nil {
			return fmt.Errorf("mock '%s' verification failed: %w", name, err)
		}
	}
	return nil
}

// ResetAll resets all registered mocks
func (mm *MockManager) ResetAll() {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	for _, mock := range mm.mocks {
		mock.Reset()
	}
}

// RestoreAll restores all mocks
func (mm *MockManager) RestoreAll() {
	mm.mu.RLock()
	defer mm.mu.RUnlock()

	for _, mock := range mm.mocks {
		mock.Restore()
	}
}

// RemoveAll removes all mocks from the manager
func (mm *MockManager) RemoveAll() {
	mm.mu.Lock()
	defer mm.mu.Unlock()

	mm.mocks = make(map[string]*Mock)
}

// Global convenience functions

// CreateMock creates a new mock using the global manager
func CreateMock(name string) *Mock {
	return GlobalMockManager.CreateMock(name)
}

// GetMock retrieves a mock by name using the global manager
func GetMock(name string) (*Mock, bool) {
	return GlobalMockManager.GetMock(name)
}

// VerifyAllMocks verifies all mocks using the global manager
func VerifyAllMocks() error {
	return GlobalMockManager.VerifyAll()
}

// ResetAllMocks resets all mocks using the global manager
func ResetAllMocks() {
	GlobalMockManager.ResetAll()
}

// RestoreAllMocks restores all mocks using the global manager
func RestoreAllMocks() {
	GlobalMockManager.RestoreAll()
}
