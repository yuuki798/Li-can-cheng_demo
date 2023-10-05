package funcs

import (
	"database/sql"
	"github.com/gin-gonic/gin"
	"golang.org/x/crypto/bcrypt"
)

type LoginRequest struct {
	Username string `json:"username"`
	Password string `json:"password"`
}

type LoginResponse struct {
	Status string `json:"status"`
	Token  string `json:"token,omitempty"`
	Error  string `json:"error,omitempty"`
}

func Login(r *gin.Engine, db *sql.DB) {
	r.POST("/login", func(c *gin.Context) {
		var req LoginRequest

		// 从请求中解析用户数据
		if err := c.BindJSON(&req); err != nil {
			c.JSON(400, LoginResponse{Status: "BadRequest", Error: "Invalid JSON format"})
			return
		}

		var hashedPassword string
		err := db.QueryRow("SELECT password FROM users WHERE username = ?", req.Username).Scan(&hashedPassword)

		if err != nil || bcrypt.CompareHashAndPassword([]byte(hashedPassword), []byte(req.Password)) != nil {
			c.JSON(404, LoginResponse{Status: "NotFound", Error: "Username not found or password mismatch"})
			return
		}

		token, err := generateToken(req.Username)
		if err != nil {
			c.JSON(500, LoginResponse{Status: "InternalServerError", Error: "Could not generate token"})
			return
		}
		c.JSON(200, LoginResponse{Status: "OK", Token: token})
	})
}
