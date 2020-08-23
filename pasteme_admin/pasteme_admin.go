package pasteme_admin

import (
	"fmt"
	"github.com/PasteUs/PasteMeGoBackend/config"
	"github.com/PasteUs/PasteMeGoBackend/http_client"
	"github.com/PasteUs/PasteMeGoBackend/util"
	"github.com/wonderivan/logger"
)

var Url string

func init() {
	Url = config.Get().AdminUrl
}

func Classify(IP string, key uint64) {
	realAddress := fmt.Sprintf("%s/api/risk/classify/%d", Url, key)
	response := http_client.Post(IP, realAddress, make(map[string]interface{}))
	logger.Info(util.LoggerInfo(IP, fmt.Sprintf("Get response = %s from %s", response, realAddress)))
}
