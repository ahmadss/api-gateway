docker build -t api-gateway:latest .


docker compose -f docker-compose-api-gateway-multi.yml up -d

docker ps