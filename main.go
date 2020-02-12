package main

import (
	"github.com/go-ini/ini"
	"log"
	check "lostcloud/check"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"time"
)

func main(){
	var (
		cfg *ini.File
		cfg_filepath string
		login_url string
		check_url string
		email string
		passwd string
	)
	//设置log句柄
	fileName:="lostcloud.log"
	logFile,err:=os.Create(fileName)
	if err!=nil{
		log.Fatalln("open file error!")
	}
	logg :=log.New(logFile,"[Debug]",log.LstdFlags)
	logg.SetPrefix("[Info]")
	logg.SetFlags(logg.Flags()|log.LstdFlags)

	current_dir, _ := os.Getwd()
	logg.Printf("current_dir========>%v",current_dir)
	cfg_filepath=current_dir+"/lostcloud.ini"
	//读取配置文件
	login_url,check_url,email,passwd=check.Read_params(cfg,cfg_filepath,logg,login_url,check_url,email,passwd)

	if login_url=="login_url"|| check_url=="check_url"|| email=="email"|| passwd=="passwd"{
		logg.Printf("ini配置文件未填写")
		os.Exit(3)//退出程序
	}

	ticker:=time.NewTicker(60*60*24* time.Second)
	defer ticker.Stop()
	for {
		select{
		case <-ticker.C:

			jar, err := cookiejar.New(nil)
			if err != nil {
				logg.Fatal("cookiejar.New  error: %v",err)
			}
			//初始化client,保持cookie
			var client = &http.Client{Transport: nil, CheckRedirect: nil, Jar: jar}
			//调用请求
			login_data := url.Values{} //初始化post参数
			login_data.Add("email", email)
			login_data.Add("passwd", passwd)
			var resp = check.Post(login_url, login_data,client,logg)
			check.LostCloud_Login(logg,resp)

			check_data := url.Values{} //初始化post参数
			var resp2=check.Post(check_url,check_data,client,logg)
			check.LostCloud_Check(logg,resp2)
		}
	}

}
