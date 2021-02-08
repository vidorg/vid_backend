package jwt

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestGenerateToken(t *testing.T) {
	userID := 1
	token, err := GenerateToken(userID, 3*time.Hour)
	if err != nil {
		assert.Error(t, err)
	}
	newUserClaim, err := ParseToken(token)
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, newUserClaim.UID, userID)
}

func TestGenerateTokenWithoutExpire(t *testing.T) {
	userID := 1
	token, err := GenerateTokenWithoutExpire(userID)
	if err != nil {
		assert.Error(t, err)
	}
	newUserClaim, err := ParseToken(token)
	if err != nil {
		assert.Error(t, err)
	}
	assert.Equal(t, newUserClaim.UID, userID)
}
