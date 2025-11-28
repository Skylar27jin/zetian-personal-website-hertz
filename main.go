package main

import (
	"zetian-personal-website-hertz/biz/config"
	SES_email "zetian-personal-website-hertz/biz/pkg/SES_email"
	"zetian-personal-website-hertz/biz/pkg/s3uploader"
	"zetian-personal-website-hertz/biz/repository"
	"zetian-personal-website-hertz/biz/repository/category_repo"
	"zetian-personal-website-hertz/biz/repository/school_repo"

	"github.com/cloudwego/hertz/pkg/app/server"
	"github.com/hertz-contrib/cors"
)

func main() {
	config.InitConfig() //初始化配置
	repository.InitPostgres() //初始化数据库
	SES_email.InitSES() //初始化SES 发邮件服务
	
	school_repo.InitSchoolCache() //初始化school缓存
	category_repo.InitCategoryCache()	//初始化category缓存
	s3uploader.InitS3Uploader() //初始化S3上传服务

	
	
	h := server.New(server.WithMaxRequestBodySize(16 << 20 * 20))
	//max 320MB (single picture 16MB, max 20 pictures)
	h.Use(
		cors.New(cors.Config{
		AllowOrigins: []string{
			"http://localhost:5173",   // 本地调试
			"https://skylar27.com",    // 线上正式域名
			"https://www.skylar27.com", // 线上正式域名带www
		},
		AllowMethods:     []string{"GET", "POST", "OPTIONS"},
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true, // 允许跨域携带 Cookie
	}))
	register(h)

	h.Spin()
}