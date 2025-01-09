package services

import (
    "mime/multipart"
    "context"
    "fmt"
    "strings"

    "github.com/aws/aws-sdk-go-v2/aws"
    "github.com/aws/aws-sdk-go-v2/service/s3"
    "github.com/aws/aws-sdk-go-v2/service/s3/types"

    log "github.com/sirupsen/logrus"
)

// BucketManager handles operations on S3 buckets
type BucketManager struct {
    client     *s3.Client
    bucketName string
}

// NewBucketManager creates a new instance of BucketManager
func NewBucketManager(client *s3.Client, bucketName string) *BucketManager {
    return &BucketManager{
        client:     client,
        bucketName: bucketName,
    }
}

// PutObject uploads an object to the specified path in the bucket
func (bm *BucketManager) PutObject(ctx context.Context, objectPath string, file multipart.File) error {
    _, err := bm.client.PutObject(ctx, &s3.PutObjectInput{
        Bucket: aws.String(bm.bucketName),
        Key:    aws.String(objectPath),
        Body:   file,
    })
    
    if err != nil {
        log.Error("failed to upload object: ", err)
        return err
    }

    return nil
}

// ClearFolder deletes all objects within the specified folder path
func (bm *BucketManager) ClearFolder(ctx context.Context, folderPath string) error {
    // Ensure folder path ends with "/"
    if !strings.HasSuffix(folderPath, "/") {
        folderPath = folderPath + "/"
    }

    // List all objects in the folder
    input := &s3.ListObjectsV2Input{
        Bucket: aws.String(bm.bucketName),
        Prefix: aws.String(folderPath),
    }

    paginator := s3.NewListObjectsV2Paginator(bm.client, input)

    for paginator.HasMorePages() {
        page, err := paginator.NextPage(ctx)
        if err != nil {
            return fmt.Errorf("failed to list objects: %v", err)
        }

        // If there are no objects, we're done
        if len(page.Contents) == 0 {
            return nil
        }

        // Create delete objects input
        objectIdentifiers := make([]types.ObjectIdentifier, len(page.Contents))
        for i, object := range page.Contents {
            objectIdentifiers[i] = types.ObjectIdentifier{
                Key: object.Key,
            }
        }

        // Delete the objects
        _, err = bm.client.DeleteObjects(ctx, &s3.DeleteObjectsInput{
            Bucket: aws.String(bm.bucketName),
            Delete: &types.Delete{
                Objects: objectIdentifiers,
                Quiet:   aws.Bool(true),
            },
        })

        if err != nil {
            return fmt.Errorf("failed to delete objects: %v", err)
        }
    }

    return nil
}

// Helper method to check if bucket exists
func (bm *BucketManager) BucketExists(ctx context.Context) (bool, error) {
    _, err := bm.client.HeadBucket(ctx, &s3.HeadBucketInput{
        Bucket: aws.String(bm.bucketName),
    })
    
    if err != nil {
        if strings.Contains(err.Error(), "NotFound") || strings.Contains(err.Error(), "NoSuchBucket") || strings.Contains(err.Error(), "404") {
            return false, nil
        }
        return false, fmt.Errorf("failed to check bucket existence: %v", err)
    }
    
    return true, nil
}
