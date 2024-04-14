package automated

import (
	"Rvvijays/logapi/storage"
	"Rvvijays/logapi/util"
	"fmt"
	"log"
	"os"
	"time"
)

// function to write logs in local file from memory in a given time interval.
func LogMapToLocalFile() {
	for {
		time.Sleep(time.Duration(util.Config.WriteLogInterval) * time.Minute)

		// fmt.Println("current file..", util.CurrentLogFile())

		mapItr := util.LogsMap.IterBuffered()
		for newlog := range mapItr {
			values, _ := newlog.Val.([]util.Log)
			err := updateLogInFile(newlog.Key, values)
			if err != nil {
				fmt.Println("error while writng log", newlog.Key)
				continue
			}
			util.LogsMap.Remove(newlog.Key)
		}
	}
}

func updateLogInFile(key string, newValue []util.Log) error {
	oldLogFile, err := os.OpenFile(util.Config.LogFileDir+"/"+key+".txt", os.O_APPEND|os.O_CREATE|os.O_WRONLY, 0666)
	if err != nil {
		fmt.Println("error while opening log file.", util.Config.LogFileDir+"/"+key+".txt")
		return err
	}

	defer oldLogFile.Close()

	for _, lo := range newValue {
		logger := log.New(oldLogFile, lo.Time+" ", log.LUTC)
		logger.Println(lo.Message)
	}

	fmt.Println("updated in local file")
	return nil
}

// function to upload local logs to storage server in given interval
func LocalFileToServer() {
	for {
		time.Sleep(time.Duration(util.Config.UploadLogInterval) * time.Minute)
		// fmt.Println("opening log dir", util.Config.LogFileDir)

		files, err := os.ReadDir(util.Config.LogFileDir)
		if err != nil {
			fmt.Println("err while opening logs dir", err)
			continue
		}

		tobeuploadFiles := []string{}
		for _, file := range files {
			// except current file as it is being written for the current hour.
			if file.Name() != util.CurrentLogFile() {
				// upload it to server
				tobeuploadFiles = append(tobeuploadFiles, file.Name())
			}
		}

		if len(tobeuploadFiles) > 0 {
			storage.UploadFiles(tobeuploadFiles)
			fmt.Println("file uploaded to server")
		}
	}

}
