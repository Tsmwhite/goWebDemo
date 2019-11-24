package controllers

import (
	"github.com/gin-gonic/gin"
	"regexp"
	"strconv"
	"strings"
	"time"
	"z3/model"
)

func ChangeHeadimg(c *gin.Context){
	headimgFile, err := c.FormFile("headimg")
	if err == nil {
		//匹配文件类型是否为图片
		imgReg := regexp.MustCompile("^image")
		//转为小写
		content_type := headimgFile.Header["Content-Type"][0]
		content_type = strings.ToLower(content_type)
		if ! imgReg.MatchString(content_type) {
			ResError("文件格式不正确",c)
			return
		}
		//保存文件
		timeUnix := strconv.FormatInt(time.Now().UnixNano(),10)
		fileName := UPLOAD_HEADIMG_PATH+timeUnix + "_" +strconv.Itoa(MemberId) + ".png"
		err = c.SaveUploadedFile(headimgFile, fileName)
		if err != nil {
			ResError(err,c)
			return
		}
		update := map[string]interface{}{
			"id":MemberId,
			"headimg":fileName,
		}
		if err:= model.Update("z3_member",update); err != nil {
			ResError(err,c)
			return
		}
	}else{
		ResError(err,c)
		return
	}
	ResSucces("头像修改完成",c)
}

//修改密码
func ChangePassword(c *gin.Context){
	old_password := c.PostForm("password")
	new_password := c.PostForm("password_")
	confirm_pass := c.PostForm("password_confirm")

	defer func() {
		if e := recover(); e != nil {
			ResError(e,c)
		}
	}()

	if old_password == ""{
		panic("请输入密码")
	}

	if new_password == ""{
		panic("请输入新密码")
	}

	if old_password == ""{
		panic("请确认密码")
	}

	if PassWrodMd5(old_password,MEMBER_INFO.Salt) != MEMBER_INFO.Password{
		panic("密码不正确")
	}

	if new_password != confirm_pass{
		panic("确认密码不正确")
	}

	salt := GetSaLt()
	update := map[string]interface{}{
		"id":MemberId,
		"salt":salt,
		"password": PassWrodMd5(new_password,salt),
	}
	err := model.Update("z3_member",update)
	if err != nil {
		panic(err)
	}
	ResSucces("修改完成",c)
}
