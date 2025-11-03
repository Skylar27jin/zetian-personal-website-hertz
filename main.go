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
	h.Use(cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",   // 本地调试
			"https://skylar27.com",    // ✅ 线上正式域名
		},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // ✅ 允许跨域携带 Cookie
	}))
	register(h)

	h.Spin()
}