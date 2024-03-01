package currencies

import (
	"github.com/DATA-DOG/go-sqlmock"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"github.com/wjoseperez20/boletia-currency-api/pkg/database"
	"github.com/wjoseperez20/boletia-currency-api/pkg/helper"
	"github.com/wjoseperez20/boletia-currency-api/pkg/models"
	"gorm.io/gorm"
	"net/http"
	"net/url"
	"testing"
)

func TestHandleCurrencyRequest_EmptyCurrency(t *testing.T) {
	// Given
	r := gin.Default()
	r.GET("/currency/:name", HandleCurrencyRequest)

	// When
	w := helper.PerformRequest(r, "GET", "/currency/", nil)
	require.Equal(t, http.StatusNotFound, w.Code)
}

func TestHandleCurrencyRequest_InvalidCurrency(t *testing.T) {
	// Given
	r := gin.Default()
	r.GET("/currency/:name", HandleCurrencyRequest)

	dbMock, gormDB := helper.SetupTestDatabase(t)
	defer dbMock.ExpectClose()
	database.DB = gormDB

	dbMock.ExpectQuery(`SELECT \* FROM "currency" WHERE name = (.+) ORDER BY "currency"."id" LIMIT (.+)`).
		WithArgs("INVALID", 1).
		WillReturnError(gorm.ErrRecordNotFound)

	// When
	w := helper.PerformRequest(r, "GET", "/currency/INVALID", nil)
	require.Equal(t, http.StatusNotFound, w.Code)

	expected := `{"error":"Currency is not valid"}`
	require.Equal(t, expected, w.Body.String())

	// Verify all expectations were met
	if err := dbMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHandleCurrencyRequest_InvalidFinitDate(t *testing.T) {
	// Given
	r := gin.Default()
	r.GET("/currency/:name", HandleCurrencyRequest)

	q := url.Values{}
	q.Add("finit", "InvalidDate")
	q.Add("fend", "2024-03-01T19:15:00")

	dbMock, gormDB := helper.SetupTestDatabase(t)
	defer dbMock.ExpectClose()
	database.DB = gormDB

	mockCurrency := models.Currency{
		ID:    1,
		Name:  "USD",
		Value: 1.0,
	}

	dbMock.ExpectQuery(`SELECT \* FROM "currency" WHERE name = (.+) ORDER BY "currency"."id" LIMIT (.+)`).
		WithArgs("USD", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "value"}).
			AddRow(mockCurrency.ID, mockCurrency.Name, mockCurrency.Value))

	// When
	w := helper.PerformRequest(r, "GET", "/currency/usd?"+q.Encode(), nil)
	require.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"error":"Invalid finit date format"}`
	require.Equal(t, expected, w.Body.String())

	// Verify all expectations were met
	if err := dbMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}

func TestHandleCurrencyRequest_InvalidFendDate(t *testing.T) {
	// Given
	r := gin.Default()
	r.GET("/currency/:name", HandleCurrencyRequest)

	q := url.Values{}
	q.Add("finit", "2024-03-01T19:15:00")
	q.Add("fend", "InvalidDate")

	dbMock, gormDB := helper.SetupTestDatabase(t)
	defer dbMock.ExpectClose()
	database.DB = gormDB

	mockCurrency := models.Currency{
		ID:    1,
		Name:  "USD",
		Value: 1.0,
	}

	dbMock.ExpectQuery(`SELECT \* FROM "currency" WHERE name = (.+) ORDER BY "currency"."id" LIMIT (.+)`).
		WithArgs("USD", 1).
		WillReturnRows(sqlmock.NewRows([]string{"id", "name", "value"}).
			AddRow(mockCurrency.ID, mockCurrency.Name, mockCurrency.Value))

	// When
	w := helper.PerformRequest(r, "GET", "/currency/usd?"+q.Encode(), nil)
	require.Equal(t, http.StatusBadRequest, w.Code)

	expected := `{"error":"Invalid fend date format"}`
	require.Equal(t, expected, w.Body.String())

	// Verify all expectations were met
	if err := dbMock.ExpectationsWereMet(); err != nil {
		t.Errorf("there were unfulfilled expectations: %s", err)
	}
}
