package paste

import (
	"bytes"
	"encoding/json"
	model "github.com/PasteUs/PasteMeGoBackend/model/paste"
	_ "github.com/PasteUs/PasteMeGoBackend/tests"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"testing"
)

func mockJSONRequest(c *gin.Context, jsonMap map[string]interface{}, method string) {
	c.Request.Method = method
	c.Request.Header.Set("Content-Type", "application/json")

	jsonBytes, err := json.Marshal(jsonMap)
	if err != nil {
		panic(err)
	}

	// the request body must be an io.ReadCloser
	// the bytes buffer though doesn't implement io.Closer,
	// so you wrap it in a no-op closer
	c.Request.Body = io.NopCloser(bytes.NewBuffer(jsonBytes))
}

func testHandler(
	ginParams map[string]string, requestBody map[string]interface{},
	mockIPPort string, method string, handler func(*gin.Context), response interface{},
) error {
	recorder := httptest.NewRecorder()
	context, _ := gin.CreateTestContext(recorder)

	var params []gin.Param
	for k, v := range ginParams {
		params = append(params, gin.Param{Key: k, Value: v})
	}
	context.Params = params

	context.Request = &http.Request{
		Header:     http.Header{},
		RemoteAddr: mockIPPort,
	}
	mockJSONRequest(context, requestBody, method)
	handler(context)
	return json.Unmarshal(recorder.Body.Bytes(), &response)
}

type Input struct {
	ginParams   map[string]string
	requestBody map[string]interface{}
	mockIPPort  string
	method      string
}

type Expect struct {
	ip      string
	status  uint
	message string
}

type testCase struct {
	name string
	Input
	Expect
}

func creatCaseGenerator() []testCase {
	var testCaseList []testCase

	for _, pasteType := range []string{"permanent", "temporary_count", "temporary_time"} {
		for _, password := range []string{"", "_with_password"} {
			s := strings.Split(pasteType, "_")
			expireType := s[len(s)-1]
			testCaseList = append(testCaseList, testCase{
				pasteType + password,
				Input{map[string]string{
					"namespace": "nobody",
				}, map[string]interface{}{
					"content":       "print('Hello World!')",
					"lang":          "python",
					"password":      password,
					"self_destruct": pasteType != "permanent",
					"expire_type":   expireType,
					"expiration":    1,
				}, "127.0.0.1:10086", "POST"},
				Expect{"127.0.0.1", 201, ""},
			})
		}
	}

	for _, name := range []string{
		"bind_failed", "empty_lang", "empty_content",
		"zero_expiration", "empty_expire_type", "other_expire_type",
		"month_expiration", "big_expiration",
	} {
		var (
			expectedStatus uint        = 400
			expiration     interface{} = model.OneMonth
			expireType                 = "time"
			content                    = "print('Hello World!')"
			lang                       = "python"
			message                    = ""
		)

		switch name {
		case "empty_lang":
			lang = ""
			message = ErrEmptyLang.Error()
		case "empty_content":
			content = ""
			message = ErrEmptyContent.Error()
		case "bind_failed":
			expiration = "1"
			message = "wrong param type"
		case "zero_expiration":
			expiration = 0
			message = ErrZeroExpiration.Error()
		case "empty_expire_type":
			expireType = ""
			message = ErrEmptyExpireType.Error()
		case "other_expire_type":
			expireType = "other"
			message = ErrInvalidExpireType.Error()
		case "month_expiration":
			expiration = model.OneMonth + 1
			message = ErrExpirationGreaterThanMonth.Error()
		case "big_expiration":
			expireType = "count"
			expiration = model.MaxCount + 1
			message = ErrExpirationGreaterThanMaxCount.Error()
		}

		testCaseList = append(testCaseList, testCase{
			name,
			Input{map[string]string{
				"namespace": "nobody",
			}, map[string]interface{}{
				"content":       content,
				"lang":          lang,
				"password":      "",
				"self_destruct": true,
				"expire_type":   expireType,
				"expiration":    expiration,
			}, "127.0.0.1:10086", "POST"},
			Expect{"127.0.0.1", expectedStatus, message},
		})
	}
	return testCaseList
}

func TestCreate(t *testing.T) {
	type Response struct {
		Message   string `json:"message"`
		Key       string `json:"key"`
		Namespace string `json:"namespace"`
		Status    uint   `json:"status"`
	}

	testCaseList := creatCaseGenerator()

	for i, c := range testCaseList {
		t.Run(c.name, func(t *testing.T) {
			response := Response{}

			if err := testHandler(c.Input.ginParams, c.Input.requestBody, c.Input.mockIPPort,
				c.Input.method, Create, &response); err != nil {
				t.Error(err.Error())
			}

			if response.Status != c.status {
				t.Errorf("test %d | check status failed | expected = %d, actual = %d, message = %s",
					i, c.status, response.Status, response.Message)
			} else if c.status == 201 && response.Namespace != c.ginParams["namespace"] {
				t.Errorf("test %d | check namespace failed | expected = %s, actual = %s, message = %s",
					i, c.Input.ginParams["namespace"], response.Namespace, response.Message)
			} else if c.status != 201 && response.Message != c.Expect.message {
				t.Errorf("test %d | check error message failed | expected = %s, actual = %s",
					i, c.Expect.message, response.Message)
			}
		})
	}
}
