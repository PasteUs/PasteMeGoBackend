package request

import (
    "bytes"
    "encoding/json"
    "fmt"
    "github.com/gin-gonic/gin"
    "io/ioutil"
    "net/http/httptest"
    "testing"
)

var baseUrl = "/api/v2/"

// Get 根据特定请求 uri，发起 get 请求返回响应
func get(t *testing.T, uri string, router *gin.Engine) []byte {
    req := httptest.NewRequest("GET", uri, nil) // 构造get请求
    w := httptest.NewRecorder()                 // 初始化响应
    router.ServeHTTP(w, req)                    // 调用相应的handler接口
    result := w.Result()                        // 提取响应
    body, err := ioutil.ReadAll(result.Body)    // 读取响应body
    if err != nil {
        t.Fatal(err)
    }
    return body
}

// requestJson 根据特定请求 uri 和参数 param，以 Json 形式传递参数，发起 post 请求返回响应
func requestJson(t *testing.T, method string, uri string, param map[string]interface{}, router *gin.Engine) []byte {
    jsonByte, err := json.Marshal(param) // 将参数转化为 json 比特流
    if err != nil {
        t.Fatal(err)
    }
    req := httptest.NewRequest(method, uri, bytes.NewReader(jsonByte)) // 构造post请求，json数据以请求body的形式传递
    w := httptest.NewRecorder()                                        // 初始化响应
    router.ServeHTTP(w, req)                                           // 调用相应的handler接口
    result := w.Result()                                               // 提取响应
    body, err := ioutil.ReadAll(result.Body)                           // 读取响应body
    if err != nil {
        t.Fatal(err)
    }
    return body
}

func Set(t *testing.T, router *gin.Engine, Key string, Lang string, Content string, Password string) []byte {
    uri := baseUrl + Key
    if Key == "" || Key == "once" {
        params := make(map[string]interface{})
        params["lang"] = Lang
        params["content"] = Content
        params["password"] = Password
        return requestJson(t, "POST", uri, params, router)
    } else {
        params := make(map[string]interface{})
        params["lang"] = Lang
        params["content"] = Content
        params["password"] = Password
        return requestJson(t, "PUT", uri, params, router)
    }
}

func Get(t *testing.T, router *gin.Engine, Key string, Password string) []byte {
    uri := fmt.Sprintf("%s%s,%s?json", baseUrl, Key, Password)
    return get(t, uri, router)
}
