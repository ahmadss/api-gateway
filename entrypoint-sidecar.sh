#!/bin/sh
set -e

echo "🔐 [1] Menunggu token Vault Agent..."

for i in $(seq 1 10); do
  if [ -s /vault-agent/rendered/api-gateway-1.env ]; then
    set -a
    . /vault-agent/rendered/api-gateway-1.env
    set +a

    if [ -n "$CONSUL_HTTP_TOKEN" ]; then
      echo "✅ Token ditemukan setelah ${i}s"
      break
    fi
  fi
  echo "⏳ Tunggu token Vault (iterasi $i)..."
  sleep 1
done

if [ -z "$CONSUL_HTTP_TOKEN" ]; then
  echo "❌ Token tidak ditemukan. Exit."
  exit 1
fi

echo "🚀 Jalankan Envoy sidecar..."
# Menunggu hingga sidecar proxy service terdaftar
echo "⏳ Menunggu sidecar proxy api-gateway-1 terdaftar di Consul..."
for i in $(seq 1 10); do
  if curl -s http://localhost:8500/v1/agent/service/api-gateway-1; then
    echo "✅ Sidecar proxy ditemukan (iterasi ke-$i)"
    break
  fi
  echo "⏳ Belum terdaftar... ($i)"
  sleep 1
done

# Gagal setelah 10 detik
if ! curl -s http://localhost:8500/v1/agent/service/api-gateway-1; then
  echo "❌ Sidecar proxy api-gateway-1 belum terdaftar. Exit."
  exit 1
fi

# Jalankan Envoy
exec consul connect envoy \
  -sidecar-for api-gateway-1 \
  -token="$CONSUL_HTTP_TOKEN" \
  -http-addr=http://127.0.0.1:8500 \
  -grpc-addr=http://127.0.0.1:8504 \
  -admin-bind=0.0.0.0:19001