package integration

import (
	"bytes"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vnnyx/golang-dot-api/model/entity"
	"github.com/vnnyx/golang-dot-api/model/web"
	"golang.org/x/crypto/bcrypt"
)

func TestLogin(t *testing.T) {
	tests := []struct {
		name               string
		payload            web.LoginRequest
		codeExpected       int
		statusCodeExpected string
	}{
		{
			name: "Login Success",
			payload: web.LoginRequest{
				Username: "username_test",
				Password: "password",
			},
			codeExpected:       http.StatusOK,
			statusCodeExpected: web.OK,
		},
		{
			name: "Login Failed",
			payload: web.LoginRequest{
				Username: "username_test",
				Password: "wrong_password",
			},
			codeExpected:       http.StatusUnauthorized,
			statusCodeExpected: web.UNAUTHORIZATION,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = userRepository.DeleteAllUser(ctx)

			password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

			dataDB := entity.User{
				UserID:    "123",
				Username:  "username_test",
				Email:     "email_test@gmail.com",
				Handphone: "08123456789",
				Password:  string(password),
			}

			_, _ = userRepository.InsertUser(ctx, dataDB)

			requestBody, _ := json.Marshal(tt.payload)

			request := httptest.NewRequest("POST", "/dot-api/login", bytes.NewBuffer(requestBody))
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)

			recorder := httptest.NewRecorder()

			app.ServeHTTP(recorder, request)
			response := recorder.Result()

			responseBody, _ := io.ReadAll(response.Body)
			webResponse := web.WebResponse{}
			json.Unmarshal(responseBody, &webResponse)
			assert.Equal(t, tt.codeExpected, webResponse.Code)
			assert.Equal(t, tt.statusCodeExpected, webResponse.Status)
		})
	}
}

func TestLogout(t *testing.T) {
	tests := []struct {
		name               string
		codeExpected       int
		statusCodeExpected string
		wantUnathorized    bool
	}{
		{
			name:               "Logout Success",
			codeExpected:       http.StatusOK,
			statusCodeExpected: web.OK,
			wantUnathorized:    false,
		},
		{
			name:               "Unathorized",
			codeExpected:       http.StatusUnauthorized,
			statusCodeExpected: web.UNAUTHORIZATION,
			wantUnathorized:    true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = userRepository.DeleteAllUser(ctx)
			_ = authRepository.FlushAll(ctx)

			password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

			dataDB := entity.User{
				UserID:    "123",
				Username:  "username_test",
				Email:     "email_test@gmail.com",
				Handphone: "08123456789",
				Password:  string(password),
			}

			_, _ = userRepository.InsertUser(ctx, dataDB)

			accessToken := getAuthorization(web.LoginRequest{Username: "username_test", Password: "password"})

			request := httptest.NewRequest("POST", "/dot-api/logout", nil)
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			if !tt.wantUnathorized {
				request.Header.Set("Authorization", "Bearer "+accessToken)
			}

			recorder := httptest.NewRecorder()

			app.ServeHTTP(recorder, request)
			response := recorder.Result()

			responseBody, _ := io.ReadAll(response.Body)
			webResponse := web.WebResponse{}
			json.Unmarshal(responseBody, &webResponse)
			assert.Equal(t, tt.codeExpected, webResponse.Code)
			assert.Equal(t, tt.statusCodeExpected, webResponse.Status)
		})
	}
}
