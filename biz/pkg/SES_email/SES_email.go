package SES_email

import (
	"context"
	"fmt"
	"log"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/sesv2"
	"github.com/aws/aws-sdk-go-v2/service/sesv2/types"
)

const (
	Sender    = "skylar27-no-reply@skylar27.com"
	AWSRegion = "us-east-2" // 你的 SES 域在 Ohio
)

var client *sesv2.Client

func InitSES() {
	cfg, err := config.LoadDefaultConfig(context.Background(), config.WithRegion(AWSRegion))
	if err != nil {
		log.Fatalf("failed to load AWS config: %v", err)
	}
	client = sesv2.NewFromConfig(cfg)
}

//send email to a custom email through aws ses
func SendEmail(ctx context.Context, to, subject, body string) error {


	// 组装请求 —— 注意：类型来自 types 包
	input := &sesv2.SendEmailInput{
		FromEmailAddress: aws.String(Sender),
		Destination: &types.Destination{
			ToAddresses: []string{to},
		},
		Content: &types.EmailContent{
			Simple: &types.Message{
				Subject: &types.Content{Data: aws.String(subject)},
				Body: &types.Body{
					Text: &types.Content{Data: aws.String(body)},
				},
			},
		},
	}

	// 发送
	out, err := client.SendEmail(ctx, input)
	if err != nil {
		return fmt.Errorf("failed to send email: %w", err)
	}
	log.Printf("✅ Email sent! Message ID: %s\n", aws.ToString(out.MessageId))
	return nil
}
