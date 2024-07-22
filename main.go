package main

import (
	"fmt"
	"github.com/Andras5014/webook/internal/repository"
	"github.com/Andras5014/webook/internal/repository/dao"
	"github.com/Andras5014/webook/internal/service"
	"github.com/Andras5014/webook/internal/web"
	"github.com/gin-contrib/cors"
	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/redis"
	"github.com/gin-gonic/gin"
	"gorm.io/driver/mysql"
	_ "gorm.io/driver/mysql"
	"gorm.io/gorm"
	"time"
)

func main() {

	db := initDB()
	u := initUser(db)
	server := initWebServer()
	u.RegisterRouters(server)
	server.Run(":8081")
}

func initUser(db *gorm.DB) *web.UserHandler {
	ud := dao.NewUserDAO(db)
	repo := repository.NewUserRepository(ud)
	svc := service.NewUserService(repo)
	uh := web.NewUserHandler(svc)
	return uh
}
func initWebServer() *gin.Engine {
	server := gin.Default()
	server.Use(cors.New(cors.Config{
		AllowCredentials: true,
		AllowMethods:     []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},

		AllowHeaders: []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
		AllowOriginFunc: func(origin string) bool {
			fmt.Println("origin", origin)
			//if strings.HasPrefix(origin, "http://127.0.0.1") {
			//	return true
			//}
			//return strings.Contains(origin, "andras.icu")
			return true
		},
		MaxAge:        12 * time.Hour,
		ExposeHeaders: []string{"x-jwt-token"},
	}))

	store, err := redis.NewStore(16, "tcp", "127.0.0.1:6379", "", []byte("secret"), []byte("secret"))
	if err != nil {
		panic(err)
	}
	server.Use(sessions.Sessions("mysession", store))
	return server
}
func initDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open("root:root@tcp(127.0.0.1:13306)/webook"))
	if err != nil {
		panic(err)
	}

	err = dao.InitTable(db)
	if err != nil {
		panic(err)
	}
	return db
}
