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


# 如何call hertz:
# 本地测试
curl "http://localhost:8888/to_binary?number=64"
# 测试EC2 服务器
curl "https://api.skylar27.com/to_binary?number=64"


# 本地起服务：
## MAC设置环境变量
export ENV=dev 
## Windows设置环境变量
$env:ENV = "dev"
go run .

# 编译:

GOOS=linux GOARCH=amd64 go build -o server-linux

# 云端(ubuntu)起服务：
./server-linux
