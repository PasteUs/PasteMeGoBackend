package paste

import (
    "bytes"
    "encoding/json"
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

    type testCase struct {
        name string
        Input
        Expect
    }

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
                }, "127.0.0.1:10086"},
                Expect{"127.0.0.1", 201},
            })
        }
    }

    for _, name := range []string{"bind_failed", "invalid_param", "db_error"} {
        var expectedStatus uint
        var expiration interface{}
        if name == "db_error" {
            expectedStatus = 500
            expiration = 1
        } else {
            expectedStatus = 400
            expiration = "1"
        }

        testCaseList = append(testCaseList, testCase{
            name,
            Input{map[string]string{
                "namespace": "nobody",
            }, map[string]interface{}{
                "content":       "print('Hello World!')",
                "lang":          "python",
                "password":      "",
                "self_destruct": true,
                "expire_type":   "",
                "expiration":    expiration,
            }, "127.0.0.1:10086"},
            Expect{"127.0.0.1", expectedStatus},
        })
    }

    type Response struct {
        Message   string `json:"message"`
        Key       string `json:"key"`
        Namespace string `json:"namespace"`
        Status    uint   `json:"status"`
    }

    for i, c := range testCaseList {
        t.Run(c.name, func(t *testing.T) {
            response := Response{}

            if err := testHandler(c.ginParams, c.requestBody, c.mockIPPort, "POST", Create, &response); err != nil {
                t.Error(err.Error())
            }

            if response.Status != c.status {
                t.Errorf("test %d | check status failed | expected = %d, actual = %d, message = %s",
                    i, c.status, response.Status, response.Message)
            } else if c.status == 201 && response.Namespace != c.ginParams["namespace"] {
                t.Errorf("test %d | check namespace failed | expected = %s, actual = %s, message = %s",
                    i, c.Input.ginParams["namespace"], response.Namespace, response.Message)
            }
        })
    }

}
