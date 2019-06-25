/*
@File: server_test.go
@Contact: lucien@lucien.ink
@Licence: (C)Copyright 2019 Lucien Shui

@Modify Time      @Author    @Version    @Description
------------      -------    --------    -----------
2019-06-21 08:37  Lucien     1.0         Init
*/
package server

import (
	"encoding/json"
	"github.com/LucienShui/PasteMeBackend/model"
	"github.com/LucienShui/PasteMeBackend/tests/request"
	"github.com/LucienShui/PasteMeBackend/util/convert"
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

	t.Logf("permanent key: %d", keyP)
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
	body := request.Set(t, router, "asdf", "plain", "Hello", "")

	type JsonResponse struct {
		Key string `json:"key"`
	}

	response := JsonResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	keyT = response.Key

	t.Logf("template key: %s", keyT)
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
	body := request.Set(t, router, "read_once", "plain", "Hello", "")

	type JsonResponse struct {
		Key string `json:"key"`
	}

	response := JsonResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	keyR = response.Key

	t.Logf("read_once key: %s", keyR)
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

	t.Logf("permanent key: %d", keyP)
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
	body := request.Set(t, router, "asdf", "plain", "Hello", "password")

	type JsonResponse struct {
		Key string `json:"key"`
	}

	response := JsonResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	keyT = response.Key

	t.Logf("template key: %s", keyT)
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
	body := request.Set(t, router, "read_once", "plain", "Hello", "password")

	type JsonResponse struct {
		Key string `json:"key"`
	}

	response := JsonResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	keyR = response.Key

	t.Logf("read_once key: %s", keyR)
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
		t.Fatalf("test read_once key: %s failed.", keyR)
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
		t.Fatalf("test read_once key: %s failed.", keyR)
	}

	TestReadOnceGet(t)
	if model.Exist(keyR) {
		t.Fatalf("test read_once key: %s failed.", keyR)
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
		t.Fatalf("test read_once key: %s failed.", keyR)
	}

	TestReadOncePasswordGet(t)
	if model.Exist(keyR) {
		t.Fatalf("test read_once key: %s failed.", keyR)
	}
}
