docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' consul-dns

docker inspect -f '{{range .NetworkSettings.Networks}}{{.IPAddress}}{{end}}' consul-server

// cek ip mac
ipconfig getifaddr en0   


docker exec -it api-gateway cat /etc/resolv.conf


docker exec -it api-gateway sh
nslookup auth-services.service.consul
