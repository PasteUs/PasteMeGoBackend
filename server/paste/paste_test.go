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

func TestCreate(t *testing.T) {
    recorder := httptest.NewRecorder()
    context, _ := gin.CreateTestContext(recorder)

    context.Params = []gin.Param{{Key: "namespace", Value: "nobody"}}

    context.Request = &http.Request{
        Header:     http.Header{},
        RemoteAddr: "127.0.0.1:10086",
    }

    requestBody := map[string]interface{}{
        "content":       "print('Hello World!')",
        "lang":          "python",
        "self_destruct": false,
        // "expire_type":   "",
        // "expiration":    "3",
    }

    mockJSONRequest(context, requestBody, "POST")

    Create(context)

    response := struct {
        Message   string `json:"message"`
        Key       string `json:"key"`
        Namespace string `json:"namespace"`
        Status    uint   `json:"status"`
    }{}

    if err := json.Unmarshal(recorder.Body.Bytes(), &response); err != nil {
        t.Fatal(err)
    }

    if response.Status != 201 || response.Namespace != "nobody" {
        t.Fatalf("%+v", response)
    }
}
