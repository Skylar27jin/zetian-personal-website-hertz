package auth_service

import (
	"context"
	"log"
	"testing"
	"time"
	"zetian-personal-website-hertz/biz/config"

	"github.com/stretchr/testify/assert"
)

func init() {
	// 统一初始化配置，避免重复
	config.InitConfig()
}

/* ------------------ Test: Generate & Parse User JWT ------------------ */

func TestGenerateUserJWT(t *testing.T) {
	ctx := context.Background()
	username := "nwang"
	email := "skyjin@bu.edu"
	id := 1
	now := time.Now().Unix()
	validDuration := int64(7 * 24 * 3600)

	token, err := GenerateUserJWT(ctx, now, id, username, email, validDuration)
	assert.NoError(t, err, "should generate JWT without error")
	assert.NotEmpty(t, token, "token should not be empty")
}

func TestParseUserJWT(t *testing.T) {
	ctx := context.Background()
	username := "nwang"
	email := "skyjin@bu.edu"
	id := 1
	now := time.Now().Unix()
	validDuration := int64(10 * 60)
	expectedExp := now + validDuration

	// 先生成
	token, err := GenerateUserJWT(ctx, now, id, username, email, validDuration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)

	// 再解析
	usernameGot, emailGot, iatGot, expGot, idGot, err := ParseUserJWT(ctx, token)
	assert.NoError(t, err, "should parse JWT successfully")
	assert.Equal(t, username, usernameGot)
	assert.Equal(t, email, emailGot)
	assert.Equal(t, id, idGot)
	assert.GreaterOrEqual(t, iatGot, now-1, "iat should be close to now")
	assert.GreaterOrEqual(t, expGot, expectedExp-5, "exp should be close to expected")
}

/* ------------------ Test: Generate & Parse Verification Email JWT ------------------ */

func TestGenerateVeriEmailJWT(t *testing.T) {
	ctx := context.Background()
	email := "skyjin@bu.edu"
	now := time.Now().Unix()
	validDuration := int64(15 * 60)

	token, err := GenerateVeriEmailJWT(ctx, now, email, validDuration)
	assert.NoError(t, err, "should generate verification JWT without error")
	assert.NotEmpty(t, token)
}

func TestParseVeriEmailJWT(t *testing.T) {
	ctx := context.Background()
	email := "sample@sample.com"
	now := time.Now().Unix()
	validDuration := int64(10 * 60)
	expectedExp := now + validDuration

	// ✅ 动态生成一个不会过期的 token
	token, err := GenerateVeriEmailJWT(ctx, now, email, validDuration)
	assert.NoError(t, err)
	assert.NotEmpty(t, token)
	log.Println(token)
	// ✅ 解析 token
	emailGot, expGot, err := ParseVeriEmailJWT(ctx, token)
	assert.NoError(t, err, "should parse valid JWT correctly")
	assert.Equal(t, email, emailGot, "email should match")
	assert.GreaterOrEqual(t, expGot, expectedExp-5, "exp should be close to expected value")

}
