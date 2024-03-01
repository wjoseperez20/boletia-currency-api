package users

import (
	"encoding/json"
	"errors"
	"github.com/DATA-DOG/go-sqlmock"
	"gorm.io/gorm"
	"net/http"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/wjoseperez20/boletia-currency-api/pkg/database"
	"github.com/wjoseperez20/boletia-currency-api/pkg/helper"
	"github.com/wjoseperez20/boletia-currency-api/pkg/models"
)

func TestLoginUser_BadRequest(t *testing.T) {
	// Given
	r := gin.Default()
	r.POST("/login", LoginUser)

	// When
	w := helper.PerformRequest(r, "POST", "/login", nil)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

	// Then
	require.NotNil(t, response)
}

func TestLoginUser_Unauthorized(t *testing.T) {
	// Given
	r := gin.Default()
	r.POST("/login", LoginUser)

	incomingUser := models.User{
		Username: "test",
		Password: "test",
	}

	dbMock, gormDB := helper.SetupTestDatabase(t)
	defer dbMock.ExpectClose()
	database.DB = gormDB

	dbMock.ExpectQuery(`SELECT \* FROM "user" WHERE username = (.+) ORDER BY "user"."id" LIMIT (.+)`).
		WithArgs("test", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// When
	w := helper.PerformRequest(r, "POST", "/login", helper.ToJSON(incomingUser))
	require.Equal(t, http.StatusUnauthorized, w.Code)

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

	// Then
	require.NotNil(t, response)
}

func TestLoginUser_InternalServerError(t *testing.T) {
	// Given
	r := gin.Default()
	r.POST("/login", LoginUser)

	incomingUser := models.User{
		Username: "Test",
		Password: "Test",
	}

	dbMock, gormDB := helper.SetupTestDatabase(t)
	defer dbMock.ExpectClose()
	database.DB = gormDB

	dbMock.ExpectQuery(`SELECT \* FROM "user" WHERE username = (.+) ORDER BY "user"."id" LIMIT (.+)`).
		WithArgs("Test", 1).
		WillReturnError(errors.New("internal error"))

	// When
	w := helper.PerformRequest(r, "POST", "/login", helper.ToJSON(incomingUser))
	require.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

	// Then
	require.NotNil(t, response)
}

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

func TestRegisterUser_BadRequest(t *testing.T) {
	// Given
	r := gin.Default()
	r.POST("/register", RegisterUser)

	// When
	w := helper.PerformRequest(r, "POST", "/register", nil)
	require.Equal(t, http.StatusBadRequest, w.Code)

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

	// Then
	require.NotNil(t, response)

}

func TestRegisterUser_InternalServerError(t *testing.T) {
	// Given
	r := gin.Default()
	r.POST("/register", RegisterUser)

	incomingUser := models.User{
		Username: "test",
		Password: "test",
	}

	dbMock, gormDB := helper.SetupTestDatabase(t)
	defer dbMock.ExpectClose()
	database.DB = gormDB

	dbMock.ExpectQuery(`SELECT \* FROM "user" WHERE username = (.+) ORDER BY "user"."id" LIMIT (.+)`).
		WithArgs("test", 1).
		WillReturnError(errors.New("internal error"))

	// When
	w := helper.PerformRequest(r, "POST", "/register", helper.ToJSON(incomingUser))
	require.Equal(t, http.StatusInternalServerError, w.Code)

	var response map[string]interface{}
	require.NoError(t, json.Unmarshal(w.Body.Bytes(), &response))

	// Then
	require.NotNil(t, response)
}
