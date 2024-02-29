package healtcheck

import (
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestHealthcheck(t *testing.T) {
	// Given
	router := gin.New()
	router.GET("/healthcheck", Healthcheck)

	// When
	w := httptest.NewRecorder()
	req, _ := http.NewRequest("GET", "/healthcheck", nil)
	router.ServeHTTP(w, req)

	// Then
	require.Equal(t, http.StatusOK, w.Code)

	expected := `{"message":"ok"}`
	require.Equal(t, expected, w.Body.String())
}
