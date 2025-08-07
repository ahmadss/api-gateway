-- build image ke docker

export GITHUB_TOKEN=

docker build \
  --network=host \
  --build-arg GITHUB_TOKEN=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 \
  -t api-gateway .

docker build -t api-gateway \
  --build-arg GITHUB_TOKEN=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 .
 

-- cara daftarkan layanana ke consul
curl --request PUT --data @- http://127.0.0.1:8500/v1/agent/service/register <<EOF
{
  "ID": "api-gateway",
  "Name": "api-gateway",
  "Address": "127.0.0.1",
  "Port": 9000,
  "Check": {
    "HTTP": "http://127.0.0.1:9000/health",
    "Interval": "10s"
  }
}
EOF


-- cek service consul yang terdaftar
curl http://127.0.0.1:8500/v1/agent/services | jq

-- testing rate-limit
for i in {1..200}; do curl -s -o /dev/null -w "%{http_code}\n" http://127.0.0.1:9000/account/health; done

-- ubah project
go mod edit -module ag-primary

run 
docker run --rm -p 9000:9000 api-gateway

stop
docker stop api-gateway && docker rm api-gateway

jalankan di terminal
docker exec -it consul-agent-account sh
consul members

docker exec -it consul-agent-account consul members -http-addr=127.0.0.1:8501

cek log
docker logs -f consul-agent-account


open container
docker exec -it consul-agent-account sh


sudo nano ~/.docker/config.json

docker build -t api-gateway .

bombardier -c 1000 -d 30s -m GET http://localhost:9001/health



docker exec -it consul-agent-api-gateway sh
consul services register /consul/api-gateway/api-gateway.hcl