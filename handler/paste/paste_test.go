package paste

import (
	"bytes"
	"encoding/json"
	"fmt"
	model "github.com/PasteUs/PasteMeGoBackend/model/paste"
	_ "github.com/PasteUs/PasteMeGoBackend/tests"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
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
	ginParams map[string]string, requestBody map[string]interface{}, header map[string]string,
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
		URL:        &url.URL{},
	}

	for k, v := range header {
		context.Request.Header[k] = []string{v}
	}

	if method == "GET" {
		var rawQueryList []string
		for k, v := range requestBody {
			rawQueryList = append(rawQueryList, fmt.Sprintf("%v=%v", k, v))
		}
		context.Request.URL.RawQuery = strings.Join(rawQueryList, "&")
	} else {
		mockJSONRequest(context, requestBody, method)
	}
	handler(context)
	return json.Unmarshal(recorder.Body.Bytes(), &response)
}

type Input struct {
	ginParams   map[string]string
	requestBody map[string]interface{}
	header      map[string]string
	mockIPPort  string
	method      string
}

type Expect struct {
	ip      string
	status  uint
	message string
	content string
	lang    string
}

type Response struct {
	Message   string `json:"message"`
	Key       string `json:"key"`
	Content   string `json:"content"`
	Lang      string `json:"lang"`
	Namespace string `json:"namespace"`
	Status    uint   `json:"status"`
}

type testCase struct {
	input    Input
	expect   Expect
	response *Response
}

var (
	createTestCaseDict map[string]testCase
	getTestCaseDict    map[string]testCase
)

func creatTestCaseGenerator() map[string]testCase {
	testCaseMap := map[string]testCase{}

	for _, pasteType := range []string{"permanent", "temporary_count", "temporary_time"} {
		for _, password := range []string{"", "_with_password"} {
			s := strings.Split(pasteType, "_")
			expireType := s[len(s)-1]
			testCaseMap[pasteType+password] = testCase{
				Input{
					map[string]string{
						"namespace": "nobody",
					},
					map[string]interface{}{
						"content":       "print('Hello World!')",
						"lang":          "python",
						"password":      password,
						"self_destruct": pasteType != "permanent",
						"expire_type":   expireType,
						"expiration":    1,
					},
					map[string]string{},
					"127.0.0.1:10086", "POST"},
				Expect{"127.0.0.1", 201, "", "", ""},
				&Response{},
			}
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

		testCaseMap[name] = testCase{
			Input{map[string]string{
				"namespace": "nobody",
			},
				map[string]interface{}{
					"content":       content,
					"lang":          lang,
					"password":      "",
					"self_destruct": true,
					"expire_type":   expireType,
					"expiration":    expiration,
				},
				map[string]string{},
				"127.0.0.1:10086", "POST"},
			Expect{"127.0.0.1", expectedStatus, message, "", ""},
			&Response{},
		}
	}
	return testCaseMap
}

func TestCreate(t *testing.T) {
	createTestCaseDict = creatTestCaseGenerator()

	for name, c := range createTestCaseDict {
		t.Run(name, func(t *testing.T) {
			if err := testHandler(c.input.ginParams, c.input.requestBody, c.input.header, c.input.mockIPPort,
				c.input.method, Create, &c.response); err != nil {
				t.Error(err.Error())
			}

			if c.response.Status != c.expect.status {
				t.Errorf("check status failed | expected = %d, actual = %d, message = %s",
					c.expect.status, c.response.Status, c.response.Message)
			} else if c.expect.status == 201 && c.response.Namespace != c.input.ginParams["namespace"] {
				t.Errorf("check namespace failed | expected = %s, actual = %s, message = %s",
					c.input.ginParams["namespace"], c.response.Namespace, c.response.Message)
			} else if c.expect.status != 201 && c.response.Message != c.expect.message {
				t.Errorf("check error message failed | expected = %s, actual = %s",
					c.expect.message, c.response.Message)
			}
		})
	}
}

func getTestCaseGenerator() map[string]testCase {
	testCaseMap := map[string]testCase{}

	for _, pasteType := range []string{"permanent", "temporary_count", "temporary_time"} {
		for _, password := range []string{"", "_with_password"} {
			name := pasteType + password
			testCaseMap[name] = testCase{
				Input{
					map[string]string{
						"namespace": "nobody",
						"key":       createTestCaseDict[name].response.Key,
					},
					map[string]interface{}{
						"password": password,
					},
					map[string]string{"Accept": "application/json"},
					"127.0.0.1:10086", "GET",
				},
				Expect{
					"127.0.0.1",
					200,
					"",
					createTestCaseDict[name].input.requestBody["content"].(string),
					createTestCaseDict[name].input.requestBody["lang"].(string),
				},
				&Response{},
			}
		}
	}

	return testCaseMap
}

func TestGet(t *testing.T) {
	getTestCaseDict = getTestCaseGenerator()

	for name, c := range getTestCaseDict {
		t.Run(name, func(t *testing.T) {
			if err := testHandler(c.input.ginParams, c.input.requestBody, c.input.header, c.input.mockIPPort,
				c.input.method, Get, &c.response); err != nil {
				t.Error(err.Error())
			}

			if c.response.Status != c.expect.status {
				t.Errorf("check status failed | expected = %d, actual = %d, message = %s",
					c.expect.status, c.response.Status, c.response.Message)
			} else if c.expect.status == 200 {
				if c.expect.lang != c.response.Lang {
					t.Errorf("check lang failed | expected = %s, actual = %s, message = %s",
						c.expect.lang, c.response.Lang, c.response.Message)
				} else if c.expect.content != c.response.Content {
					t.Errorf("check content failed | expected = %s, actual = %s, message = %s",
						c.expect.content, c.response.Content, c.response.Message)
				}
			} else if c.expect.status != 200 && c.response.Message != c.expect.message {
				t.Errorf("check error message failed | expected = %s, actual = %s",
					c.expect.message, c.response.Message)
			}
		})
	}
}
