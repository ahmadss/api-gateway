# 构建部分
```shell
sudo docker build --target builder -t test .
```

# 清理系统中不再使用的镜像
```shell
docker images prune
```

# 增加docker代理访问dockerhub
1. 创建 dockerd 相关的 systemd 目录，这个目录下的配置将覆盖 dockerd 的默认配置
```sh
sudo mkdir -p /etc/systemd/system/docker.service.d
```
2. 新建配置文件 /etc/systemd/system/docker.service.d/http-proxy.conf，这个文件中将包含环境变量
```sh
[Service]
Environment="HTTP_PROXY=http://proxy.example.com"
Environment="HTTPS_PROXY=http://proxy.example.com"
Environment="NO_PROXY=your-registry.com,10.10.10.10,*.example.com"
```
多个 NO_PROXY 变量的值用逗号分隔，而且可以使用通配符（*），极端情况下，如果 NO_PROXY=*，那么所有请求都将不通过代理服务器。
3. 重新加载配置文件，重启 dockerd
```sh
sudo systemctl daemon-reload
sudo systemctl restart dockersh
```
5. 检查确认环境变量已经正确配置：
```sh
sudo systemctl show --property=Environment docker
```
6. 或可从 docker info 的结果中查看配置项。
```sh
HTTP Proxy: 192.168.137.1:1090
HTTPS Proxy: 192.168.137.1:1090
```