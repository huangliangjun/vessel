package middleware

import (
	"gopkg.in/macaron.v1"
)

// InitLog log init
func InitLog(runmode, path string) {
	//log := module.GetLog()
	//
	//log.SetLevelByString("info")
	//
	//if runmode == "dev" {
	//	log.SetLevelByString("debug")
	//}
	//
	//log.SetOutputByName(path)
}

func logger(runmode string) macaron.Handler {
	//log := module.GetLog()
	//return func(ctx *macaron.Context) {
	//	if runmode == "dev" {
	//		log.Debug("")
	//		log.Debug("----------------------------------------------------------------------------------")
	//	}
	//
	//	log.Infof("[%s] [%s]", ctx.Req.Method, ctx.Req.RequestURI)
	//	log.Infof("[Header] %v", ctx.Req.Header)
	//}
	return func(ctx *macaron.Context) {}
}
