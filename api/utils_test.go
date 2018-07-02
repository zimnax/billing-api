package api

import (
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestIsCorrectMoneyFormat(t *testing.T) {
	incomingMoneyFormat := "8.88"

	is, err := isCorrectMoneyFormat(incomingMoneyFormat)

	assert.NoError(t, err)
	assert.True(t, is)
}

func TestIsCorrectMoneyFormat_Comma(t *testing.T) {
	incomingMoneyFormat := "8,88"

	is, err := isCorrectMoneyFormat(incomingMoneyFormat)

	assert.NoError(t, err)
	assert.False(t, is)
}

func TestIsCorrectMoneyFormat_WrongNumberFormat(t *testing.T) {
	incomingMoneyFormat := "8.8"

	is, err := isCorrectMoneyFormat(incomingMoneyFormat)

	assert.NoError(t, err)
	assert.False(t, is)
}

func TestIsCorrectMoneyFormat_WrongNumberFormat2(t *testing.T) {
	incomingMoneyFormat := "8"

	is, err := isCorrectMoneyFormat(incomingMoneyFormat)

	assert.NoError(t, err)
	assert.False(t, is)
}
