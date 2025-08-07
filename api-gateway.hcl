service {
  name    = "api-gateway"
  id      = "api-gateway-1"
  address = "172.27.0.31"  # ✅ ganti dari 127.0.0.1

  port    = 9000

  connect {
    sidecar_service {
        proxy {
            upstreams = [
             {
                destination_name = "auth-public-gateway"
                local_bind_port  = 21003
             }
            ]
      }
    }  # ✅ cukup ini
  }

  check {
    id       = "http-check"
    name     = "HTTP Health Check"
    http     = "http://172.27.0.31:9000/health"  # ✅ cocokkan
    interval = "10s"
    timeout  = "2s"
  }
}
