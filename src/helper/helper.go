package helper

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"pg_bridge_go/global_var"
	"pg_bridge_go/logger"
	"time"

	"github.com/gofiber/fiber/v2"
	"github.com/skip2/go-qrcode"
	"go.uber.org/zap"
	"golang.org/x/crypto/bcrypt"
)

func SendResponse(statusCode int, message interface{}, result interface{}, c *fiber.Ctx) error {
	response := global_var.TRequestResponse{
		Message: message,
		Result:  result,
	}
	if response.Message == nil || response.Message == "" {
		response.Message = http.StatusText(statusCode)
	}
	if statusCode == http.StatusInternalServerError {
		logger.Error("Internal Server Error :"+fmt.Sprintf("%v", response.Message), zap.Int("status_code", statusCode), zap.Any("result", result))
	}
	return c.Status(statusCode).JSON(&response)
}

// GetUsernameFiber extracts the username from Fiber's BasicAuth middleware
func GetUsernameFiber(c *fiber.Ctx) string {
	username := c.Locals("username")
	if usernameStr, ok := username.(string); ok {
		return usernameStr
	}
	return ""
}

// GetAuth extracts basic auth credentials from Fiber context
func GetAuth(c *fiber.Ctx) (string, bool) {
	auth := c.Get("Authorization")
	if auth == "" {
		return "", false
	}

	// For Fiber, we need to parse the basic auth manually or use middleware
	// This is a simplified version - in production use proper basic auth parsing
	username := c.Locals("username")
	if usernameStr, ok := username.(string); ok && usernameStr != "" {
		return usernameStr, true
	}
	return "", false
}

func GenerateOrderID(PGID string) string {
	timestamp := time.Now().Unix()
	return fmt.Sprintf("%s-%d", PGID, timestamp)
}

func GenerateQRCodeBase64(url string) (string, error) {
	var png []byte
	png, err := qrcode.Encode(url, qrcode.Medium, 256)
	if err != nil {
		return "", err
	}
	base64String := base64.StdEncoding.EncodeToString(png)
	return "data:image/png;base64," + base64String, nil
}

func VerifyPassword(password, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil
}

func HashPassword(password string) (string, error) {
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return "", err
	}
	return string(hashedPassword), nil
}

type AuthType int

const (
	AuthNone AuthType = iota
	AuthBasic
	AuthBearer
)

type RequestOptions struct {
	Method      string
	URL         string
	QueryParams map[string]string
	Headers     map[string]string
	Body        interface{}
	AuthType    AuthType
	Username    string
	Password    string
	BearerToken string
	ContentType string
	Client      *http.Client
	Timeout     time.Duration
}

func SendRequest(opt RequestOptions) (interface{}, int, http.Header, error) {
	reqURL, err := url.Parse(opt.URL)
	if err != nil {
		return nil, 0, nil, err
	}

	if opt.QueryParams != nil {
		query := reqURL.Query()
		for key, value := range opt.QueryParams {
			query.Set(key, value)
		}
		reqURL.RawQuery = query.Encode()
	}

	var reqBody io.Reader
	if opt.Body != nil {
		bodyBytes, err := json.Marshal(opt.Body)
		if err != nil {
			return nil, 0, nil, err
		}
		reqBody = bytes.NewBuffer(bodyBytes)
	}

	req, err := http.NewRequest(opt.Method, reqURL.String(), reqBody)
	if err != nil {
		return nil, 0, nil, err
	}

	if opt.Headers != nil {
		for key, value := range opt.Headers {
			req.Header.Set(key, value)
		}
	}

	if opt.ContentType != "" {
		req.Header.Set("Content-Type", opt.ContentType)
	} else if opt.Body != nil {
		req.Header.Set("Content-Type", "application/json")
	}

	switch opt.AuthType {
	case AuthBasic:
		req.SetBasicAuth(opt.Username, opt.Password)
	case AuthBearer:
		req.Header.Set("Authorization", "Bearer "+opt.BearerToken)
	}

	client := opt.Client
	if client == nil {
		client = &http.Client{}
	}

	if opt.Timeout > 0 {
		client.Timeout = opt.Timeout
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, 0, nil, err
	}
	defer resp.Body.Close()

	responseBody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, resp.StatusCode, resp.Header, err
	}

	var result interface{}
	if len(responseBody) > 0 {
		err = json.Unmarshal(responseBody, &result)
		if err != nil {
			return string(responseBody), resp.StatusCode, resp.Header, nil
		}
	}

	return result, resp.StatusCode, resp.Header, nil
}
