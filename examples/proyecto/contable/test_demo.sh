#!/bin/bash

echo "🧪 Testing Sistema Contable LATAM Web Demo"
echo "=========================================="

# Start server in background
echo "🚀 Starting server..."
go run main.go examples/proyecto/contable/contable_web_simple.r2 &
SERVER_PID=$!

# Wait for server to start
sleep 3

echo ""
echo "✅ Testing root endpoint..."
curl -s http://localhost:8080/ | grep -q "Sistema Contable LATAM" && echo "✓ Root page works!" || echo "✗ Root page failed"

echo ""
echo "✅ Testing API regiones..."
curl -s http://localhost:8080/api/regiones | grep -q "COL" && echo "✓ API regiones works!" || echo "✗ API regiones failed"

echo ""
echo "✅ Testing demo page..."
curl -s http://localhost:8080/demo | grep -q "Demo Completo" && echo "✓ Demo page works!" || echo "✗ Demo page failed"

echo ""
echo "✅ Testing POST venta..."
curl -s -X POST http://localhost:8080/procesar-venta -d "region=COL&importe=100000" | grep -q "COMPROBANTE" && echo "✓ POST venta works!" || echo "✗ POST venta failed"

echo ""
echo "✅ Testing POST compra..."
curl -s -X POST http://localhost:8080/procesar-compra -d "region=MX&importe=50000" | grep -q "COMPROBANTE" && echo "✓ POST compra works!" || echo "✗ POST compra failed"

echo ""
echo "✅ Testing API transacciones..."
curl -s http://localhost:8080/api/transacciones | grep -q "total" && echo "✓ API transacciones works!" || echo "✗ API transacciones failed"

# Kill server
echo ""
echo "🛑 Stopping server..."
kill $SERVER_PID

echo ""
echo "✅ Demo test completed!"