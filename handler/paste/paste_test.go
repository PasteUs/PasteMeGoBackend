package paste

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/PasteUs/PasteMeGoBackend/handler/common"
	"github.com/PasteUs/PasteMeGoBackend/handler/token"
	model "github.com/PasteUs/PasteMeGoBackend/model/paste"
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

	acceptType := ""

	for k, v := range header {
		context.Request.Header[k] = []string{v}
		if k == "Accept" && strings.Contains(v, "json") {
			acceptType = "json"
		}
	}

	if method == "GET" {
		var rawQueryList []string
		for k, v := range requestBody {
			rawQueryList = append(rawQueryList, fmt.Sprintf("%v=%v", k, v))
		}
		context.Request.URL.RawQuery = strings.Join(rawQueryList, "&")
	} else {
		mockJSONRequest(context, requestBody, method)
		context.Set(token.IdentityKey, ginParams["username"])
	}
	handler(context)
	if method == "GET" && acceptType != "json" {
		p := response.(**Response)
		(*p).Content = recorder.Body.String()
		(*p).Code = http.StatusOK
		return nil
	}
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
	status  int
	message string
	content string
	lang    string
}

type Response struct {
	Message  string `json:"message"`
	Key      string `json:"key"`
	Content  string `json:"content"`
	Lang     string `json:"lang"`
	Username string `json:"username"`
	Code     int    `json:"code"`
}

type testCase struct {
	name     string
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

	for _, pasteType := range []string{"permanent", "temporary"} {
		for _, password := range []string{"", "_with_password"} {
			username := token.Nobody
			if pasteType == "permanent" {
				username = "unittest"
			}
			name := pasteType + password
			testCaseMap[name] = testCase{
				name,
				Input{
					map[string]string{
						"username": username,
					},
					map[string]interface{}{
						"content":       "print('Hello World!')",
						"lang":          "python",
						"password":      password,
						"self_destruct": pasteType != "permanent",
						"expire_second": 60,
						"expire_count":  1,
					},
					map[string]string{},
					"127.0.0.1:10086", "POST"},
				Expect{"127.0.0.1", http.StatusCreated, "", "", ""},
				&Response{},
			}
		}
	}

	for _, name := range []string{
		"bind_failed", "empty_lang", "empty_content",
		"zero_expire_count", "zero_expire_second",
		"month_expiration", "big_expiration", "invalid_lang", // "db_locked",
	} {
		var (
			expectedStatus             = -1
			ExpireSecond   interface{} = model.OneMonth
			expireCount                = 1
			content                    = "print('Hello World!')"
			lang                       = "python"
			message                    = ""
		)

		switch name {
		case "empty_lang":
			lang = ""
			expectedStatus = common.ErrEmptyLang.Code
			message = common.ErrEmptyLang.Error()
		case "empty_content":
			content = ""
			expectedStatus = common.ErrEmptyContent.Code
			message = common.ErrEmptyContent.Error()
		case "bind_failed":
			ExpireSecond = "1"
			expectedStatus = common.ErrWrongParamType.Code
			message = common.ErrWrongParamType.Error()
		case "zero_expire_count":
			expireCount = 0
			expectedStatus = common.ErrZeroExpireCount.Code
			message = common.ErrZeroExpireCount.Error()
		case "zero_expire_second":
			ExpireSecond = 0
			expectedStatus = common.ErrZeroExpireSecond.Code
			message = common.ErrZeroExpireSecond.Error()
		case "month_expiration":
			ExpireSecond = model.OneMonth + 1
			expectedStatus = common.ErrExpireSecondGreaterThanMonth.Code
			message = common.ErrExpireSecondGreaterThanMonth.Error()
		case "big_expiration":
			expireCount = model.MaxCount + 1
			expectedStatus = common.ErrExpireCountGreaterThanMaxCount.Code
			message = common.ErrExpireCountGreaterThanMaxCount.Error()
		case "db_locked":
			expectedStatus = common.ErrQueryDBFailed.Code
			message = common.ErrQueryDBFailed.Error()
		case "invalid_lang":
			lang = "none"
			expectedStatus = common.ErrInvalidLang.Code
			message = common.ErrInvalidLang.Error()
		}

		testCaseMap[name] = testCase{
			name,
			Input{map[string]string{
				"username": token.Nobody,
			},
				map[string]interface{}{
					"content":       content,
					"lang":          lang,
					"password":      "",
					"self_destruct": true,
					"expire_second": ExpireSecond,
					"expire_count":  expireCount,
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

			if c.response.Code != c.expect.status {
				t.Errorf("check status failed | expected = %d, actual = %d, message = %s",
					c.expect.status, c.response.Code, c.response.Message)
			} else if c.expect.status != http.StatusCreated && c.response.Message != c.expect.message {
				t.Errorf("check error message failed | expected = %s, actual = %s",
					c.expect.message, c.response.Message)
			}
		})
	}
}

func getTestCaseGenerator() map[string]testCase {
	testCaseMap := map[string]testCase{}

	for _, pasteType := range []string{"permanent", "temporary"} {
		passwordList := []string{"", "_with_password"}

		if pasteType == "permanent" {
			passwordList = append(passwordList, "_wrong_password")
		}

		for _, password := range passwordList {
			var (
				name        = pasteType + password
				status  int = http.StatusOK
				message     = ""
			)

			if password == "_wrong_password" {
				status = common.ErrWrongPassword.Code
				message = common.ErrWrongPassword.Error()
				createTestCaseDict[name] = createTestCaseDict[pasteType+"_with_password"]
			}

			testCaseMap[name] = testCase{
				name,
				Input{
					map[string]string{
						"username": createTestCaseDict[name].input.ginParams["username"],
						"key":      createTestCaseDict[name].response.Key,
					},
					map[string]interface{}{
						"password": password,
					},
					map[string]string{"Accept": "application/json"},
					"127.0.0.1:10086", "GET",
				},
				Expect{
					"127.0.0.1",
					status,
					message,
					createTestCaseDict[name].input.requestBody["content"].(string),
					createTestCaseDict[name].input.requestBody["lang"].(string),
				},
				&Response{},
			}
		}
	}

	for _, name := range []string{
		"not_found", "invalid_key_length",
		"invalid_key_format", "raw_content", // "db_locked",
	} {
		var (
			key      string
			status   int
			message  string
			header   = map[string]string{"Accept": "application/json"}
			username = token.Nobody
			content  string
		)

		switch name {
		case "not_found":
			key = "12345678"
			status = common.ErrRecordNotFound.Code
			message = common.ErrRecordNotFound.Error()
		case "invalid_key_length":
			key = "123456789"
			status = common.ErrInvalidKeyLength.Code
			message = common.ErrInvalidKeyLength.Error()
		case "invalid_key_format":
			key = "123__456"
			status = common.ErrInvalidKeyFormat.Code
			message = common.ErrInvalidKeyFormat.Error()
		case "raw_content":
			key = createTestCaseDict["permanent"].response.Key
			content = createTestCaseDict["permanent"].input.requestBody["content"].(string)
			username = createTestCaseDict["permanent"].input.ginParams["username"]
			status = http.StatusOK
			header = map[string]string{}
		case "db_locked":
			key = createTestCaseDict["permanent"].response.Key
			status = common.ErrQueryDBFailed.Code
			message = common.ErrQueryDBFailed.Error()
		}

		testCaseMap[name] = testCase{
			name,
			Input{
				map[string]string{
					"username": username,
					"key":      key,
				},
				map[string]interface{}{
					"password": "",
				},
				header,
				"127.0.0.1:10086", "GET",
			},
			Expect{
				"127.0.0.1",
				status,
				message,
				content,
				"",
			},
			&Response{},
		}
	}

	return testCaseMap
}

// TODO get expired temporary testing
func TestGet(t *testing.T) {
	getTestCaseDict = getTestCaseGenerator()

	for name, c := range getTestCaseDict {
		t.Run(name, func(t *testing.T) {
			if err := testHandler(c.input.ginParams, c.input.requestBody, c.input.header, c.input.mockIPPort,
				c.input.method, Get, &c.response); err != nil {
				t.Error(err.Error())
			}

			if c.response.Code != c.expect.status {
				t.Errorf("check status failed | expected = %d, actual = %d, message = %s",
					c.expect.status, c.response.Code, c.response.Message)
			} else if c.expect.status == http.StatusOK {
				if c.expect.lang != c.response.Lang {
					t.Errorf("check lang failed | expected = %s, actual = %s, message = %s",
						c.expect.lang, c.response.Lang, c.response.Message)
				} else if c.expect.content != c.response.Content {
					t.Errorf("check content failed | expected = %s, actual = %s, message = %s",
						c.expect.content, c.response.Content, c.response.Message)
				}
			} else if c.expect.status != http.StatusOK && c.response.Message != c.expect.message {
				t.Errorf("check error message failed | expected = %s, actual = %s",
					c.expect.message, c.response.Message)
			}
		})
	}
}
