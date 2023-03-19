package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
)

const bucketsFilename = "out_buckets.txt"
const bucketObjectsFilename = "out_bucket_objects.txt"
const bucketObjectsNamesFilename = "out_bucket_objects_names.txt"

func GetFakeBuckets() ([]minio.BucketInfo, error) {
	cfg, err := GetConfig()
	if err != nil {
		return []minio.BucketInfo{}, err
	}

	srcDir := cfg.Params.InputDir
	srcDir = strings.Replace(srcDir, "input", "src", 1)
	fileContent, err := ioutil.ReadFile(filepath.Join(srcDir, bucketsFilename))
	if err != nil {
		log.Fatal(err)
	}

	var buckets []minio.BucketInfo
	err = json.Unmarshal(fileContent, &buckets)
	if err != nil {
		return nil, err
	}

	for k, v := range buckets {
		fmt.Println(k, v)
	}

	return buckets, nil
}

func GetFakeBucketObjectsNames() ([]string, error) {
	cfg, err := GetConfig()
	if err != nil {
		return []string{}, err
	}

	srcDir := cfg.Params.InputDir
	srcDir = strings.Replace(srcDir, "input", "src", 1)
	fileContent, err := ioutil.ReadFile(filepath.Join(srcDir, bucketObjectsNamesFilename))
	if err != nil {
		log.Fatal(err)
	}

	var bucketObjectsNames []string
	err = json.Unmarshal(fileContent, &bucketObjectsNames)
	if err != nil {
		return nil, err
	}

	for k, v := range bucketObjectsNames {
		fmt.Println(k, v)
	}

	return bucketObjectsNames, nil
}

func GetFakeBucketObjects() ([]FileInfo, error) {
	cfg, err := GetConfig()
	if err != nil {
		return []FileInfo{}, err
	}

	srcDir := cfg.Params.InputDir
	srcDir = strings.Replace(srcDir, "input", "src", 1)
	fileContent, err := ioutil.ReadFile(filepath.Join(srcDir, bucketObjectsFilename))
	if err != nil {
		log.Fatal(err)
	}

	var bucketObjects []FileInfo
	err = json.Unmarshal(fileContent, &bucketObjects)
	if err != nil {
		return nil, err
	}

	for k, v := range bucketObjects {
		fmt.Println(k, v)
	}

	return bucketObjects, nil
}
