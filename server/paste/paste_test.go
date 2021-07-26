package paste

import (
	"bytes"
	"encoding/json"
	"github.com/gin-gonic/gin"
	"io"
	"net/http"
	"net/http/httptest"
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

func TestCreate(t *testing.T) {
	type Input struct {
		ginParams   map[string]string
		requestBody map[string]interface{}
		mockIPPort  string
	}
	type Expect struct {
		ip     string
		status uint
	}
	var testCases = []struct {
		Input
		Expect
	}{
		{
			Input{map[string]string{
				"namespace": "nobody",
			}, map[string]interface{}{
				"content":       "print('Hello World!')",
				"lang":          "python",
				"self_destruct": false,
			}, "127.0.0.1:10086"},
			Expect{"127.0.0.1", 201},
		},
		{
			Input{map[string]string{
				"namespace": "lucien",
			}, map[string]interface{}{
				"content":       "print('Hello World!')",
				"lang":          "python",
				"self_destruct": true,
				"expire_type":   "",
				"expiration":    "3",
			}, "127.0.0.1:10086"},
			Expect{"127.0.0.1", 400},
		},
	}

	type Response struct {
		Message   string `json:"message"`
		Key       string `json:"key"`
		Namespace string `json:"namespace"`
		Status    uint   `json:"status"`
	}

	for i, testCase := range testCases {
		response := Response{}

		if err := testHandler(testCase.ginParams, testCase.requestBody, testCase.mockIPPort, "POST", Create, &response); err != nil {
			t.Fatal(err.Error())
		}

		if response.Status != testCase.status || response.Namespace != testCase.ginParams["namespace"] {
			t.Fatalf("Test %d | Input: %+v, Expected: %+v, Output: %+v\n", i, testCase.Input, testCase.Expect, response)
		}
	}

	
}

func TestMain(m *testing.M) {
	m.Run()
}
