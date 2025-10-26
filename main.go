package main

import (
	"zetian-personal-website-hertz/biz/config"
	SES_email "zetian-personal-website-hertz/biz/pkg/SES_email"
	"zetian-personal-website-hertz/biz/repository"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/cors"
)

func main() {
	config.InitConfig() //初始化配置
	repository.InitPostgres() //初始化数据库
	SES_email.InitSES() //初始化SES 发邮件服务
	
	h := server.Default()
    h.Use(cors.Default())
	register(h)

	h.Spin()
}