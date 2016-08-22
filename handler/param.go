package handler

import (
	"net/http"

	"gopkg.in/macaron.v1"
)

// ParamPOSTJSON param JSON data from HTTP POST
type ParamPOSTJSON struct {
}

// V1POSTParam handler for HTTP POST
func V1POSTParam(ctx *macaron.Context, param ParamPOSTJSON) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1GETListParams list handler for HTTP GET
func V1GETListParams(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1PUTParam handler for HTTP PUT
func V1PUTParam(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1GETParam handler for HTTP GET
func V1GETParam(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1DELETEParam handler for HTTP DELETE
func V1DELETEParam(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}
