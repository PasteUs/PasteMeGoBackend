package meta

import "github.com/wonderivan/logger"

var Version = "3.3.0"
var ValidConfigVersion = []string{"3.3.0", ""}

func init() {
	logger.Info("PasteMe Go Backend Version \"%s\"", Version)
}
