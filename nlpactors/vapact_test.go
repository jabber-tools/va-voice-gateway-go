package nlpactors

import (
	"github.com/stretchr/testify/assert"
	"testing"
	"time"
)

func TestIsExpired1(t *testing.T) {
	TokenCreatedTime := time.Now()
	time.Sleep(time.Second * 5)
	TimeNow := time.Now()
	assert.Equal(t, true, IsExpired(TimeNow, TokenCreatedTime, 3), "Should be expired")
}

func TestIsExpired2(t *testing.T) {
	TimeNow := time.Now()
	time.Sleep(time.Second * 5)
	TokenCreatedTime := time.Now()
	assert.Equal(t, true, IsExpired(TimeNow, TokenCreatedTime, 3), "Should be expired")
}

func TestIsExpired3(t *testing.T) {
	TokenCreatedTime := time.Now()
	time.Sleep(time.Second * 3)
	TimeNow := time.Now()
	assert.Equal(t, false, IsExpired(TimeNow, TokenCreatedTime, 10), "Should not be expired")
}
