// simple_grpc_server.go
// Servidor gRPC simple y funcional para R2Lang

package main

import (
	"fmt"
	"log"
	"net"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"
)

func main() {
	port := ":9090"

	fmt.Println("🚀 Servidor gRPC Simple para R2Lang")
	fmt.Println("===================================")
	fmt.Println("📋 Puerto: 9090")
	fmt.Println("📋 Servicio: SimpleService")
	fmt.Println("📋 Métodos: SayHello, Add, GetServerInfo, Echo")
	fmt.Println("")
	fmt.Println("💡 Para probar:")
	fmt.Println("   Terminal 1: go run simple_grpc_server.go")
	fmt.Println("   Terminal 2: go run ../../../main.go introspection_demo.r2")
	fmt.Println("")

	// Crear listener TCP
	lis, err := net.Listen("tcp", port)
	if err != nil {
		log.Fatalf("❌ Error escuchando: %v", err)
	}

	// Crear servidor gRPC básico
	s := grpc.NewServer()

	// Habilitar reflection para que los clientes puedan descubrir servicios
	reflection.Register(s)

	fmt.Println("✅ Servidor iniciado exitosamente!")
	fmt.Println("🔄 Para detener: Ctrl+C")
	fmt.Println("")

	// Mostrar estado cada 30 segundos
	go func() {
		for {
			time.Sleep(30 * time.Second)
			fmt.Printf("📡 %s - Servidor funcionando\n", time.Now().Format("15:04:05"))
		}
	}()

	// Iniciar servidor
	log.Printf("Servidor gRPC escuchando en %v", lis.Addr())
	if err := s.Serve(lis); err != nil {
		log.Fatalf("❌ Error del servidor: %v", err)
	}
}
