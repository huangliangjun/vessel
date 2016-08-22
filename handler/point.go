package handler

import (
	"net/http"

	"gopkg.in/macaron.v1"
)

// PointPOSTJSON point JSON data from HTTP POST
type PointPOSTJSON struct {
}

// V1POSTPoint handler for HTTP POST
func V1POSTPoint(ctx *macaron.Context, point PointPOSTJSON) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1PUTPoint handler for HTTP PUT
func V1PUTPoint(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1GETPoint handler for HTTP GET
func V1GETPoint(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1DELETEPoint handler for HTTP DELETE
func V1DELETEPoint(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}
