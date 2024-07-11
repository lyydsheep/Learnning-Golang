package main

import (
	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
	"net/http"
	"strings"
)

type employee struct {
	Id         int
	Name       string
	Department string
}

func main() {
	//dsn := "root:405356790@tcp(127.0.0.1:3306)/exercise"
	//db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
	//	DryRun: true,
	//})
	//if err != nil {
	//	panic("failed to connect database")
	//}
	//employees := []*employee{
	//	{Name: "蝙蝠侠", Department: "哥谭"},
	//	{Name: "超人", Department: "大都会"},
	//}
	//db = db.Debug()
	//db.Create(&employees)
	//for _, e := range employees {
	//	fmt.Printf("%+v\n", e)
	//}

	server := gin.Default()

	server.Use(cors.New(cors.Config{
		AllowHeaders:     []string{"Content-Type", "Authorization"},
		AllowCredentials: true,
		AllowOriginFunc: func(origin string) bool {
			if strings.HasPrefix(origin, "http://localhost") {
				return true
			}
			return strings.HasPrefix(origin, "company.com")
		},
	}))
	server.POST("/hello", func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "hello TIC")
	})

	server.Run(":8080")
}
