package users

import (
	"encoding/json"
	"net/http"
	"testing"
	"time"

	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/wjoseperez20/boletia-currency-api/pkg/database"
	"github.com/wjoseperez20/boletia-currency-api/pkg/helper"
	"github.com/wjoseperez20/boletia-currency-api/pkg/models"
)

func TestLoginUser_Success(t *testing.T) {
	// Given
	r := gin.Default()
	r.POST("/login", LoginUser)

	parseTime, err := time.Parse(time.RFC3339Nano, "2024-02-19T15:30:45.123456Z")
	require.NoError(t, err)

	incomingUser := models.User{
		Username: "test",
		Password: "test",
	}

	dbMock, gormDB := helper.SetupTestDatabase(t)
	defer dbMock.ExpectClose()
	database.DB = gormDB

	mockUser := models.User{
		ID:        1,
		Username:  "test",
		Password:  "$2a$14$q6TbZ6LL71UjKldZheALMu5jS6AA3/BbFyB6AviKCO9B5LQJ4WMcq",
		CreatedAt: parseTime,
		UpdatedAt: parseTime,
	}
	dbMock.ExpectQuery(`SELECT \* FROM "user" WHERE username = (.+) ORDER BY "user"."id" LIMIT (.+)`).
		WithArgs("test", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "username", "password", "created_at", "updated_at"}).
			AddRow(mockUser.ID, mockUser.Username, mockUser.Password, mockUser.CreatedAt, mockUser.UpdatedAt))

	// When
	w := helper.PerformRequest(r, "POST", "/login", helper.ToJSON(incomingUser))
	require.Equal(t, http.StatusOK, w.Code)

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

	// Then
	require.NotNil(t, response["token"])
}

func TestRegisterUser_Fail(t *testing.T) {
	// Given
	r := gin.Default()
	r.POST("/register", RegisterUser)

	parseTime, err := time.Parse(time.RFC3339Nano, "2024-02-19T15:30:45.123456Z")
	require.NoError(t, err)

	incomingUser := models.User{
		Username: "test",
		Password: "test",
	}

	hashedPassword := "$2a$14$7z17lzN8ckCiGEQQdbQ2c.XsnJYDunu8SQ1H9BG9EqT4FpVwez68K"

	dbMock, gormDB := helper.SetupTestDatabase(t)
	defer dbMock.ExpectClose()
	database.DB = gormDB

	dbMock.ExpectExec(`INSERT INTO "user" (.+) VALUES (.+)`).
		WithArgs("test", hashedPassword, parseTime, parseTime).
		WillReturnError(nil)

	// When
	w := helper.PerformRequest(r, "POST", "/register", helper.ToJSON(incomingUser))
	require.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

	// Then
	require.NotNil(t, response)
}
