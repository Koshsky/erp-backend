package main

import (
	"log"

	"github.com/gin-gonic/gin"
	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"

	_ "github.com/Koshsky/erp-backend/docs/swagger"
)

//	@title			Enterprise Resource Planning
//	@version		1.0
//	@description	For managing the enterprise's universal resources
//	@termsOfService	http://swagger.io/terms/

//	@contact.name	Shmonov Matvey
//	@contact.url	https://t.me/Koshsky
//	@contact.email	shmonov.mv@gmail.com

//	@license.name	Apache 2.0
//	@license.url	http://www.apache.org/licenses/LICENSE-2.0.html

//	@host		localhost:8080
//	@BasePath	/api/v1

//	@securityDefinitions.apikey	ApiKeyAuth
//	@in							header
//	@name						Authorization
//	@description				"Введите JWT токен в формате: Bearer {token}"

//	@externalDocs.description	Документация ERP (заглушка)
//	@externalDocs.url			https://swagger.io/resources/open-api/

func main() {
	router := gin.Default()
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))
	err := router.Run(":8080")
	if err != nil {
		log.Fatal(err)
	}
}
