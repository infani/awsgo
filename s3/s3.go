package s3

import (
	"context"
	"io"
	"log"
	"strings"
	"time"

	"github.com/infani/awsgo/config/awsConfig"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/s3/manager"
	awsS3 "github.com/aws/aws-sdk-go-v2/service/s3"
)

func CopyFile(bucket string, key string, newKey string) error {

	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		return err
	}
	client := awsS3.NewFromConfig(cfg)

	input := &awsS3.CopyObjectInput{
		Bucket:     aws.String(bucket),
		CopySource: aws.String(bucket + "/" + key),
		Key:        aws.String(newKey),
	}

	_, err = client.CopyObject(context.TODO(), input)

	if err != nil {
		return err
	}

	return nil
}

func GetFile(bucket string, key string) ([]byte, error) {
	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		log.Println(err)
		panic("configuration error, " + err.Error())
	}

	client := awsS3.NewFromConfig(cfg)
	downloader := manager.NewDownloader(client)

	writeAtBuf := manager.NewWriteAtBuffer([]byte{})
	_, err = downloader.Download(context.TODO(), writeAtBuf, &awsS3.GetObjectInput{Bucket: &bucket, Key: &key})

	if err != nil {
		return nil, err
	}
	// log.Println((string)(writeAtBuf.Bytes()))

	return writeAtBuf.Bytes(), nil
}

func PutFile(bucket string, key string, body io.Reader) error {
	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		log.Println(err)
		panic("configuration error, " + err.Error())
	}
	client := awsS3.NewFromConfig(cfg)
	uploader := manager.NewUploader(client)
	contentType := ""
	if strings.Contains(key, ".m3u8") {
		contentType = "application/x-mpegURL"
	}
	_, err = uploader.Upload(context.TODO(), &awsS3.PutObjectInput{Bucket: &bucket, Key: &key, Body: body, ContentType: &contentType})

	if err != nil {
		return err
	}

	return nil
}

func DeleteFolder(bucket string, folder string) error {
	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		return err
	}

	client := awsS3.NewFromConfig(cfg)

	input := &awsS3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(folder),
	}

	ctx := context.Background()
	resp, err := client.ListObjectsV2(ctx, input)
	if err != nil {
		return err
	}

	for _, item := range resp.Contents {
		input := &awsS3.DeleteObjectInput{
			Bucket: aws.String(bucket),
			Key:    aws.String(*item.Key),
		}
		_, err = client.DeleteObject(ctx, input)
		if err != nil {
			return err
		}
	}

	return nil
}

func SignedURL(bucket string, key string, expires time.Duration) (string, error) {
	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		return "", err
	}

	client := awsS3.NewFromConfig(cfg)
	psClient := awsS3.NewPresignClient(client)

	input := &awsS3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	resp, err := psClient.PresignGetObject(context.TODO(), input, func(options *awsS3.PresignOptions) { options.Expires = expires })

	if err != nil {
		return "", err
	}
	return resp.URL, nil
}

func GetSize(bucket string, folder string, filter string) (int64, error) {
	cfg, err := awsConfig.LoadAWSDefaultConfig()
	if err != nil {
		return 0, err
	}

	client := awsS3.NewFromConfig(cfg)

	input := &awsS3.ListObjectsV2Input{
		Bucket: aws.String(bucket),
		Prefix: aws.String(folder),
	}

	ctx := context.Background()
	resp, err := client.ListObjectsV2(ctx, input)
	if err != nil {
		return 0, err
	}
	var size int64
	for _, item := range resp.Contents {
		if strings.Contains(*item.Key, filter) {
			size += item.Size
		}
	}
	return size, nil
}	