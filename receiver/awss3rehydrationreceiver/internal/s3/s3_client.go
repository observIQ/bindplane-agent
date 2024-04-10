// Copyright observIQ, Inc.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package s3 //import "github.com/observiq/bindplane-agent/receiver/awss3rehydrationreceiver/internal/s3"

import (
	"context"
	"fmt"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials/stscreds"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/aws/aws-sdk-go/service/s3/s3manager"
)

// ObjectInfo contains necessary info to process S3 objects
type ObjectInfo struct {
	Name string
	size int64
}

// S3Client provides a client for S3 object operations
//
//go:generate mockery --name S3Client --output ./mocks --with-expecter --filename mock_s3_client.go --structname MockS3Client
type S3Client interface {
	// ListObjects returns a list of ObjectInfo objects present in the bucket with the given prefix
	ListObjects(ctx context.Context, bucket string, prefix, continuationToken *string) ([]*ObjectInfo, *string, error)

	// DownloadObject downloads the contents of the object into the buffer.
	DownloadObject(ctx context.Context, bucket, key string, buf []byte) (int64, error)

	// DeleteObjects deletes the keys in the specified bucket
	DeleteObjects(ctx context.Context, bucket string, keys []string) error
}

// AWSClient is an implementation of the S3Client for AWS
type AWSClient struct {
	sessionCfg *aws.Config
	roleArn    string
}

// NewAWSClient creates a new AWS Client
func NewAWSClient(region, roleArn string) S3Client {
	sessionConfig := &aws.Config{
		Region: aws.String(region),
	}

	return &AWSClient{
		sessionCfg: sessionConfig,
		roleArn:    roleArn,
	}
}

func (a *AWSClient) ListObjects(ctx context.Context, bucket string, prefix, continuationToken *string) ([]*ObjectInfo, *string, error) {
	sess, err := getSession(a.sessionCfg, a.roleArn)
	if err != nil {
		return nil, nil, fmt.Errorf("get session: %w", err)
	}

	svc := s3.New(sess)

	input := &s3.ListObjectsV2Input{
		Bucket:            aws.String(bucket),
		ContinuationToken: continuationToken,
		Prefix:            prefix,
	}

	output, err := svc.ListObjectsV2WithContext(ctx, input)
	if err != nil {
		return nil, nil, fmt.Errorf("list objects: %w", err)
	}

	nextToken := output.ContinuationToken

	objects := make([]*ObjectInfo, len(output.Contents))

	for i, object := range output.Contents {
		// All fields are pointers so check all pointers needed to process are present
		if object.Key == nil || object.Size == nil {
			continue
		}

		objects[i] = &ObjectInfo{
			Name: *object.Key,
			size: *object.Size,
		}
	}

	return objects, nextToken, nil
}

// DownloadObject downloads the contents of the object.
func (a *AWSClient) DownloadObject(ctx context.Context, bucket, key string, buf []byte) (int64, error) {
	sess, err := getSession(a.sessionCfg, a.roleArn)
	if err != nil {
		return 0, fmt.Errorf("get session: %w", err)
	}

	downloader := s3manager.NewDownloader(sess)

	input := &s3.GetObjectInput{
		Bucket: aws.String(bucket),
		Key:    aws.String(key),
	}

	buffer := aws.NewWriteAtBuffer(buf)
	n, err := downloader.DownloadWithContext(ctx, buffer, input)
	if err != nil {
		return 0, fmt.Errorf("download: %w", err)
	}

	return n, err
}

// DeleteObjects deletes the keys in the specified bucket
func (a *AWSClient) DeleteObjects(ctx context.Context, bucket string, keys []string) error {
	sess, err := getSession(a.sessionCfg, a.roleArn)
	if err != nil {
		return fmt.Errorf("get session: %w", err)
	}

	deleter := s3manager.NewBatchDelete(sess, func(bd *s3manager.BatchDelete) { bd.BatchSize = len(keys) })

	objects := make([]s3manager.BatchDeleteObject, len(keys))
	for i, key := range keys {
		objects[i] = s3manager.BatchDeleteObject{
			Object: &s3.DeleteObjectInput{
				Bucket: aws.String(bucket),
				Key:    aws.String(key),
			},
		}
	}

	if err := deleter.Delete(ctx, &s3manager.DeleteObjectsIterator{
		Objects: objects,
	}); err != nil {
		return fmt.Errorf("delete: %w", err)
	}

	return nil
}

func getSession(sessionConfig *aws.Config, roleArn string) (*session.Session, error) {
	sess, err := session.NewSession(sessionConfig)
	if err != nil {
		return nil, fmt.Errorf("new session: %w", err)
	}

	if roleArn != "" {
		credentials := stscreds.NewCredentials(sess, roleArn)
		sess.Config.Credentials = credentials
	}

	return sess, nil
}
