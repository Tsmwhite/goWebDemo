package controllers

//用户登陆注册

import (
	"net/http"
	time2 "time"
	"z3/model"
	"github.com/gin-gonic/gin"
	"regexp"
)

const (
	MobileReg = `^1([38][0-9]|14[579]|5[^4]|16[6]|7[1-35-8]|9[189])\d{8}$`
)

//登录权限验证
func AuthVerify() gin.HandlerFunc{
	return func(c *gin.Context) {
		defer func() {
			if catch := recover();catch != nil{
				ResError(catch,c)
			}
		}()
		res := CheckUser(c.Request)
		if res != ""{
			panic(res)
		}
	}
}

//请求头验证userkey
func CheckUser(req *http.Request) string{
	header := req.Header
	if header == nil || header["User_key"] == nil || header["User_key"][0] == "" {
		return  "暂无权限进行此操作"
	}

	if err := ResUserKeyMsg(header["User_key"][0]);err != "" {
		return err
	}

	return ""
}



//用户登录
func LoginAuth(c *gin.Context){
	var (
		member  *model.Member
		err error
	)
	mobile := c.PostForm("mobile")
	password := c.PostForm("password")
	if mobile == "" {
		ResError("请输入11位手机号",c)
		return
	}

	rgx := regexp.MustCompile(MobileReg)
	if ! rgx.MatchString(mobile) {
		ResError("手机号码格式不正确",c)
		return
	}

	if password == "" {
		ResError("请输入密码",c)
		return
	}

	//查询用户信息
	member,err = model.GetMemberInfo("mobile",mobile)
	if err == nil{
		if member.Id > 0{
			//验证密码
			password = PassWrodMd5(password,member.Salt)
			if password == member.Password{
				AfterLoginRes(*member,c)
				return
			}
			ResError("密码不正确",c)
			return
		}else{
			ResError("账户不存在",c)
			return
		}
	}else{
		ResError("账户不存在",c)
		return
	}
}

func AfterLoginRes(member model.Member,c *gin.Context){
	resData := struct {
		User_key 	string
		Member_id 	int
		Nickname  	string
	}{
		User_key: GetUserKey(member.Id),
		Member_id:member.Id,
		Nickname:member.Name,
	}
	ResSucces(resData,c)
}

//用户注册
func Register (c *gin.Context) {
	mobile := c.PostForm("mobile")
	password := c.PostForm("password")
	password_ := c.PostForm("password_")
	if mobile == "" {
		ResError("请输入账户名",c)
		return
	}

	rgx := regexp.MustCompile(MobileReg)
	if ! rgx.MatchString(mobile) {
		ResError("账户格式不正确（请输入11位手机号）",c)
		return
	}

	if password == "" {
		ResError("请输入密码",c)
		return
	}

	if password_ == "" {
		ResError("请确认密码",c)
		return
	}

	if password != password_ {
		ResError("两次密码输入不一致",c)
		return
	}

	//查询手机号是否已存在
	_,unique := model.GetMemberInfo("mobile",mobile)
	if unique == nil{
		ResError("该手机号已注册",c)
		return
	}

	//当前时间戳
	time := time2.Now().Unix()
	//当前ip
	ip := c.Request.RemoteAddr

	//默认手机号位昵称
	name := mobile[0 : 3]+"****"+mobile[7: ]
	salt := GetSaLt()
	password = PassWrodMd5(password,salt)
	memberInfo := &model.Member{
		Name:name,
		Salt:salt,
		Mobile:mobile,
		Password:password,
		RegisterTime:time,
		LoginTime:time,
		RegisterIp:ip,
		LoginIp:ip,
	}
	err := memberInfo.Insert()
	if err == nil{
		//查询用户信息并
		memberInfo,_ := model.GetMemberInfo("mobile",mobile)
		AfterLoginRes(*memberInfo,c)
		return
	}else{
		ResError("注册失败，请重试-",c)
	}
	return
}