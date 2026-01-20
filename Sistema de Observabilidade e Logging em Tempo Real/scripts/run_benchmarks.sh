#!/bin/bash

echo "🚀 Running Observability System Benchmarks"
echo "=========================================="

# Verificar se a infraestrutura está rodando
if ! docker ps | grep -q redis; then
    echo "⚠️  Redis não está rodando. Iniciando infraestrutura..."
    docker-compose up -d
    sleep 5
fi

echo ""
echo "📊 Running Benchmarks..."
echo ""

# Executar benchmarks
go test -bench=. -benchmem -benchtime=10s ./benchmarks/ | tee benchmark_results.txt

echo ""
echo "🔥 Running Stress Tests..."
echo ""

# Executar testes de stress
go test -v -run=TestStress ./benchmarks/ | tee -a benchmark_results.txt

echo ""
echo "📈 Generating Memory Profile..."
go test -bench=BenchmarkMemoryAllocation -memprofile=mem.prof ./benchmarks/
go tool pprof -text mem.prof > memory_profile.txt

echo ""
echo "⚡ Generating CPU Profile..."
go test -bench=BenchmarkMetricsCollection -cpuprofile=cpu.prof ./benchmarks/
go tool pprof -text cpu.prof > cpu_profile.txt

echo ""
echo "✅ Benchmarks completed!"
echo "Results saved to:"
echo "  - benchmark_results.txt"
echo "  - memory_profile.txt"
echo "  - cpu_profile.txt"
