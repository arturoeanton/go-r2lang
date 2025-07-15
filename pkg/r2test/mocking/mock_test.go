package mocking

import (
	"testing"
)

func TestNewMock(t *testing.T) {
	mock := NewMock("test_mock")

	if mock.name != "test_mock" {
		t.Errorf("Expected mock name 'test_mock', got '%s'", mock.name)
	}

	if mock.expectations == nil {
		t.Error("Expected expectations map to be initialized")
	}

	if mock.calls == nil {
		t.Error("Expected calls slice to be initialized")
	}

	if !mock.autoRestore {
		t.Error("Expected auto restore to be enabled by default")
	}
}

func TestMockWhen(t *testing.T) {
	mock := NewMock("test_mock")

	expectation := mock.When("testFunction", "arg1", "arg2")

	if expectation.FunctionName != "testFunction" {
		t.Errorf("Expected function name 'testFunction', got '%s'", expectation.FunctionName)
	}

	if len(expectation.Args) != 2 {
		t.Errorf("Expected 2 arguments, got %d", len(expectation.Args))
	}

	if expectation.Args[0] != "arg1" || expectation.Args[1] != "arg2" {
		t.Errorf("Expected args ['arg1', 'arg2'], got %v", expectation.Args)
	}

	if expectation.MinCalls != 1 {
		t.Errorf("Expected min calls to be 1, got %d", expectation.MinCalls)
	}
}

func TestExpectationReturns(t *testing.T) {
	mock := NewMock("test_mock")
	expectation := mock.When("testFunction")

	expectation.Returns("result1", "result2")

	if len(expectation.ReturnValues) != 2 {
		t.Errorf("Expected 2 return values, got %d", len(expectation.ReturnValues))
	}

	if expectation.ReturnValues[0] != "result1" || expectation.ReturnValues[1] != "result2" {
		t.Errorf("Expected return values ['result1', 'result2'], got %v", expectation.ReturnValues)
	}
}

func TestExpectationTimes(t *testing.T) {
	mock := NewMock("test_mock")
	expectation := mock.When("testFunction")

	expectation.Times(3)

	if expectation.MinCalls != 3 {
		t.Errorf("Expected min calls to be 3, got %d", expectation.MinCalls)
	}

	if expectation.MaxCalls != 3 {
		t.Errorf("Expected max calls to be 3, got %d", expectation.MaxCalls)
	}
}

func TestMockCall(t *testing.T) {
	mock := NewMock("test_mock")
	mock.When("testFunction", "arg1").Returns("result1")

	results, err := mock.Call("testFunction", "arg1")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if len(results) != 1 {
		t.Errorf("Expected 1 result, got %d", len(results))
	}

	if results[0] != "result1" {
		t.Errorf("Expected result 'result1', got '%v'", results[0])
	}

	// Check that call was recorded
	calls := mock.GetCalls()
	if len(calls) != 1 {
		t.Errorf("Expected 1 recorded call, got %d", len(calls))
	}

	if calls[0].FunctionName != "testFunction" {
		t.Errorf("Expected function name 'testFunction', got '%s'", calls[0].FunctionName)
	}
}

func TestMockCallUnexpected(t *testing.T) {
	mock := NewMock("test_mock")

	_, err := mock.Call("unexpectedFunction")

	if err == nil {
		t.Error("Expected error for unexpected function call")
	}

	verificationErr, ok := err.(*MockVerificationError)
	if !ok {
		t.Errorf("Expected MockVerificationError, got %T", err)
	}

	expectedMsg := "unexpected call to unexpectedFunction with args []"
	if verificationErr.Message != expectedMsg {
		t.Errorf("Expected error message '%s', got '%s'", expectedMsg, verificationErr.Message)
	}
}

func TestMockVerify(t *testing.T) {
	mock := NewMock("test_mock")
	mock.When("testFunction").Returns("result")

	// Should fail verification before calling
	err := mock.Verify()
	if err == nil {
		t.Error("Expected verification to fail before calling function")
	}

	// Call the function
	mock.Call("testFunction")

	// Should pass verification after calling
	err = mock.Verify()
	if err != nil {
		t.Errorf("Unexpected verification error: %v", err)
	}
}

func TestMockVerifyTooManyCalls(t *testing.T) {
	mock := NewMock("test_mock")
	mock.When("testFunction").Returns("result").Times(1)

	// Call once (should be fine)
	mock.Call("testFunction")

	// Call again (should fail)
	_, err := mock.Call("testFunction")
	if err == nil {
		t.Error("Expected error for too many calls")
	}
}

func TestMockWasCalled(t *testing.T) {
	mock := NewMock("test_mock")
	mock.When("testFunction", "arg1").Returns("result")

	// Should return false before calling
	if mock.WasCalled("testFunction", "arg1") {
		t.Error("Expected WasCalled to return false before calling")
	}

	// Call the function
	mock.Call("testFunction", "arg1")

	// Should return true after calling
	if !mock.WasCalled("testFunction", "arg1") {
		t.Error("Expected WasCalled to return true after calling")
	}

	// Should return false for different args
	if mock.WasCalled("testFunction", "arg2") {
		t.Error("Expected WasCalled to return false for different args")
	}
}

func TestMockWasCalledTimes(t *testing.T) {
	mock := NewMock("test_mock")
	mock.When("testFunction").Returns("result").AtMost(3)

	// Should return true for 0 times initially
	if !mock.WasCalledTimes("testFunction", 0) {
		t.Error("Expected WasCalledTimes(0) to return true initially")
	}

	// Call twice
	mock.Call("testFunction")
	mock.Call("testFunction")

	// Should return true for 2 times
	if !mock.WasCalledTimes("testFunction", 2) {
		t.Error("Expected WasCalledTimes(2) to return true after 2 calls")
	}

	// Should return false for 1 time
	if mock.WasCalledTimes("testFunction", 1) {
		t.Error("Expected WasCalledTimes(1) to return false after 2 calls")
	}
}

func TestMockReset(t *testing.T) {
	mock := NewMock("test_mock")
	mock.When("testFunction").Returns("result")
	mock.Call("testFunction")

	// Verify there are calls and expectations
	if len(mock.GetCalls()) != 1 {
		t.Error("Expected 1 call before reset")
	}

	if len(mock.expectations) != 1 {
		t.Error("Expected 1 expectation before reset")
	}

	// Reset
	mock.Reset()

	// Verify everything is cleared
	if len(mock.GetCalls()) != 0 {
		t.Error("Expected 0 calls after reset")
	}

	if len(mock.expectations) != 0 {
		t.Error("Expected 0 expectations after reset")
	}
}

func TestGlobalMockManager(t *testing.T) {
	// Create a mock through the global manager
	mock := CreateMock("global_test_mock")

	if mock.name != "global_test_mock" {
		t.Errorf("Expected mock name 'global_test_mock', got '%s'", mock.name)
	}

	// Retrieve the mock
	retrievedMock, exists := GetMock("global_test_mock")
	if !exists {
		t.Error("Expected to find mock in global manager")
	}

	if retrievedMock != mock {
		t.Error("Expected retrieved mock to be the same instance")
	}

	// Clean up
	GlobalMockManager.RemoveAll()
}

func TestMockWithCallback(t *testing.T) {
	mock := NewMock("test_mock")
	callbackCalled := false

	mock.When("testFunction", "arg1").WithCallback(func(args ...interface{}) []interface{} {
		callbackCalled = true
		if len(args) != 1 || args[0] != "arg1" {
			t.Errorf("Expected callback args ['arg1'], got %v", args)
		}
		return []interface{}{"callback_result"}
	})

	results, err := mock.Call("testFunction", "arg1")

	if err != nil {
		t.Errorf("Unexpected error: %v", err)
	}

	if !callbackCalled {
		t.Error("Expected callback to be called")
	}

	if len(results) != 1 || results[0] != "callback_result" {
		t.Errorf("Expected callback result ['callback_result'], got %v", results)
	}
}
