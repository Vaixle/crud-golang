package midleware

import (
	"github.com/Vaixle/crud-golang/internal/config"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"net/http/httptest"
	"testing"
)

func TestBasicAuth(t *testing.T) {
	testTable := []struct {
		name               string
		login              string
		password           string
		expectedStatusCode uint16
	}{
		{
			name:               "OK",
			login:              "admin",
			password:           "admin",
			expectedStatusCode: 200,
		},
		{
			name:               "UNAUTHORIZED",
			expectedStatusCode: 401,
		},
	}

	for _, testCase := range testTable {
		t.Run(testCase.name, func(t *testing.T) {
			err := config.Init()
			if err != nil {
				return
			}

			r := gin.New()
			r.GET("/", BasicAuth(), nil)
			w := httptest.NewRecorder()
			req := httptest.NewRequest("GET", "/", nil)
			req.SetBasicAuth(testCase.login, testCase.password)

			r.ServeHTTP(w, req)

			require.NoError(t, err)
			assert.Equal(t, w.Code, testCase.expectedStatusCode)
		})
	}
}
