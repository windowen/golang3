package test

import (
	"context"
	"log"
	"testing"

	upload2 "serverApi/pkg/tools/upload"
)

func TestUploadS3(t *testing.T) {
	InitConfig()

	upload := upload2.NewUploader()

	// 上传文件
	fileData := []byte("your file content")
	err := upload.Upload(context.Background(), "file/test02.txt", fileData)
	if err != nil {
		log.Fatalf("Failed to upload file: %v", err)
	}

}
