package http_client

import (
	"bytes"
	"encoding/json"
	"github.com/PasteUs/PasteMeGoBackend/util"
	"github.com/wonderivan/logger"
	"io/ioutil"
	"net/http"
	"net/url"
)

func Post(IP string, url string, param map[string]interface{}) string {
	if data, err := json.Marshal(param); err != nil {
		logger.Error(util.LoggerInfo(IP, "JSON parse failed: "+err.Error()))
	} else {
		if response, err := http.Post(url, "application/json", bytes.NewReader(data)); err != nil {
			logger.Error(util.LoggerInfo(IP, "Post failed: "+err.Error()))
		} else {
			if responseBody, err := ioutil.ReadAll(response.Body); err != nil {
				logger.Error(util.LoggerInfo(IP, "Read response failed: "+err.Error()))
			} else {
				logger.Info(util.LoggerInfo(IP, "http.Post success"))
				return string(responseBody)
			}
		}
	}

	return "Exception"
}

func Get(IP string, rawUrl string, params map[string]string) string {
	if urlObject, err := url.Parse(rawUrl); err != nil {
		logger.Error(util.LoggerInfo(IP, "url.Parse exception, err = "+err.Error()))
	} else {
		values := url.Values{}
		for k, v := range params {
			values.Set(k, v)
		}

		urlObject.RawQuery = values.Encode()
		urlPath := urlObject.String()
		if response, err := http.Get(urlPath); err != nil {
			logger.Error(util.LoggerInfo(IP, "http.Get exception, err = "+err.Error()))
		} else {
			if body, err := ioutil.ReadAll(response.Body); err != nil {
				logger.Error(util.LoggerInfo(IP, "ioutil.ReadAll exception, err = "+err.Error()))
			} else {
				logger.Info(util.LoggerInfo(IP, "http.Get success"))
				return string(body)
			}
		}
	}
	return "Exception"
}
