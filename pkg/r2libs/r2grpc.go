package r2libs

import (
	"context"
	"crypto/tls"
	"crypto/x509"
	"fmt"
	"io"
	"io/ioutil"
	"path/filepath"
	"strconv"
	"strings"
	"sync"
	"time"

	"github.com/arturoeanton/go-r2lang/pkg/r2core"
	"github.com/golang/protobuf/proto"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/dynamic/grpcdynamic"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/metadata"
	"google.golang.org/grpc/status"
	"google.golang.org/protobuf/types/descriptorpb"
)

// r2grpc.go: Dynamic gRPC client library for R2Lang with Protocol Buffers parsing

// GRPCClient represents a dynamic gRPC client
type GRPCClient struct {
	ProtoFile   string
	ServerAddr  string
	FileDesc    *desc.FileDescriptor
	Services    map[string]*GRPCService
	Connection  *grpc.ClientConn
	Timeout     time.Duration
	Metadata    map[string]string
	TLSConfig   *tls.Config
	UseTLS      bool
	Auth        *GRPCAuth
	Compression string
	mu          sync.RWMutex
}

// GRPCAuth represents authentication configuration for gRPC
type GRPCAuth struct {
	Type     string // "bearer", "basic", "mtls", "custom"
	Token    string
	Username string
	Password string
	CertFile string
	KeyFile  string
	CAFile   string
	Metadata map[string]string
}

// GRPCService represents a gRPC service with its methods
type GRPCService struct {
	Name    string
	Methods map[string]*GRPCMethod
	Desc    *desc.ServiceDescriptor
}

// GRPCMethod represents a gRPC method with its details
type GRPCMethod struct {
	Name            string
	InputType       string
	OutputType      string
	ClientStreaming bool
	ServerStreaming bool
	Desc            *desc.MethodDescriptor
}

// GRPCStream represents a gRPC streaming connection
type GRPCStream struct {
	serverStream *grpcdynamic.ServerStream
	clientStream *grpcdynamic.ClientStream
	bidiStream   *grpcdynamic.BidiStream
	methodDesc   *desc.MethodDescriptor
	onReceive    func(interface{})
	onError      func(error)
	onClose      func()
	isClient     bool
	isServer     bool
	isBidi       bool
	ctx          context.Context
	cancel       context.CancelFunc
	mu           sync.RWMutex
	closed       bool
}

// RegisterGRPC registers gRPC functions in R2Lang environment
func RegisterGRPC(env *r2core.Environment) {
	functions := map[string]r2core.BuiltinFunction{
		"grpcClient": r2core.BuiltinFunction(func(args ...interface{}) interface{} {
			if len(args) < 2 {
				panic("grpcClient requires (protoFile, serverAddr, [metadata])")
			}

			protoFile, ok1 := args[0].(string)
			serverAddr, ok2 := args[1].(string)

			if !ok1 || !ok2 {
				panic("grpcClient: protoFile and serverAddr must be strings")
			}

			// Parse optional metadata
			var customMetadata map[string]interface{}
			if len(args) > 2 {
				if metadata, ok := args[2].(map[string]interface{}); ok {
					customMetadata = metadata
				} else {
					panic("grpcClient: metadata must be a map")
				}
			}

			client, err := createGRPCClient(protoFile, serverAddr, customMetadata)
			if err != nil {
				errorMsg := fmt.Sprintf("grpcClient: failed to create client from '%s' to '%s'", protoFile, serverAddr)
				if strings.Contains(err.Error(), "connection refused") {
					errorMsg += " - Server is not running or not accepting connections."
				} else if strings.Contains(err.Error(), "no such host") {
					errorMsg += " - DNS resolution failed. Check the server address."
				} else if strings.Contains(err.Error(), "timeout") {
					errorMsg += " - Connection timeout. The server may be slow or overloaded."
				} else if strings.Contains(err.Error(), "proto") {
					errorMsg += " - Protocol buffer file parsing failed. Check the .proto file syntax."
				}
				errorMsg += fmt.Sprintf(" Error: %v", err)
				panic(errorMsg)
			}

			return grpcClientToMap(client)
		}),
	}

	RegisterModule(env, "grpc", functions)
}

// createGRPCClient creates a new gRPC client from proto file
func createGRPCClient(protoFile, serverAddr string, customMetadata map[string]interface{}) (*GRPCClient, error) {
	// Parse proto file
	parser := protoparse.Parser{
		ImportPaths: []string{filepath.Dir(protoFile)},
	}

	fileDescs, err := parser.ParseFiles(filepath.Base(protoFile))
	if err != nil {
		return nil, fmt.Errorf("failed to parse proto file: %v", err)
	}

	if len(fileDescs) == 0 {
		return nil, fmt.Errorf("no file descriptors found in proto file")
	}

	fileDesc := fileDescs[0]

	// Parse services and methods
	services := make(map[string]*GRPCService)
	for _, serviceDesc := range fileDesc.GetServices() {
		service := &GRPCService{
			Name:    serviceDesc.GetName(),
			Methods: make(map[string]*GRPCMethod),
			Desc:    serviceDesc,
		}

		for _, methodDesc := range serviceDesc.GetMethods() {
			method := &GRPCMethod{
				Name:            methodDesc.GetName(),
				InputType:       methodDesc.GetInputType().GetName(),
				OutputType:      methodDesc.GetOutputType().GetName(),
				ClientStreaming: methodDesc.IsClientStreaming(),
				ServerStreaming: methodDesc.IsServerStreaming(),
				Desc:            methodDesc,
			}
			service.Methods[methodDesc.GetName()] = method
		}

		services[serviceDesc.GetName()] = service
	}

	// Initialize default metadata
	defaultMetadata := map[string]string{
		"user-agent": "R2Lang-gRPC-Client/1.0",
		"accept":     "application/grpc",
	}

	// Override with custom metadata if provided
	if customMetadata != nil {
		for key, value := range customMetadata {
			if strValue, ok := value.(string); ok {
				defaultMetadata[key] = strValue
			}
		}
	}

	client := &GRPCClient{
		ProtoFile:   protoFile,
		ServerAddr:  serverAddr,
		FileDesc:    fileDesc,
		Services:    services,
		Timeout:     30 * time.Second,
		Metadata:    defaultMetadata,
		TLSConfig:   &tls.Config{},
		UseTLS:      false,
		Auth:        nil,
		Compression: "",
	}

	return client, nil
}

// connect establishes a connection to the gRPC server
func (client *GRPCClient) connect() error {
	client.mu.Lock()
	defer client.mu.Unlock()

	if client.Connection != nil {
		return nil // Already connected
	}

	var opts []grpc.DialOption

	// Configure TLS/credentials
	if client.UseTLS {
		if client.Auth != nil && client.Auth.Type == "mtls" {
			// mTLS configuration
			cert, err := tls.LoadX509KeyPair(client.Auth.CertFile, client.Auth.KeyFile)
			if err != nil {
				return fmt.Errorf("failed to load client certificate: %v", err)
			}

			var caCert []byte
			if client.Auth.CAFile != "" {
				caCert, err = ioutil.ReadFile(client.Auth.CAFile)
				if err != nil {
					return fmt.Errorf("failed to read CA certificate: %v", err)
				}
			}

			caCertPool := x509.NewCertPool()
			if caCert != nil {
				caCertPool.AppendCertsFromPEM(caCert)
			}

			tlsConfig := &tls.Config{
				Certificates: []tls.Certificate{cert},
				RootCAs:      caCertPool,
			}

			opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(tlsConfig)))
		} else {
			// Standard TLS
			opts = append(opts, grpc.WithTransportCredentials(credentials.NewTLS(client.TLSConfig)))
		}
	} else {
		opts = append(opts, grpc.WithTransportCredentials(insecure.NewCredentials()))
	}

	// Set timeout
	opts = append(opts, grpc.WithTimeout(client.Timeout))

	// Set compression
	if client.Compression != "" {
		opts = append(opts, grpc.WithDefaultCallOptions(grpc.UseCompressor(client.Compression)))
	}

	// Connect to server
	conn, err := grpc.Dial(client.ServerAddr, opts...)
	if err != nil {
		return fmt.Errorf("failed to connect to server: %v", err)
	}

	client.Connection = conn
	return nil
}

// grpcClientToMap converts GRPCClient to R2Lang map with methods
func grpcClientToMap(client *GRPCClient) map[string]interface{} {
	clientMap := make(map[string]interface{})

	// Client properties
	clientMap["protoFile"] = client.ProtoFile
	clientMap["serverAddr"] = client.ServerAddr
	clientMap["useTLS"] = client.UseTLS

	// List available services
	clientMap["listServices"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		services := make([]interface{}, 0, len(client.Services))
		for name := range client.Services {
			services = append(services, name)
		}
		return services
	})

	// List methods for a service
	clientMap["listMethods"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("listMethods requires (serviceName)")
		}

		serviceName, ok := args[0].(string)
		if !ok {
			panic("listMethods: serviceName must be a string")
		}

		service, exists := client.Services[serviceName]
		if !exists {
			panic(fmt.Sprintf("listMethods: service '%s' not found", serviceName))
		}

		methods := make([]interface{}, 0, len(service.Methods))
		for name := range service.Methods {
			methods = append(methods, name)
		}
		return methods
	})

	// Get method info
	clientMap["getMethodInfo"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("getMethodInfo requires (serviceName, methodName)")
		}

		serviceName, ok1 := args[0].(string)
		methodName, ok2 := args[1].(string)

		if !ok1 || !ok2 {
			panic("getMethodInfo: serviceName and methodName must be strings")
		}

		service, exists := client.Services[serviceName]
		if !exists {
			panic(fmt.Sprintf("getMethodInfo: service '%s' not found", serviceName))
		}

		method, exists := service.Methods[methodName]
		if !exists {
			panic(fmt.Sprintf("getMethodInfo: method '%s' not found in service '%s'", methodName, serviceName))
		}

		return map[string]interface{}{
			"name":            method.Name,
			"inputType":       method.InputType,
			"outputType":      method.OutputType,
			"clientStreaming": method.ClientStreaming,
			"serverStreaming": method.ServerStreaming,
		}
	})

	// Call unary RPC method
	clientMap["call"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 3 {
			panic("call requires (serviceName, methodName, request)")
		}

		serviceName, ok1 := args[0].(string)
		methodName, ok2 := args[1].(string)
		request, ok3 := args[2].(map[string]interface{})

		if !ok1 || !ok2 || !ok3 {
			panic("call: serviceName and methodName must be strings, request must be a map")
		}

		return client.callUnary(serviceName, methodName, request)
	})

	// Call unary RPC method (simplified response)
	clientMap["callSimple"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 3 {
			panic("callSimple requires (serviceName, methodName, request)")
		}

		serviceName, ok1 := args[0].(string)
		methodName, ok2 := args[1].(string)
		request, ok3 := args[2].(map[string]interface{})

		if !ok1 || !ok2 || !ok3 {
			panic("callSimple: serviceName and methodName must be strings, request must be a map")
		}

		response := client.callUnary(serviceName, methodName, request)
		if responseMap, ok := response.(map[string]interface{}); ok {
			if success, exists := responseMap["success"]; exists && success == true {
				if result, exists := responseMap["result"]; exists {
					return result
				}
			} else {
				if err, exists := responseMap["error"]; exists {
					return err
				}
			}
		}
		return response
	})

	// Call server streaming RPC method
	clientMap["callServerStream"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 3 {
			panic("callServerStream requires (serviceName, methodName, request)")
		}

		serviceName, ok1 := args[0].(string)
		methodName, ok2 := args[1].(string)
		request, ok3 := args[2].(map[string]interface{})

		if !ok1 || !ok2 || !ok3 {
			panic("callServerStream: serviceName and methodName must be strings, request must be a map")
		}

		return client.callServerStream(serviceName, methodName, request)
	})

	// Call client streaming RPC method
	clientMap["callClientStream"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("callClientStream requires (serviceName, methodName)")
		}

		serviceName, ok1 := args[0].(string)
		methodName, ok2 := args[1].(string)

		if !ok1 || !ok2 {
			panic("callClientStream: serviceName and methodName must be strings")
		}

		return client.callClientStream(serviceName, methodName)
	})

	// Call bidirectional streaming RPC method
	clientMap["callBidirectionalStream"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 2 {
			panic("callBidirectionalStream requires (serviceName, methodName)")
		}

		serviceName, ok1 := args[0].(string)
		methodName, ok2 := args[1].(string)

		if !ok1 || !ok2 {
			panic("callBidirectionalStream: serviceName and methodName must be strings")
		}

		return client.callBidirectionalStream(serviceName, methodName)
	})

	// Set timeout
	clientMap["setTimeout"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("setTimeout requires (seconds)")
		}

		seconds, ok := args[0].(float64)
		if !ok {
			panic("setTimeout: seconds must be a number")
		}

		client.Timeout = time.Duration(seconds) * time.Second
		return nil
	})

	// Set metadata
	clientMap["setMetadata"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("setMetadata requires (key, value) or (metadataMap)")
		}

		// If single argument, treat as metadata map
		if len(args) == 1 {
			if metadataMap, ok := args[0].(map[string]interface{}); ok {
				client.mu.Lock()
				for key, value := range metadataMap {
					if strValue, ok := value.(string); ok {
						client.Metadata[key] = strValue
					}
				}
				client.mu.Unlock()
				return nil
			}
		}

		// If two arguments, treat as key-value pair
		if len(args) >= 2 {
			key, ok1 := args[0].(string)
			value, ok2 := args[1].(string)

			if !ok1 || !ok2 {
				panic("setMetadata: key and value must be strings")
			}

			client.mu.Lock()
			client.Metadata[key] = value
			client.mu.Unlock()
			return nil
		}

		panic("setMetadata: invalid arguments")
	})

	// Get metadata
	clientMap["getMetadata"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		client.mu.RLock()
		defer client.mu.RUnlock()

		metadataMap := make(map[string]interface{})
		for key, value := range client.Metadata {
			metadataMap[key] = value
		}
		return metadataMap
	})

	// Set TLS configuration
	clientMap["setTLSConfig"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("setTLSConfig requires (configMap)")
		}

		configMap, ok := args[0].(map[string]interface{})
		if !ok {
			panic("setTLSConfig: configMap must be a map")
		}

		client.mu.Lock()
		defer client.mu.Unlock()

		// Enable TLS
		if enabled, exists := configMap["enabled"]; exists {
			if enable, ok := enabled.(bool); ok {
				client.UseTLS = enable
			}
		}

		// Skip certificate verification
		if skipVerify, exists := configMap["skipVerify"]; exists {
			if skip, ok := skipVerify.(bool); ok {
				client.TLSConfig.InsecureSkipVerify = skip
			}
		}

		// Set server name
		if serverName, exists := configMap["serverName"]; exists {
			if name, ok := serverName.(string); ok {
				client.TLSConfig.ServerName = name
			}
		}

		// Set minimum TLS version
		if minVersion, exists := configMap["minVersion"]; exists {
			if version, ok := minVersion.(string); ok {
				switch version {
				case "1.0":
					client.TLSConfig.MinVersion = tls.VersionTLS10
				case "1.1":
					client.TLSConfig.MinVersion = tls.VersionTLS11
				case "1.2":
					client.TLSConfig.MinVersion = tls.VersionTLS12
				case "1.3":
					client.TLSConfig.MinVersion = tls.VersionTLS13
				default:
					client.TLSConfig.MinVersion = tls.VersionTLS12
				}
			}
		}

		return nil
	})

	// Set authentication
	clientMap["setAuth"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("setAuth requires (authConfig)")
		}

		authConfig, ok := args[0].(map[string]interface{})
		if !ok {
			panic("setAuth: authConfig must be a map")
		}

		client.mu.Lock()
		defer client.mu.Unlock()

		auth := &GRPCAuth{
			Metadata: make(map[string]string),
		}

		if authType, exists := authConfig["type"]; exists {
			if typeStr, ok := authType.(string); ok {
				auth.Type = typeStr
			}
		}

		if token, exists := authConfig["token"]; exists {
			if tokenStr, ok := token.(string); ok {
				auth.Token = tokenStr
			}
		}

		if username, exists := authConfig["username"]; exists {
			if userStr, ok := username.(string); ok {
				auth.Username = userStr
			}
		}

		if password, exists := authConfig["password"]; exists {
			if passStr, ok := password.(string); ok {
				auth.Password = passStr
			}
		}

		if certFile, exists := authConfig["certFile"]; exists {
			if certStr, ok := certFile.(string); ok {
				auth.CertFile = certStr
			}
		}

		if keyFile, exists := authConfig["keyFile"]; exists {
			if keyStr, ok := keyFile.(string); ok {
				auth.KeyFile = keyStr
			}
		}

		if caFile, exists := authConfig["caFile"]; exists {
			if caStr, ok := caFile.(string); ok {
				auth.CAFile = caStr
			}
		}

		if metadata, exists := authConfig["metadata"]; exists {
			if metadataMap, ok := metadata.(map[string]interface{}); ok {
				for key, value := range metadataMap {
					if strValue, ok := value.(string); ok {
						auth.Metadata[key] = strValue
					}
				}
			}
		}

		client.Auth = auth
		return nil
	})

	// Set compression
	clientMap["setCompression"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("setCompression requires (algorithm)")
		}

		algorithm, ok := args[0].(string)
		if !ok {
			panic("setCompression: algorithm must be a string")
		}

		client.mu.Lock()
		client.Compression = algorithm
		client.mu.Unlock()
		return nil
	})

	// Close connection
	clientMap["close"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		client.mu.Lock()
		defer client.mu.Unlock()

		if client.Connection != nil {
			err := client.Connection.Close()
			client.Connection = nil
			if err != nil {
				panic(fmt.Sprintf("close: %v", err))
			}
		}
		return nil
	})

	return clientMap
}

// callUnary performs a unary RPC call
func (client *GRPCClient) callUnary(serviceName, methodName string, request map[string]interface{}) interface{} {
	// Get service and method
	service, exists := client.Services[serviceName]
	if !exists {
		panic(fmt.Sprintf("service '%s' not found", serviceName))
	}

	method, exists := service.Methods[methodName]
	if !exists {
		panic(fmt.Sprintf("method '%s' not found in service '%s'", methodName, serviceName))
	}

	if method.ClientStreaming || method.ServerStreaming {
		panic(fmt.Sprintf("method '%s' is streaming, use appropriate streaming method", methodName))
	}

	// Connect if not connected
	if err := client.connect(); err != nil {
		panic(fmt.Sprintf("connection failed: %v", err))
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout)
	defer cancel()

	// Add metadata to context
	if len(client.Metadata) > 0 || client.Auth != nil {
		md := metadata.New(client.Metadata)

		// Add authentication metadata
		if client.Auth != nil {
			switch client.Auth.Type {
			case "bearer":
				if client.Auth.Token != "" {
					md.Set("authorization", "Bearer "+client.Auth.Token)
				}
			case "basic":
				if client.Auth.Username != "" && client.Auth.Password != "" {
					md.Set("authorization", "Basic "+encodeBasicAuth(client.Auth.Username, client.Auth.Password))
				}
			case "custom":
				for key, value := range client.Auth.Metadata {
					md.Set(key, value)
				}
			}
		}

		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	// Create dynamic message from request
	inputMsg := dynamic.NewMessage(method.Desc.GetInputType())
	if err := populateMessage(inputMsg, request); err != nil {
		panic(fmt.Sprintf("failed to populate request message: %v", err))
	}

	// Create dynamic stub
	stub := grpcdynamic.NewStub(client.Connection)

	// Invoke method
	response, err := stub.InvokeRpc(ctx, method.Desc, inputMsg)
	if err != nil {
		if st, ok := status.FromError(err); ok {
			return map[string]interface{}{
				"success": false,
				"error": map[string]interface{}{
					"code":    st.Code().String(),
					"message": st.Message(),
					"details": st.Details(),
				},
			}
		}
		return map[string]interface{}{
			"success": false,
			"error":   err.Error(),
		}
	}

	// Convert response to map
	var result map[string]interface{}
	if dynMsg, ok := response.(*dynamic.Message); ok {
		result, err = messageToMap(dynMsg)
		if err != nil {
			panic(fmt.Sprintf("failed to convert response: %v", err))
		}
	} else {
		panic(fmt.Sprintf("unexpected response type: %T", response))
	}

	return map[string]interface{}{
		"success": true,
		"result":  result,
	}
}

// callServerStream performs a server streaming RPC call
func (client *GRPCClient) callServerStream(serviceName, methodName string, request map[string]interface{}) interface{} {
	// Get service and method
	service, exists := client.Services[serviceName]
	if !exists {
		panic(fmt.Sprintf("service '%s' not found", serviceName))
	}

	method, exists := service.Methods[methodName]
	if !exists {
		panic(fmt.Sprintf("method '%s' not found in service '%s'", methodName, serviceName))
	}

	if !method.ServerStreaming {
		panic(fmt.Sprintf("method '%s' is not server streaming", methodName))
	}

	// Connect if not connected
	if err := client.connect(); err != nil {
		panic(fmt.Sprintf("connection failed: %v", err))
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout)

	// Add metadata to context
	if len(client.Metadata) > 0 || client.Auth != nil {
		md := metadata.New(client.Metadata)

		// Add authentication metadata
		if client.Auth != nil {
			switch client.Auth.Type {
			case "bearer":
				if client.Auth.Token != "" {
					md.Set("authorization", "Bearer "+client.Auth.Token)
				}
			case "basic":
				if client.Auth.Username != "" && client.Auth.Password != "" {
					md.Set("authorization", "Basic "+encodeBasicAuth(client.Auth.Username, client.Auth.Password))
				}
			case "custom":
				for key, value := range client.Auth.Metadata {
					md.Set(key, value)
				}
			}
		}

		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	// Create dynamic message from request
	inputMsg := dynamic.NewMessage(method.Desc.GetInputType())
	if err := populateMessage(inputMsg, request); err != nil {
		panic(fmt.Sprintf("failed to populate request message: %v", err))
	}

	// Create dynamic stub
	stub := grpcdynamic.NewStub(client.Connection)

	// Invoke streaming method
	stream, err := stub.InvokeRpcServerStream(ctx, method.Desc, inputMsg)
	if err != nil {
		panic(fmt.Sprintf("failed to invoke server stream: %v", err))
	}

	// Create stream wrapper
	grpcStream := &GRPCStream{
		serverStream: stream,
		methodDesc:   method.Desc,
		isServer:     true,
		ctx:          ctx,
		cancel:       cancel,
	}

	// Start receiving messages in background
	go grpcStream.receiveMessages()

	return grpcStreamToMap(grpcStream)
}

// callClientStream performs a client streaming RPC call
func (client *GRPCClient) callClientStream(serviceName, methodName string) interface{} {
	// Get service and method
	service, exists := client.Services[serviceName]
	if !exists {
		panic(fmt.Sprintf("service '%s' not found", serviceName))
	}

	method, exists := service.Methods[methodName]
	if !exists {
		panic(fmt.Sprintf("method '%s' not found in service '%s'", methodName, serviceName))
	}

	if !method.ClientStreaming {
		panic(fmt.Sprintf("method '%s' is not client streaming", methodName))
	}

	// Connect if not connected
	if err := client.connect(); err != nil {
		panic(fmt.Sprintf("connection failed: %v", err))
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout)

	// Add metadata to context
	if len(client.Metadata) > 0 || client.Auth != nil {
		md := metadata.New(client.Metadata)

		// Add authentication metadata
		if client.Auth != nil {
			switch client.Auth.Type {
			case "bearer":
				if client.Auth.Token != "" {
					md.Set("authorization", "Bearer "+client.Auth.Token)
				}
			case "basic":
				if client.Auth.Username != "" && client.Auth.Password != "" {
					md.Set("authorization", "Basic "+encodeBasicAuth(client.Auth.Username, client.Auth.Password))
				}
			case "custom":
				for key, value := range client.Auth.Metadata {
					md.Set(key, value)
				}
			}
		}

		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	// Create dynamic stub
	stub := grpcdynamic.NewStub(client.Connection)

	// Invoke streaming method
	stream, err := stub.InvokeRpcClientStream(ctx, method.Desc)
	if err != nil {
		panic(fmt.Sprintf("failed to invoke client stream: %v", err))
	}

	// Create stream wrapper
	grpcStream := &GRPCStream{
		clientStream: stream,
		methodDesc:   method.Desc,
		isClient:     true,
		ctx:          ctx,
		cancel:       cancel,
	}

	return grpcStreamToMap(grpcStream)
}

// callBidirectionalStream performs a bidirectional streaming RPC call
func (client *GRPCClient) callBidirectionalStream(serviceName, methodName string) interface{} {
	// Get service and method
	service, exists := client.Services[serviceName]
	if !exists {
		panic(fmt.Sprintf("service '%s' not found", serviceName))
	}

	method, exists := service.Methods[methodName]
	if !exists {
		panic(fmt.Sprintf("method '%s' not found in service '%s'", methodName, serviceName))
	}

	if !method.ClientStreaming || !method.ServerStreaming {
		panic(fmt.Sprintf("method '%s' is not bidirectional streaming", methodName))
	}

	// Connect if not connected
	if err := client.connect(); err != nil {
		panic(fmt.Sprintf("connection failed: %v", err))
	}

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), client.Timeout)

	// Add metadata to context
	if len(client.Metadata) > 0 || client.Auth != nil {
		md := metadata.New(client.Metadata)

		// Add authentication metadata
		if client.Auth != nil {
			switch client.Auth.Type {
			case "bearer":
				if client.Auth.Token != "" {
					md.Set("authorization", "Bearer "+client.Auth.Token)
				}
			case "basic":
				if client.Auth.Username != "" && client.Auth.Password != "" {
					md.Set("authorization", "Basic "+encodeBasicAuth(client.Auth.Username, client.Auth.Password))
				}
			case "custom":
				for key, value := range client.Auth.Metadata {
					md.Set(key, value)
				}
			}
		}

		ctx = metadata.NewOutgoingContext(ctx, md)
	}

	// Create dynamic stub
	stub := grpcdynamic.NewStub(client.Connection)

	// Invoke streaming method
	stream, err := stub.InvokeRpcBidiStream(ctx, method.Desc)
	if err != nil {
		panic(fmt.Sprintf("failed to invoke bidirectional stream: %v", err))
	}

	// Create stream wrapper
	grpcStream := &GRPCStream{
		bidiStream: stream,
		methodDesc: method.Desc,
		isClient:   true,
		isServer:   true,
		isBidi:     true,
		ctx:        ctx,
		cancel:     cancel,
	}

	// Start receiving messages in background
	go grpcStream.receiveMessages()

	return grpcStreamToMap(grpcStream)
}

// grpcStreamToMap converts GRPCStream to R2Lang map with methods
func grpcStreamToMap(stream *GRPCStream) map[string]interface{} {
	streamMap := make(map[string]interface{})

	// Stream properties
	streamMap["isClient"] = stream.isClient
	streamMap["isServer"] = stream.isServer
	streamMap["isBidi"] = stream.isBidi

	// Send message (for client streaming)
	streamMap["send"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if !stream.isClient {
			panic("send: not a client streaming method")
		}

		if len(args) < 1 {
			panic("send requires (message)")
		}

		message, ok := args[0].(map[string]interface{})
		if !ok {
			panic("send: message must be a map")
		}

		stream.mu.RLock()
		if stream.closed {
			stream.mu.RUnlock()
			panic("send: stream is closed")
		}
		stream.mu.RUnlock()

		// Create dynamic message
		inputMsg := dynamic.NewMessage(stream.methodDesc.GetInputType())
		if err := populateMessage(inputMsg, message); err != nil {
			panic(fmt.Sprintf("failed to populate message: %v", err))
		}

		// Send message
		var sendErr error
		if stream.clientStream != nil {
			sendErr = stream.clientStream.SendMsg(inputMsg)
		} else if stream.bidiStream != nil {
			sendErr = stream.bidiStream.SendMsg(inputMsg)
		} else {
			panic("send: invalid stream type")
		}

		if sendErr != nil {
			panic(fmt.Sprintf("failed to send message: %v", sendErr))
		}

		return nil
	})

	// Close send (for client streaming)
	streamMap["closeSend"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if !stream.isClient {
			panic("closeSend: not a client streaming method")
		}

		// Note: grpcdynamic streams don't have a direct CloseSend method
		// The stream is closed when the parent context is cancelled
		stream.mu.Lock()
		if !stream.closed {
			stream.closed = true
			stream.cancel()
		}
		stream.mu.Unlock()

		return nil
	})

	// Set receive callback
	streamMap["onReceive"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if !stream.isServer {
			panic("onReceive: not a server streaming method")
		}

		if len(args) < 1 {
			panic("onReceive requires (callback)")
		}

		callback, ok := args[0].(r2core.BuiltinFunction)
		if !ok {
			panic("onReceive: callback must be a function")
		}

		stream.mu.Lock()
		stream.onReceive = func(msg interface{}) {
			callback(msg)
		}
		stream.mu.Unlock()

		return nil
	})

	// Set error callback
	streamMap["onError"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("onError requires (callback)")
		}

		callback, ok := args[0].(r2core.BuiltinFunction)
		if !ok {
			panic("onError: callback must be a function")
		}

		stream.mu.Lock()
		stream.onError = func(err error) {
			callback(err.Error())
		}
		stream.mu.Unlock()

		return nil
	})

	// Set close callback
	streamMap["onClose"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		if len(args) < 1 {
			panic("onClose requires (callback)")
		}

		callback, ok := args[0].(r2core.BuiltinFunction)
		if !ok {
			panic("onClose: callback must be a function")
		}

		stream.mu.Lock()
		stream.onClose = func() {
			callback()
		}
		stream.mu.Unlock()

		return nil
	})

	// Close stream
	streamMap["close"] = r2core.BuiltinFunction(func(args ...interface{}) interface{} {
		stream.mu.Lock()
		defer stream.mu.Unlock()

		if stream.closed {
			return nil
		}

		stream.closed = true
		stream.cancel()

		if stream.onClose != nil {
			stream.onClose()
		}

		return nil
	})

	return streamMap
}

// receiveMessages receives messages from server streaming
func (stream *GRPCStream) receiveMessages() {
	defer func() {
		stream.mu.Lock()
		stream.closed = true
		stream.mu.Unlock()

		if stream.onClose != nil {
			stream.onClose()
		}
	}()

	for {
		select {
		case <-stream.ctx.Done():
			return
		default:
			// Receive message based on stream type
			var msg proto.Message
			var err error

			if stream.serverStream != nil {
				msg, err = stream.serverStream.RecvMsg()
			} else if stream.bidiStream != nil {
				msg, err = stream.bidiStream.RecvMsg()
			} else {
				if stream.onError != nil {
					stream.onError(fmt.Errorf("invalid stream type for receiving"))
				}
				return
			}

			if err != nil {
				if err == io.EOF {
					return // End of stream
				}
				if stream.onError != nil {
					stream.onError(err)
				}
				return
			}

			// Convert message to map
			var result map[string]interface{}
			if dynMsg, ok := msg.(*dynamic.Message); ok {
				result, err = messageToMap(dynMsg)
				if err != nil {
					if stream.onError != nil {
						stream.onError(err)
					}
					return
				}
			} else {
				if stream.onError != nil {
					stream.onError(fmt.Errorf("unexpected message type: %T", msg))
				}
				return
			}

			// Call receive callback
			if stream.onReceive != nil {
				stream.onReceive(result)
			}
		}
	}
}

// populateMessage populates a dynamic message from a map
func populateMessage(msg *dynamic.Message, data map[string]interface{}) error {
	msgDesc := msg.GetMessageDescriptor()

	for fieldName, value := range data {
		field := msgDesc.FindFieldByName(fieldName)
		if field == nil {
			continue // Skip unknown fields
		}

		convertedValue, err := convertValueToProto(value, field)
		if err != nil {
			return fmt.Errorf("failed to convert field '%s': %v", fieldName, err)
		}

		msg.SetField(field, convertedValue)
	}

	return nil
}

// convertValueToProto converts R2Lang value to protobuf value
func convertValueToProto(value interface{}, field *desc.FieldDescriptor) (interface{}, error) {
	if field == nil {
		return nil, fmt.Errorf("field descriptor is nil")
	}

	switch field.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		if str, ok := value.(string); ok {
			return str, nil
		}
		return fmt.Sprintf("%v", value), nil

	case descriptorpb.FieldDescriptorProto_TYPE_INT32, descriptorpb.FieldDescriptorProto_TYPE_SINT32, descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
		if num, ok := value.(float64); ok {
			return int32(num), nil
		}
		if str, ok := value.(string); ok {
			if i, err := strconv.ParseInt(str, 10, 32); err == nil {
				return int32(i), nil
			}
		}
		return nil, fmt.Errorf("cannot convert %v to int32", value)

	case descriptorpb.FieldDescriptorProto_TYPE_INT64, descriptorpb.FieldDescriptorProto_TYPE_SINT64, descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
		if num, ok := value.(float64); ok {
			return int64(num), nil
		}
		if str, ok := value.(string); ok {
			if i, err := strconv.ParseInt(str, 10, 64); err == nil {
				return i, nil
			}
		}
		return nil, fmt.Errorf("cannot convert %v to int64", value)

	case descriptorpb.FieldDescriptorProto_TYPE_UINT32, descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
		if num, ok := value.(float64); ok {
			return uint32(num), nil
		}
		if str, ok := value.(string); ok {
			if i, err := strconv.ParseUint(str, 10, 32); err == nil {
				return uint32(i), nil
			}
		}
		return nil, fmt.Errorf("cannot convert %v to uint32", value)

	case descriptorpb.FieldDescriptorProto_TYPE_UINT64, descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
		if num, ok := value.(float64); ok {
			return uint64(num), nil
		}
		if str, ok := value.(string); ok {
			if i, err := strconv.ParseUint(str, 10, 64); err == nil {
				return i, nil
			}
		}
		return nil, fmt.Errorf("cannot convert %v to uint64", value)

	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		if num, ok := value.(float64); ok {
			return float32(num), nil
		}
		if str, ok := value.(string); ok {
			if f, err := strconv.ParseFloat(str, 32); err == nil {
				return float32(f), nil
			}
		}
		return nil, fmt.Errorf("cannot convert %v to float32", value)

	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		if num, ok := value.(float64); ok {
			return num, nil
		}
		if str, ok := value.(string); ok {
			if f, err := strconv.ParseFloat(str, 64); err == nil {
				return f, nil
			}
		}
		return nil, fmt.Errorf("cannot convert %v to float64", value)

	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		if b, ok := value.(bool); ok {
			return b, nil
		}
		if str, ok := value.(string); ok {
			if b, err := strconv.ParseBool(str); err == nil {
				return b, nil
			}
		}
		return nil, fmt.Errorf("cannot convert %v to bool", value)

	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		if str, ok := value.(string); ok {
			return []byte(str), nil
		}
		return nil, fmt.Errorf("cannot convert %v to bytes", value)

	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		if msgMap, ok := value.(map[string]interface{}); ok {
			nestedMsg := dynamic.NewMessage(field.GetMessageType())
			if err := populateMessage(nestedMsg, msgMap); err != nil {
				return nil, err
			}
			return nestedMsg, nil
		}
		return nil, fmt.Errorf("cannot convert %v to message", value)

	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		if str, ok := value.(string); ok {
			enumDesc := field.GetEnumType()
			enumValue := enumDesc.FindValueByName(str)
			if enumValue != nil {
				return enumValue.GetNumber(), nil
			}
		}
		if num, ok := value.(float64); ok {
			return int32(num), nil
		}
		return nil, fmt.Errorf("cannot convert %v to enum", value)

	default:
		return value, nil
	}
}

// messageToMap converts a dynamic message to a map
func messageToMap(msg *dynamic.Message) (map[string]interface{}, error) {
	result := make(map[string]interface{})
	msgDesc := msg.GetMessageDescriptor()

	for _, field := range msgDesc.GetFields() {
		if !msg.HasField(field) {
			continue
		}

		value := msg.GetField(field)
		convertedValue, err := convertProtoValueToR2Lang(value, field)
		if err != nil {
			return nil, fmt.Errorf("failed to convert field '%s': %v", field.GetName(), err)
		}

		result[field.GetName()] = convertedValue
	}

	return result, nil
}

// convertProtoValueToR2Lang converts protobuf value to R2Lang value
func convertProtoValueToR2Lang(value interface{}, field *desc.FieldDescriptor) (interface{}, error) {
	if value == nil {
		return nil, nil
	}

	switch field.GetType() {
	case descriptorpb.FieldDescriptorProto_TYPE_STRING:
		return value, nil

	case descriptorpb.FieldDescriptorProto_TYPE_INT32, descriptorpb.FieldDescriptorProto_TYPE_SINT32, descriptorpb.FieldDescriptorProto_TYPE_SFIXED32:
		if i, ok := value.(int32); ok {
			return float64(i), nil
		}
		return value, nil

	case descriptorpb.FieldDescriptorProto_TYPE_INT64, descriptorpb.FieldDescriptorProto_TYPE_SINT64, descriptorpb.FieldDescriptorProto_TYPE_SFIXED64:
		if i, ok := value.(int64); ok {
			return float64(i), nil
		}
		return value, nil

	case descriptorpb.FieldDescriptorProto_TYPE_UINT32, descriptorpb.FieldDescriptorProto_TYPE_FIXED32:
		if i, ok := value.(uint32); ok {
			return float64(i), nil
		}
		return value, nil

	case descriptorpb.FieldDescriptorProto_TYPE_UINT64, descriptorpb.FieldDescriptorProto_TYPE_FIXED64:
		if i, ok := value.(uint64); ok {
			return float64(i), nil
		}
		return value, nil

	case descriptorpb.FieldDescriptorProto_TYPE_FLOAT:
		if f, ok := value.(float32); ok {
			return float64(f), nil
		}
		return value, nil

	case descriptorpb.FieldDescriptorProto_TYPE_DOUBLE:
		return value, nil

	case descriptorpb.FieldDescriptorProto_TYPE_BOOL:
		return value, nil

	case descriptorpb.FieldDescriptorProto_TYPE_BYTES:
		if b, ok := value.([]byte); ok {
			return string(b), nil
		}
		return value, nil

	case descriptorpb.FieldDescriptorProto_TYPE_MESSAGE:
		if msg, ok := value.(*dynamic.Message); ok {
			return messageToMap(msg)
		}
		return value, nil

	case descriptorpb.FieldDescriptorProto_TYPE_ENUM:
		if enumVal, ok := value.(int32); ok {
			enumDesc := field.GetEnumType()
			enumValue := enumDesc.FindValueByNumber(enumVal)
			if enumValue != nil {
				return enumValue.GetName(), nil
			}
			return float64(enumVal), nil
		}
		return value, nil

	default:
		return value, nil
	}
}

// encodeBasicAuth encodes username and password for basic authentication
func encodeBasicAuth(username, password string) string {
	auth := username + ":" + password
	return "Basic " + base64EncodeGRPC(auth)
}

// base64EncodeGRPC encodes a string to base64 for gRPC
func base64EncodeGRPC(s string) string {
	// Simple base64 encoding implementation
	const charset = "ABCDEFGHIJKLMNOPQRSTUVWXYZabcdefghijklmnopqrstuvwxyz0123456789+/"
	var result strings.Builder

	for i := 0; i < len(s); i += 3 {
		// Get 3 bytes (24 bits)
		var b1, b2, b3 byte
		b1 = s[i]
		if i+1 < len(s) {
			b2 = s[i+1]
		}
		if i+2 < len(s) {
			b3 = s[i+2]
		}

		// Convert to 4 6-bit values
		result.WriteByte(charset[b1>>2])
		result.WriteByte(charset[((b1&0x03)<<4)|(b2>>4)])

		if i+1 < len(s) {
			result.WriteByte(charset[((b2&0x0F)<<2)|(b3>>6)])
		} else {
			result.WriteByte('=')
		}

		if i+2 < len(s) {
			result.WriteByte(charset[b3&0x3F])
		} else {
			result.WriteByte('=')
		}
	}

	return result.String()
}
