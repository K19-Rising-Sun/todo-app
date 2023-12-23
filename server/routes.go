package server

import (
	"net/http"
	"todo-app/component"

	"github.com/gin-gonic/gin"
)

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()
	r.HTMLRender = &TemplRender{}

	r.GET("/", s.HomePage)
	r.GET("/register", s.RegisterPage)
	r.GET("/logout", s.LogoutHandler)
	r.POST("/auth/login", s.LoginHandler)
	r.POST("/auth/register", s.RegisterHandler)
	r.GET("/todo/get/all", s.TodoGetHandler)
	r.POST("/todo/add", s.TodoAddHandler)
	r.GET("/todo/delete/:id", s.TodoDeleteHandler)

	return r
}
func (s *Server) HomePage(c *gin.Context) {
	_, err := c.Cookie("username")
	if err != nil {
		c.HTML(http.StatusOK, "", component.Login())
		return
	}
	c.HTML(http.StatusOK, "", component.Home())
}
func (s *Server) RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register", component.Register())
}
func (s *Server) LogoutHandler(c *gin.Context) {
	c.SetCookie("username", "", -1, "/", "", false, false)
	c.SetCookie("password", "", -1, "/", "", false, false)
	c.Redirect(http.StatusSeeOther, "/")
}
