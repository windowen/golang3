package upload

import (
	"bytes"
	"context"
	"fmt"
	"log"
	"mime/multipart"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/gin-gonic/gin"

	"serverApi/pkg/common/config"
	siteRes "serverApi/pkg/response/site"
	"serverApi/pkg/tools/apiresp"
	"serverApi/pkg/tools/errs"
	"serverApi/pkg/tools/utils"
	"serverApi/pkg/zlogger"
)

// 默认 ACL
const defaultAcl = "private"

// Uploader 结构体包含上传所需的所有配置
type Uploader struct {
	bucket   string // 存储桶
	acl      string // 权限
	metadata map[string]*string
	s3Client *s3.S3
}

// NewUploader 创建一个新的 Uploader，并应用 Option 配置
func NewUploader(opts ...Option) *Uploader {
	s3Cfg := config.Config.S3

	s3Config := &aws.Config{
		Credentials:      credentials.NewStaticCredentials(s3Cfg.AccessKeyID, s3Cfg.SecretAccessKey, ""),
		Region:           aws.String(s3Cfg.Region),
		Endpoint:         aws.String(s3Cfg.Endpoint),
		S3ForcePathStyle: aws.Bool(true),
	}

	sess, err := session.NewSession(s3Config)
	if err != nil {
		log.Fatalf("Failed to create AWS session: %v", err)
	}

	u := &Uploader{
		bucket:   s3Cfg.Bucket,
		acl:      defaultAcl,
		metadata: nil,
		s3Client: s3.New(sess),
	}

	for _, opt := range opts {
		opt(u)
	}

	return u
}

// CheckFile 检查文件的类型和大小
func (u *Uploader) CheckFile(file *multipart.FileHeader) ([]byte, error) {
	allowedTypes := []string{".jpg", ".png", ".gif", ".jpeg"}
	fileExt := strings.ToLower(filepath.Ext(file.Filename))

	if !utils.SliceHas(allowedTypes, fileExt) {
		return nil, errs.ErrUpload.WithDetail("upload_type_limit")
	}

	maxSize := int64(2 << 20) // 2MB
	if file.Size > maxSize {
		return nil, errs.ErrUpload.WithDetail("upload_exceed_limit")
	}

	src, err := file.Open()
	if err != nil {
		return nil, fmt.Errorf("CheckFile failed to open file | err=%w", err)
	}
	defer func(src multipart.File) {
		err := src.Close()
		if err != nil {
			zlogger.Errorf("CheckFile failed to close file: %v", err)
		}
	}(src)

	buf := bytes.NewBuffer(nil)
	if _, err := buf.ReadFrom(src); err != nil {
		return nil, fmt.Errorf("CheckFile failed to read from file: %w", err)
	}

	return buf.Bytes(), nil
}

// Upload 上传文件到 S3
func (u *Uploader) Upload(ctx context.Context, fileName string, fileData []byte) error {
	_, err := u.s3Client.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket:   aws.String(u.bucket),
		Key:      aws.String(fileName),
		Body:     bytes.NewReader(fileData),
		ACL:      aws.String(u.acl),
		Metadata: u.metadata,
	})
	if err != nil {
		return fmt.Errorf("failed to upload file to S3: %w", err)
	}
	return nil
}

// 生成唯一的文件名
func generateFileName(file *multipart.FileHeader) string {
	timestamp := time.Now().Unix()
	randomStr := utils.RandString(10)
	fileExt := filepath.Ext(file.Filename)
	return fmt.Sprintf("file/%v-%s%s", timestamp, randomStr, fileExt)
}

// GinHandler 上传文件的 Gin 处理器
func GinHandler(c *gin.Context) {
	file, err := c.FormFile("file")
	if err != nil {
		apiresp.GinError(c, errs.ErrUpload.WithDetail("upload_please_select_file"))
		return
	}

	// 初始化上传器
	upload := NewUploader()

	// 检查文件的类型和大小
	fileData, err := upload.CheckFile(file)
	if err != nil {
		apiresp.GinError(c, err)
		return
	}

	// 生成唯一文件名
	filename := generateFileName(file)

	// 上传文件到 S3
	if err := upload.Upload(c, filename, fileData); err != nil {
		apiresp.GinError(c, err)
		return
	}

	apiresp.GinSuccess(c, siteRes.UploadResp{Uri: filename})
}

// Option 是一个函数类型，用于配置 Uploader
type Option func(*Uploader)

// WithBucket 设置上传的 S3 存储桶
func WithBucket(bucket string) Option {
	return func(u *Uploader) {
		u.bucket = bucket
	}
}

// WithACL 设置文件的访问控制列表（ACL）
func WithACL(acl string) Option {
	return func(u *Uploader) {
		u.acl = acl
	}
}

// WithMetadata 设置文件的元数据
func WithMetadata(metadata map[string]*string) Option {
	return func(u *Uploader) {
		u.metadata = metadata
	}
}
