package s3util

import (
	"context"
	"fmt"
	"log"
	"net/url"
	"time"

	"github.com/google/uuid"
	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/minio/minio-go/v7/pkg/tags"
)

func CreateMinioClient(endpoint string, accessKey string, secretKey string, useSSL bool) (*minio.Client, error) {
	if endpoint == "" || accessKey == "" || secretKey == "" {
		return nil, fmt.Errorf("MINIO_ENDPOINT, MINIO_USERNAME, and MINIO_PASSWORD must be set")
	}
	client, err := minio.New(endpoint, &minio.Options{
		Creds:  credentials.NewStaticV4(accessKey, secretKey, ""),
		Secure: useSSL,
	})
	if err != nil {
		return nil, err
	}
	return client, nil
}

func GeneratePresignedPutURL(client *minio.Client, bucketName string, expires time.Duration, baseURL string) (string, string, error) {
    uniqueID := uuid.New().String()

    presignedURL, err := client.PresignedPutObject(context.Background(), bucketName, uniqueID, expires)
    if err != nil {
        return "", "", fmt.Errorf("error generating presigned PUT URL: %w", err)
    }

    urlStr := presignedURL.String()
    log.Printf("Generated presigned PUT URL: %s for file ID: %s", urlStr, uniqueID)
    urlStr, err = replaceHostWithBaseURL(urlStr, baseURL)
    log.Printf("Replaced host in presigned URL: %s", urlStr)
    if err != nil {
        return "", "", fmt.Errorf("error replacing host with base URL: %w", err)
    }

    return uniqueID, urlStr, nil
}

// GetPermanentObjectURL returns a permanent URL for an object
func GetPermanentObjectURL(bucketName, objectKey string, baseURL string) string {
    if baseURL == "" {
        // Fallback to a default format if S3_BASE_URL is not set
        return fmt.Sprintf("/api/files/%s/%s", bucketName, objectKey)
    }
    
    // Ensure baseURL doesn't end with slash
    if baseURL[len(baseURL)-1] == '/' {
        baseURL = baseURL[:len(baseURL)-1]
    }
    
    return fmt.Sprintf("%s/%s/%s", baseURL, bucketName, objectKey)
}

// ChangeObjectStatusToPermanent changes the status tag of an object from temporary to permanent
func ChangeObjectStatusToPermanent(client *minio.Client, bucketName, objectName string) error {
    // Get current tags if any
    t, err := client.GetObjectTagging(context.Background(), bucketName, objectName, minio.GetObjectTaggingOptions{})
    if err != nil {
        return fmt.Errorf("error getting object tags: %w", err)
    }

    // Create or update the status tag
    tagsMap := t.ToMap()
    if tagsMap == nil {
        tagsMap = make(map[string]string)
    }
    tagsMap["status"] = "permanent"

    // Apply the updated tags
    newTags, err := tags.NewTags(tagsMap, false)
    if err != nil {
        return fmt.Errorf("error creating new tags: %w", err)
    }

    // Use PutObjectTagging with the required PutObjectTaggingOptions parameter
    err = client.PutObjectTagging(context.Background(), bucketName, objectName, newTags, minio.PutObjectTaggingOptions{})
    if err != nil {
        return fmt.Errorf("error setting object tags: %w", err)
    }

    return nil
}

// IsObjectTemporary checks if an object has the temporary status tag
func IsObjectTemporary(client *minio.Client, bucketName, objectName string) (bool, error) {
    t, err := client.GetObjectTagging(context.Background(), bucketName, objectName, minio.GetObjectTaggingOptions{})
    if err != nil {
        return false, fmt.Errorf("error getting object tags: %w", err)
    }

    tagsMap := t.ToMap()
    return tagsMap["status"] == "temporary", nil
}

// replaceHostWithBaseURL replaces the host part of a URL with the S3_BASE_URL if set
func replaceHostWithBaseURL(originalURL string, baseURL string) (string, error) {
    if baseURL == "" {
        return originalURL, nil // Return original URL if S3_BASE_URL is not set
    }

    u, err := url.Parse(originalURL)
    if err != nil {
        return "", fmt.Errorf("error parsing URL: %w", err)
    }

    baseU, err := url.Parse(baseURL)
    if err != nil {
        return "", fmt.Errorf("error parsing base URL: %w", err)
    }

    // Replace the scheme, host, and port with those from the base URL
    u.Scheme = baseU.Scheme
    u.Host = baseU.Host

    // Preserve the path from the base URL if it exists and append the original path
    if baseU.Path != "" && baseU.Path != "/" {
        // Ensure base path doesn't end with slash and original path doesn't start with slash
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