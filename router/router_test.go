package router

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"net/http/httptest"
	"strings"
	"testing"
)

func request(t *testing.T, method string, uri string, param map[string]interface{},
	header map[string]string) (result map[string]interface{}) {
	var body io.Reader = nil

	if method == "GET" {
		var rawQueryList []string
		for k, v := range param {
			rawQueryList = append(rawQueryList, fmt.Sprintf("%v=%v", k, v))
		}
		uri = uri + "?" + strings.Join(rawQueryList, "&")
	} else {
		if jsonByte, err := json.Marshal(param); err != nil {
			t.Error(err)
			return
		} else {
			body = bytes.NewReader(jsonByte)
		}
	}

	req := httptest.NewRequest(method, uri, body)

	for k, v := range header {
		req.Header.Set(k, v)
	}

	w := httptest.NewRecorder()
	router.ServeHTTP(w, req)
	if data, err := ioutil.ReadAll(w.Result().Body); err != nil {
		t.Error(err)
		return
	} else {
		if err := json.Unmarshal(data, &result); err != nil {
			t.Error(err)
			return
		}
	}
	return
}

type testCase struct {
	name     string
	method   string
	uri      string
	param    map[string]interface{}
	response map[string]interface{}
	expect   map[string]interface{}
}

// TODO permanent test case

func makeCreateTestCase() (result map[string]*testCase) {
	result = make(map[string]*testCase)
	// create paste
	for _, pasteType := range []string{"temporary"} {
		var (
			method = "POST"
			status = 201
		)

		if pasteType == "permanent" {
			status = 401
		}

		for _, password := range []string{"", "_with_password"} {
			name := method + "_" + pasteType + password
			result[name] = &testCase{
				name, method, "/api/v3/paste/",
				map[string]interface{}{
					"lang":     "plain",
					"content":  "Hello World!",
					"password": password,
				},
				map[string]interface{}{},
				map[string]interface{}{
					"code": status,
				},
			}

			if pasteType == "temporary" {
				result[name].param["self_destruct"] = true
				result[name].param["expire_second"] = 5
				result[name].param["expire_count"] = 1
			}
		}
	}
	return
}

func makeGetTestCase(createCaseList map[string]*testCase) (result map[string]*testCase) {
	result = make(map[string]*testCase)
	// get paste
	for _, pasteType := range []string{"temporary"} {
		method := "GET"
		for _, password := range []string{"", "_with_password"} {
			name := method + "_" + pasteType + password
			previousName := "POST_" + pasteType + password
			result[name] = &testCase{
				name, method, "/api/v3/paste/" + (createCaseList[previousName].response["key"]).(string),
				map[string]interface{}{
					"password": password,
				},
				map[string]interface{}{},
				map[string]interface{}{
					"code":  200,
					"lang":    "plain",
					"content": "Hello World!",
				},
			}
		}
	}
	return
}

func equal(expect interface{}, value interface{}) bool {
	return fmt.Sprintf("%v", expect) == fmt.Sprintf("%v", value)
}

func test(t *testing.T, caseList map[string]*testCase) {
	for name, c := range caseList {
		t.Run(name, func(t *testing.T) {
			c.response = request(t, c.method, c.uri, c.param, map[string]string{"Accept": "application/json"})

			for k, v := range c.expect {
				if !equal(v, c.response[k]) {
					t.Errorf("check field \"%s\" failed, expect %v, got %v", k, v, c.response[k])
				}
			}
		})
	}
}

func Test(t *testing.T) {
	createCaseList := makeCreateTestCase()
	test(t, createCaseList)
	getCaseList := makeGetTestCase(createCaseList)
	test(t, getCaseList)
}

func TestMain(m *testing.M) {
	m.Run()
}
