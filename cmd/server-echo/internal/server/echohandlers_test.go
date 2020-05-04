package server

import (
	"fmt"
	"github.com/dgrijalva/jwt-go"
	"github.com/gorilla/websocket"
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

type WsHandler struct {
	handler echo.HandlerFunc
}

func (h *WsHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	e := echo.New()
	c := e.NewContext(r, w)

	token := jwt.New(jwt.SigningMethodHS256)

	claims := token.Claims.(jwt.MapClaims)
	claims["name"] = "ndjordjevic"
	claims["admin"] = true
	claims["exp"] = time.Now().Add(time.Hour * 72).Unix()

	ts, err := token.SignedString([]byte("secret"))
	if err != nil {
		fmt.Println(err)
	}

	to, err := jwt.Parse(ts, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, fmt.Errorf("Unexpected signing method: %v", token.Header["alg"])
		}
		return []byte("secret"), nil
	})

	c.Set("user", to)

	var wsClients sync.Map
	cc := &CustomContext{c, &wsClients}

	forever := make(chan struct{})
	h.handler(cc)
	<-forever
}

func Test_wsEndpoint_successful(t *testing.T) {
	h := WsHandler{handler: wsEndpoint}
	server := httptest.NewServer(http.HandlerFunc(h.ServeHTTP))
	defer server.Close()

	wsURL := "ws" + strings.TrimPrefix(server.URL, "http") + "/restricted/ws"
	_, _, err := websocket.DefaultDialer.Dial(wsURL, nil)
	assert.Nil(t, err, err)
}
