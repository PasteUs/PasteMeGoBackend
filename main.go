package main

import (
	"github.com/PasteUs/PasteMeGoBackend/common/config"
	"github.com/PasteUs/PasteMeGoBackend/router"
)

// @title PasteMe API
// @version 3.4.1
// @description PasteMe Go Backend API
// @termsOfService https://github.com/LucienShui/PasteMe#%E5%85%8D%E8%B4%A3%E5%A3%B0%E6%98%8E

// @contact.name Lucien
// @contact.url https://blog.lucien.ink
// @contact.email lucien@lucien.ink

// @license.name GNU General Public License v3.0
// @license.url https://github.com/PasteUs/PasteMeGoBackend/blob/main/LICENSE

// @host
// @BasePath /api/v3

func main() {
	router.Run(config.Config.Address, config.Config.Port)
}
