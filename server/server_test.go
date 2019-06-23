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
	"github.com/LucienShui/PasteMeBackend/tests/request"
	"github.com/LucienShui/PasteMeBackend/util"
	"testing"
)

var keyD int
var keyS string

func TestPermanentPost(t *testing.T) {
	body := request.Set(t, router, "", "plain", "<h1>Hello!</h1>", "")

	type JsonResponse struct {
		Key int `json:"key"`
	}

	response := JsonResponse{}
	if err := json.Unmarshal(body, &response); err != nil {
		t.Fatal(err)
	}

	keyD = response.Key

	t.Log(keyD)
}

func TestPermanentGet(t *testing.T) {
	body := request.Get(t, router, util.Uint2string(uint64(keyD)), "")

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

	keyS = response.Key

	t.Log(keyS)
}

func TestTemporaryGet(t *testing.T) {
	body := request.Get(t, router, keyS, "")

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

	keyS = response.Key

	t.Log(keyS)
}

func TestReadOnceGet(t *testing.T) {
	body := request.Get(t, router, keyS, "")

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
