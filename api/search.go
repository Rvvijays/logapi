package api

import (
	"Rvvijays/logapi/storage"
	"Rvvijays/logapi/util"
	"fmt"
	"net/http"
	"time"

	"github.com/valyala/fasthttp"
)

type searchReq struct {
	SearchKeyword string `json:"searchKeyword"`
	From          string `json:"from"`
	To            string `json:"to"`
}

// Returns logs between time frame with search key word
func SearchLog(ctx *fasthttp.RequestCtx) {

	body := searchReq{}

	util.JSONDecode(ctx.PostBody(), &body)

	// fmt.Println("from time....", body.From)

	startTime, err := time.Parse("2006-01-02", body.From)
	if err != nil {
		// fmt.Println("error while parsing time")
		util.JSONResponse(ctx, nil, &util.APIError{
			Code:     http.StatusBadRequest,
			Status:   http.StatusText(http.StatusBadRequest),
			Internal: err.Error(),
			Message:  "Invalid time format. Please provide time in 2006-01-02 format.",
		})
		return
	}

	endTime, err := time.Parse("2006-01-02", body.To)
	if err != nil {
		// fmt.Println("error while parsing time")
		util.JSONResponse(ctx, nil, &util.APIError{
			Code:     http.StatusBadRequest,
			Status:   http.StatusText(http.StatusBadRequest),
			Internal: err.Error(),
			Message:  "Invalid time format. Please provide time in 2006-01-02 format.",
		})
		return
	}

	//

	// get all folders where logs would be stored and a string concat of all folder
	folders, folderString := util.GenerateAllFolders(startTime, endTime)

	// fmt.Println("searching in folders", folderString)
	var logs []string

	if util.Config.ServerStorage {
		logs, err = storage.Search(folders, body.SearchKeyword)
		if err != nil {
			fmt.Println("error while search", err)
			util.JSONResponse(ctx, nil, &util.APIError{
				Code:     http.StatusInternalServerError,
				Status:   http.StatusText(http.StatusInternalServerError),
				Message:  "Error searching keyword.",
				Internal: err.Error(),
			})
			return
		}
	}

	// search from local if exists.
	searchedLogs := util.SearchLogsFromLocal(folderString, body.SearchKeyword)
	logs = append(logs, searchedLogs...)

	// search in log map if exists.
	searchedLogs = util.SearchLogsFromMap(folderString, body.SearchKeyword)
	logs = append(logs, searchedLogs...)

	resp := make(map[string]interface{})
	resp["logs"] = logs

	respBytes, _ := util.JSONEncode(resp)

	util.JSONResponse(ctx, respBytes, nil)
}
