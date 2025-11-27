package s3uploader

import (
    "context"
    "fmt"
    "io"
    "net/url"
    "path/filepath"
    "strings"
    "time"

    "github.com/aws/aws-sdk-go-v2/aws"
    aws_config "github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/s3"

    "zetian-personal-website-hertz/biz/config"
)

type Uploader interface {
    // 上传帖子图片，返回完整 URL（S3 或 CDN）
    UploadPostMedia(ctx context.Context, userID int64, filename string, r io.Reader) (string, error)

    // 根据完整 URL 删除 S3 对象
    DeleteByURL(ctx context.Context, rawURL string) error
}

type uploader struct {
    client    *s3.Client
    bucket    string
    cdnDomain string
    region    string
}

var S3uploader Uploader

// 在你的 main/init 里调用：s3uploader.InitS3Uploader()
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
    cfg, err := aws_config.LoadDefaultConfig(
        context.Background(),
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

    // 有 CDN 域名就优先走 CDN
    if u.cdnDomain != "" {
        return fmt.Sprintf("https://%s/%s", u.cdnDomain, key), nil
    }
    // 默认走 S3 域名
    return fmt.Sprintf("https://%s.s3.%s.amazonaws.com/%s", u.bucket, u.region, key), nil
}

// 从 URL 里提取 key，比如：
// https://project-talk-media.s3.us-east-1.amazonaws.com/post_media/6/xxx.png
// https://cdn.skylar27.com/post_media/6/xxx.png
func (u *uploader) extractKeyFromURL(rawURL string) (string, error) {
    parsed, err := url.Parse(rawURL)
    if err != nil {
        return "", fmt.Errorf("parse url: %w", err)
    }

    // path 形如 /post_media/6/xxx.png
    path := strings.TrimPrefix(parsed.Path, "/")
    if path == "" {
        return "", fmt.Errorf("empty path in url: %s", rawURL)
    }
    return path, nil
}

// DeleteByURL 根据完整 URL 删除 S3 对象
func (u *uploader) DeleteByURL(ctx context.Context, rawURL string) error {
    key, err := u.extractKeyFromURL(rawURL)
    if err != nil {
        return err
    }

    _, err = u.client.DeleteObject(ctx, &s3.DeleteObjectInput{
        Bucket: aws.String(u.bucket),
        Key:    aws.String(key),
    })
    if err != nil {
        return fmt.Errorf("delete object %s failed: %w", key, err)
    }
    return nil
}
