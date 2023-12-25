package server

import (
	"fmt"
	"net/http"
	"strconv"
	"todo-app/component"
	"todo-app/database"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var secret = []byte("todo-app")

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.HTMLRender = &TemplRender{}
	r.Use(sessions.Sessions("auth", cookie.NewStore(secret)))

	r.GET("/register", s.RegisterPage)
	r.GET("/logout", s.LogoutHandler)
	r.POST("/auth/login", s.LoginHandler)
	r.POST("/auth/register", s.RegisterHandler)
	r.GET("/", s.HomePage)

	todo_routes := r.Group("/todo")
	todo_routes.Use(AuthRequired)
	{
		todo_routes.GET("/", s.TodoQueryHandler)
		todo_routes.POST("/", s.TodoAddHandler)
		todo_routes.DELETE("/:id", s.TodoDeleteHandler)
	}

	return r
}

func AuthRequired(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("username")
	if user == nil {
		c.AbortWithStatus(http.StatusUnauthorized)
		return
	}
	c.Next()
}

func (s *Server) HomePage(c *gin.Context) {
	session := sessions.Default(c)
	user := session.Get("username")
	if user == nil {
		c.HTML(http.StatusOK, "", component.Login())
		return
	}

	c.HTML(http.StatusOK, "", component.Home())
}
func (s *Server) RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register", component.Register())
}

func (s *Server) LogoutHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username")
	if username == nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid session token"})
		return
	}

	session.Delete("username")

	if err := session.Save(); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to save session"})
		return
	}

	c.Redirect(http.StatusSeeOther, "/")
}
func (s *Server) LoginHandler(c *gin.Context) {
	user := database.User{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	}

	is_user_valid, err := user.IsValid(s.db)
	if err != nil {
		fmt.Println(err)
		c.HTML(http.StatusUnauthorized, "", component.Error("Invalid username or password"))
		return
	}
	if !is_user_valid {
		c.HTML(http.StatusUnauthorized, "", component.Error("Invalid username or password"))
		return
	}

	session := sessions.Default(c)
	session.Set("username", user.Username)
	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "", component.Error("Failed to save session"))
		return
	}

    c.Header("HX-Redirect", "/")
}
func (s *Server) RegisterHandler(c *gin.Context) {
	user := database.User{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	}

	err := user.Add(s.db)
	if err != nil {
		fmt.Println(err)
		c.HTML(http.StatusConflict, "", component.Error("Username already exist"))
		return
	}

	session := sessions.Default(c)
	session.Set("username", user.Username)
	if err := session.Save(); err != nil {
		c.HTML(http.StatusInternalServerError, "", component.Error("Failed to save session"))
		return
	}

    c.Header("HX-Redirect", "/")
}

func (s *Server) TodoQueryHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username").(string)

	url_parameter := c.Request.URL.Query()

	todos, err := database.QueryTodo(s.db, username, url_parameter.Get("category"), url_parameter.Get("title"))
	if err != nil {
		fmt.Println(err)
		// c.HTML(http.StatusUnauthorized, "", component.Error("Failed to get todo list"))
        c.Status(http.StatusUnauthorized)
		return
	}

	// c.HTML(http.StatusOK, "", component.TodoList(todos))
    c.JSON(http.StatusOK, todos)
}
func (s *Server) TodoAddHandler(c *gin.Context) {
	session := sessions.Default(c)
	username := session.Get("username").(string)

    is_done, err := strconv.ParseBool(c.PostForm("is_done"))
	todo := database.Todo{
		Username:    username,
		Category:    c.PostForm("category"),
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		IsDone:       is_done,
	}

	created_todo, err := todo.Create(s.db)
	if err != nil {
		fmt.Println(err)
		// c.HTML(http.StatusInternalServerError, "", component.Error("Failed to add new todo"))
		return
	}

	// todos, err := database.QueryTodo(s.db, username, "", "")
	// if err != nil {
	// 	fmt.Println(err)
	// 	c.HTML(http.StatusUnauthorized, "", component.Error("Failed to get new todo list"))
	// 	return
	// }

	// c.HTML(http.StatusOK, "", component.TodoList(todos))
    c.JSON(http.StatusOK, created_todo)
}
func (s *Server) TodoDeleteHandler(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println(err)
		// c.HTML(http.StatusBadRequest, "", component.Error("Invalid todo id"))
        c.Status(http.StatusBadRequest)
		return
	}

	session := sessions.Default(c)
	username := session.Get("username").(string)

	if err := database.DeleteTodo(s.db, id, username); err != nil {
		fmt.Println(err)
		// c.HTML(http.StatusUnauthorized, "", component.Error("User does not own this todo"))
        c.Status(http.StatusUnauthorized)
		return
	}

    c.Status(http.StatusOK)
}
