package server

import (
	"context"
	"database/sql"
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"todo-app/component"
	"todo-app/internal/database"
	"todo-app/view"

	"github.com/gin-contrib/sessions"
	"github.com/gin-contrib/sessions/cookie"
	"github.com/gin-gonic/gin"
)

var secret = []byte("todo-app")

func (s *Server) RegisterRoutes() http.Handler {
	r := gin.Default()

	r.Static("/static", "./static")

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
		todo_routes.PUT("/", s.TodoPutHandler)
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
		c.HTML(http.StatusOK, "", view.Login())
		return
	}

	c.HTML(http.StatusOK, "", view.Home(user.(string)))
}
func (s *Server) RegisterPage(c *gin.Context) {
	c.HTML(http.StatusOK, "register", view.Register())
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
	ctx := context.Background()

	user, err := database.New(s.db).GetUser(ctx, c.PostForm("username"))
	if err != nil {
		fmt.Println(err)
		c.HTML(http.StatusUnauthorized, "", component.Error("Invalid username"))
		return
	}

	if user.Password != c.PostForm("password") {
		c.HTML(http.StatusUnauthorized, "", component.Error("Invalid password"))
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
	ctx := context.Background()

	user, err := database.New(s.db).CreateUser(ctx, database.CreateUserParams{
		Username: c.PostForm("username"),
		Password: c.PostForm("password"),
	})
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
	ctx := context.Background()

	session := sessions.Default(c)
	username := session.Get("username").(string)

	url_parameter := c.Request.URL.Query()
	category := url_parameter.Get("category")
	title := url_parameter.Get("title")

	var todos []database.Todo
	var err error

	if category != "" && title != "" {
		todos, err = database.New(s.db).SearchTodos(ctx, database.SearchTodosParams{
			Username: username,
			Query: sql.NullString{
				String: fmt.Sprintf("category: %s and title: %s", category, title),
				Valid:  true,
			},
		})
	} else if category != "" {
		todos, err = database.New(s.db).SearchTodos(ctx, database.SearchTodosParams{
			Username: username,
			Query: sql.NullString{
				String: fmt.Sprintf("category: %s", category),
				Valid:  true,
			},
		})
	} else if title != "" {
		todos, err = database.New(s.db).SearchTodos(ctx, database.SearchTodosParams{
			Username: username,
			Query: sql.NullString{
				String: fmt.Sprintf("title: %s", title),
				Valid:  true,
			},
		})
	} else {
		todos, err = database.New(s.db).GetTodos(ctx, username)
	}

	if err != nil {
		fmt.Println(err)
		// c.HTML(http.StatusUnauthorized, "", component.Error("Failed to get todo list"))
		c.Status(http.StatusUnauthorized)
		return
	}

	// c.HTML(http.StatusOK, "", component.TodoList(todos))
    fmt.Println(todos)
	c.JSON(http.StatusOK, todos)
}
func (s *Server) TodoAddHandler(c *gin.Context) {
	ctx := context.Background()

	session := sessions.Default(c)
	username := session.Get("username").(string)

	is_done, err := strconv.Atoi(c.PostForm("is_done"))

	created_todo, err := database.New(s.db).CreateTodo(ctx, database.CreateTodoParams{
		Username:    username,
		Category:    c.PostForm("category"),
		Title:       c.PostForm("title"),
		Description: c.PostForm("description"),
		IsDone:      int64(is_done),
	})
	if err != nil {
		fmt.Println(err)
		// c.HTML(http.StatusInternalServerError, "", component.Error("Failed to add new todo"))
		return
	}

	// c.HTML(http.StatusOK, "", component.TodoList(todos))
	c.JSON(http.StatusOK, created_todo)
}

func (s *Server) TodoPutHandler(c *gin.Context) {
	ctx := context.Background()

	session := sessions.Default(c)
	username := session.Get("username").(string)

	id, err := strconv.Atoi(c.PostForm("id"))
	if err != nil {
		c.Status(http.StatusBadRequest)
		return
	}

	category := sql.NullString{}
	if c.PostForm("category") != "" {
		category = sql.NullString{
			String: c.PostForm("category"),
			Valid:  true,
		}
	}
	title := sql.NullString{}
	if c.PostForm("title") != "" {
		title = sql.NullString{
			String: c.PostForm("title"),
			Valid:  true,
		}
	}
	description := sql.NullString{}
	if c.PostForm("description") != "" {
		description = sql.NullString{
			String: c.PostForm("description"),
			Valid:  true,
		}
	}
	is_done_raw, err := strconv.ParseBool(c.PostForm("is_done"))
	is_done := sql.NullInt64{}
	if is_done_raw == true {
		is_done = sql.NullInt64{
			Int64: int64(1),
			Valid: true,
		}
	} else if err == nil {
		is_done = sql.NullInt64{
			Int64: int64(0),
			Valid: true,
		}
	}

	updated_todo, err := database.New(s.db).UpdateTodo(ctx, database.UpdateTodoParams{
		ID:          int64(id),
		Username:    sql.NullString{String: username, Valid: true},
		Category:    category,
		Title:       title,
		Description: description,
		IsDone:      is_done,
	})
	if err != nil {
		fmt.Println(err)
		c.Status(http.StatusInternalServerError)
		// c.HTML(http.StatusInternalServerError, "", component.Error("Failed to add new todo"))
		return
	}

	// c.HTML(http.StatusOK, "", component.TodoList(todos))
	c.JSON(http.StatusOK, updated_todo)
}
func (s *Server) TodoDeleteHandler(c *gin.Context) {
	ctx := context.Background()

	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		fmt.Println(err)
		// c.HTML(http.StatusBadRequest, "", component.Error("Invalid todo id"))
		c.Status(http.StatusBadRequest)
		return
	}

	session := sessions.Default(c)
	username := session.Get("username").(string)

	if err := database.New(s.db).DeleteTodo(ctx, database.DeleteTodoParams{
		ID:       int64(id),
		Username: username,
	}); err != nil {
		fmt.Println(err)
		// c.HTML(http.StatusUnauthorized, "", component.Error("User does not own this todo"))
		c.Status(http.StatusUnauthorized)
		return
	}

	c.Status(http.StatusOK)
}
