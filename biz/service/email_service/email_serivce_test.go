package email_service

import (
	"testing"
	"zetian-personal-website-hertz/biz/pkg/SES_email"

	"github.com/stretchr/testify/assert"
)

func init() {
	SES_email.InitSES()
}
func TestSendVeriCodeEmailTo1(t *testing.T) {
	err := SendVeriCodeEmailTo(t.Context(), "skyjin@bu.edu", "114514", "idk")
	assert.Nil(t, err)
}

func TestSendVeriCodeEmailTo2(t *testing.T) {
	err := SendVeriCodeEmailTo(t.Context(), "skyjin0127@gmail.com", "1433223", "idk")
	assert.Nil(t, err)
}