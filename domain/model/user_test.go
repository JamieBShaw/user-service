package model

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestUser_Validate(t *testing.T) {

	tt := []struct {
		name     string
		user     *User
		errorMsg string
		ok       bool
	}{
		{
			name: "valid user",
			user: &User{
				ID:       1,
				Username: "James",
				Admin:    true,
			},
			errorMsg: "nil",
			ok:       true,
		},
		{
			name: "invalid user id",
			user: &User{
				ID:       0,
				Username: "James",
				Admin:    true,
			},
			errorMsg: "user id is invalid",
			ok:       false,
		},
		{
			name: "invalid user name, empty username",
			user: &User{
				ID:       1,
				Username: "",
				Admin:    true,
			},
			errorMsg: "username is empty",
			ok:       false,
		},
		{
			name: "invalid user name, username too long",
			user: &User{
				ID:       1,
				Username: "JamesNameOver10",
				Admin:    true,
			},
			errorMsg: "username is to long",
			ok:       false,
		},
		{
			name: "invalid user name, username too short",
			user: &User{
				ID:       1,
				Username: "Ja",
				Admin:    true,
			},
			errorMsg: "username is too short",
			ok:       false,
		},
		{
			name:     "invalid user, user is nil",
			user:     nil,
			errorMsg: "user is empty",
			ok:       false,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			err := tc.user.Validate()
			if err != nil {
				assert.Equal(t, tc.errorMsg, err.Error())
				assert.Equal(t, tc.ok, false)
				return
			}
			//Valid User
			assert.Equal(t, tc.ok, true)

		})
	}
}

func TestUser_IsAdmin(t *testing.T) {
	tt := []struct {
		name     string
		user     *User
		expected bool
	}{
		{
			name: "User is Admin",
			user: &User{
				Admin: true,
			},
			expected: true,
		},
		{
			name: "User is not an Admin",
			user: &User{
				Admin: false,
			},
			expected: false,
		},
		{
			name:     "User is nil",
			user:     nil,
			expected: false,
		},
	}

	for _, tc := range tt {
		tc := tc
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			actual := tc.user.IsAdmin()
			assert.Equal(t, actual, tc.expected)
		})
	}
}
