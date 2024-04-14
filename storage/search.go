package storage

import (
	"Rvvijays/logapi/util"
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/s3"
)

// search given keyword in all files inside all given folders
func Search(folders []string, searchValue string) ([]string, error) {
	logs := []string{}

	newSession, err := s3Configure(util.Config.Credentials)
	if err != nil {
		// fmt.Println("errror while creating session with s3", err)
		return nil, err
	}

	newClient := s3.New(newSession)

	for _, folder := range folders {

		folder += "/"

		// fmt.Println("checking in folder", folder)

		searchedLogs, err := searchInFolder(newClient, util.Config.BucketName, folder, searchValue)
		if err != nil {
			fmt.Println("errow in search folder", err)
			continue
		}

		logs = append(logs, searchedLogs...)
	}

	return logs, nil
}

// read all files in a folder and returns matching lines with the keywords.
func searchInFolder(newClient *s3.S3, bucket, folder, searchText string) ([]string, error) {
	resp, err := newClient.ListObjectsV2(&s3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(folder),
	})

	if err != nil {
		fmt.Println("error while listing files.", err)
		return nil, err
	}

	// Iterate through each object (file) in the folder
	logs := []string{}
	for _, obj := range resp.Contents {

		// Create a reader to stream the file content from S3
		fileContentReader, err := newClient.GetObject(&s3.GetObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(*obj.Key),
		})

		if err != nil {
			fmt.Println("Error creating file content reader:", err)
			continue
		}

		// fmt.Println("checking in file", *obj.Key)
		searchedLogs, err := getLineWithText(fileContentReader.Body, searchText)
		if err != nil {
			fmt.Println("error while reading file", err)
			continue
		}

		logs = append(logs, searchedLogs...)

	}

	return logs, nil
}

// returns all the lines in which given text matches
func getLineWithText(reader io.Reader, searchText string) ([]string, error) {

	lines := []string{}
	scanner := bufio.NewScanner(reader)
	for scanner.Scan() {
		line := scanner.Text()
		if strings.Contains(strings.ToLower(line), strings.ToLower(searchText)) {
			lines = append(lines, line)
		}
	}
	if err := scanner.Err(); err != nil {
		return lines, err
	}
	return lines, nil
}

// Function to retrieve file content from S3
// func getFileContent(svc *s3.S3, bucket, key string) (string, error) {
// 	resp, err := svc.GetObject(&s3.GetObjectInput{
// 		Bucket: aws.String(bucket),
// 		Key:    aws.String(key),
// 	})
// 	if err != nil {
// 		return "", err
// 	}
// 	defer resp.Body.Close()

// 	var content strings.Builder
// 	_, err = io.Copy(&content, resp.Body)
// 	if err != nil {
// 		return "", err
// 	}

// 	return content.String(), nil
// }

// func containsText(reader io.ReadCloser, searchText string) bool {
// 	var buf strings.Builder
// 	_, err := io.Copy(&buf, reader)
// 	if err != nil {
// 		log.Println("Error reading file content:", err)
// 		return false
// 	}
// 	return strings.Contains(buf.String(), searchText)
// }
