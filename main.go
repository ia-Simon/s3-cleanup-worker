package main

import (
	"context"
	"fmt"
	"sync"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"

	"s3_cleanup_worker/core"
)

func main() {
	fileNames := make(chan string, 25)
	go fetchFileNames(fileNames)

	wg := &sync.WaitGroup{}

	for fileName := range fileNames {
		wg.Add(1)
		go cleanupTask(wg, fileName)
	}

	wg.Wait()
}

func fetchFileNames(out chan string) {
	defer close(out)

	var continuationToken *string
	for {
		data, err := core.S3Client.ListObjectsV2(
			context.TODO(),
			&s3.ListObjectsV2Input{
				Bucket:            aws.String(core.Settings.AWS_FILES_BUCKET),
				MaxKeys:           50,
				ContinuationToken: continuationToken,
			},
		)
		if err != nil {
			panic(err)
		}

		for _, item := range data.Contents {
			out <- aws.ToString(item.Key)
		}

		if !data.IsTruncated {
			return
		}
		continuationToken = data.NextContinuationToken
	}
}

func cleanupTask(wg *sync.WaitGroup, filename string) {
	defer wg.Done()

	fileInUse := checkFileUsageInDB(
		map[string][]string{
			"attachment":   {"file"},
			"professional": {"avatar"},
			"exercise":     {"image", "audio"},
			"answer":       {"image"},
		},
		filename,
	)

	if !fileInUse {
		fmt.Printf("File <[%s]> is not in use, and will be deleted.\n", filename)
		_, err := core.S3Client.DeleteObject(
			context.TODO(),
			&s3.DeleteObjectInput{
				Bucket: aws.String(core.Settings.AWS_FILES_BUCKET),
				Key:    aws.String(filename),
			},
		)
		if err != nil {
			panic(err)
		}
	} else {
		fmt.Printf("File <[%s]> is in use\n", filename)
	}
}

func checkFileUsageInDB(searchMap map[string][]string, filename string) bool {
	jobData := make(chan bool)
	stopJobs := make(chan bool)
	jobCount := 0

	for tableName, searchFields := range searchMap {
		go searchTableForFile(jobData, stopJobs, tableName, searchFields, filename)
		jobCount++
	}

	fileInUse := false
	for fileFound := range jobData {
		jobCount--
		if fileFound {
			fileInUse = true
			break
		} else if jobCount == 0 {
			fileInUse = false
			break
		}
	}

	return fileInUse
}

func searchTableForFile(out chan bool, quit chan bool, tableName string, searchFields []string, filename string) {
	var results []map[string]interface{}
	core.DbSession.Table(tableName).Find(&results)

	for _, item := range results {
		select {
		case <-quit:
			return
		default:
			for _, field := range searchFields {
				if filename == item[field] {
					out <- true
					return
				}
			}
		}
	}
	out <- false
}
