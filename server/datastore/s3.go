package datastore

import (
	"context"
	"fmt"
	"io"
	"path/filepath"
	"strings"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/credentials"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/s3"
	"github.com/karlseguin/ccache/v2"
)

type S3 struct {
	bucket string
	cache  *ccache.Cache
	ctx    context.Context
	conn   *s3.S3
}

type S3ConfigOptions struct {
	Bucket    string
	Endpoint  string
	Region    string
	AccessKey string
	SecretKey string
}

func (opts S3ConfigOptions) init() (result DS, err error) {
	conf := aws.NewConfig()
	conf.WithRegion(opts.Region)
	conf.WithEndpoint(opts.Endpoint)
	conf.WithCredentials(credentials.NewStaticCredentials(opts.AccessKey, opts.SecretKey, ""))

	sess := session.Must(session.NewSession(conf))

	conn := s3.New(sess)

	ctx := context.Background()

	cache := ccache.New(ccache.Configure().MaxSize(1000))

	result = &S3{
		bucket: opts.Bucket,
		cache:  cache,
		conn:   conn,
		ctx:    ctx,
	}

	return result, nil
}

func (s *S3) Put(key string, value string) error {
	ctx, cancelFn := context.WithTimeout(s.ctx, time.Second*5)
	defer cancelFn()

	_, err := s.conn.PutObjectWithContext(ctx, &s3.PutObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filepath.Join("artifacts", key)),
		Body:   strings.NewReader(value),
	})
	return err
}

func (s *S3) Get(key string) (string, error) {
	item, err := s.cache.Fetch(key, time.Hour*24*365, func() (interface{}, error) {
		ctx, cancelFn := context.WithTimeout(s.ctx, time.Second*5)
		defer cancelFn()

		obj, err := s.conn.GetObjectWithContext(ctx, &s3.GetObjectInput{
			Bucket: aws.String(s.bucket),
			Key:    aws.String(filepath.Join("artifacts", key)),
		})
		if err != nil {
			return "", err
		}

		buf := new(strings.Builder)
		_, err = io.Copy(buf, obj.Body)
		return buf, err
	})

	return fmt.Sprintf("%v", item.Value()), err
}

func (s *S3) Delete(key string) error {
	ctx, cancelFn := context.WithTimeout(s.ctx, time.Second*5)
	defer cancelFn()

	_, err := s.conn.DeleteObjectWithContext(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.bucket),
		Key:    aws.String(filepath.Join("artifacts", key)),
	})

	s.cache.Delete(key)
	return err
}
