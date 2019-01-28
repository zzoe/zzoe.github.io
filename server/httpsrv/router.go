package httpsrv

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func router() *gin.Engine {
	r := gin.Default()

	r.Static("/", viper.GetString("http.dist"))

	return r
}
