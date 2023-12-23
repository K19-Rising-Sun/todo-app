package main

import (
	"fmt"
	"net/http"
	"strconv"
	"todo-app/component"
	"todo-app/database"
	"todo-app/service"

	"github.com/gin-gonic/gin"
)

var domain string = ""

func main() {
	db, err := database.Init()
	if err != nil {
		return
	}
	r := gin.Default()
	r.HTMLRender = &TemplRender{}

	r.GET("/", func(c *gin.Context) {
		_, err := c.Cookie("username")
		if err != nil {
			c.HTML(http.StatusOK, "", component.Home(false))
            return
		}
		c.HTML(http.StatusOK, "", component.Home(true))
	})
	r.GET("/register", func(c *gin.Context) {
		c.HTML(http.StatusOK, "register", component.Register())
	})
    r.GET("/logout", func(c *gin.Context) {
        c.SetCookie("username", "", -1, "/", domain, false, false)
        c.SetCookie("password", "", -1, "/", domain, false, false)
		c.Redirect(http.StatusSeeOther, "/")
    })
	r.POST("/auth/register", func(c *gin.Context) {
		user := database.User{
			Username: c.PostForm("username"),
			Password: c.PostForm("password"),
		}

		fmt.Println("user:", user)
		service.AddUser(db, &user)
		c.Redirect(http.StatusSeeOther, "/")
	})
	r.POST("/auth/login", func(c *gin.Context) {
		user := database.User{
			Username: c.PostForm("username"),
			Password: c.PostForm("password"),
		}

		is_user_exist, err := service.IsValidUser(db, &user)
		if err != nil {
			fmt.Println(err)
			return
		}
		if is_user_exist {
			c.SetCookie("username", user.Username, 0, "/", domain, false, false)
			c.SetCookie("password", user.Password, 0, "/", domain, false, false)
		}
		c.Redirect(http.StatusSeeOther, "/")
	})
	r.GET("/todo/get/all", func(c *gin.Context) {
		username, err := c.Cookie("username")
		if err != nil {
			fmt.Println(err)
			c.Redirect(http.StatusOK, "/")
			return
		}

		todos, err := service.GetTodos(db, username)
		c.HTML(http.StatusOK, "", component.TodoList(todos))
	})
	r.POST("/todo/add", func(c *gin.Context) {
		username, err := c.Cookie("username")
		if err != nil {
			fmt.Println(err)
			c.Redirect(http.StatusOK, "/")
			return
		}
		todo := database.Todo{
			Username:    username,
			Category:    c.PostForm("category"),
			Title:       c.PostForm("title"),
			Description: c.PostForm("description"),
			State:       c.PostForm("state"),
		}
		service.CreateTodo(db, &todo)
		todos, _ := service.GetTodos(db, username)
		c.HTML(http.StatusOK, "", component.TodoList(todos))
	})
	r.GET("/todo/delete/:id", func(c *gin.Context) {
		id, err := strconv.Atoi(c.Param("id"))
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusForbidden, "")
		}
		fmt.Println(id)
		_, err = c.Cookie("username")
		if err != nil {
			fmt.Println(err)
			c.JSON(http.StatusForbidden, "")
		}
		if err := service.DeleteTodo(db, id); err != nil {
			fmt.Println(err)
			c.JSON(http.StatusForbidden, "")
		}
	})
	fmt.Println("Listening on localhost:3000")
	r.Run(":3000")
}
