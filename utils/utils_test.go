package utils

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestNormalizeAWB(t *testing.T) {
	assert.Equal(t, "it is 12342", NormalizeAWB("it is one two three four two"), "expected does not match actual")
	assert.Equal(t, "It is 12something 345", NormalizeAWB("It is one two something three four five"), "expected does not match actual")
	assert.Equal(t, "转移到代理 ", NormalizeAWB("转移到代理"), "expected does not match actual")
}

func TestRemoveNonAlphaNumericChars(t *testing.T) {
	assert.Equal(t, "It is 123456 ", RemoveNonAlphaNumericChars("It is---,./ 123456 #$%^&"), "expected does not match actual")
	assert.Equal(t, "Adam was here", RemoveNonAlphaNumericChars("Adam was here"), "expected does not match actual")
	assert.Equal(t, "转移到代理", RemoveNonAlphaNumericChars("转移到代理"), "expected does not match actual")
	assert.Equal(t, "转移到代理", RemoveNonAlphaNumericChars("转移---./,./,/.到代理"), "expected does not match actual")
}