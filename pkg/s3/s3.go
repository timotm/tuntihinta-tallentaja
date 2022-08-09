package s3

import (
	"bytes"
	"context"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

func PutFileToS3(contents []byte, region string, bucketName string, key string) {
	cfg, err := config.LoadDefaultConfig(context.TODO(), config.WithClientLogMode(aws.LogDeprecatedUsage))
	if err != nil {
		log.Fatalf("Error creating config: %s", err)
	}
	client := s3.NewFromConfig(cfg)
	_, err = client.PutObject(context.TODO(), &s3.PutObjectInput{
		Bucket: aws.String(bucketName),
		Key:    aws.String(key),
		Body:   bytes.NewReader(contents)})

	if err != nil {
		log.Fatalf("Error uploading file %s to bucket %s in %s: %s\n", key, bucketName, region, err)
	}
}
