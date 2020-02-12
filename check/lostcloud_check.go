package check

import (
	"bytes"
	"encoding/json"
	"github.com/go-ini/ini"
	"os"

	"io"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"
)


type Login_params struct {
	Email string
	password string
}

//
type resp_struct struct {
	Code    int               `json:"code"`
	Message string            `json:"message"`
	Data    map[string]string `json:"data"`
}


//response body处理方法
func resp_text(resp *http.Response,logg *log.Logger) string {
	var buffer [512]byte             //缓冲块，重复读写
	contents := bytes.NewBuffer(nil) //byte类型的缓冲区
	for {
		n, err := resp.Body.Read(buffer[0:])
		contents.Write(buffer[0:n])
		if err != nil && err == io.EOF {
			break
		} else if err != nil {
			logg.Printf("resp.Body.Read error =========>%v",err)
			panic(err)
		}
	}
	//println(contents.String())
	return contents.String()
}

//response body处理方法
func resp_json(resp *http.Response,logg *log.Logger)(msg resp_struct){

			//判断json是否decode成功
			if err := json.NewDecoder(resp.Body).Decode(&msg); err != nil {
				logg.Fatal("response body json decode error")
			}
			return
		}


//初始化cookie
func Cookies(logg *log.Logger)(*http.Client){

	jar, err := cookiejar.New(nil)
	if err != nil {
		logg.Fatal("cookiejar.New error : %v",err)
	}
	//初始化client,保持cookie
	client := &http.Client{Transport: nil, CheckRedirect: nil, Jar: jar}
	return client
}
//只请求，不处理返回
func Post(post_url string,params url.Values,client *http.Client,logg *log.Logger)(*http.Response){

	post_data := strings.NewReader(params.Encode()) //格式化post参数
	//登录request
	req, _ := http.NewRequest("POST", post_url, post_data)
	//添加头部
	req.Header.Add("User-Agent", "Mozilla/5.0 (X11; Linux x86_64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/73.0.3683.86 Safari/537.36")
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")

	resp, err := (client).Do(req)
	if err != nil {
		logg.Printf("client Do error :%v",err)
	}

	if err != nil {
		logg.Printf("query topic failed %v", err.Error())
	}


	return resp
}

//read lastcloud.ini file
func Read_params(cfg *ini.File,cfg_filepath string,logg  *log.Logger,login_url,check_url,email,passwd string )(string,string,string,string){
	var err error
	cfg,err=ini.Load(cfg_filepath)
	if err!=nil{
		logg.Fatalf("加载ini配置文件失败:%v",err)
	}
	login_url=cfg.Section("lastcloud").Key("login_url").MustString("login_url")
	check_url=cfg.Section("lastcloud").Key("check_url").MustString("check_url")
	email=cfg.Section("lastcloud").Key("email").MustString("email")
	passwd=cfg.Section("lastcloud").Key("passwd").MustString("passwd")
	return login_url,check_url,email,passwd
}

func Log_to_file()(logg *log.Logger){
	fileName:="lostcloud.log"
	logFile,err:=os.Create(fileName)
	if err!=nil{
		log.Fatalln("open file error!")
	}
	debugLog:=log.New(logFile,"[Debug]",log.LstdFlags)
	debugLog.SetPrefix("[Info]")
	debugLog.SetFlags(debugLog.Flags()|log.LstdFlags)
	return debugLog
}

func LostCloud_Login(logg *log.Logger,resp *http.Response){

	//开始处理返回
	content_type := resp.Header["Content-Type"]

	if len(content_type) == 0 {
		logg.Printf("content_type is none")
		contents := resp_text(resp,logg)
		logg.Printf("contents string =====>%v", contents)
	} else {
		logg.Printf("content_type=======>%v", content_type[0])
		//有content-type，但不一一定是application/json
		if (strings.Contains(content_type[0], "json")) {
			//按照json处理数据
			contents := resp_json(resp,logg)
			logg.Printf("contents json =====>%v", contents)
		}else{
			//非json，按照文本处理
			contents := resp_text(resp,logg)
			logg.Printf("contents string =====>%v", contents)
		}
	}

	//最后才关闭resp
	defer resp.Body.Close()
}

func LostCloud_Check(logg *log.Logger,resp *http.Response){


	logg.Printf("code======>%v", resp.StatusCode)
	//开始处理返回
	content_type := resp.Header["Content-Type"]

	if len(content_type) == 0 {
		logg.Printf("content_type is none")
		contents := resp_text(resp,logg)
		logg.Printf("contents string =====>%v", contents)
	} else {
		logg.Printf(" content_type=======>%v", content_type[0])
		//有content-type，但不一一定是application/json
		if (strings.Contains(content_type[0], "json")) {
			//按照json处理数据
			contents := resp_json(resp,logg)
			logg.Printf("contents json =====>%v", contents)
		}else{
			//非json，按照文本处理
			contents := resp_text(resp,logg)
			logg.Printf("contents string =====>%v", contents)
		}
	}

	//最后才关闭resp
	defer resp.Body.Close()
}

