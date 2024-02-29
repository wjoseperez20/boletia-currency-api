package currencies

import (
	"encoding/json"
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/wjoseperez20/boletia-currency-api/pkg/database"
	"github.com/wjoseperez20/boletia-currency-api/pkg/helper"
	"github.com/wjoseperez20/boletia-currency-api/pkg/models"
	"net/http"
	"testing"
	"time"
)

func TestFindCurrency_Success(t *testing.T) {
	// Given
	router := gin.Default()
	router.GET("/currencies/:name", FindCurrency)

	parseTime, err := time.Parse(time.RFC3339Nano, "2024-02-19T15:30:45.123456Z")
	require.NoError(t, err)

	dbMock, gormDB := helper.SetupTestDatabase(t)
	database.DB = gormDB

	mockCurrency := models.Currency{ID: 1, Name: "USD", Code: "USD", Value: 1.0, CreatedAt: parseTime}
	dbMock.ExpectQuery(`SELECT \* FROM "currency" WHERE name = (.+) ORDER BY "currency"."id" LIMIT (.+)`).
		WithArgs("USD", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "code", "value", "created_at"}).
			AddRow(mockCurrency.ID, mockCurrency.Name, mockCurrency.Code, mockCurrency.Value, mockCurrency.CreatedAt))

	// When
	w := helper.PerformRequest(router, "GET", "/currencies/USD", nil)
	require.Equal(t, http.StatusOK, w.Code)

	var expected models.Currency
	err = json.Unmarshal(w.Body.Bytes(), &expected)

	// Then
	require.NoError(t, err)
	require.Equal(t, mockCurrency.ID, expected.ID)
	require.Equal(t, mockCurrency.Name, expected.Name)
	require.Equal(t, mockCurrency.Code, expected.Code)
	require.Equal(t, mockCurrency.Value, expected.Value)

	// Verify all expectations were met
	if err := dbMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestFindCurrency_NotFound(t *testing.T) {
	// Given
	router := gin.Default()
	router.GET("/currencies/:name", FindCurrency)

	dbMock, gormDB := helper.SetupTestDatabase(t)
	database.DB = gormDB

	// When
	w := helper.PerformRequest(router, "GET", "/currencies/USD", nil)
	require.Equal(t, http.StatusNotFound, w.Code)

	// Verify all expectations were met
	if err := dbMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
