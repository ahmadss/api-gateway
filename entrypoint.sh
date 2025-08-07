#!/bin/sh
set -e

echo "🚀 [1] Mulai startup api-gateway-1..."

echo "🔐 [0] Load ENV dari Vault Agent..."
if [ -s /vault-agent/rendered/api-gateway-1.env ]; then
  set -a
  . /vault-agent/rendered/api-gateway-1.env
  set +a
  echo "✅ ENV Loaded"
else
  echo "❌ File ENV tidak ditemukan!"
  exit 1
fi

echo "📄 [2] Render config Consul Agent..."
envsubst < /consul/config/agent-config-api-gateway.hcl.tmpl > /consul/config/agent.hcl


echo "📡 [1.2] Jalankan Consul Agent lokal..."
if [ -f /consul/config/agent.hcl ]; then
  consul agent -config-dir=/consul/config &
  echo "✅ Consul Agent running in background"
else
  echo "❌ agent.hcl tidak ditemukan!"
  exit 1
fi

echo "📡 [1.5] Konsul Addr: 172.27.0.31:8500
"
CONSUL_ADDR=http://172.27.0.31:8500

echo "📡 [1.2] Jalankan Consul Agent lokal..."
if [ -f /consul/config/agent.hcl ]; then
  consul agent -config-dir=/consul/config &
  echo "✅ Consul Agent running in background"
else
  echo "❌ agent.hcl tidak ditemukan!"
  exit 1
fi

# 🕒 Tunggu maksimal 10 detik sampai port 8500 terbuka
echo "⏳ [1.3] Menunggu Consul Agent siap di ${CONSUL_ADDR} ..."
for i in $(seq 1 10); do
  if curl -s "${CONSUL_ADDR}/v1/status/leader" | grep -q '"'; then
    echo "✅ Consul Agent siap (iterasi ke-$i)"
    break
  fi
  echo "⏳ Menunggu... ($i)"
  sleep 1
done

# Jika masih gagal
if ! curl -s "${CONSUL_ADDR}/v1/status/leader" | grep -q '"'; then
  echo "❌ Consul Agent tidak merespons pada ${CONSUL_ADDR}!"
  exit 1
fi

echo "🔍 [2] Mendaftarkan service ke Consul..."
if ! consul services register -http-addr="${CONSUL_ADDR}" /consul/api-gateway/api-gateway.hcl; then
  echo "❌ Registrasi ke Consul gagal!"
  exit 1
fi

echo "📦 [3] Apply service-defaults (Connect Mesh)..."

echo "✅ [4] Registrasi berhasil. Menjalankan binary utama..."
exec /api-gateway
