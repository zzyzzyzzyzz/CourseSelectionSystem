package main

import (
	"CourseSelectionSystem/db"
	"CourseSelectionSystem/routers"
)

func main() {
	db.InitDb()
	db.InitRedis()
	routers.InitRouter()
}
