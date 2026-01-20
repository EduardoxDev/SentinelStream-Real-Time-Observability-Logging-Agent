#!/bin/bash

echo "🔧 Generating gRPC code from proto files..."

# Criar diretório de saída
mkdir -p proto/gen

# Gerar código Go
protoc --go_out=proto/gen --go_opt=paths=source_relative \
       --go-grpc_out=proto/gen --go-grpc_opt=paths=source_relative \
       proto/metrics.proto

echo "✅ gRPC code generated successfully!"
echo "Files created in proto/gen/"
