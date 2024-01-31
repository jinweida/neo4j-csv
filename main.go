package main

import (
	"fmt"
	"neo4j-csv/docs"
	"neo4j-csv/routers"

	"github.com/gin-gonic/gin"
	swaggerfiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

var (
	model  string
	dir    string
	delete bool
)

// 模拟一些私人数据
var secrets = gin.H{
	"foo":    gin.H{"email": "foo@bar.com", "phone": "123433"},
	"austin": gin.H{"email": "austin@example.com", "phone": "666"},
	"lena":   gin.H{"email": "lena@guapa.com", "phone": "523443"},
}

func main() {
	docs.SwaggerInfo.Title = "用户360关系"
	docs.SwaggerInfo.Description = "用户360关系"
	docs.SwaggerInfo.Version = "1.0"
	// docs.SwaggerInfo.Host = "localhost:8080"
	docs.SwaggerInfo.BasePath = "/api"
	gin.SetMode(gin.DebugMode)

	action := &routers.Graph{}
	// 禁止日志的颜色
	// gin.DisableConsoleColor()
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{
			"message": "pong",
		})
	})
	r.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerfiles.Handler))

	api := r.Group("/api", gin.BasicAuth(gin.Accounts{
		"foo":    "bar",
		"austin": "1234",
		"lena":   "hello2",
		"manu":   "4321",
	}))
	api.GET("/target_graph", action.Targetgraph)
	r.Run("0.0.0.0:8080") // 监听并在 0.0.0.0:8080 上启动服务

	fmt.Println("主线程执行完毕")

}
