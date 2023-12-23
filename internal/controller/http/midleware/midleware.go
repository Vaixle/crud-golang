package midleware

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func BasicAuth() gin.HandlerFunc {
	return gin.BasicAuth(gin.Accounts{viper.GetString("auth.basic.username"): viper.GetString("auth.basic.password")})
}
