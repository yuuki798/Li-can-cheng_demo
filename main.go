package main

import (
	"awesomeProject/funcs"
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"path/filepath"
)

func initDB() (db *sql.DB, err error) {

	//打开并记录数据库
	//db, err = sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/todoDB")
	db, err = sql.Open("mysql", "root:1234@tcp(db:3306)/tododb")

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
	r.Static("/static", "./static")

	// 加载template文件夹下的所有HTML文件
	r.LoadHTMLGlob("./template/*")

	//注册路由——登录页面
	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})

	//遍历所有的htmlFile并注册路由
	htmlFiles, _ := filepath.Glob("./template/*.html")
	for _, file := range htmlFiles {
		route := "/" + filepath.Base(file)
		r.GET(route, getRouteHandler(filepath.Base(file)))
	}

	//不需要身份验证的路由
	funcs.Register(r, db)
	funcs.Login(r, db)

	//需要身份验证的路由分组
	authorized := r.Group("/")
	authorized.Use(funcs.AuthMiddleware())
	{
		// 在此分组中，注册需要身份验证的路由
		funcs.Rooting(authorized, db)
	}

	fmt.Printf("Starting server at http://127.0.0.1:8080\n")
	//跑进程
	r.Run(":8080")
}
func getRouteHandler(filename string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, filename, nil)
	}
}
