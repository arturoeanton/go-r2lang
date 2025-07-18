package main

// greeter_server.go
//
// -----------------------------------------------------------------------------
// Paso 0. Instrucción go:generate
// -----------------------------------------------------------------------------
// La primera vez (o cuando cambie el .proto) ejecuta:
//      go generate
// Esto llamará a `protoc` y generará los archivos *_pb.go y *_grpc.pb.go
// justo al lado del .proto, usando rutas relativas.
//

/*go:generate protoc \
    --go_out=. --go_opt=paths=source_relative \
    --go-grpc_out=. --go-grpc_opt=paths=source_relative \
    greeter.proto
	//*/
// -----------------------------------------------------------------------------
// Paso 1. Paquete e imports
// -----------------------------------------------------------------------------
// Este paquete es un ejemplo de implementación de un servidor gRPC

import (
	"context"
	"fmt"
	"log"
	"net"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	// El paquete generado a partir del .proto (mismo directorio → nombre helloworld)
	pb "github.com/arturoeanton/go-r2lang/examples/grpc/example2/helloworld"
)

// -----------------------------------------------------------------------------
// Paso 2. Implementar el servicio Greeter
// -----------------------------------------------------------------------------
type greeterServer struct {
	pb.UnimplementedGreeterServer // Embebemos la impl‑vacía para compatibilidad
}

// SayHello responde con “hello {name}”
func (s *greeterServer) SayHello(
	ctx context.Context,
	req *pb.HelloRequest,
) (*pb.HelloReply, error) {

	name := req.GetName()
	log.Printf("▶️  Received: %s", name)

	// Construimos la respuesta
	return &pb.HelloReply{
		Message: fmt.Sprintf("hello %s", name),
	}, nil
}

// (Opcional) Si más RPCs fuesen añadidas al .proto, se implementan aquí.

// -----------------------------------------------------------------------------
// Paso 3. Punto de entrada (main)
// -----------------------------------------------------------------------------
func main() {
	const addr = ":9090"

	// 3a. Creamos un listener TCP
	lis, err := net.Listen("tcp", addr)
	if err != nil {
		log.Fatalf("❌  listen: %v", err)
	}

	// 3b. Instanciamos el servidor gRPC
	s := grpc.NewServer()

	// 3c. Registramos nuestra implementación
	pb.RegisterGreeterServer(s, &greeterServer{})

	// 3d. Activamos reflection para poder usar herramientas como grpcurl
	reflection.Register(s)

	log.Printf("✅  gRPC Greeter server listening on %s", addr)

	// 3e. Bloqueante: arranca el loop del servidor
	if err := s.Serve(lis); err != nil {
		log.Fatalf("❌  serve: %v", err)
	}
}
