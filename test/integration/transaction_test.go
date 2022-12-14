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

func TestCreateTransaction(t *testing.T) {
	tests := []struct {
		name                string
		payload             web.TransactionCreateRequest
		codeExpected        int
		statusCodeExpected  string
		wantUnauthorized    bool
		wantErrUserNotFound bool
	}{
		{
			name: "Create Transaction Success",
			payload: web.TransactionCreateRequest{
				Name: "product_test",
			},
			codeExpected:        http.StatusCreated,
			statusCodeExpected:  web.CREATED,
			wantUnauthorized:    false,
			wantErrUserNotFound: false,
		},
		{
			name: "User Not Found",
			payload: web.TransactionCreateRequest{
				Name: "product_test",
			},
			codeExpected:        http.StatusNotFound,
			statusCodeExpected:  web.NOT_FOUND,
			wantUnauthorized:    false,
			wantErrUserNotFound: true,
		},
		{
			name: "Blank Field",
			payload: web.TransactionCreateRequest{
				Name: "",
			},
			codeExpected:        http.StatusBadRequest,
			statusCodeExpected:  web.BAD_REQUEST,
			wantUnauthorized:    false,
			wantErrUserNotFound: false,
		},
		{
			name: "Unauthorized",
			payload: web.TransactionCreateRequest{
				Name: "product_test",
			},
			codeExpected:        http.StatusUnauthorized,
			statusCodeExpected:  web.UNAUTHORIZATION,
			wantUnauthorized:    true,
			wantErrUserNotFound: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = transactionRepository.DeleteAllTransaction(ctx)
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

			requestBody, _ := json.Marshal(tt.payload)

			var accessToken string
			if !tt.wantUnauthorized {
				accessToken = getAuthorization(web.LoginRequest{Username: dataDB.Username, Password: "password"})
			}

			var request *http.Request
			if !tt.wantErrUserNotFound {
				request = httptest.NewRequest("POST", "/dot-api/transaction?user_id="+dataDB.UserID, bytes.NewBuffer(requestBody))
			} else {
				request = httptest.NewRequest("POST", "/dot-api/transaction?user_id=wrong_id", bytes.NewBuffer(requestBody))
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

func TestGetTransactionById(t *testing.T) {
	tests := []struct {
		name               string
		codeExpected       int
		statusCodeExpected string
		wantUnauthorized   bool
		wantErrNotFound    bool
	}{
		{
			name:               "Get Transaction By Id Success",
			codeExpected:       http.StatusOK,
			statusCodeExpected: web.OK,
			wantUnauthorized:   false,
			wantErrNotFound:    false,
		},
		{
			name:               "Transaction Not Found",
			codeExpected:       http.StatusNotFound,
			statusCodeExpected: web.NOT_FOUND,
			wantUnauthorized:   false,
			wantErrNotFound:    true,
		},
		{
			name:               "Unathorized",
			codeExpected:       http.StatusUnauthorized,
			statusCodeExpected: web.UNAUTHORIZATION,
			wantUnauthorized:   true,
			wantErrNotFound:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = transactionRepository.DeleteAllTransaction(ctx)
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

			transactionDB := entity.Transaction{
				TransactionID: "123",
				Name:          "product_test",
				UserID:        "123",
			}

			_, _ = transactionRepository.InsertTransaction(ctx, transactionDB)

			var accessToken string
			if !tt.wantUnauthorized {
				accessToken = getAuthorization(web.LoginRequest{Username: dataDB.Username, Password: "password"})
			}

			var request *http.Request
			if !tt.wantErrNotFound {
				request = httptest.NewRequest("GET", "/dot-api/transaction/"+transactionDB.TransactionID, nil)
			} else {
				request = httptest.NewRequest("GET", "/dot-api/transaction/wrong_id", nil)
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

func TestGetAllTransaction(t *testing.T) {
	tests := []struct {
		name               string
		codeExpected       int
		statusCodeExpected string
		wantUnauthorized   bool
	}{
		{
			name:               "Get All Transaction Success",
			codeExpected:       http.StatusOK,
			statusCodeExpected: web.OK,
			wantUnauthorized:   false,
		},
		{
			name:               "Unauthorized",
			codeExpected:       http.StatusUnauthorized,
			statusCodeExpected: web.UNAUTHORIZATION,
			wantUnauthorized:   true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = transactionRepository.DeleteAllTransaction(ctx)
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

			transactionDB := entity.Transaction{
				TransactionID: "123",
				Name:          "product_test",
				UserID:        "123",
			}

			_, _ = transactionRepository.InsertTransaction(ctx, transactionDB)

			var accessToken string
			if !tt.wantUnauthorized {
				accessToken = getAuthorization(web.LoginRequest{Username: dataDB.Username, Password: "password"})
			}

			request := httptest.NewRequest("GET", "/dot-api/transaction", nil)
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

func TestGetTransactionByUserId(t *testing.T) {
	tests := []struct {
		name                string
		codeExpected        int
		statusCodeExpected  string
		wantUnauthorized    bool
		wantErrUserNotFound bool
	}{
		{
			name:                "Get Transaction By UserId Success",
			codeExpected:        http.StatusOK,
			statusCodeExpected:  web.OK,
			wantUnauthorized:    false,
			wantErrUserNotFound: false,
		},
		{
			name:                "User Not Found",
			codeExpected:        http.StatusNotFound,
			statusCodeExpected:  web.NOT_FOUND,
			wantUnauthorized:    false,
			wantErrUserNotFound: true,
		},
		{
			name:                "Unathorized",
			codeExpected:        http.StatusUnauthorized,
			statusCodeExpected:  web.UNAUTHORIZATION,
			wantUnauthorized:    true,
			wantErrUserNotFound: false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = transactionRepository.DeleteAllTransaction(ctx)
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

			transactionDB := entity.Transaction{
				TransactionID: "123",
				Name:          "product_test",
				UserID:        "123",
			}

			_, _ = transactionRepository.InsertTransaction(ctx, transactionDB)

			var accessToken string
			if !tt.wantUnauthorized {
				accessToken = getAuthorization(web.LoginRequest{Username: dataDB.Username, Password: "password"})
			}

			var request *http.Request
			if !tt.wantErrUserNotFound {
				request = httptest.NewRequest("GET", "/dot-api/transaction/user?user_id="+dataDB.UserID, nil)
			} else {
				request = httptest.NewRequest("GET", "/dot-api/transaction/user?user_id=wrong_id", nil)
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

func TestUpdateTransaction(t *testing.T) {
	tests := []struct {
		name               string
		payload            web.TransactionUpdateRequest
		codeExpected       int
		statusCodeExpected string
		wantUnauthorized   bool
		wantErrNotFound    bool
	}{
		{
			name: "Create Transaction Success",
			payload: web.TransactionUpdateRequest{
				Name: "product_update_test",
			},
			codeExpected:       http.StatusOK,
			statusCodeExpected: web.OK,
			wantUnauthorized:   false,
			wantErrNotFound:    false,
		},
		{
			name: "User Not Found",
			payload: web.TransactionUpdateRequest{
				Name: "product_update_test",
			},
			codeExpected:       http.StatusNotFound,
			statusCodeExpected: web.NOT_FOUND,
			wantUnauthorized:   false,
			wantErrNotFound:    true,
		},
		{
			name: "Blank Field",
			payload: web.TransactionUpdateRequest{
				Name: "",
			},
			codeExpected:       http.StatusBadRequest,
			statusCodeExpected: web.BAD_REQUEST,
			wantUnauthorized:   false,
			wantErrNotFound:    false,
		},
		{
			name: "Unauthorized",
			payload: web.TransactionUpdateRequest{
				Name: "product_update_test",
			},
			codeExpected:       http.StatusUnauthorized,
			statusCodeExpected: web.UNAUTHORIZATION,
			wantUnauthorized:   true,
			wantErrNotFound:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = transactionRepository.DeleteAllTransaction(ctx)
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

			transactionDB := entity.Transaction{
				TransactionID: "123",
				Name:          "product_test",
				UserID:        "123",
			}

			_, _ = transactionRepository.InsertTransaction(ctx, transactionDB)

			requestBody, _ := json.Marshal(tt.payload)

			var accessToken string
			if !tt.wantUnauthorized {
				accessToken = getAuthorization(web.LoginRequest{Username: dataDB.Username, Password: "password"})
			}

			var request *http.Request
			if !tt.wantErrNotFound {
				request = httptest.NewRequest("PATCH", "/dot-api/transaction/"+dataDB.UserID, bytes.NewBuffer(requestBody))
			} else {
				request = httptest.NewRequest("PATCH", "/dot-api/transaction/wrong_id", bytes.NewBuffer(requestBody))
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

func TestDeleteTransaction(t *testing.T) {
	tests := []struct {
		name               string
		codeExpected       int
		statusCodeExpected string
		wantUnauthorized   bool
		wantErrNotFound    bool
	}{
		{
			name:               "Delete Transaction Success",
			codeExpected:       http.StatusOK,
			statusCodeExpected: web.OK,
			wantUnauthorized:   false,
			wantErrNotFound:    false,
		},
		{
			name:               "Transaction Not Found",
			codeExpected:       http.StatusNotFound,
			statusCodeExpected: web.NOT_FOUND,
			wantUnauthorized:   false,
			wantErrNotFound:    true,
		},
		{
			name:               "Unathorized",
			codeExpected:       http.StatusUnauthorized,
			statusCodeExpected: web.UNAUTHORIZATION,
			wantUnauthorized:   true,
			wantErrNotFound:    false,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_ = transactionRepository.DeleteAllTransaction(ctx)
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

			transactionDB := entity.Transaction{
				TransactionID: "123",
				Name:          "product_test",
				UserID:        "123",
			}

			_, _ = transactionRepository.InsertTransaction(ctx, transactionDB)

			var accessToken string
			if !tt.wantUnauthorized {
				accessToken = getAuthorization(web.LoginRequest{Username: dataDB.Username, Password: "password"})
			}

			var request *http.Request
			if !tt.wantErrNotFound {
				request = httptest.NewRequest("DELETE", "/dot-api/transaction/"+transactionDB.TransactionID, nil)
			} else {
				request = httptest.NewRequest("DELETE", "/dot-api/transaction/wrong_id", nil)
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
