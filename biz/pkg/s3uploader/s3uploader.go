package s3uploader

import (
    "context"
    "fmt"
    "io"
    "path/filepath"
    "time"

    "github.com/aws/aws-sdk-go-v2/aws"
    aws_config "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"

    "zetian-personal-website-hertz/biz/config"
)

type Uploader interface {
    UploadPostMedia(ctx context.Context, userID int64, filename string, r io.Reader) (string, error)
}

type uploader struct {
    client    *s3.Client
    bucket    string
    cdnDomain string
    region    string
}


var S3uploader Uploader


func InitS3Uploader() {
    var err error
    S3uploader, err = New(
        config.GetSpecificConfig().S3Bucket,
        config.GetSpecificConfig().AWSRegion,
        config.GetSpecificConfig().CDNDomain,
    )
    if err != nil {
        panic(fmt.Sprintf("failed to init s3 uploader: %v", err))
    }
}

func New(bucket, region, cdnDomain string) (Uploader, error) {
    cfg, err := aws_config.LoadDefaultConfig(context.Background(),
        aws_config.WithRegion(region),
    )
    if err != nil {
        return nil, fmt.Errorf("load aws config: %w", err)
    }

    return &uploader{
        client:    s3.NewFromConfig(cfg),
        bucket:    bucket,
        cdnDomain: cdnDomain,
        region:    region,
    }, nil
}

func (u *uploader) UploadPostMedia(
    ctx context.Context,
    userID int64,
    filename string,
    r io.Reader,
) (string, error) {
    ext := filepath.Ext(filename)
    if ext == "" {
        ext = ".jpg"
    }
    key := fmt.Sprintf("post_media/%d/%d%s", userID, time.Now().UnixNano(), ext)

    _, err := u.client.PutObject(ctx, &s3.PutObjectInput{
        Bucket: aws.String(u.bucket),
        Key:    aws.String(key),
        Body:   r,
    })
    if err != nil {
        return "", fmt.Errorf("put object: %w", err)
    }

    if u.cdnDomain != "" {
        return fmt.Sprintf("https://%s/%s", u.cdnDomain, key), nil
    }
    return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.bucket, u.region, key), nil
}
