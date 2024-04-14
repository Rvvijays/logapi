package storage

import (
	"Rvvijays/logapi/util"
	"fmt"
	"os"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// uploads files to storage server.
func UploadFiles(files []string) {

	// create new session
	newSession, err := s3Configure(util.Config.Credentials)
	if err != nil {
		fmt.Println("err", err)
		return
	}

	newClient := s3.New(newSession)

	for _, filePath := range files {

		// read files
		file, err := os.Open(util.Config.LogFileDir + "/" + filePath)
		if err != nil {
			fmt.Println("error while opening file", filePath)
			continue
		}
		defer file.Close()

		// format it for folders
		key := strings.Replace(filePath, "_", "/", 1)

		_, err = newClient.PutObject(&s3.PutObjectInput{
			Bucket: aws.String(util.Config.BucketName),
			Key:    aws.String(key),
			Body:   file,
		})

		if err != nil {
			fmt.Println("error while uploading file", err)
			continue
		}

		err = os.Remove(util.Config.LogFileDir + "/" + filePath)
		if err != nil {
			fmt.Println("error while removing file", err)
		}

		// fmt.Println("deleted from local", filePath)
	}

}
