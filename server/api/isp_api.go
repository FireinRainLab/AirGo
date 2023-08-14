package api

import (
	"AirGo/global"
	"AirGo/model"
	"AirGo/service"
	"AirGo/utils/encrypt_plugin"
	"AirGo/utils/isp_plugin"
	"AirGo/utils/jwt_plugin"
	"AirGo/utils/response"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

// 获取监控
func GetMonitorByUserID(ctx *gin.Context) {
	uID, ok := ctx.Get("uID")
	if !ok || uID == nil {
		response.Fail("获取信息,uID参数错误", nil, ctx)
		return
	}
	uIDInt := uID.(int)
	isp, err := service.GetMonitorByUserID(uIDInt)
	if err == gorm.ErrRecordNotFound {
		//创建新的
		var ispNew = model.ISP{
			UserID: uIDInt,
			UnicomConfig: model.UnicomConfig{
				APPID: encrypt_plugin.RandomString(160),
			},
		}
		service.NewMonitor(&ispNew)
		isp, _ = service.GetMonitorByUserID(uIDInt)
	}
	response.OK("获取成功", isp, ctx)

}

// 发送验证码
func SendCode(ctx *gin.Context) {
	var isp model.ISP
	err := ctx.ShouldBind(&isp)
	if err != nil {
		global.Logrus.Error("运营商参数错误:", err)
		response.Fail("运营商参数错误", nil, ctx)
		return
	}
	//处理mobile
	mb, _ := encrypt_plugin.RSAEnCrypt(isp.Mobile, isp_plugin.UnicomPublicKey)
	isp.UnicomConfig.Mobile = mb
	//
	var resp string
	switch isp.ISPType {
	case "unicom":
		resp, err = isp_plugin.UnicomCode(&isp)
	case "telecom":
	}
	if err != nil {
		global.Logrus.Error("发送验证码错误:", err)
		response.Fail("发送验证码错误", nil, ctx)
		return
	}
	//判断是否为空
	respMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(resp), &respMap)
	if err != nil {
		global.Logrus.Error("resp解析错误:", err)
		response.Fail("resp解析错误", nil, ctx)
		return
	}
	if respMap["rsp_code"] != "0000" {
		msg := respMap["rsp_desc"].(string)
		if msg == "" {
			msg = "发送验证码失败"
		}
		response.Fail(msg, nil, ctx)
		return
	}
	response.OK(respMap["rsp_desc"].(string), resp, ctx)

}

// 登录运营商
func ISPLogin(ctx *gin.Context) {
	uID, ok := ctx.Get("uID")
	if !ok || uID == nil {
		response.Fail("获取信息,uID参数错误", nil, ctx)
		return
	}
	uIDInt := uID.(int)

	var isp model.ISP
	err := ctx.ShouldBind(&isp)
	if err != nil {
		global.Logrus.Error("运营商参数错误:", err)
		response.Fail("运营商参数错误", nil, ctx)
		return
	}
	if isp.ISPType == "loginAgain" {
		//清空手机号信息，重新登录
		//fmt.Println("清空手机号信息，重新登录")
		isp1, _ := service.GetMonitorByUserID(uIDInt)
		isp1.UnicomConfig.Cookie = ""
		isp1.UnicomConfig.Password = ""
		isp1.UnicomConfig.Mobile = ""
		isp1.Status = false
		go service.UpdateMonitor(isp1)
		response.OK("获取成功", isp, ctx)
		return
	}

	//fmt.Println("登录运营商", isp)
	//处理mobile,appid,验证码
	mb, _ := encrypt_plugin.RSAEnCrypt(isp.Mobile, isp_plugin.UnicomPublicKey)
	//appid, _ := encrypt_plugin.RSAEnCrypt(encrypt_plugin.RandomString(160), isp_plugin.UnicomPublicKey)
	pw, _ := encrypt_plugin.RSAEnCrypt(isp.UnicomConfig.Password, isp_plugin.UnicomPublicKey)
	isp.UnicomConfig.Mobile = mb
	//isp.UnicomConfig.APPID = appid
	isp.UnicomConfig.Password = pw
	//尝试登录
	var resp, cookie string
	switch isp.ISPType {
	case "unicom":
		resp, cookie, err = isp_plugin.UnicomCodeLogin(isp.UnicomConfig.Password, isp.UnicomConfig.Mobile, isp.UnicomConfig.APPID)
	case "telecom":
	}
	if err != nil {
		global.Logrus.Error("尝试登录错误:", err)
		response.Fail("尝试登录错误", nil, ctx)
		return
	}
	//fmt.Println("登录resp", resp)
	//fmt.Println("登录cookie", cookie)

	if resp == "" || cookie == "" {
		global.Logrus.Error("尝试登录错误,resp,cookie为空")
		response.Fail("尝试登录错误", nil, ctx)
		return
	}
	//判断响应
	respMap := make(map[string]interface{})
	err = json.Unmarshal([]byte(resp), &respMap)

	if err != nil {
		global.Logrus.Error("resp解析错误:", err)
		response.Fail("resp解析错误", nil, ctx)
		return
	}
	if respMap["code"] != "0" {
		response.Fail(respMap["dsc"].(string), nil, ctx)
		return
	}

	//处理cookie，保存isp信息
	isp.UserID = uIDInt
	isp.UnicomConfig.Cookie = cookie
	isp.Status = true
	go service.UpdateMonitor(&isp)

	//response.OK("登录成功", resp, ctx)
	response.OK("登录成功", string(json.RawMessage(resp)), ctx)
}

// 套餐查询
func QueryPackage(ctx *gin.Context) {
	token := ctx.Query("id")
	claims, err := jwt_plugin.ParseToken(token)
	if err != nil {
		response.Fail(err.Error(), nil, ctx)
		return
	}
	//log.Println("token解析后 claims.ID：", claims.ID)
	//设置user id
	uID := claims.ID
	//查询monitor
	isp, err := service.GetMonitorByUserID(uID)
	if err != nil {
		ctx.JSON(200, gin.H{
			"packageName": "查询流量失败，请重新登录",
		})
		return
	}
	//查询套餐
	resp, err := isp_plugin.UnicomQueryTraffic(isp)
	if err != nil {
		//修改monitor状态
		isp.Status = false
		go service.UpdateMonitor(isp)
		ctx.JSON(200, gin.H{
			"packageName": "查询流量失败，请重新登录",
			"mobile":      err.Error(),
		})
		return
	}
	ctx.String(200, resp) //将响应原样返回，避免json编码出现斜杠

}
