package storage

import (
	"bytes"
	"context"
	"errors"
	"fmt"
	"io"
	"mime"
	"net/url"
	"path"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/credentials"
	"github.com/aws/aws-sdk-go-v2/service/s3"
	"github.com/aws/aws-sdk-go-v2/service/s3/types"
	"github.com/mauriciorobertodev/whappy-go/internal/infra/config"
)

type S3Storage struct {
	cfg    *config.StorageConfig
	client *s3.Client
}

func NewS3Storage(cfg *config.StorageConfig) *S3Storage {
	awsCfg := aws.Config{
		Credentials: aws.NewCredentialsCache(credentials.NewStaticCredentialsProvider(cfg.Key, cfg.Secret, "")),
		Region:      cfg.Region,
	}

	if cfg.Endpoint != "" {
		awsCfg.EndpointResolverWithOptions = aws.EndpointResolverWithOptionsFunc(
			func(service, region string, _ ...interface{}) (aws.Endpoint, error) {
				return aws.Endpoint{
					URL:               cfg.Endpoint,
					HostnameImmutable: true,
				}, nil
			},
		)
	}

	client := s3.NewFromConfig(awsCfg, func(o *s3.Options) {
		o.UsePathStyle = cfg.PathStyle
	})

	return &S3Storage{
		cfg:    cfg,
		client: client,
	}
}

func (s *S3Storage) Save(ctx context.Context, key string, r io.Reader) error {
	mimeType := mime.TypeByExtension(path.Ext(key))
	_, err := s.client.PutObject(ctx, &s3.PutObjectInput{
		Bucket:      aws.String(s.cfg.Bucket),
		Key:         aws.String(key),
		Body:        r,
		ContentType: &mimeType,
	})

	return err
}

func (s *S3Storage) Load(ctx context.Context, key string) (io.ReadCloser, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	return out.Body, nil
}

func (s *S3Storage) Get(ctx context.Context, key string) ([]byte, error) {
	out, err := s.client.GetObject(ctx, &s3.GetObjectInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		return nil, err
	}

	defer out.Body.Close()

	buf := new(bytes.Buffer)
	if _, err := io.Copy(buf, out.Body); err != nil {
		return nil, err
	}

	return buf.Bytes(), nil
}

func (s *S3Storage) Put(ctx context.Context, key string, data []byte) error {
	return s.Save(ctx, key, bytes.NewReader(data))
}

func (s *S3Storage) Delete(ctx context.Context, key string) error {
	_, err := s.client.DeleteObject(ctx, &s3.DeleteObjectInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(key),
	})

	return err
}

func (s *S3Storage) URL(ctx context.Context, key string) (string, error) {
	if s.cfg.URL != "" {
		u, err := url.JoinPath(s.cfg.URL, s.cfg.Bucket, key)
		if err != nil {
			return "", err
		}
		return u, nil
	}

	return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", s.cfg.Bucket, s.cfg.Region, key), nil
}

func (s *S3Storage) Exists(ctx context.Context, key string) (bool, error) {
	_, err := s.client.HeadObject(ctx, &s3.HeadObjectInput{
		Bucket: aws.String(s.cfg.Bucket),
		Key:    aws.String(key),
	})

	if err != nil {
		var nfe *types.NotFound
		if ok := errors.As(err, &nfe); ok {
			return false, nil
		}
		return false, err
	}

	return true, nil
}

func (s *S3Storage) EnsureBucket(ctx context.Context) error {
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.cfg.Bucket),
	})

	if err == nil {
		return nil
	}

	var nfe *types.NotFound
	if ok := errors.As(err, &nfe); !ok {
		return fmt.Errorf("erro checando bucket: %w", err)
	}

	_, err = s.client.CreateBucket(ctx, &s3.CreateBucketInput{
		Bucket: aws.String(s.cfg.Bucket),
	})

	if err != nil {
		return fmt.Errorf("erro criando bucket: %w", err)
	}

	return nil
}

func (s *S3Storage) Healthy(ctx context.Context) error {
	_, err := s.client.HeadBucket(ctx, &s3.HeadBucketInput{
		Bucket: aws.String(s.cfg.Bucket),
	})
	if err != nil {
		return fmt.Errorf("health failed: %w", err)
	}
	return nil
}
