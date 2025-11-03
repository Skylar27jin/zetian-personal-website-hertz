# What is Hertz?
1. Hertz is a high performance go network framework by Bytedance. It is quite similar to Gin, or spring boot for Java.
2. Hertz's handler(or "controller" in spring boot term) is defined in a IDL(interface define language) file. (This is why Hertz is "Document-Oriendted Programming")
3. After IDL is defined, you could run a command and let Hertz generates correponding codes for you.

#Thrift
Thrift is a IDL, protobuf is another IDL. In this project, I used thrift as I have experiences writing thrift. In real life, protobuf seems to be more popular.
## Init the Hertz by Force(DO NOT do this): 

```
// 不在 `$GOPATH` 下的项目通过工具提供的 `-module` 命令指定一个自定义 module 名称即可：
hz new -module zetian-personal-website-hertz -idl idl/base.thrift -force
go mod tidy
//查看 go.mod 中 github.com/apache/thrift 版本是否为 v0.13.0，如果不是则继续执行 2.2 小节剩余代码
go mod edit -replace github.com/apache/thrift=github.com/apache/thrift@v0.13.0
go mod tidy
```

## Update thrift (if you update the thrift, do this)
```
hz update -idl idl/base.thrift
```



# Official Gudiance
## How to Initialize a Hertz project 
https://www.cloudwego.io/zh/docs/hertz/tutorials/toolkit/usage-thrift/
## Hertz：Basic Feature(set/get cookie, etc)
https://www.cloudwego.io/docs/hertz/tutorials/example/


# Setup AWS CLI:
We need AWS CLI primarily because we use AWS's Simple Email Service(SES).
To send emails from our Hertz project, we need to set up the AWS account locally. 

## How
go download aws cli（for send verification code to user） 
https://docs.aws.amazon.com/cli/latest/userguide/cli-chap-getting-started.html
then set up IAM locally(ask zetian for more details)
```
aws configure
```

# How to call hertz:
Test Locally
```
curl "http://localhost:8888/to_binary?number=64"
```
or test in Postman

Test EC2 Instance
```
curl "https://api.skylar27.com/to_binary?number=64"
```


# Run the project：
Windows
```
$env:ENV = "dev"
go run .
```
Mac
```
ENV=dev go run .
```

# Complie into Linux:
windows
```
$env:GOOS = "linux"
$env:GOARCH = "amd64"
go build -o server-linux
```
mac
```
GOOS=linux GOARCH=amd64 go build -o server-linux
```


# 云端EC2(ubuntu)升级服务：
目前网站在关闭terminal后会持续运行（已被“持久化”）, so you should first stop the initial old service, and run and 持久化 the new service.

1. 找到旧进程 PID
```
ps aux | grep server-linux
```
You will see a row like：
```
ubuntu      4919  0.0  1.1 1240800 11496 ?       Sl   Oct19   1:14 ./server-linux
```

2. Kill the old process（assume its PID is 12345）
```
kill 12345
```

3. Compile on your machine


4. 更新权限
```
chmod +x server-linux
```

5. 运行，并持久化运行服务：
```
nohup env ENV=prod ./server-linux > server.log 2>&1 &
```

7. 检查运行状态

```
ps aux | grep server-linux 
tail -f server.log
```

