package users

import (
	"c/pkg/database"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/wjoseperez20/boletia-currency-api/pkg/helper"
	"github.com/wjoseperez20/boletia-currency-api/pkg/models"
	"net/http"
	"testing"
	"time"
)

func TestLoginUser(t *testing.T) {
	// Given
	r := gin.Default()
	r.POST("/login", LoginUser)

	parseTime, err := time.Parse(time.RFC3339Nano, "2023-11-25T15:30:45.123456Z")
	incomingUser := models.User{
		Username: "test",
		Password: "test",
	}

	dbMock, gormDB := setupTestDatabase(t)
	database.DB = gormDB
	mockUser := models.User{Username: "test", Password: "$2a$14$7z17lzN8ckCiGEQQdbQ2c.XsnJYDunu8SQ1H9BG9EqT4FpVwez68K", CreatedAt: parseTime, UpdatedAt: parseTime}
	dbMock.ExpectQuery(`SELECT \* FROM "users" WHERE username = (.+) ORDER BY "users"."username" LIMIT 1`).
		WithArgs("test").
		WillReturnRows(sqlmock.NewRows([]string{"username", "password", "created_at", "updated_at"}).
			AddRow(mockUser.Username, mockUser.Password, mockUser.CreatedAt, mockUser.UpdatedAt))

	// When
	w := helper.PerformRequest(r, "POST", "/login", toJSON(incomingUser))
	require.Equal(t, http.StatusOK, w.Code)

	var expected map[string]interface{}
	err = json.Unmarshal(w.Body.Bytes(), &expected)

	// Then
	require.NoError(t, err)
	require.NotNil(t, expected["token"])
}
