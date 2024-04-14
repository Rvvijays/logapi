package api

import (
	"Rvvijays/logapi/util"
	"net/http"

	"github.com/valyala/fasthttp"
)

func GetCurrentLogs(ctx *fasthttp.RequestCtx) {

	logs := util.LogsMap.Items()

	jsonBytes, err := util.JSONEncode(logs)
	if err != nil {
		// fmt.Println("error while encoding logs")
		util.JSONResponse(ctx, nil, &util.APIError{
			Code:     http.StatusInternalServerError,
			Status:   http.StatusText(http.StatusInternalServerError),
			Internal: err.Error(),
			Message:  "Error while encoding logs.",
		})
		return
	}

	util.JSONResponse(ctx, jsonBytes, nil)

}
