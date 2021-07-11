package server

import (
	"encoding/json"
	"github.com/PasteUs/PasteMeGoBackend/model"
	"github.com/PasteUs/PasteMeGoBackend/tests/request"
	"github.com/PasteUs/PasteMeGoBackend/util/convert"
	"testing"
)

var keyP uint64
var keyT, keyR string

func TestPermanentPost(t *testing.T) {
	body := request.Set(t, router, "", "plain", "Hello", "")

	type JsonResponse struct {
		Key uint64 `json:"key"`
	}

	response := JsonResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	keyP = response.Key
}

func TestPermanentGet(t *testing.T) {
	body := request.Get(t, router, convert.Uint2string(uint64(keyP)), "")

	type JsonResponse struct {
		Content string `json:"content"`
		Lang    string `json:"lang"`
	}

	response := JsonResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	if response.Content != "Hello" {
		t.Fatalf("content not equal: %s", response.Content)
	}
}

func TestTemporaryPost(t *testing.T) {
	body := request.Set(t, router, "example", "plain", "Hello", "")

	type JsonResponse struct {
		Key string `json:"key"`
	}

	response := JsonResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	keyT = response.Key
}

func TestTemporaryGet(t *testing.T) {
	body := request.Get(t, router, keyT, "")

	type JsonResponse struct {
		Content string `json:"content"`
		Lang    string `json:"lang"`
	}

	response := JsonResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	if response.Content != "Hello" {
		t.Fatalf("content not equal: %s", response.Content)
	}
}

func TestReadOncePost(t *testing.T) {
	body := request.Set(t, router, "once", "plain", "Hello", "")

	type JsonResponse struct {
		Key string `json:"key"`
	}

	response := JsonResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	keyR = response.Key
}

func TestReadOnceGet(t *testing.T) {
	body := request.Get(t, router, keyR, "")

	type JsonResponse struct {
		Content string `json:"content"`
		Lang    string `json:"lang"`
	}

	response := JsonResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	if response.Content != "Hello" {
		t.Fatalf("content not equal: %s", response.Content)
	}
}

func TestPermanentPasswordPost(t *testing.T) {
	body := request.Set(t, router, "", "plain", "Hello", "password")

	type JsonResponse struct {
		Key uint64 `json:"key"`
	}

	response := JsonResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	keyP = response.Key
}

func TestPermanentPasswordGet(t *testing.T) {
	body := request.Get(t, router, convert.Uint2string(uint64(keyP)), "password")

	type JsonResponse struct {
		Content string `json:"content"`
		Lang    string `json:"lang"`
	}

	response := JsonResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	if response.Content != "Hello" {
		t.Fatalf("content not equal: %s", response.Content)
	}
}

func TestTemporaryPasswordPost(t *testing.T) {
	body := request.Set(t, router, "example", "plain", "Hello", "password")

	type JsonResponse struct {
		Key string `json:"key"`
	}

	response := JsonResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	keyT = response.Key
}

func TestTemporaryPasswordGet(t *testing.T) {
	body := request.Get(t, router, keyT, "password")

	type JsonResponse struct {
		Content string `json:"content"`
		Lang    string `json:"lang"`
	}

	response := JsonResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	if response.Content != "Hello" {
		t.Fatalf("content not equal: %s", response.Content)
	}
}

func TestReadOncePasswordPost(t *testing.T) {
	body := request.Set(t, router, "once", "plain", "Hello", "password")

	type JsonResponse struct {
		Key string `json:"key"`
	}

	response := JsonResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	keyR = response.Key
}

func TestReadOncePasswordGet(t *testing.T) {
	body := request.Get(t, router, keyR, "password")

	type JsonResponse struct {
		Content string `json:"content"`
		Lang    string `json:"lang"`
	}

	response := JsonResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	if response.Content != "Hello" {
		t.Fatalf("content not equal: %s", response.Content)
	}
}

func TestExist(t *testing.T) {
	if model.Exist(keyT) {
		t.Fatalf("test temporary key: %s failed.", keyT)
	}

	if model.Exist(keyR) {
		t.Fatalf("test once key: %s failed.", keyR)
	}

	TestTemporaryPost(t)
	if !model.Exist(keyT) {
		t.Fatalf("test temporary key: %s failed.", keyT)
	}

	TestTemporaryGet(t)
	if model.Exist(keyT) {
		t.Fatalf("test temporary key: %s failed.", keyT)
	}

	TestReadOncePost(t)
	if !model.Exist(keyR) {
		t.Fatalf("test once key: %s failed.", keyR)
	}

	TestReadOnceGet(t)
	if model.Exist(keyR) {
		t.Fatalf("test once key: %s failed.", keyR)
	}

	TestTemporaryPasswordPost(t)
	if !model.Exist(keyT) {
		t.Fatalf("test temporary key: %s failed.", keyT)
	}

	TestTemporaryPasswordGet(t)
	if model.Exist(keyT) {
		t.Fatalf("test temporary key: %s failed.", keyT)
	}

	TestReadOncePasswordPost(t)
	if !model.Exist(keyR) {
		t.Fatalf("test once key: %s failed.", keyR)
	}

	TestReadOncePasswordGet(t)
	if model.Exist(keyR) {
		t.Fatalf("test once key: %s failed.", keyR)
	}
}
