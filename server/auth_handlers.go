package server

import (
	"fmt"
	"net/http"
	"todo-app/database"

	"github.com/gin-gonic/gin"
)

func (s *Server) LoginHandler(c *gin.Context) {
	user := database.User{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	}
	is_user_valid, err := user.IsValid(s.db)
	if err != nil {
		fmt.Println(err)
		return
	}
	if is_user_valid {
		c.SetCookie("username", user.Username, 0, "/", "", false, false)
		c.SetCookie("password", user.Password, 0, "/", "", false, false)
	}
	c.Redirect(http.StatusSeeOther, "/")
}
func (s *Server) RegisterHandler(c *gin.Context) {
	user := database.User{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	}

    user.Add(s.db)
	c.Redirect(http.StatusSeeOther, "/")
}
