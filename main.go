package main

import (
	"awesomeProject/funcs"
	"database/sql"
	"github.com/gin-gonic/gin"
	_ "github.com/go-sql-driver/mysql"
	"log"
	"net/http"
	"path/filepath"
)

//var todos []TODO

func initDB() (db *sql.DB, err error) {

	// username:password@tcp(host:port)/dbname
	//db, err = sql.Open("mysql", "root:1234@tcp(127.0.0.1:3306)/todoDB")

	db, err = sql.Open("mysql", "user:password@tcp(127.0.0.1:3306)/todoDB")

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

	// 配置CORS (如果需要的话)
	//r.Use(cors.Default())

	// 设置静态文件目录
	r.Static("/static", "./static")

	// 加载HTML文件
	r.LoadHTMLGlob("./template/*") // 加载template文件夹下的所有HTML文件

	r.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "login.html", nil)
	})
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

	//跑进程
	r.Run(":8080")
}
func getRouteHandler(filename string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.HTML(http.StatusOK, filename, nil)
	}
}
