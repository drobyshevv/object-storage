package storage

import (
	"context"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/service/s3"
)

// EnsureBucket проверяет, есть ли bucket, и создаёт его, если нет
func EnsureBucket(ctx context.Context, client *s3.Client, bucket string) error {
	// Проверяем наличие
	_, err := client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(bucket),
	})

	if err == nil {
		// Bucket существует
		return nil
	}

	// Создаём bucket
	_, err = client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(bucket),
	})

	return err
}
