package adduser_test

import (
	"bytes"
	"encoding/json"
	"fmt"
	"net/http"
	"net/http/httptest"
	"strconv"
	"testing"

	"github.com/m1al04949/avito-tech-service/internal/http-server/handlers/adduser"
	"github.com/m1al04949/avito-tech-service/internal/http-server/handlers/adduser/mocks"
	"github.com/m1al04949/avito-tech-service/pkg/slogdiscard"
	"github.com/stretchr/testify/mock"
	"github.com/stretchr/testify/require"
)

func TestAddUserHandler(t *testing.T) {
	cases := []struct {
		name      string
		user      string
		respError string
		mockError error
	}{
		{
			name: "Success",
			user: "123",
		},
		{
			name:      "Empty user",
			user:      "",
			respError: "user_id is empty",
		},
		{
			name:      "Invalid USER",
			user:      "someinvalidURL",
			respError: "failed to save user",
		},
	}

	for _, tc := range cases {
		tc := tc

		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()

			userSaverMock := mocks.NewUserSaver(t)

			if tc.respError == "" || tc.mockError != nil {
				userSaverMock.On("SaveUser", tc.user, mock.AnythingOfType("string")).
					Return(int64(1), tc.mockError).
					Once()
			}

			handler := adduser.AddUser(slogdiscard.NewDiscardLogger(), userSaverMock)

			var input string
			user, err := strconv.Atoi(tc.user)
			if err != nil {
				input = fmt.Sprintf(`{"user_id": "%s"}`, tc.user)
			} else {
				input = fmt.Sprintf(`{"user_id": "%d"}`, user)
			}

			req, err := http.NewRequest(http.MethodPost, "/users", bytes.NewReader([]byte(input)))
			require.NoError(t, err)

			rr := httptest.NewRecorder()
			handler.ServeHTTP(rr, req)

			require.Equal(t, rr.Code, http.StatusOK)

			body := rr.Body.String()

			var resp adduser.Response

			require.NoError(t, json.Unmarshal([]byte(body), &resp))

			require.Equal(t, tc.respError, resp.Error)

			// add more checks
		})
	}
}
