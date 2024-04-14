package api

import (
	"Rvvijays/logapi/util"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

// Function to all incoming logs
func AddLog(ctx *fasthttp.RequestCtx) {

	var log util.Log

	err := util.JSONDecode(ctx.PostBody(), &log)
	if err != nil {
		// fmt.Println("error while reading body")
		util.JSONResponse(ctx, nil, &util.APIError{
			Code:     http.StatusBadRequest,
			Status:   http.StatusText(http.StatusBadRequest),
			Internal: err.Error(),
			Message:  "Bad request. Unable to decode the request body.",
		})
		return
	}

	log.Time = time.Now().UTC().Format("2006-01-02 15:04:05")

	logMapKey := util.CurrentLogFile()

	// fmt.Println("log key..", logMapKey)

	// checking if present in the concurrent map or not
	if logMapValue, ok := util.LogsMap.Get(logMapKey); ok {

		logs, _ := logMapValue.([]util.Log)
		// fmt.Println("witohut inter", logMapValue)

		// fmt.Println("existing logs...", logs, err)

		logs = append(logs, log)

		// fmt.Println("after appending logs..", logs)

		util.LogsMap.Set(logMapKey, logs)

		// fmt.Println("added log in existing current hour log")

		util.SuccessResponse(ctx)
		return
	}

	logs := []util.Log{
		log,
	}

	util.LogsMap.Set(logMapKey, logs)
	util.SuccessResponse(ctx)

	// fmt.Println("added log in new hour.")

}
