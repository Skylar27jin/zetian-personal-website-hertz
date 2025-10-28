# 需要更新thrift: 
# 不在 `$GOPATH` 下的项目通过工具提供的 `-module` 命令指定一个自定义 module 名称即可：
hz new -module zetian-personal-website-hertz -idl idl/base.thrift -force
go mod tidy
# 查看 go.mod 中 github.com/apache/thrift 版本是否为 v0.13.0，如果不是则继续执行 2.2 小节剩余代码
go mod edit -replace github.com/apache/thrift=github.com/apache/thrift@v0.13.0
go mod tidy
# 更新 thrift
hz update -idl idl/base.thrift



# 官方指导
## 如何用thrift开始开发
https://www.cloudwego.io/zh/docs/hertz/tutorials/toolkit/usage-thrift/
## 开发指南：Basic Feature
https://www.cloudwego.io/docs/hertz/tutorials/example/


# 设置AWS CLI:

Why we need this: 
to let hertz send verification code to user, aws needs to verify that the action is done by an account with permission.
the cli helps aws verify the account.
Let's download aws cli
https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html

then set up IAM locally(ask zetian for more details)

# 获取config:
check biz/config, we can see config.go is trying to read many sensitive secret keys from some .yaml files.
Those files are git ignored, so ask zetian for more details

# 如何call hertz:
# 本地测试
curl "http://localhost:8888/to_binary?number=64"
# 测试EC2 服务器
curl "https://api.skylar27.com/to_binary?number=64"


# 本地起服务：

Windows
$env:ENV = "dev"
go run .

Mac
ENV=dev go run .


# 编译:

windows
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o server-linux

mac
GOOS=linux GOARCH=amd64 go build -o server-linux


# 云端EC2(ubuntu)起服务：
./server-linux


# 云端EC2(ubuntu)升级服务：
## 目前网站在关闭terminal后会持续运行（已被“持久化”）

1. 找到旧进程 PID
ps aux | grep server-linux
一帮长这样：
ubuntu      4919  0.0  1.1 1240800 11496 ?       Sl   Oct19   1:14 ./server-linux


2. 停止旧进程（假设 PID 为 12345）
kill 12345

3. 在自己电脑上编译

windows
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o server-linux

mac
GOOS=linux GOARCH=amd64 go build -o server-linux


4. 更新权限
chmod +x server-linux

5. 持久化运行服务：
nohup env ENV=prod ./server-linux > server.log 2>&1 &

6. 检查运行状态

ps aux | grep server-linux
tail -f server.log


