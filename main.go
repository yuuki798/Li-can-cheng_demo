package main

import (
	"awesomeProject/funcs"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
)

//var todos []TODO

func initDB() (db *sql.DB, err error) {

	// username:password@tcp(host:port)/dbname
	db, err = sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/todoDB")
	if err != nil {
		log.Fatalf("Error opening database: %q", err)
	}
	return db, err
}
func main() {
	var db *sql.DB
	//初始化数据库
	db, err := initDB()
	if err != nil {
		log.Fatalf("Error initializing database: %q", err)
	}

	//进程结束后确保关闭
	defer db.Close()
	//调用gin框架
	r := gin.Default()

	// 设置静态文件目录
	r.Static("/static", "./frontend/js") // 这里假设您的JavaScript文件在 /frontend/js

	// 加载HTML文件
	r.LoadHTMLFiles("./frontend/index.html") // 路径应根据您的实际情况进行调整

	r.GET("/", func(c *gin.Context) {
		c.HTML(200, "index.html", nil)
	})
	//路由实现
	funcs.Register(r, db)
	funcs.Login(r, db)

	r.Use(funcs.AuthMiddleware())
	funcs.Rooting(r, db)
	//跑进程
	r.Run(":8080")

}
