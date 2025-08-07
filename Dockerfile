# ⬆️ Stage 1: Build binary (Go) untuk arsitektur arm64
FROM arm64v8/golang:1.24 AS builder
WORKDIR /app

ARG GITHUB_TOKEN
ENV GO111MODULE=on
ENV CGO_ENABLED=0

RUN git config --global url."https://${GITHUB_TOKEN}:x-oauth-basic@github.com/".insteadOf "https://github.com/"

COPY go.mod go.sum ./
RUN go mod download

COPY . .
RUN cd cmd && go build -o /api-gateway

# ⬇️ Stage 2: Final image untuk arm64
FROM arm64v8/alpine:3.18
WORKDIR /

# 1. Install semua tools yang dibutuhkan
RUN apk add --no-cache curl unzip gettext bash

# 2. Tambahkan Consul ARM64
ARG CONSUL_VERSION=1.21.1
RUN curl -Lo /tmp/consul.zip https://releases.hashicorp.com/consul/${CONSUL_VERSION}/consul_${CONSUL_VERSION}_linux_arm64.zip \
    && unzip /tmp/consul.zip -d /usr/local/bin/ \
    && chmod +x /usr/local/bin/consul \
    && rm -rf /tmp/*


# 3. Salin binary dan entrypoint
COPY --from=builder /api-gateway /api-gateway
COPY entrypoint.sh /entrypoint.sh
RUN chmod +x /entrypoint.sh

EXPOSE 9000
ENTRYPOINT ["/entrypoint.sh"]
