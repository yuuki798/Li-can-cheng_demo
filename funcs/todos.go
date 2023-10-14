package funcs

import (
	"database/sql"
	"fmt"
	"github.com/gin-gonic/gin"
	"strconv"
)

type TODO struct {
	ID      int    `json:"id"`
	Content string `json:"content"`
	Done    bool   `json:"done"`
}

func Rooting(r *gin.RouterGroup, db *sql.DB) {
	//添加
	addToTodo(r, db)
	//删除
	deleteTodo(r, db)
	//修改
	changeTodo(r, db)
	//获取
	getAllTodo(r, db)
	//真查询
	searchForTodo(r, db)
}

func getAllTodo(r *gin.RouterGroup, db *sql.DB) {
	// 获取所有TODO
	r.GET("/todo", func(c *gin.Context) {
		rows, err := db.Query("SELECT * FROM todos")
		if err != nil {
			fmt.Println("Error:", err) // 输出具体错误到控制台
			c.JSON(500, gin.H{"status": "InternalServerError", "error": err.Error()})
			return
		}
		defer rows.Close()

		var todos []TODO
		for rows.Next() {
			var todo TODO
			if err := rows.Scan(&todo.ID, &todo.Content, &todo.Done); err != nil {
				fmt.Println("Error:", err) // 输出具体错误到控制台
				c.JSON(500, gin.H{"status": "InternalServerError", "error": err.Error()})
				return
			}
			todos = append(todos, todo)
		}

		c.JSON(200, todos)
	})
}

// 定义一个新的结构体只包含 done 字段
type UpdateTODO struct {
	Done bool `json:"done"`
}

func changeTodo(r *gin.RouterGroup, db *sql.DB) {
	// 修改TODO
	r.PUT("/todo/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"status": "BadRequest", "error": "Invalid ID"})
			return
		}

		var update UpdateTODO
		if err := c.BindJSON(&update); err != nil {
			c.JSON(400, gin.H{"status": "BadRequest", "error": "Invalid JSON format"})
			return
		}

		// 更新语句只设置 done 字段
		result, err := db.Exec("UPDATE todos SET done = ? WHERE id = ?", update.Done, id)
		if err != nil {
			fmt.Println("Error:", err)
			c.JSON(500, gin.H{"status": "InternalServerError", "error": err.Error()})
			return
		}

		count, err := result.RowsAffected()
		if err != nil {
			fmt.Println("Error getting rows affected:", err)
			return
		}
		if count == 0 {
			fmt.Println("Warning: No rows were updated.")
			c.JSON(404, gin.H{"status": "NotFound", "error": "TODO not found"})
		}

		c.JSON(200, gin.H{"status": "OK"})
	})
}

func deleteTodo(r *gin.RouterGroup, db *sql.DB) {
	// 删除TODO
	r.DELETE("/todo/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"status": "BadRequest", "error": "Invalid ID"})
			return
		}

		_, err = db.Exec("DELETE FROM todos WHERE id = ?", id)
		if err != nil {
			fmt.Println("Error:", err) // 输出具体错误到控制台
			c.JSON(500, gin.H{"status": "InternalServerError", "error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"status": "OK"})
	})
}

func addToTodo(r *gin.RouterGroup, db *sql.DB) {
	//添加TODO
	r.POST("/todo", func(c *gin.Context) {
		var todo TODO
		// 从请求体中解析用户数据
		if err := c.BindJSON(&todo); err != nil {
			fmt.Println("Error:", err) // 输出具体错误到控制台
			c.JSON(400, gin.H{"status": "BadRequest", "error": err.Error()})
			return
		}

		_, err := db.Exec("INSERT INTO todos(content, done) VALUES(?, ?)", todo.Content, todo.Done)
		if err != nil {
			c.JSON(500, gin.H{"status": "InternalServerError"})
			return
		}
		c.JSON(200, gin.H{"status": "OK"})
	})
}

func searchForTodo(r *gin.RouterGroup, db *sql.DB) {
	r.GET("/todo/search/:query", func(c *gin.Context) {
		query := c.Param("query")
		rows, err := db.Query("SELECT * FROM todos WHERE content LIKE ?", "%"+query+"%")
		if err != nil {
			fmt.Println("Error:", err)
			c.JSON(500, gin.H{"status": "InternalServerError", "error": err.Error()})
			return
		}
		defer rows.Close()

		var todos []TODO
		for rows.Next() {
			var todo TODO
			if err := rows.Scan(&todo.ID, &todo.Content, &todo.Done); err != nil {
				fmt.Println("Error:", err)
				c.JSON(500, gin.H{"status": "InternalServerError", "error": err.Error()})
				return
			}
			todos = append(todos, todo)
		}

		c.JSON(200, todos)
	})
}
