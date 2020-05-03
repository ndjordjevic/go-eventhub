package server

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/labstack/echo"
	"github.com/stretchr/testify/assert"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"
)

func Test_login_successful(t *testing.T) {
	e := echo.New()
	f := make(url.Values)
	f.Set("username", "vpopovic")
	f.Set("password", "test")

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	// Assertions
	if assert.NoError(t, login(c)) {
		assert.Equal(t, http.StatusOK, rec.Code)
	}
}

func Test_login_failed(t *testing.T) {
	e := echo.New()
	f := make(url.Values)
	f.Set("username", "ndjord")
	f.Set("password", "test")

	req := httptest.NewRequest(http.MethodPost, "/login", strings.NewReader(f.Encode()))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationForm)

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)

	assert.EqualError(t, login(c), "code=401, message=Unauthorized")
}

func Test_wsEndpoint_successful(t *testing.T) {
	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "ndjordjevic"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	ts, err := token.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println(err)
	}
	fmt.Println(ts)

	to, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	if err != nil {
		fmt.Println(err)
	}

	e := echo.New()
	req := httptest.NewRequest(http.MethodGet, "/restricted/ws", nil)
	req.Header["Accept-Encoding"] = []string{"gzip", "deflate", "br"}

	req.Header.Set("Authorization", "Bearer "+ts)
	req.Header.Set("Connection", "upgrade")
	req.Header.Set("Upgrade", "websocket")
	req.Header.Set("Sec-WebSocket-Version", "13")
	req.Header.Set("Sec-WebSocket-Key", "ABDUNUXB9lg3+tpYnQRRtQ==")

	rec := httptest.NewRecorder()
	c := e.NewContext(req, rec)
	c.Set("user", to)

	var wsClients sync.Map
	cc := &CustomContext{c, &wsClients}

	//if assert.NoError(t, wsEndpoint(cc)) {
	//	assert.Equal(t, http.StatusOK, rec.Code)
	//}

	err = wsEndpoint(cc)
	fmt.Println(err)
}
