package handler

import (
	"net/http"

	"gopkg.in/macaron.v1"
)

// StagePOSTJSON stage JSON data for HTTP POST
type StagePOSTJSON struct {
}

// V1POSTStage handler for HTTP POST
func V1POSTStage(ctx *macaron.Context, stage StagePOSTJSON) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1PUTStage handler for HTTP PUT
func V1PUTStage(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1GETStage handler for HTTP GET
func V1GETStage(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}

// V1DELETEStage handler for HTTP DELETE
func V1DELETEStage(ctx *macaron.Context) (int, []byte) {
	return http.StatusOK, []byte("")
}
