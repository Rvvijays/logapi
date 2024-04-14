package util

import (
	"bufio"
	"encoding/json"
	"fmt"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/valyala/fasthttp"

	cmap "github.com/orcaman/concurrent-map"
)

var LogsMap cmap.ConcurrentMap

// JSONEncode ...
func JSONEncode(data interface{}) ([]byte, error) {
	if data == nil {
		return []byte{}, nil
	}
	jsonBytes, err := json.MarshalIndent(data, "", " ")
	if err != nil {
		return []byte{}, err
	}
	return jsonBytes, nil
}

// JSONDecode ...
func JSONDecode(jsonBytes []byte, data interface{}) error {
	if len(jsonBytes) <= 0 || jsonBytes == nil {
		return nil
	}
	err := json.Unmarshal(jsonBytes, &data)
	if err != nil {
		return err
	}

	return nil

}

// SuccessResponse ...
func SuccessResponse(ctx *fasthttp.RequestCtx) {
	ctx.Response.Header.SetContentType("application/json")
	ctx.Response.SetStatusCode(http.StatusOK)
	data := []byte(`{"message":"success"}`)
	ctx.Response.SetBody(data)

}

func JSONResponse(ctx *fasthttp.RequestCtx, data []byte, customErr *APIError) error {

	ctx.Response.Header.SetContentType("application/json")

	if customErr != nil {

		ctx.Response.SetStatusCode(customErr.Code)
		errBytes, err := JSONEncode(customErr)
		ctx.Response.SetBody(errBytes)
		return err
	}
	ctx.Response.SetStatusCode(http.StatusOK)
	ctx.Response.SetBody(data)
	return nil
}

// Returns YYYY-MM-DD_HH
func CurrentLogFile() string {
	t := time.Now().UTC()
	name := t.Format(Config.TimeLayout) + "_" + strconv.Itoa(t.Hour())
	return name
}

// Returns all folders between startTime and endTime
func GenerateAllFolders(startTime time.Time, endTime time.Time) ([]string, string) {

	folders := []string{}

	folderString := ""

	currentTime := startTime
	endTime = endTime.AddDate(0, 0, 1)

	for currentTime.Before(endTime) {
		timeString := currentTime.Format(Config.TimeLayout)

		folderString += timeString + ","
		// fmt.Println("file : ", timeString)
		folders = append(folders, timeString)

		currentTime = currentTime.AddDate(0, 0, 1)
	}

	return folders, folderString

}

func SearchLogsFromMap(folderString string, searchValue string) []string {

	logs := []string{}
	mapItr := LogsMap.IterBuffered()
	for newlog := range mapItr {
		values, _ := newlog.Val.([]Log)

		logKey := strings.Split(newlog.Key, "_")[0]

		if strings.Contains(folderString, logKey) {
			for _, val := range values {
				if strings.Contains(strings.ToLower(val.Message), strings.ToLower(searchValue)) {
					logs = append(logs, val.Time+" "+val.Message)
				}
			}
		}

	}

	return logs

}

func SearchLogsFromLocal(folderString, searchValue string) []string {
	logs := []string{}
	files, err := os.ReadDir(Config.LogFileDir)
	if err != nil {
		fmt.Println("err", err)
		return logs
	}

	for _, file := range files {
		logKey := strings.Split(file.Name(), "_")[0]

		fmt.Println("log key...", logKey)

		if strings.Contains(folderString, logKey) {
			fmt.Println("matched")
			searchedLogs, err := searchInLocalFile(Config.LogFileDir+"/"+file.Name(), searchValue)
			if err != nil {
				fmt.Println("error while seaching in local file...")
				continue
			}
			logs = append(logs, searchedLogs...)
		}
	}

	return logs

}

// Function to search for a value in a file
func searchInLocalFile(filePath, searchValue string) ([]string, error) {

	fmt.Println("opening local file to search", filePath)
	logs := []string{}
	file, err := os.Open(filePath)
	if err != nil {
		fmt.Println("eror whi le opening ", err)
		return logs, err
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)
	for scanner.Scan() {
		line := scanner.Text()
		fmt.Println("line...", line)
		if strings.Contains(strings.ToLower(line), strings.ToLower(searchValue)) {
			logs = append(logs, line)
		}
	}

	if err := scanner.Err(); err != nil {
		return logs, err
	}

	return logs, nil
}
