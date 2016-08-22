package web

import (
	"github.com/containerops/vessel/middleware"
	"github.com/containerops/vessel/router"
	"gopkg.in/macaron.v1"
)

// SetVesselMacaron set middle wares and routers to macaron
func SetVesselMacaron(m *macaron.Macaron) {
	//Setting Middleware
	middleware.SetMiddlewares(m)

	//Setting Router
	router.SetRouters(m)
}
