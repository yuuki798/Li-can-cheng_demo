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

func Rooting(r *gin.Engine, db *sql.DB) {
	//添加
	add_to_TODO(r, db)
	//删除
	delete_TODO(r, db)
	//修改
	change_TODO(r, db)
	//获取
	get_all_TODO(r, db)
	//查询
	query_for_TODO(r, db)
}

func query_for_TODO(r *gin.Engine, db *sql.DB) {
	// 查询单个TODO
	r.GET("/todo/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"status": "BadRequest", "error": "Invalid ID"})
			return
		}

		var todo TODO
		err = db.QueryRow("SELECT * FROM todos WHERE id = ?", id).Scan(&todo.ID, &todo.Content, &todo.Done)
		if err != nil {
			if err == sql.ErrNoRows {
				c.JSON(404, gin.H{"status": "NotFound", "error": "TODO not found"})
			} else {
				fmt.Println("Error:", err) // 输出具体错误到控制台
				c.JSON(500, gin.H{"status": "InternalServerError", "error": err.Error()})
			}
			return
		}

		c.JSON(200, todo)
	})
}

func get_all_TODO(r *gin.Engine, db *sql.DB) {
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

func change_TODO(r *gin.Engine, db *sql.DB) {
	// 修改TODO
	r.PUT("/todo/:id", func(c *gin.Context) {
		idStr := c.Param("id")
		id, err := strconv.Atoi(idStr)
		if err != nil {
			c.JSON(400, gin.H{"status": "BadRequest", "error": "Invalid ID"})
			return
		}

		var todo TODO
		if err := c.BindJSON(&todo); err != nil {
			c.JSON(400, gin.H{"status": "BadRequest", "error": "Invalid JSON format"})
			return
		}

		_, err = db.Exec("UPDATE todos SET content = ?, done = ? WHERE id = ?", todo.Content, todo.Done, id)
		if err != nil {
			fmt.Println("Error:", err) // 输出具体错误到控制台
			c.JSON(500, gin.H{"status": "InternalServerError", "error": err.Error()})
			return
		}

		c.JSON(200, gin.H{"status": "OK"})
	})
}

func delete_TODO(r *gin.Engine, db *sql.DB) {
	//删除TODO

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

func add_to_TODO(r *gin.Engine, db *sql.DB) {
	//添加TODO
	r.POST("/todo", func(c *gin.Context) {
		var todo TODO

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
