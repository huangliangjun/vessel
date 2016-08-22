package middleware

import (
	"gopkg.in/macaron.v1"
)

// SetMiddlewares set middlewares to macaron
func SetMiddlewares(m *macaron.Macaron) {
	// InitLog(setting.RunMode, setting.LogPath)

	//Set recovery handler to returns a middleware that recovers from any panics
	m.Use(macaron.Recovery())
}
