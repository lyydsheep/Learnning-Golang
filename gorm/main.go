package main

import (
	"fmt"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type employee struct {
	Id         int
	Name       string
	Department string
}

func main() {
	dsn := "root:405356790@tcp(127.0.0.1:3306)/exercise"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{
		DryRun: true,
	})
	if err != nil {
		panic("failed to connect database")
	}
	employees := []*employee{
		{Name: "蝙蝠侠", Department: "哥谭"},
		{Name: "超人", Department: "大都会"},
	}
	db = db.Debug()
	db.Create(&employees)
	for _, e := range employees {
		fmt.Printf("%+v\n", e)
	}
}
