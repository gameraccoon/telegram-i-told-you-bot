package processing

import (
	"github.com/stretchr/testify/require"
	"testing"
)

func TestParseBetTime(t *testing.T) {
	assert := require.New(t)

	{
		time, isSuccessful, errorMessage := ParseBetTime("2d15h")
		assert.True(isSuccessful, errorMessage)
		assert.Equal(2 * 24 + 15, int(time.Hours()))
	}

	{
		time, isSuccessful, errorMessage := ParseBetTime("1y2m3d4h")
		assert.True(isSuccessful, errorMessage)
		assert.Equal(1 * 365 * 24 + 2 * 30 * 24 + 3 * 24 + 4, int(time.Hours()))
	}

	{
		time, isSuccessful, errorMessage := ParseBetTime("20m")
		assert.True(isSuccessful, errorMessage)
		assert.Equal(20 * 30 * 24, int(time.Hours()))
	}

	{
		_, isSuccessful, _ := ParseBetTime("5h2m")
		assert.False(isSuccessful)
	}

	{
		_, isSuccessful, _ := ParseBetTime("")
		assert.False(isSuccessful)
	}
}
