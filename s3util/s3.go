package s3util

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/tags"
)

type S3Config struct {
	Endpoint  string
	AccessKey string
	SecretKey string
	UseSSL    bool
	BaseURL   string
}

type Client struct {
	minio   *minio.Client
	baseURL string
}

func New(cfg S3Config) (*Client, error) {
	if cfg.Endpoint == "" || cfg.AccessKey == "" || cfg.SecretKey == "" {
		return nil, fmt.Errorf("endpoint, accessKey and secretKey must be set")
	}
	client, err := minio.New(cfg.Endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(cfg.AccessKey, cfg.SecretKey, ""),
		Secure: cfg.UseSSL,
	})
	if err != nil {
		return nil, err
	}
	return &Client{minio: client, baseURL: cfg.BaseURL}, nil
}

func (c *Client) GeneratePresignedPutURL(ctx context.Context, bucket string, expires time.Duration) (string, string, error) {
	uniqueID := uuid.New().String()
	presignedURL, err := c.minio.PresignedPutObject(ctx, bucket, uniqueID, expires)
	if err != nil {
		return "", "", fmt.Errorf("generate presigned PUT url: %w", err)
	}
	urlStr, err := replaceHostWithBaseURL(presignedURL.String(), c.baseURL)
	if err != nil {
		return "", "", err
	}
	return uniqueID, urlStr, nil
}

func (c *Client) GetPermanentObjectURL(bucket, object string) string {
	if c.baseURL == "" {
		return fmt.Sprintf("/api/files/%s/%s", bucket, object)
	}
	if c.baseURL[len(c.baseURL)-1] == '/' {
		return fmt.Sprintf("%s%s/%s", c.baseURL[:len(c.baseURL)-1], bucket, object)
	}
	return fmt.Sprintf("%s/%s/%s", c.baseURL, bucket, object)
}

func (c *Client) MarkObjectPermanent(ctx context.Context, bucket, object string) error {
	t, err := c.minio.GetObjectTagging(ctx, bucket, object, minio.GetObjectTaggingOptions{})
	if err != nil {
		return err
	}
	tagsMap := t.ToMap()
	if tagsMap == nil {
		tagsMap = make(map[string]string)
	}
	tagsMap["status"] = "permanent"

	newTags, err := tags.NewTags(tagsMap, false)
	if err != nil {
		return err
	}
	return c.minio.PutObjectTagging(ctx, bucket, object, newTags, minio.PutObjectTaggingOptions{})
}

func (c *Client) IsObjectTemporary(ctx context.Context, bucket, object string) (bool, error) {
	t, err := c.minio.GetObjectTagging(ctx, bucket, object, minio.GetObjectTaggingOptions{})
	if err != nil {
		return false, err
	}
	return t.ToMap()["status"] == "temporary", nil
}

func replaceHostWithBaseURL(originalURL string, baseURL string) (string, error) {
	if baseURL == "" {
		return originalURL, nil
	}
	u, err := url.Parse(originalURL)
	if err != nil {
		return "", err
	}
	baseU, err := url.Parse(baseURL)
	if err != nil {
		return "", err
	}
	u.Scheme, u.Host = baseU.Scheme, baseU.Host
	if baseU.Path != "" && baseU.Path != "/" {
		basePath := baseU.Path
		if basePath[len(basePath)-1] == '/' {
			basePath = basePath[:len(basePath)-1]
		}
		origPath := u.Path
		if len(origPath) > 0 && origPath[0] == '/' {
			origPath = origPath[1:]
		}
		u.Path = basePath + "/" + origPath
	}
	return u.String(), nil
}
