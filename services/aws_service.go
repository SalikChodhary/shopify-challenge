package services

import (
	"bytes"
	"errors"
	"mime/multipart"
	"net/http"
	"os"
	"sync"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
)

var instance *session.Session 
var err error

func InitS3Instance() (error) {
	var once sync.Once
	once.Do(initS3Session)
	if err != nil {
		return errors.New("Failed to get aws instance.")
	}
	return nil
}

func initS3Session()  {
	instance, err = session.NewSession(&aws.Config{
		Region: aws.String("us-east-2"),
		Credentials: credentials.NewStaticCredentials(
			os.Getenv("AWS_ID"), // id
			os.Getenv("AWS_SECRET"),   // secret
			""),  // token can be left blank for now
	})
}

func UploadFileToS3(file multipart.File, fileHeader *multipart.FileHeader, key string) error {
	size := fileHeader.Size
  buffer := make([]byte, size)
	file.Read(buffer)
	println(http.DetectContentType(buffer))
	
	_, err := s3.New(instance).PutObject(&s3.PutObjectInput{
		Bucket:               aws.String("salik-test-bucket"),
		Key:                  aws.String(key),
		ACL:                  aws.String("public-read"),// could be private if you want it to be access by only authorized users
		Body:                 bytes.NewReader(buffer),
		ContentLength:        aws.Int64(int64(size)),
		ContentType:        	aws.String(http.DetectContentType(buffer)),
		ServerSideEncryption: aws.String("AES256"),
		StorageClass:         aws.String("INTELLIGENT_TIERING"),
  })

	return err
}