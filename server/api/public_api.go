package api

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/service"
	"AirGo/utils/encode_plugin"
	"AirGo/utils/mail_plugin"
	"AirGo/utils/response"
	"fmt"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 邮箱验证码
func GetMailCode(ctx *gin.Context) {

	var u model.User
	err := ctx.ShouldBind(&u)
	if err != nil {
		response.Fail("邮箱验证码参数错误"+err.Error(), nil, ctx)
		return
	}
	_, ok := global.LocalCache.Get(u.UserName + "emailcode")
	if ok {
		response.Fail("邮箱验证码获取频繁", nil, ctx)
		return
	}
	//用户是否存在且是否有效
	user, err := service.FindUserByEmail(&u)
	if err == gorm.ErrRecordNotFound || user == nil {
		response.Fail("用户不存在"+err.Error(), nil, ctx)
		return
	} else if !user.Enable {
		response.Fail("用户被封禁", nil, ctx)
		return
	}
	//生成验证码
	randomStr := encode_plugin.RandomString(4) //4位随机数
	mail_plugin.SendEmail(u.UserName, randomStr)
	//验证码存入local cache
	global.LocalCache.Set(u.UserName+"emailcode", randomStr, 60000000000) //60秒过期
	//发送邮件
	err = mail_plugin.SendEmail(user.UserName, randomStr)
	if err != nil {
		fmt.Println("验证码获取失败:", err.Error())
		response.OK("验证码获取失败"+err.Error(), nil, ctx)
		return
	}
	response.OK("验证码获取成功", nil, ctx)
}
