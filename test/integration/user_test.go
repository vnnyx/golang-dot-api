package integration

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/vnnyx/golang-dot-api/model/entity"
	"github.com/vnnyx/golang-dot-api/model/web"
	"golang.org/x/crypto/bcrypt"
)

func TestCreateUser(t *testing.T) {
	tests := []struct {
		name               string
		payload            web.UserCreateRequest
		codeExpected       int
		statusCodeExpected string
	}{
		{
			name: "Create User Success",
			payload: web.UserCreateRequest{
				Username:             fmt.Sprintf("username_test_1%d", time.Now().UnixMilli()),
				Email:                fmt.Sprintf("integration_1%d@email.com", time.Now().UnixMilli()),
				Handphone:            "08123456789",
				Password:             "password",
				PasswordConfirmation: "password",
			},
			codeExpected:       http.StatusCreated,
			statusCodeExpected: web.CREATED,
		},
		{
			name: "Password Not Match",
			payload: web.UserCreateRequest{
				Username:             fmt.Sprintf("username_test_2%d", time.Now().UnixMilli()),
				Email:                fmt.Sprintf("integration_2%d@email.com", time.Now().UnixMilli()),
				Handphone:            "08123456789",
				Password:             "password",
				PasswordConfirmation: "wrong_password",
			},
			codeExpected:       http.StatusBadRequest,
			statusCodeExpected: web.BAD_REQUEST,
		},
		{
			name: "Some Field Empty",
			payload: web.UserCreateRequest{
				Username:             fmt.Sprintf("username_test_2%d", time.Now().UnixMilli()),
				Email:                fmt.Sprintf("integration_2%d@email.com", time.Now().UnixMilli()),
				Handphone:            "",
				Password:             "password",
				PasswordConfirmation: "password",
			},
			codeExpected:       http.StatusBadRequest,
			statusCodeExpected: web.BAD_REQUEST,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := transactionRepository.DeleteAllTransaction(ctx)
			assert.NoError(t, err)
			err = userRepository.DeleteAllUser(ctx)
			assert.NoError(t, err)

			requestBody, _ := json.Marshal(tt.payload)

			request := httptest.NewRequest("POST", "/dot-api/user", bytes.NewBuffer(requestBody))
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

func TestGetUserById(t *testing.T) {
	tests := []struct {
		name               string
		codeExpected       int
		statusCodeExpected string
		wantErr            bool
	}{
		{
			name:               "Get User By Id Success",
			codeExpected:       http.StatusOK,
			statusCodeExpected: web.OK,
			wantErr:            false,
		},
		{
			name:               "User Not Found",
			codeExpected:       http.StatusNotFound,
			statusCodeExpected: web.NOT_FOUND,
			wantErr:            true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := transactionRepository.DeleteAllTransaction(ctx)
			assert.NoError(t, err)
			err = userRepository.DeleteAllUser(ctx)
			assert.NoError(t, err)

			password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

			dataDB := entity.User{
				UserID:    "123",
				Username:  "username_test",
				Email:     "email_test@gmail.com",
				Handphone: "08123456789",
				Password:  string(password),
			}

			_, err = userRepository.InsertUser(ctx, dataDB)
			assert.NoError(t, err)

			var request *http.Request
			if !tt.wantErr {
				request = httptest.NewRequest("GET", "/dot-api/user/"+dataDB.UserID, nil)
			} else {
				request = httptest.NewRequest("GET", "/dot-api/user/wrong_id", nil)
			}
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

func TestGetAllUserId(t *testing.T) {
	tests := []struct {
		name               string
		codeExpected       int
		statusCodeExpected string
		wantErr            bool
	}{
		{
			name:               "Get All User Success",
			codeExpected:       http.StatusOK,
			statusCodeExpected: web.OK,
			wantErr:            false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := transactionRepository.DeleteAllTransaction(ctx)
			assert.NoError(t, err)
			err = userRepository.DeleteAllUser(ctx)
			assert.NoError(t, err)

			password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

			dataDB := entity.User{
				UserID:    "123",
				Username:  "username_test",
				Email:     "email_test@gmail.com",
				Handphone: "08123456789",
				Password:  string(password),
			}

			_, err = userRepository.InsertUser(ctx, dataDB)
			assert.NoError(t, err)

			var request *http.Request
			if !tt.wantErr {
				request = httptest.NewRequest("GET", "/dot-api/user/"+dataDB.UserID, nil)
			} else {
				request = httptest.NewRequest("GET", "/dot-api/user/wrong_id", nil)
			}
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

func TestUpdateUserProfile(t *testing.T) {
	tests := []struct {
		name               string
		payload            web.UserUpdateProfileRequest
		codeExpected       int
		statusCodeExpected string
		wanErrNotFound     bool
		wantUnauthorized   bool
	}{
		{
			name: "Update User Success",
			payload: web.UserUpdateProfileRequest{
				Username:  fmt.Sprintf("username_test_1%d", time.Now().UnixMilli()),
				Email:     fmt.Sprintf("integration_1%d@email.com", time.Now().UnixMilli()),
				Handphone: "08123456789",
			},
			codeExpected:       http.StatusOK,
			statusCodeExpected: web.OK,
			wanErrNotFound:     false,
			wantUnauthorized:   false,
		},
		{
			name: "Some Field Empty",
			payload: web.UserUpdateProfileRequest{
				Username:  fmt.Sprintf("username_test_2%d", time.Now().UnixMilli()),
				Email:     fmt.Sprintf("integration_2%d@email.com", time.Now().UnixMilli()),
				Handphone: "",
			},
			codeExpected:       http.StatusBadRequest,
			statusCodeExpected: web.BAD_REQUEST,
			wanErrNotFound:     false,
			wantUnauthorized:   false,
		},
		{
			name: "User Not Found",
			payload: web.UserUpdateProfileRequest{
				Username:  fmt.Sprintf("username_test_3%d", time.Now().UnixMilli()),
				Email:     fmt.Sprintf("integration_3%d@email.com", time.Now().UnixMilli()),
				Handphone: "08123456789",
			},
			codeExpected:       http.StatusNotFound,
			statusCodeExpected: web.NOT_FOUND,
			wanErrNotFound:     true,
			wantUnauthorized:   false,
		},
		{
			name: "Unauthorized",
			payload: web.UserUpdateProfileRequest{
				Username:  fmt.Sprintf("username_test_4%d", time.Now().UnixMilli()),
				Email:     fmt.Sprintf("integration_4%d@email.com", time.Now().UnixMilli()),
				Handphone: "08123456789",
			},
			codeExpected:       http.StatusUnauthorized,
			statusCodeExpected: web.UNAUTHORIZATION,
			wanErrNotFound:     false,
			wantUnauthorized:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := transactionRepository.DeleteAllTransaction(ctx)
			assert.NoError(t, err)
			err = userRepository.DeleteAllUser(ctx)
			assert.NoError(t, err)

			password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

			dataDB := entity.User{
				UserID:    uuid.NewString(),
				Username:  fmt.Sprintf("username_test_%d", time.Now().UnixMilli()),
				Email:     fmt.Sprintf("integration%d@email.com", time.Now().UnixMilli()),
				Handphone: "08123456789",
				Password:  string(password),
			}

			_, err = userRepository.InsertUser(ctx, dataDB)
			assert.NoError(t, err)

			requestBody, _ := json.Marshal(tt.payload)

			var accessToken string
			if !tt.wantUnauthorized {
				accessToken = getAuthorization(web.LoginRequest{Username: dataDB.Username, Password: "password"})
			}

			var request *http.Request
			if !tt.wanErrNotFound {
				request = httptest.NewRequest("PUT", "/dot-api/user/"+dataDB.UserID, bytes.NewBuffer(requestBody))
			} else {
				request = httptest.NewRequest("PUT", "/dot-api/user/wrong_id", bytes.NewBuffer(requestBody))
			}
			request.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			request.Header.Set("Authorization", "Bearer "+accessToken)

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

func TestRemoveUser(t *testing.T) {
	tests := []struct {
		name               string
		codeExpected       int
		statusCodeExpected string
		wanErrNotFound     bool
	}{
		{
			name:               "Remove User Success",
			codeExpected:       http.StatusOK,
			statusCodeExpected: web.OK,
			wanErrNotFound:     false,
		},
		{
			name:               "User Not Found",
			codeExpected:       http.StatusNotFound,
			statusCodeExpected: web.NOT_FOUND,
			wanErrNotFound:     true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := transactionRepository.DeleteAllTransaction(ctx)
			assert.NoError(t, err)
			err = userRepository.DeleteAllUser(ctx)
			assert.NoError(t, err)

			password, _ := bcrypt.GenerateFromPassword([]byte("password"), bcrypt.DefaultCost)

			dataDB := entity.User{
				UserID:    "123",
				Username:  "username_test",
				Email:     "email_test@gmail.com",
				Handphone: "08123456789",
				Password:  string(password),
			}

			_, err = userRepository.InsertUser(ctx, dataDB)
			assert.NoError(t, err)

			var request *http.Request
			if !tt.wanErrNotFound {
				request = httptest.NewRequest("DELETE", "/dot-api/user/"+dataDB.UserID, nil)
			} else {
				request = httptest.NewRequest("DELETE", "/dot-api/user/wrong_id", nil)
			}
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
