package picture_upload_service

import (
    "context"
    "fmt"
    "mime/multipart"

    "zetian-personal-website-hertz/biz/pkg/s3uploader"
)

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
