package middleware_test

import (
	"net/http"
	test "net/http/httptest"
	"net/url"
	"testing"

	"github.com/labstack/echo/v4"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"github.com/usual2970/retell/internal/rest/middleware"
)

func TestCORS(t *testing.T) {
	e := echo.New()
	req := test.NewRequest(echo.GET, "/", nil)
	res := test.NewRecorder()
	c := e.NewContext(req, res)

	h := middleware.CORS(echo.HandlerFunc(func(c echo.Context) error {
		return c.NoContent(http.StatusOK)
	}))

	err := h(c)
	require.NoError(t, err)
	assert.Equal(t, "*", res.Header().Get("Access-Control-Allow-Origin"))
}

func TestUrl(t *testing.T) {
	uri := "?orderUri=/play/user-list"
	uriData, err := url.Parse(uri)
	if err != nil {
		t.Errorf("Error parsing URL: %v", err)
		return
	}
	query := uriData.Query()
	orderUri := query.Get("orderUri")
	if orderUri == "" {
		t.Errorf("orderUri not found in URL")
		return
	}
	t.Logf("orderUri: %s", orderUri)

	t.Logf("Path: %s", uriData.Path)
}
