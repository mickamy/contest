package jwt_test

import (
	"testing"

	"github.com/google/uuid"
	"github.com/stretchr/testify/assert"

	"github.com/mickamy/contest/example/internal/lib/jwt"
)

func TestJWT_New(t *testing.T) {
	t.Parallel()

	// arrange
	id := uuid.NewString()

	// act
	token, err := jwt.New(id)

	// assert
	assert.NoError(t, err)
	assert.NotEmpty(t, token.Access.Value)
	assert.NotEmpty(t, token.Access.ExpiresAt)
	assert.NotEmpty(t, token.Refresh.Value)
	assert.NotEmpty(t, token.Refresh.ExpiresAt)
}

func TestJWT_Verify(t *testing.T) {
	t.Parallel()

	// arrange
	id := uuid.NewString()
	token, err := jwt.New(id)
	assert.NoError(t, err)

	// act
	accessClaim, err := jwt.Verify(token.Access.Value)
	assert.NoError(t, err)
	refreshClaim, err := jwt.Verify(token.Refresh.Value)
	assert.NoError(t, err)

	// assert
	assert.Equal(t, accessClaim["id"], id)
	assert.NotEmpty(t, accessClaim["exp"])
	assert.Equal(t, refreshClaim["id"], id)
	assert.NotEmpty(t, refreshClaim["exp"], token.Access.ExpiresAt)
}

func TestJWT_ExtractID(t *testing.T) {
	t.Parallel()

	// arrange
	userID := uuid.NewString()
	token, err := jwt.New(userID)
	assert.NoError(t, err)

	tcs := []struct {
		name     string
		token    jwt.Token
		expected string
	}{
		{
			name:     "from access token",
			token:    token.Access,
			expected: userID,
		},
		{
			name:     "from refresh token",
			token:    token.Refresh,
			expected: userID,
		},
	}
	for _, c := range tcs {
		t.Run(c.name, func(t *testing.T) {
			c := c
			t.Parallel()

			// act
			actual, err := jwt.ExtractID(c.token.Value)

			// assert
			assert.NoError(t, err)
			assert.Equal(t, c.expected, actual)
		})
	}
}
