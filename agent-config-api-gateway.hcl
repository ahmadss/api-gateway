node_name   = "api-gateway-node"
datacenter  = "dc1"
data_dir    = "/consul/data"

bind_addr   = "0.0.0.0"
client_addr = "0.0.0.0"
advertise_addr = "172.27.0.31"

retry_join = ["172.27.0.5", "172.27.0.15", "172.27.0.25"]

ports {
  http = 8500
  grpc = 8504
}

connect {
  enabled = true
}

enable_script_checks = true
enable_central_service_config = true
