package helper

import (
	"bytes"
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"gorm.io/driver/postgres"
	"gorm.io/gorm"
	"net/http"
	"net/http/httptest"
	"testing"
)

// SetupTestDatabase sets up a mock database for testing.
func SetupTestDatabase(t *testing.T) (sqlmock.Sqlmock, *gorm.DB) {
	// Create a mock database for testing
	db, dbMock, err := sqlmock.New()
	require.NoError(t, err)

	// Replace the actual database with the mock database for testing
	gormDB, err := gorm.Open(postgres.New(postgres.Config{
		Conn: db,
	}), &gorm.Config{})
	require.NoError(t, err)

	return dbMock, gormDB
}

// PerformRequest performs an HTTP request and returns the response recorder.
func PerformRequest(router *gin.Engine, method, path string, requestBody ...[]byte) *httptest.ResponseRecorder {
	var reqBody []byte
	if len(requestBody) > 0 {
		reqBody = requestBody[0]
	}

	w := httptest.NewRecorder()
	req, _ := http.NewRequest(method, path, bytes.NewBuffer(reqBody))
	req.Header.Set("Content-Type", "application/json")
	router.ServeHTTP(w, req)
	return w
}

// ToJSON converts a value to JSON.
func ToJSON(v interface{}) []byte {
	result, _ := json.Marshal(v)
	return result
}
