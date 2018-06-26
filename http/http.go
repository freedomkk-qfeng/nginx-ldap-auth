package http

import (
	"github.com/astaxie/beego"

	"github.com/freedomkk-qfeng/nginx-ldap-auth/g"
	"github.com/freedomkk-qfeng/nginx-ldap-auth/http/controllers"
)

func Start() {
	beego.BConfig.CopyRequestBody = true
	beego.BConfig.WebConfig.Session.SessionOn = true
	beego.BConfig.WebConfig.Session.SessionName = "sessionID"
	beego.SetStaticPath("/static", "static")

	if !g.Config().Http.Debug {
		beego.SetLevel(beego.LevelInformational)
	}
	ConfigRouters()
	beego.ErrorController(&controllers.ErrorController{})
	beego.Run(g.Config().Http.Listen)
}
