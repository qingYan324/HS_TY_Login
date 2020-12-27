package main

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"net"
	"os"

	"github.com/kirinlabs/HttpRequest"
	"io"
	"log"
	"strconv"
	"strings"
	"time"



)

var clientip string
var nasip string
var mac string

var secret="Eshore!@#"
var username string
var password string
var timestamp=""
type send01 struct {
	iswifi string `json:"challenge"`
	clientip string `json:"resinfo"`
	nasip string `json:""`
	mac string `json:""`
	timestamp string `json:""`
	authenticator string `json:""`
	username string `json:""`


}
type ret01 struct {
	Challenge string `json:"challenge"`
	Resinfo string `json:"resinfo"`
	Rescode string `json:"rescode"`
}
type ret02 struct {
	Resinfo string `json:"resinfo"`
	Rescode string `json:"rescode"`
}
func gettimestamp()string{
	timeUnix:=time.Now().Unix()*1000
	string:=strconv.FormatInt(timeUnix,10)

	return string
}
func md5V3(str string) string {
	w := md5.New()
	io.WriteString(w, str)
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	return md5str
}

func getVerifyCodeString()string{
	var code string
	var timestamp=gettimestamp()
	fmt.Println("第一个时间戳")
	fmt.Println(timestamp)
	var tmp = clientip +nasip+ mac +timestamp+ secret
	var md5string=strings.ToUpper(md5V3(tmp))
	fmt.Println(md5string)
	//实例化
	req := HttpRequest.NewRequest()
    //定义头
	req.SetHeaders(map[string]string{
		"Content-Type":"application/json",
		//"Content-Length":"XXX",
	})

	/*
	postData := map[string]interface{}{
		"id":    1,
		"title": "xxx",
	}
	*/
	postData :=map[string]interface{}{
		"iswifi":"4060",
		"clientip":clientip,
		"nasip":nasip,
		"mac":mac,
		"timestamp":timestamp,
		"authenticator":md5string,
		"username":username,
	}
	//转换json数据
	mjson, err := json.Marshal(postData)

	mString :=string(mjson)
	fmt.Println("第一次传出去的数据")
	fmt.Println(mString)

	resp,err := req.Post("http://61.140.12.23:10001/client/challenge",mString)

	if err != nil {
		fmt.Println(err)
		return ""
	}

	if resp.StatusCode() == 200 {
		body, err := resp.Body()

		if err != nil {
			log.Println(err)
			return ""
		}

		fmt.Println(string(body))
		var a=string(body)

		//
		println(a)

		var reet01 ret01
		 err2 :=json.Unmarshal(body,&reet01)
		if err2 !=nil {
			fmt.Println(err2)

		}
		code=reet01.Challenge



	}

	//body,err := res.Body()

	//fmt.Println(res)
	//fmt.Println(err)
	//fmt.Println(postData)

	return code
}
func login(vertifyCode string) {
	var timestamp = gettimestamp()
	fmt.Println("第二个时间戳")
	fmt.Println(timestamp)
	var tmp= clientip + nasip + mac + timestamp +vertifyCode+ secret
	fmt.Println("看看这里")
	fmt.Println(tmp)
	var md5string = strings.ToUpper(md5V3(tmp))
	fmt.Println(md5string)
	postData := map[string]interface{}{
		"password":password,
		"verificationcode":vertifyCode,
		"iswifi":"4060",
		"clientip":clientip,
		"nasip":nasip,
		"mac": mac,
		"timestamp":timestamp,
		"authenticator":md5string,
		"username":username,
	}
	println(postData)
	//实例化
	//转换json数据
	mjson, err := json.Marshal(postData)

	mString :=string(mjson)
	fmt.Println("第二次传出去的数据")
	fmt.Println(mString)
	req := HttpRequest.NewRequest()
	//定义头
	req.SetHeaders(map[string]string{
		"Content-Type": "application/json",
		//"Content-Length":"XXX",
	})
	resp,err := req.Post("http://61.140.12.23:10001/client/login",mString)


	if err != nil {
		fmt.Println(err)
		return
	}

	if resp.StatusCode() == 200 {
		body, err := resp.Body()

		if err != nil {
			log.Println(err)
			return
		}
		//fmt.Println(body)
		fmt.Println("第二次服务器回复")
		fmt.Println(string(body))

	}
}
// 创建一个错误处理函数，避免过多的 if err != nil{} 出现
func dropErr(e error) {
	if e != nil {
		panic(e)
	}
}
func init(){
	//初始化第一个发送数据


}
func getmac()string{
	// 获取本机的MAC地址
	interfaces, err := net.Interfaces()
	if err != nil {
		panic("Poor soul, here is what you got: " + err.Error())
	}
	//for _, inter := range interfaces {
	//fmt.Println(inter.Name)
	inter := interfaces[0]
	mac := inter.HardwareAddr.String() //获取本机MAC地址

	fmt.Println("MAC = ", mac)
	return mac
}
func getclientip()string{

	socket, err := net.DialUDP("udp",nil,&net.UDPAddr{
		IP:  net.IPv4(1,1,1,8) ,
		Port: 3850,
	})
	if err != nil{
		fmt.Println("连接失败，err：",err)
		return "未获取到本地ip"
	}

	var tmp string=socket.LocalAddr().String()
	//var ind  =strings.Index(tmp, ":")
	tmp=tmp[:strings.Index(tmp, ":")]
	fmt.Printf("客户端链接的地址及端口是：%v\n",tmp)

	defer socket.Close()

	return tmp
}
//参入参数顺序为 账号 密码 nasip clientip mac typeA/typeB
//类型A 为本机登录  传入参数 为  账号 密码 nasip
//类型B 为协助法自定义登录例如 适用于在二级路由器下的登录或者自定义mac nasip mac 的登录
func main() {
	getclientip()
	fmt.Printf("传入的参数：\n" )
	var tmp [8] string
	for idx,args:=range os.Args{

		fmt.Printf("参数"+strconv.Itoa(idx)+":",args)
		tmp[idx]=args


	}
	//本机登录
	if tmp[1]=="typeA" {
		username = tmp[2]
		password = tmp[3]
		clientip = getclientip()
		nasip =tmp[4]
		mac =getmac()
	}
	//协助登录（协助路由器登录、自定义参数登录）
	if tmp[1]=="typeB" {
		username = tmp[2]
		password = tmp[3]
		nasip = tmp[4]
		clientip =tmp[5]
		mac =tmp[6]
	}



	fmt.Printf("传入参数结束\n")

	var code =getVerifyCodeString()
	println(code)
	login(code)


	//login(getVerifyCodeString())

}
