package server

import (
    "encoding/json"
    "github.com/PasteUs/PasteMeGoBackend/model"
    _ "github.com/PasteUs/PasteMeGoBackend/tests"
    "github.com/PasteUs/PasteMeGoBackend/tests/request"
    "github.com/PasteUs/PasteMeGoBackend/util"
    "testing"
)

var keyP uint64
var keyT, keyR string

func checkGetResponse(t *testing.T, body []byte) {
    type JsonResponse struct {
        Content string `json:"content"`
        Lang    string `json:"lang"`
    }

    response := JsonResponse{}
    if err := json.Unmarshal(body, &response); err != nil {
        t.Error(err)
    }

    if response.Content != "Hello" {
        t.Errorf("content not equal: \"%s\"", response.Content)
    }

    if response.Lang != "plain" {
        t.Errorf("lang not equal: \"%s\"", response.Lang)
    }
}

func TestPermanentPost(t *testing.T) {
    body := request.Set(t, router, "", "plain", "Hello", "")

    type JsonResponse struct {
        Key uint64 `json:"key"`
    }

    response := JsonResponse{}
    if err := json.Unmarshal(body, &response); err != nil {
        t.Error(err)
    }

    keyP = response.Key
}

func TestPermanentGet(t *testing.T) {
    body := request.Get(t, router, util.Uint2string(uint64(keyP)), "")
    checkGetResponse(t, body)
}

func TestTemporaryPost(t *testing.T) {
    body := request.Set(t, router, "example", "plain", "Hello", "")
    type JsonResponse struct {
        Key string `json:"key"`
    }

    response := JsonResponse{}
    if err := json.Unmarshal(body, &response); err != nil {
        t.Error(err)
    }

    keyT = response.Key
}

func TestTemporaryGet(t *testing.T) {
    body := request.Get(t, router, keyT, "")
    checkGetResponse(t, body)
}

func TestReadOncePost(t *testing.T) {
    body := request.Set(t, router, "once", "plain", "Hello", "")

    type JsonResponse struct {
        Key string `json:"key"`
    }

    response := JsonResponse{}
    if err := json.Unmarshal(body, &response); err != nil {
        t.Error(err)
    }

    keyR = response.Key
}

func TestReadOnceGet(t *testing.T) {
    body := request.Get(t, router, keyR, "")
    checkGetResponse(t, body)
}

func TestPermanentPasswordPost(t *testing.T) {
    body := request.Set(t, router, "", "plain", "Hello", "password")

    type JsonResponse struct {
        Key uint64 `json:"key"`
    }

    response := JsonResponse{}
    if err := json.Unmarshal(body, &response); err != nil {
        t.Error(err)
    }

    keyP = response.Key
}

func TestPermanentPasswordGet(t *testing.T) {
    body := request.Get(t, router, util.Uint2string(uint64(keyP)), "password")
    checkGetResponse(t, body)
}

func TestTemporaryPasswordPost(t *testing.T) {
    body := request.Set(t, router, "example", "plain", "Hello", "password")

    type JsonResponse struct {
        Key string `json:"key"`
    }

    response := JsonResponse{}
    if err := json.Unmarshal(body, &response); err != nil {
        t.Error(err)
    }

    keyT = response.Key
}

func TestTemporaryPasswordGet(t *testing.T) {
    body := request.Get(t, router, keyT, "password")
    checkGetResponse(t, body)
}

func TestReadOncePasswordPost(t *testing.T) {
    body := request.Set(t, router, "once", "plain", "Hello", "password")

    type JsonResponse struct {
        Key string `json:"key"`
    }

    response := JsonResponse{}
    if err := json.Unmarshal(body, &response); err != nil {
        t.Error(err)
    }

    keyR = response.Key
}

func TestReadOncePasswordGet(t *testing.T) {
    body := request.Get(t, router, keyR, "password")
    checkGetResponse(t, body)
}

func TestExist(t *testing.T) {
    if model.Exist(keyT) {
        t.Errorf("test temporary key: %s failed.", keyT)
    }

    if model.Exist(keyR) {
        t.Errorf("test once key: %s failed.", keyR)
    }

    TestTemporaryPost(t)
    if !model.Exist(keyT) {
        t.Errorf("test temporary key: %s failed.", keyT)
    }

    TestTemporaryGet(t)
    if model.Exist(keyT) {
        t.Errorf("test temporary key: %s failed.", keyT)
    }

    TestReadOncePost(t)
    if !model.Exist(keyR) {
        t.Errorf("test once key: %s failed.", keyR)
    }

    TestReadOnceGet(t)
    if model.Exist(keyR) {
        t.Errorf("test once key: %s failed.", keyR)
    }

    TestTemporaryPasswordPost(t)
    if !model.Exist(keyT) {
        t.Errorf("test temporary key: %s failed.", keyT)
    }

    TestTemporaryPasswordGet(t)
    if model.Exist(keyT) {
        t.Errorf("test temporary key: %s failed.", keyT)
    }

    TestReadOncePasswordPost(t)
    if !model.Exist(keyR) {
        t.Errorf("test once key: %s failed.", keyR)
    }

    TestReadOncePasswordGet(t)
    if model.Exist(keyR) {
        t.Errorf("test once key: %s failed.", keyR)
    }
}
