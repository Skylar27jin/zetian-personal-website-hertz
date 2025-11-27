package picture_upload_service

import (
	"context"
	"encoding/json"
	"fmt"
	"mime/multipart"
	"strings"

	"zetian-personal-website-hertz/biz/pkg/s3uploader"
)


func UploadAvatar(ctx context.Context, userID int64, fileHeader *multipart.FileHeader) (string, error) {
    f, err := fileHeader.Open()
    if err != nil {
        return "", fmt.Errorf("open file %s failed: %w", fileHeader.Filename, err)
    }

    url, err := s3uploader.S3uploader.UploadAvatar(ctx, userID, fileHeader.Filename, f)
    f.Close()

    if err != nil {
        return "", fmt.Errorf("upload file %s failed: %w", fileHeader.Filename, err)
    }
    return url, nil
}

// UploadPostImages uploads one or more images for a post and returns S3 URLs.
func UploadPostImages(ctx context.Context, userID int64, files []*multipart.FileHeader) ([]string, error) {
    urls := make([]string, 0, len(files))

    for _, fh := range files {
        // TODO: 可以在这里加大小类型校验，例如 fh.Size、fh.Header.Get("Content-Type")

        f, err := fh.Open()
        if err != nil {
            return nil, fmt.Errorf("open file %s failed: %w", fh.Filename, err)
        }

        // 注意：这里不用 defer，直接用完就关，避免很多文件时把 fd 堆在一起
        url, err := s3uploader.S3uploader.UploadPostMedia(ctx, userID, fh.Filename, f)
        f.Close()

        if err != nil {
            return nil, fmt.Errorf("upload file %s failed: %w", fh.Filename, err)
        }

        urls = append(urls, url)
    }

    return urls, nil
}


func DeletePostImagesByURLs(ctx context.Context, urls []string) {
    for _, u := range urls {
        if u == "" {
            continue
        }
        if err := s3uploader.S3uploader.DeleteByURL(ctx, u); err != nil {
            fmt.Printf("⚠ delete s3 media failed: %s, err=%v\n", u, err)
        }
    }
}


// DeletePostImagesJSON parses JSON string from DB and deletes all S3 objects.
// mediaJSON 形如：["https://xxx","https://yyy"]
func DeletePostImagesJSON(ctx context.Context, mediaJSON string) {
    mediaJSON = strings.TrimSpace(mediaJSON)
    if mediaJSON == "" || mediaJSON == "[]" || mediaJSON == "null" {
        return
    }

    var urls []string
    if err := json.Unmarshal([]byte(mediaJSON), &urls); err != nil {
        fmt.Printf("⚠ parse media urls json failed, json=%s, err=%v\n", mediaJSON, err)
        return
    }

    DeletePostImagesByURLs(ctx, urls)
}