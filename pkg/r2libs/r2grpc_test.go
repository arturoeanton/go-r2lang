package r2libs

import (
	"context"
	"fmt"
	"io"
	"io/ioutil"
	"net"
	"os"
	"path/filepath"
	"sync"
	"testing"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"google.golang.org/grpc"
)

// Test proto file content
const testProtoContent = `syntax = "proto3";

package testservice;

service TestService {
    rpc GetUser(GetUserRequest) returns (GetUserResponse);
    rpc ListUsers(ListUsersRequest) returns (stream GetUserResponse);
    rpc CreateUsers(stream CreateUserRequest) returns (CreateUserResponse);
    rpc Chat(stream ChatMessage) returns (stream ChatMessage);
}

message GetUserRequest {
    string user_id = 1;
}

message GetUserResponse {
    string user_id = 1;
    string name = 2;
    string email = 3;
    int32 age = 4;
    bool active = 5;
}

message ListUsersRequest {
    int32 page_size = 1;
    string page_token = 2;
}

message CreateUserRequest {
    string name = 1;
    string email = 2;
    int32 age = 3;
}

message CreateUserResponse {
    int32 created_count = 1;
    string message = 2;
}

message ChatMessage {
    string user_id = 1;
    string message = 2;
    int64 timestamp = 3;
}

message RepeatedFieldsTest {
    repeated string tags = 1;
    repeated int32 scores = 2;
    repeated GetUserResponse users = 3;
    map<string, string> labels = 4;
    map<string, int32> counts = 5;
}`

// Mock gRPC server for testing
type mockTestServer struct {
	users map[string]*dynamic.Message
	fd    *desc.FileDescriptor
}

func newMockTestServer(fd *desc.FileDescriptor) *mockTestServer {
	return &mockTestServer{
		users: make(map[string]*dynamic.Message),
		fd:    fd,
	}
}

func (s *mockTestServer) GetUser(ctx context.Context, req *dynamic.Message) (*dynamic.Message, error) {
	userID := req.GetFieldByName("user_id").(string)

	if user, exists := s.users[userID]; exists {
		return user, nil
	}

	// Create a mock user
	responseDesc := s.fd.FindMessage("testservice.GetUserResponse")
	response := dynamic.NewMessage(responseDesc)
	response.SetFieldByName("user_id", userID)
	response.SetFieldByName("name", "Test User")
	response.SetFieldByName("email", "test@example.com")
	response.SetFieldByName("age", int32(25))
	response.SetFieldByName("active", true)

	s.users[userID] = response
	return response, nil
}

func (s *mockTestServer) ListUsers(req *dynamic.Message, stream grpc.ServerStream) error {
	// Send multiple users
	for i := 0; i < 3; i++ {
		responseDesc := s.fd.FindMessage("testservice.GetUserResponse")
		response := dynamic.NewMessage(responseDesc)
		response.SetFieldByName("user_id", fmt.Sprintf("user_%d", i))
		response.SetFieldByName("name", fmt.Sprintf("User %d", i))
		response.SetFieldByName("email", fmt.Sprintf("user%d@example.com", i))
		response.SetFieldByName("age", int32(20+i))
		response.SetFieldByName("active", true)

		if err := stream.SendMsg(response); err != nil {
			return err
		}
	}
	return nil
}

func (s *mockTestServer) CreateUsers(stream grpc.ServerStream) error {
	count := 0
	for {
		requestDesc := s.fd.FindMessage("testservice.CreateUserRequest")
		req := dynamic.NewMessage(requestDesc)

		err := stream.RecvMsg(req)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		count++
	}

	responseDesc := s.fd.FindMessage("testservice.CreateUserResponse")
	response := dynamic.NewMessage(responseDesc)
	response.SetFieldByName("created_count", int32(count))
	response.SetFieldByName("message", "Users created successfully")

	return stream.SendMsg(response)
}

func (s *mockTestServer) Chat(stream grpc.ServerStream) error {
	for {
		messageDesc := s.fd.FindMessage("testservice.ChatMessage")
		msg := dynamic.NewMessage(messageDesc)

		err := stream.RecvMsg(msg)
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}

		// Echo the message back
		if err := stream.SendMsg(msg); err != nil {
			return err
		}
	}
	return nil
}

// Test setup functions
func createTestProtoFile(tb testing.TB) string {
	tmpDir, err := ioutil.TempDir("", "r2grpc_test")
	if err != nil {
		tb.Fatalf("Failed to create temp dir: %v", err)
	}

	protoFile := filepath.Join(tmpDir, "test.proto")
	err = ioutil.WriteFile(protoFile, []byte(testProtoContent), 0644)
	if err != nil {
		tb.Fatalf("Failed to write proto file: %v", err)
	}

	return protoFile
}

func startMockServer(t *testing.T, fd *desc.FileDescriptor) (*grpc.Server, string) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}

	server := grpc.NewServer()
	_ = newMockTestServer(fd)

	// Register the mock server dynamically
	serviceDesc := fd.FindService("testservice.TestService")
	if serviceDesc == nil {
		t.Fatalf("Service not found in proto file")
	}

	go server.Serve(listener)

	return server, listener.Addr().String()
}

// startFullMockServer registers a real, working grpc.ServiceDesc for
// TestService (unlike startMockServer, whose mockTestServer is constructed
// and discarded without ever being registered on the grpc.Server), so
// client-streaming and bidi-streaming calls actually reach a handler that
// exercises the real wire protocol end-to-end.
func startFullMockServer(t *testing.T, fd *desc.FileDescriptor) (*grpc.Server, string) {
	listener, err := net.Listen("tcp", "localhost:0")
	if err != nil {
		t.Fatalf("Failed to create listener: %v", err)
	}

	impl := newMockTestServer(fd)
	server := grpc.NewServer()

	sd := grpc.ServiceDesc{
		ServiceName: "testservice.TestService",
		HandlerType: (*any)(nil),
		Streams: []grpc.StreamDesc{
			{
				StreamName:    "CreateUsers",
				ClientStreams: true,
				Handler: func(srv interface{}, stream grpc.ServerStream) error {
					return impl.CreateUsers(stream)
				},
			},
			{
				StreamName:    "Chat",
				ClientStreams: true,
				ServerStreams: true,
				Handler: func(srv interface{}, stream grpc.ServerStream) error {
					return impl.Chat(stream)
				},
			},
		},
	}

	server.RegisterService(&sd, impl)
	go server.Serve(listener)

	return server, listener.Addr().String()
}

func TestGRPCClientCreation(t *testing.T) {
	protoFile := createTestProtoFile(t)
	defer os.RemoveAll(filepath.Dir(protoFile))

	// Test successful client creation
	client, err := createGRPCClient(protoFile, "localhost:50051", nil)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}

	if client.ProtoFile != protoFile {
		t.Errorf("Expected ProtoFile %s, got %s", protoFile, client.ProtoFile)
	}

	if client.ServerAddr != "localhost:50051" {
		t.Errorf("Expected ServerAddr localhost:50051, got %s", client.ServerAddr)
	}

	// Test that services are parsed correctly
	if len(client.Services) != 1 {
		t.Errorf("Expected 1 service, got %d", len(client.Services))
	}

	service, exists := client.Services["TestService"]
	if !exists {
		t.Errorf("TestService not found in parsed services")
	}

	if len(service.Methods) != 4 {
		t.Errorf("Expected 4 methods, got %d", len(service.Methods))
	}

	// Test method parsing
	expectedMethods := []string{"GetUser", "ListUsers", "CreateUsers", "Chat"}
	for _, methodName := range expectedMethods {
		if _, exists := service.Methods[methodName]; !exists {
			t.Errorf("Method %s not found", methodName)
		}
	}
}

func TestGRPCClientRegistration(t *testing.T) {
	env := r2core.NewEnvironment()
	RegisterGRPC(env)

	// Test that grpc module is registered
	grpcModule, exists := env.Get("grpc")
	if !exists {
		t.Errorf("grpc module not registered")
	}

	grpcMap, ok := grpcModule.(map[string]interface{})
	if !ok {
		t.Errorf("grpc is not a module map")
	}

	// Test that grpcClient function is in the grpc module
	grpcClientFunc, exists := grpcMap["grpcClient"]
	if !exists {
		t.Errorf("grpcClient function not registered in grpc module")
	}

	if _, ok := grpcClientFunc.(r2core.BuiltinFunction); !ok {
		t.Errorf("grpcClient is not a BuiltinFunction")
	}
}

func TestGRPCClientMapGeneration(t *testing.T) {
	protoFile := createTestProtoFile(t)
	defer os.RemoveAll(filepath.Dir(protoFile))

	client, err := createGRPCClient(protoFile, "localhost:50051", nil)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}

	clientMap := grpcClientToMap(client)

	// Test that all expected methods are present
	expectedMethods := []string{
		"listServices", "listMethods", "getMethodInfo", "call", "callSimple",
		"callServerStream", "callClientStream", "callBidirectionalStream",
		"setTimeout", "setMetadata", "getMetadata", "setTLSConfig", "setAuth",
		"setCompression", "close",
	}

	for _, method := range expectedMethods {
		if _, exists := clientMap[method]; !exists {
			t.Errorf("Method %s not found in client map", method)
		}
	}

	// Test properties
	if clientMap["protoFile"] != protoFile {
		t.Errorf("Expected protoFile %s, got %v", protoFile, clientMap["protoFile"])
	}

	if clientMap["serverAddr"] != "localhost:50051" {
		t.Errorf("Expected serverAddr localhost:50051, got %v", clientMap["serverAddr"])
	}
}

func TestListServices(t *testing.T) {
	protoFile := createTestProtoFile(t)
	defer os.RemoveAll(filepath.Dir(protoFile))

	client, err := createGRPCClient(protoFile, "localhost:50051", nil)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}

	clientMap := grpcClientToMap(client)
	listServicesFunc := clientMap["listServices"].(r2core.BuiltinFunction)

	services := listServicesFunc()
	servicesSlice, ok := services.([]interface{})
	if !ok {
		t.Errorf("listServices should return []interface{}")
	}

	if len(servicesSlice) != 1 {
		t.Errorf("Expected 1 service, got %d", len(servicesSlice))
	}

	if servicesSlice[0] != "TestService" {
		t.Errorf("Expected TestService, got %v", servicesSlice[0])
	}
}

func TestListMethods(t *testing.T) {
	protoFile := createTestProtoFile(t)
	defer os.RemoveAll(filepath.Dir(protoFile))

	client, err := createGRPCClient(protoFile, "localhost:50051", nil)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}

	clientMap := grpcClientToMap(client)
	listMethodsFunc := clientMap["listMethods"].(r2core.BuiltinFunction)

	methods := listMethodsFunc("TestService")
	methodsSlice, ok := methods.([]interface{})
	if !ok {
		t.Errorf("listMethods should return []interface{}")
	}

	if len(methodsSlice) != 4 {
		t.Errorf("Expected 4 methods, got %d", len(methodsSlice))
	}

	expectedMethods := map[string]bool{
		"GetUser": false, "ListUsers": false, "CreateUsers": false, "Chat": false,
	}

	for _, method := range methodsSlice {
		methodName, ok := method.(string)
		if !ok {
			t.Errorf("Method name should be string, got %T", method)
			continue
		}

		if _, exists := expectedMethods[methodName]; !exists {
			t.Errorf("Unexpected method %s", methodName)
		} else {
			expectedMethods[methodName] = true
		}
	}

	for method, found := range expectedMethods {
		if !found {
			t.Errorf("Method %s not found", method)
		}
	}
}

func TestGetMethodInfo(t *testing.T) {
	protoFile := createTestProtoFile(t)
	defer os.RemoveAll(filepath.Dir(protoFile))

	client, err := createGRPCClient(protoFile, "localhost:50051", nil)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}

	clientMap := grpcClientToMap(client)
	getMethodInfoFunc := clientMap["getMethodInfo"].(r2core.BuiltinFunction)

	// Test GetUser method (unary)
	methodInfo := getMethodInfoFunc("TestService", "GetUser")
	infoMap, ok := methodInfo.(map[string]interface{})
	if !ok {
		t.Errorf("getMethodInfo should return map[string]interface{}")
	}

	if infoMap["name"] != "GetUser" {
		t.Errorf("Expected name GetUser, got %v", infoMap["name"])
	}

	if infoMap["inputType"] != "GetUserRequest" {
		t.Errorf("Expected inputType GetUserRequest, got %v", infoMap["inputType"])
	}

	if infoMap["outputType"] != "GetUserResponse" {
		t.Errorf("Expected outputType GetUserResponse, got %v", infoMap["outputType"])
	}

	if infoMap["clientStreaming"] != false {
		t.Errorf("Expected clientStreaming false, got %v", infoMap["clientStreaming"])
	}

	if infoMap["serverStreaming"] != false {
		t.Errorf("Expected serverStreaming false, got %v", infoMap["serverStreaming"])
	}

	// Test ListUsers method (server streaming)
	methodInfo = getMethodInfoFunc("TestService", "ListUsers")
	infoMap = methodInfo.(map[string]interface{})

	if infoMap["serverStreaming"] != true {
		t.Errorf("Expected ListUsers serverStreaming true, got %v", infoMap["serverStreaming"])
	}

	// Test CreateUsers method (client streaming)
	methodInfo = getMethodInfoFunc("TestService", "CreateUsers")
	infoMap = methodInfo.(map[string]interface{})

	if infoMap["clientStreaming"] != true {
		t.Errorf("Expected CreateUsers clientStreaming true, got %v", infoMap["clientStreaming"])
	}

	// Test Chat method (bidirectional streaming)
	methodInfo = getMethodInfoFunc("TestService", "Chat")
	infoMap = methodInfo.(map[string]interface{})

	if infoMap["clientStreaming"] != true {
		t.Errorf("Expected Chat clientStreaming true, got %v", infoMap["clientStreaming"])
	}

	if infoMap["serverStreaming"] != true {
		t.Errorf("Expected Chat serverStreaming true, got %v", infoMap["serverStreaming"])
	}
}

func TestMetadataManagement(t *testing.T) {
	protoFile := createTestProtoFile(t)
	defer os.RemoveAll(filepath.Dir(protoFile))

	client, err := createGRPCClient(protoFile, "localhost:50051", nil)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}

	clientMap := grpcClientToMap(client)
	setMetadataFunc := clientMap["setMetadata"].(r2core.BuiltinFunction)
	getMetadataFunc := clientMap["getMetadata"].(r2core.BuiltinFunction)

	// Test setting metadata with key-value pair
	setMetadataFunc("x-custom-header", "test-value")

	metadata := getMetadataFunc()
	metadataMap, ok := metadata.(map[string]interface{})
	if !ok {
		t.Errorf("getMetadata should return map[string]interface{}")
	}

	if metadataMap["x-custom-header"] != "test-value" {
		t.Errorf("Expected x-custom-header test-value, got %v", metadataMap["x-custom-header"])
	}

	// Test setting metadata with map
	metadataToSet := map[string]interface{}{
		"x-api-key": "api-key-value",
		"x-version": "1.0",
	}
	setMetadataFunc(metadataToSet)

	metadata = getMetadataFunc()
	metadataMap = metadata.(map[string]interface{})

	if metadataMap["x-api-key"] != "api-key-value" {
		t.Errorf("Expected x-api-key api-key-value, got %v", metadataMap["x-api-key"])
	}

	if metadataMap["x-version"] != "1.0" {
		t.Errorf("Expected x-version 1.0, got %v", metadataMap["x-version"])
	}
}

func TestGRPCTLSConfiguration(t *testing.T) {
	protoFile := createTestProtoFile(t)
	defer os.RemoveAll(filepath.Dir(protoFile))

	client, err := createGRPCClient(protoFile, "localhost:50051", nil)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}

	clientMap := grpcClientToMap(client)
	setTLSConfigFunc := clientMap["setTLSConfig"].(r2core.BuiltinFunction)

	// Test TLS configuration
	tlsConfig := map[string]interface{}{
		"enabled":    true,
		"skipVerify": true,
		"serverName": "test-server",
		"minVersion": "1.2",
	}

	setTLSConfigFunc(tlsConfig)

	if !client.UseTLS {
		t.Errorf("Expected UseTLS true, got %v", client.UseTLS)
	}

	if !client.TLSConfig.InsecureSkipVerify {
		t.Errorf("Expected InsecureSkipVerify true, got %v", client.TLSConfig.InsecureSkipVerify)
	}

	if client.TLSConfig.ServerName != "test-server" {
		t.Errorf("Expected ServerName test-server, got %v", client.TLSConfig.ServerName)
	}
}

func TestAuthConfiguration(t *testing.T) {
	protoFile := createTestProtoFile(t)
	defer os.RemoveAll(filepath.Dir(protoFile))

	client, err := createGRPCClient(protoFile, "localhost:50051", nil)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}

	clientMap := grpcClientToMap(client)
	setAuthFunc := clientMap["setAuth"].(r2core.BuiltinFunction)

	// Test Bearer token authentication
	authConfig := map[string]interface{}{
		"type":  "bearer",
		"token": "test-jwt-token",
	}

	setAuthFunc(authConfig)

	if client.Auth.Type != "bearer" {
		t.Errorf("Expected auth type bearer, got %v", client.Auth.Type)
	}

	if client.Auth.Token != "test-jwt-token" {
		t.Errorf("Expected token test-jwt-token, got %v", client.Auth.Token)
	}

	// Test Basic authentication
	authConfig = map[string]interface{}{
		"type":     "basic",
		"username": "testuser",
		"password": "testpass",
	}

	setAuthFunc(authConfig)

	if client.Auth.Type != "basic" {
		t.Errorf("Expected auth type basic, got %v", client.Auth.Type)
	}

	if client.Auth.Username != "testuser" {
		t.Errorf("Expected username testuser, got %v", client.Auth.Username)
	}

	if client.Auth.Password != "testpass" {
		t.Errorf("Expected password testpass, got %v", client.Auth.Password)
	}

	// Test mTLS authentication
	authConfig = map[string]interface{}{
		"type":     "mtls",
		"certFile": "/path/to/cert.pem",
		"keyFile":  "/path/to/key.pem",
		"caFile":   "/path/to/ca.pem",
	}

	setAuthFunc(authConfig)

	if client.Auth.Type != "mtls" {
		t.Errorf("Expected auth type mtls, got %v", client.Auth.Type)
	}

	if client.Auth.CertFile != "/path/to/cert.pem" {
		t.Errorf("Expected certFile /path/to/cert.pem, got %v", client.Auth.CertFile)
	}

	// Test custom metadata authentication
	authConfig = map[string]interface{}{
		"type": "custom",
		"metadata": map[string]interface{}{
			"x-api-key": "custom-api-key",
			"x-tenant":  "tenant-123",
		},
	}

	setAuthFunc(authConfig)

	if client.Auth.Type != "custom" {
		t.Errorf("Expected auth type custom, got %v", client.Auth.Type)
	}

	if client.Auth.Metadata["x-api-key"] != "custom-api-key" {
		t.Errorf("Expected x-api-key custom-api-key, got %v", client.Auth.Metadata["x-api-key"])
	}
}

func TestGRPCTimeout(t *testing.T) {
	protoFile := createTestProtoFile(t)
	defer os.RemoveAll(filepath.Dir(protoFile))

	client, err := createGRPCClient(protoFile, "localhost:50051", nil)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}

	clientMap := grpcClientToMap(client)
	setTimeoutFunc := clientMap["setTimeout"].(r2core.BuiltinFunction)

	// Test setting timeout
	setTimeoutFunc(60.0)

	expectedTimeout := 60 * time.Second
	if client.Timeout != expectedTimeout {
		t.Errorf("Expected timeout %v, got %v", expectedTimeout, client.Timeout)
	}
}

func TestCompression(t *testing.T) {
	protoFile := createTestProtoFile(t)
	defer os.RemoveAll(filepath.Dir(protoFile))

	client, err := createGRPCClient(protoFile, "localhost:50051", nil)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}

	clientMap := grpcClientToMap(client)
	setCompressionFunc := clientMap["setCompression"].(r2core.BuiltinFunction)

	// Test setting compression
	setCompressionFunc("gzip")

	if client.Compression != "gzip" {
		t.Errorf("Expected compression gzip, got %v", client.Compression)
	}
}

func TestGRPCErrorHandling(t *testing.T) {
	// Test with invalid proto file
	_, err := createGRPCClient("/nonexistent/file.proto", "localhost:50051", nil)
	if err == nil {
		t.Errorf("Expected error for nonexistent proto file")
	}

	// Test with invalid server address format
	protoFile := createTestProtoFile(t)
	defer os.RemoveAll(filepath.Dir(protoFile))

	client, err := createGRPCClient(protoFile, "invalid-address", nil)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}

	// Test connection error (should happen on first call)
	clientMap := grpcClientToMap(client)
	callFunc := clientMap["call"].(r2core.BuiltinFunction)

	// Call should return a map with success=false and error info
	result := callFunc("TestService", "GetUser", map[string]interface{}{
		"user_id": "test123",
	})

	resultMap, ok := result.(map[string]interface{})
	if !ok {
		t.Errorf("Expected map result, got %T", result)
		return
	}

	success, exists := resultMap["success"]
	if !exists || success.(bool) {
		t.Errorf("Expected success=false for invalid server address")
	}

	errorInfo, exists := resultMap["error"]
	if !exists {
		t.Errorf("Expected error information in result")
	} else {
		errorMap := errorInfo.(map[string]interface{})
		if errorMap["message"] == nil {
			t.Errorf("Expected error message")
		}
	}
}

func TestValueConversion(t *testing.T) {
	// This test is a placeholder for value conversion functionality
	// Real testing would require proper field descriptors with type information
	// For now, we just test that the function doesn't panic with nil input

	// Test with nil field descriptor (should handle gracefully)
	_, err := convertValueToProto("test", nil)
	if err == nil {
		t.Errorf("Expected error for nil field descriptor")
	}

	// Note: Comprehensive value conversion tests would require
	// mock field descriptors with proper type information from protobuf
}

// TestRepeatedAndMapFieldConversion guards against a regression where
// convertValueToProto/convertProtoValueToR2Lang only handled scalar fields:
// field.GetType() was switched on directly without checking IsRepeated()/
// IsMap() first, so a repeated field's R2Lang array value (represented as
// []interface{}) fell through to the scalar branch (e.g. TYPE_MESSAGE tried
// a map[string]interface{} type assertion on the array and failed, TYPE_
// STRING silently stringified the array via fmt.Sprintf instead of
// converting each element). Map fields were entirely unhandled the same way.
func TestRepeatedAndMapFieldConversion(t *testing.T) {
	protoFile := createTestProtoFile(t)
	defer os.RemoveAll(filepath.Dir(protoFile))

	client, err := createGRPCClient(protoFile, "localhost:50051", nil)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}

	msgDesc := client.FileDesc.FindMessage("testservice.RepeatedFieldsTest")
	if msgDesc == nil {
		t.Fatalf("RepeatedFieldsTest message not found")
	}

	request := map[string]interface{}{
		"tags":   []interface{}{"a", "b", "c"},
		"scores": []interface{}{1.0, 2.0, 3.0},
		"users": []interface{}{
			map[string]interface{}{"user_id": "u1", "name": "Alice"},
			map[string]interface{}{"user_id": "u2", "name": "Bob"},
		},
		"labels": map[string]interface{}{
			"env": "prod",
		},
		"counts": map[string]interface{}{
			"x": 1.0,
			"y": 2.0,
		},
	}

	msg := dynamic.NewMessage(msgDesc)
	if err := populateMessage(msg, request); err != nil {
		t.Fatalf("populateMessage failed: %v", err)
	}

	result, err := messageToMap(msg)
	if err != nil {
		t.Fatalf("messageToMap failed: %v", err)
	}

	tags, ok := result["tags"].([]interface{})
	if !ok || len(tags) != 3 || tags[0] != "a" || tags[1] != "b" || tags[2] != "c" {
		t.Errorf("Expected tags [a b c], got %#v", result["tags"])
	}

	scores, ok := result["scores"].([]interface{})
	if !ok || len(scores) != 3 {
		t.Fatalf("Expected 3 scores, got %#v", result["scores"])
	}
	for i, want := range []float64{1, 2, 3} {
		got, ok := scores[i].(float64)
		if !ok || got != want {
			t.Errorf("scores[%d] = %#v (%T), want float64 %v", i, scores[i], scores[i], want)
		}
	}

	users, ok := result["users"].([]interface{})
	if !ok || len(users) != 2 {
		t.Fatalf("Expected 2 users, got %#v", result["users"])
	}
	u0, ok := users[0].(map[string]interface{})
	if !ok || u0["user_id"] != "u1" || u0["name"] != "Alice" {
		t.Errorf("users[0] = %#v, want map with user_id=u1 name=Alice", users[0])
	}

	labels, ok := result["labels"].(map[string]interface{})
	if !ok || labels["env"] != "prod" {
		t.Errorf("Expected labels map with env=prod, got %#v", result["labels"])
	}

	counts, ok := result["counts"].(map[string]interface{})
	if !ok {
		t.Fatalf("Expected counts map, got %#v", result["counts"])
	}
	if x, ok := counts["x"].(float64); !ok || x != 1 {
		t.Errorf("counts[x] = %#v, want float64 1", counts["x"])
	}
	if y, ok := counts["y"].(float64); !ok || y != 2 {
		t.Errorf("counts[y] = %#v, want float64 2", counts["y"])
	}
}

// TestClientStreamCloseAndReceive is an end-to-end test (real grpc.Server +
// real network connection) guarding against a regression where closeSend()
// for a pure client-streaming call just cancelled the RPC's context instead
// of calling grpcdynamic's ClientStream.CloseAndReceive(). Cancelling
// aborts the RPC outright (the server sees a cancelled context and the
// response is never produced), so the R2Lang script's onReceive callback
// could never fire for a client-streaming call -- the feature was
// unusable for its core purpose of returning the server's aggregated
// response.
func TestClientStreamCloseAndReceive(t *testing.T) {
	protoFile := createTestProtoFile(t)
	defer os.RemoveAll(filepath.Dir(protoFile))

	client, err := createGRPCClient(protoFile, "localhost:50051", nil)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}

	server, addr := startFullMockServer(t, client.FileDesc)
	defer server.Stop()

	client.ServerAddr = addr

	stream := client.callClientStream("TestService", "CreateUsers").(map[string]interface{})
	sendFn := stream["send"].(r2core.BuiltinFunction)
	closeSendFn := stream["closeSend"].(r2core.BuiltinFunction)
	onReceiveFn := stream["onReceive"].(r2core.BuiltinFunction)
	onErrorFn := stream["onError"].(r2core.BuiltinFunction)

	received := make(chan map[string]interface{}, 1)
	streamErrs := make(chan string, 1)

	onReceiveFn(r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		received <- args[0].(map[string]interface{})
		return nil
	}))
	onErrorFn(r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		streamErrs <- fmt.Sprintf("%v", args[0])
		return nil
	}))

	sendFn(map[string]interface{}{"name": "Alice", "email": "alice@example.com", "age": 30.0})
	sendFn(map[string]interface{}{"name": "Bob", "email": "bob@example.com", "age": 25.0})
	closeSendFn()

	select {
	case result := <-received:
		if result["created_count"] != float64(2) {
			t.Errorf("Expected created_count 2, got %v", result["created_count"])
		}
	case errMsg := <-streamErrs:
		t.Fatalf("closeSend produced an error instead of a response: %s", errMsg)
	case <-time.After(5 * time.Second):
		t.Fatal("closeSend/CloseAndReceive never delivered a response via onReceive")
	}
}

// TestBidiStreamCloseSendKeepsReceiving is an end-to-end test guarding
// against a regression where closeSend() on a bidi stream cancelled the
// entire RPC context instead of half-closing the send side via
// BidiStream.CloseSend(). The mock Chat handler echoes each received
// message back, so if closeSend() prematurely kills the stream, the echo
// for the last message sent before closeSend() would be lost.
func TestBidiStreamCloseSendKeepsReceiving(t *testing.T) {
	protoFile := createTestProtoFile(t)
	defer os.RemoveAll(filepath.Dir(protoFile))

	client, err := createGRPCClient(protoFile, "localhost:50051", nil)
	if err != nil {
		t.Fatalf("Failed to create gRPC client: %v", err)
	}

	server, addr := startFullMockServer(t, client.FileDesc)
	defer server.Stop()

	client.ServerAddr = addr

	stream := client.callBidirectionalStream("TestService", "Chat").(map[string]interface{})
	sendFn := stream["send"].(r2core.BuiltinFunction)
	closeSendFn := stream["closeSend"].(r2core.BuiltinFunction)
	onReceiveFn := stream["onReceive"].(r2core.BuiltinFunction)
	onCloseFn := stream["onClose"].(r2core.BuiltinFunction)

	var mu sync.Mutex
	var echoes []string
	closed := make(chan struct{})

	onReceiveFn(r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		mu.Lock()
		echoes = append(echoes, args[0].(map[string]interface{})["message"].(string))
		mu.Unlock()
		return nil
	}))
	onCloseFn(r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		close(closed)
		return nil
	}))

	sendFn(map[string]interface{}{"user_id": "u1", "message": "hello"})
	sendFn(map[string]interface{}{"user_id": "u1", "message": "world"})
	closeSendFn()

	select {
	case <-closed:
	case <-time.After(5 * time.Second):
		t.Fatal("bidi stream never closed after server finished echoing")
	}

	mu.Lock()
	defer mu.Unlock()
	if len(echoes) != 2 || echoes[0] != "hello" || echoes[1] != "world" {
		t.Errorf("Expected echoes [hello world], got %v (closeSend must not drop in-flight responses)", echoes)
	}
}

func TestR2LangIntegration(t *testing.T) {
	// Test the R2Lang integration
	env := r2core.NewEnvironment()
	RegisterGRPC(env)

	// Create a test proto file
	protoFile := createTestProtoFile(t)
	defer os.RemoveAll(filepath.Dir(protoFile))

	// Test accessing grpc module and calling grpcClient
	grpcModule, exists := env.Get("grpc")
	if !exists {
		t.Fatalf("grpc module not found")
	}

	grpcMap := grpcModule.(map[string]interface{})
	grpcClientFunc, exists := grpcMap["grpcClient"]
	if !exists {
		t.Fatalf("grpcClient function not found in grpc module")
	}

	clientFunc := grpcClientFunc.(r2core.BuiltinFunction)
	clientMap := clientFunc(protoFile, "localhost:50051")
	clientMapTyped, ok := clientMap.(map[string]interface{})
	if !ok {
		t.Errorf("grpcClient should return map[string]interface{}")
	}

	// Test that all expected methods are available
	expectedMethods := []string{
		"listServices", "listMethods", "getMethodInfo", "call",
		"callSimple", "callServerStream", "callClientStream",
		"callBidirectionalStream", "setTimeout", "setMetadata",
	}

	for _, method := range expectedMethods {
		if _, exists := clientMapTyped[method]; !exists {
			t.Errorf("Method %s not available in R2Lang client", method)
		}
	}
}

func TestBase64Encoding(t *testing.T) {
	// Test our custom base64 implementation
	testCases := []struct {
		input    string
		expected string
	}{
		{"hello", "aGVsbG8="},
		{"test", "dGVzdA=="},
		{"a", "YQ=="},
		{"", ""},
	}

	for _, tc := range testCases {
		result := base64EncodeGRPC(tc.input)
		if result != tc.expected {
			t.Errorf("base64EncodeGRPC(%s) = %s, expected %s", tc.input, result, tc.expected)
		}
	}

	// Test basic auth encoding
	auth := encodeBasicAuth("user", "pass")
	expected := "Basic dXNlcjpwYXNz"
	if auth != expected {
		t.Errorf("encodeBasicAuth(user, pass) = %s, expected %s", auth, expected)
	}
}

// Benchmark tests
func BenchmarkGRPCClientCreation(b *testing.B) {
	protoFile := createTestProtoFile(b)
	defer os.RemoveAll(filepath.Dir(protoFile))

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		client, err := createGRPCClient(protoFile, "localhost:50051", nil)
		if err != nil {
			b.Fatalf("Failed to create gRPC client: %v", err)
		}
		_ = client
	}
}

// TestStreamCloseReentrantOnCloseDoesNotDeadlock guards against a regression
// where streamMap["close"] invoked stream.onClose() while still holding
// stream.mu (via `defer stream.mu.Unlock()`). Since sync.RWMutex is not
// reentrant, an onClose callback that calls stream.close() again (a natural
// "make sure it's closed" pattern in a script) used to self-deadlock instead
// of hitting the `if stream.closed { return nil }` fast path.
func TestStreamCloseReentrantOnCloseDoesNotDeadlock(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	stream := &GRPCStream{
		isClient: true,
		isServer: true,
		ctx:      ctx,
		cancel:   cancel,
	}

	m := grpcStreamToMap(stream)
	closeFn := m["close"].(r2core.BuiltinFunction)
	onCloseFn := m["onClose"].(r2core.BuiltinFunction)

	reentered := false
	onCloseFn(r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		// Reentrant call: must not deadlock, and must be a no-op since the
		// stream is already marked closed.
		closeFn()
		reentered = true
		return nil
	}))

	done := make(chan struct{})
	go func() {
		closeFn()
		close(done)
	}()

	select {
	case <-done:
		if !reentered {
			t.Fatal("onClose callback never ran")
		}
	case <-time.After(3 * time.Second):
		t.Fatal("DEADLOCK: close() did not return within 3s when onClose re-entered close()")
	}
}

// TestStreamClosePanickingOnCloseIsRecovered ensures that, consistent with
// receiveMessages(), a panicking onClose callback invoked from the close()
// builtin is recovered via safeInvokeCallback rather than crashing the
// calling goroutine.
func TestStreamClosePanickingOnCloseIsRecovered(t *testing.T) {
	ctx, cancel := context.WithCancel(context.Background())
	stream := &GRPCStream{
		isClient: true,
		isServer: true,
		ctx:      ctx,
		cancel:   cancel,
	}

	m := grpcStreamToMap(stream)
	closeFn := m["close"].(r2core.BuiltinFunction)
	onCloseFn := m["onClose"].(r2core.BuiltinFunction)

	onCloseFn(r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		panic("boom: bad script logic in onClose")
	}))

	defer func() {
		if r := recover(); r != nil {
			t.Fatalf("close() should not propagate a panic from onClose, got: %v", r)
		}
	}()

	closeFn()
}

func BenchmarkMetadataOperations(b *testing.B) {
	protoFile := createTestProtoFile(b)
	defer os.RemoveAll(filepath.Dir(protoFile))

	client, err := createGRPCClient(protoFile, "localhost:50051", nil)
	if err != nil {
		b.Fatalf("Failed to create gRPC client: %v", err)
	}

	clientMap := grpcClientToMap(client)
	setMetadataFunc := clientMap["setMetadata"].(r2core.BuiltinFunction)
	getMetadataFunc := clientMap["getMetadata"].(r2core.BuiltinFunction)

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		setMetadataFunc("x-test-header", "test-value")
		getMetadataFunc()
	}
}
