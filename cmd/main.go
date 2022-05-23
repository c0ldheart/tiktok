package main

import (
	"github.com/gin-gonic/gin"
	"tikapp/common/config"
	"tikapp/common/db"
	"tikapp/common/log"
	"tikapp/common/oss"
)

func init() {
	config.ReadCfg()
	config.Init()
	log.Init()
	db.Init()
	oss.Init()
}

func main() {

	r := gin.Default()

	handle(r)

	r.Run(config.AppCfg.Host + ":" + config.AppCfg.Port)
}
