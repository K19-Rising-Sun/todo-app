package server

import (
	"fmt"
	"net/http"
	"strconv"
	"todo-app/component"
	"todo-app/database"

	"github.com/gin-gonic/gin"
)

func (s *Server) TodoGetHandler(c *gin.Context) {
	username, err := c.Cookie("username")
	if err != nil {
		fmt.Println(err)
		c.Redirect(http.StatusOK, "/")
		return
	}

	todos, err := database.GetTodos(s.db, username)
	c.HTML(http.StatusOK, "", component.TodoList(todos))
}
func (s *Server) TodoAddHandler(c *gin.Context) {
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
	todo.Create(s.db)
	todos, _ := database.GetTodos(s.db, username)
	c.HTML(http.StatusOK, "", component.TodoList(todos))
}
func (s *Server) TodoDeleteHandler(c *gin.Context) {
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
	if err := database.DeleteTodo(s.db, id); err != nil {
		fmt.Println(err)
		c.JSON(http.StatusForbidden, "")
	}

}
