#!/bin/bash

echo "🚀 EJECUTANDO DEMO SISTEMA CONTABLE LATAM PARA SIIGO"
echo "==================================================="
echo ""

# Limpiar puerto
echo "🧹 Limpiando puerto 8080..."
lsof -ti :8080 | xargs kill -9 2>/dev/null
sleep 2

# Ejecutar demo
echo "✅ Iniciando servidor..."
echo ""
echo "📋 URLs disponibles:"
echo "   - http://localhost:8080 - Página principal"
echo "   - http://localhost:8080/demo - Demo automática"
echo "   - http://localhost:8080/api/transacciones - API REST"
echo ""
echo "🎯 Value Proposition Siigo:"
echo "   • 18 meses → 2 meses"
echo "   • $500K → $150K"
echo "   • 7 sistemas → 1 DSL"
echo "   • ROI: 1,020%"
echo ""
echo "Presiona Ctrl+C para detener el servidor"
echo ""

# Ejecutar R2Lang
go run main.go examples/proyecto/contable/demo_final.r2