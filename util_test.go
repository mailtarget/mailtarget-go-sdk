package layang

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsEmail_Success(t *testing.T) {
	res := IsEmail("email@example.com")
	assert.Equal(t, true, res)
}

func TestIsEmail_Error(t *testing.T) {
	res := IsEmail("not_email")
	assert.Equal(t, false, res)
}

func TestIsHTML_Success(t *testing.T) {
	res := IsHTML("<!DOCTYPE html><nothtml>&#<body><a href=\"unclosed tag\"</html>")
	assert.Equal(t, true, res)
}

func TestIsHTML_Error(t *testing.T) {
	res := IsHTML("not_html")
	assert.Equal(t, false, res)
}
