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
	"github.com/LucienShui/PasteMeBackend/util"
	"testing"
)

var keyP uint64
var keyT, keyR string

func TestPermanentPost(t *testing.T) {
	body := request.Set(t, router, "", "plain", "<h1>Hello!</h1>", "")

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
	body := request.Get(t, router, util.Uint2string(uint64(keyP)), "")

	type JsonResponse struct {
		Content string `json:"content"`
		Lang    string `json:"lang"`
	}

	response := JsonResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	content, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(content))
}

func TestTemporaryPost(t *testing.T) {
	body := request.Set(t, router, "asdf", "plain", "<h1>Hello!</h1>", "")

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

	content, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(content))
}

func TestReadOncePost(t *testing.T) {
	body := request.Set(t, router, "read_once", "plain", "<h1>Hello!</h1>", "")

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

	content, err := json.Marshal(response)
	if err != nil {
		t.Fatal(err)
	}

	t.Log(string(content))
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
}
