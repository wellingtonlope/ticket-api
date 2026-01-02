package steps

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"

	"github.com/labstack/echo/v4"
	"github.com/wellingtonlope/ticket-api/internal/app/repository"
	"github.com/wellingtonlope/ticket-api/internal/app/usecase/user"
	"github.com/wellingtonlope/ticket-api/internal/domain"
	infrahttp "github.com/wellingtonlope/ticket-api/internal/infra/http"
	"github.com/wellingtonlope/ticket-api/internal/infra/jwt"
	"github.com/wellingtonlope/ticket-api/internal/infra/memory"
)

type BDDTestContext struct {
	app         *echo.Echo
	userRepo    repository.UserRepository
	auth        *jwt.Authenticator
	controller  *infrahttp.UserController
	response    *httptest.ResponseRecorder
	requestBody map[string]interface{}
	rawBody     string
}

func NewBDDTestContext(t *testing.T) *BDDTestContext {
	app := echo.New()
	userRepo := &memory.UserRepository{}
	auth := jwt.NewAuthenticator("secret", 24*time.Hour)

	controller := &infrahttp.UserController{
		UCRegister:        user.NewRegister(userRepo),
		UCLogin:           user.NewLogin(userRepo),
		UCGetAllOperators: user.NewGetAllOperators(userRepo),
		Authenticator:     auth,
	}

	app.POST("/users", func(c echo.Context) error {
		bodyBytes, err := io.ReadAll(c.Request().Body)
		if err != nil {
			return c.JSON(http.StatusBadRequest, map[string]string{"message": err.Error()})
		}
		defer c.Request().Body.Close()

		response := controller.Register(infrahttp.Request{
			Body: string(bodyBytes),
		})

		c.Response().Header().Set(infrahttp.ContentTypeHeader, infrahttp.ContentTypeJSON)
		return c.String(response.HttpCode, response.Body)
	})

	return &BDDTestContext{
		app:        app,
		userRepo:   userRepo,
		auth:       auth,
		controller: controller,
	}
}

func (ctx *BDDTestContext) CreateTestUser(email, name, password string) error {
	now := time.Now()
	user, err := domain.UserRegister(name, email, password, now)
	if err != nil {
		return err
	}

	_, err = ctx.userRepo.Insert(*user)
	return err
}

func (ctx *BDDTestContext) SendPostRequest(path string) error {
	var body []byte
	var err error

	if ctx.rawBody != "" {
		body = []byte(ctx.rawBody)
		ctx.rawBody = ""
	} else {
		if ctx.requestBody == nil {
			ctx.requestBody = make(map[string]interface{})
		}

		body, err = json.Marshal(ctx.requestBody)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %w", err)
		}
	}

	req := httptest.NewRequest(http.MethodPost, path, bytes.NewReader(body))
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	ctx.response = httptest.NewRecorder()

	ctx.app.ServeHTTP(ctx.response, req)
	return nil
}

func (ctx *BDDTestContext) SetRawBody(body string) {
	ctx.rawBody = body
}

func (ctx *BDDTestContext) GetStatusCode() int {
	return ctx.response.Code
}

func (ctx *BDDTestContext) GetResponseBody() map[string]interface{} {
	var body map[string]interface{}
	err := json.Unmarshal(ctx.response.Body.Bytes(), &body)
	if err != nil {
		return nil
	}
	return body
}

func (ctx *BDDTestContext) GetToken() (string, error) {
	body := ctx.GetResponseBody()
	if body == nil {
		return "", fmt.Errorf("response body is empty")
	}

	token, ok := body["token"].(string)
	if !ok {
		return "", fmt.Errorf("token not found in response")
	}

	return token, nil
}

func (ctx *BDDTestContext) SetRequestBody(key string, value interface{}) {
	if ctx.requestBody == nil {
		ctx.requestBody = make(map[string]interface{})
	}
	ctx.requestBody[key] = value
}

func (ctx *BDDTestContext) ClearRequestBody() {
	ctx.requestBody = make(map[string]interface{})
}

type ctxKey struct{}

func SetTestContext(ctx context.Context, testCtx *BDDTestContext) context.Context {
	return context.WithValue(ctx, ctxKey{}, testCtx)
}

func GetTestContext(ctx context.Context) *BDDTestContext {
	testCtx, ok := ctx.Value(ctxKey{}).(*BDDTestContext)
	if !ok {
		return nil
	}
	return testCtx
}
