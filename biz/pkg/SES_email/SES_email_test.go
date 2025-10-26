package SES_email

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestSendEmail(t *testing.T) {
	InitSES()
	err := SendEmail(context.Background(), "skyjin@bu.edu", "Hello From website!", "hello world, this is the body")
	assert.Nil(t, err)

}