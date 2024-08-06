package main

import (
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

type TodoTask struct {
	ID        int    `json:"id"`
	Title     string `json:"title"`
	Completed bool   `json:"completed"`
}

type TodoApp struct {
	Tasks []TodoTask `json:"tasks"`

	nextId int
}

var globalApp = CreateTodoApp()

func CreateTodoApp() *TodoApp {
	return &TodoApp{
		Tasks:  []TodoTask{},
		nextId: 1,
	}
}

func (app *TodoApp) GetTaskOne(id int) *TodoTask {
	for _, task := range app.Tasks {
		if task.ID == id {
			return &task
		}
	}

	/* Return nil if the task is not found */
	return nil
}

func (app *TodoApp) GetTaskAll() []TodoTask {
	return app.Tasks
}

func (app *TodoApp) CreateTask(name string) TodoTask {
	newTask := TodoTask{
		ID:        app.nextId,
		Title:     name,
		Completed: false,
	}
	app.nextId++

	app.Tasks = append(app.Tasks, newTask)
	app.nextId++

	return newTask
}

func (app *TodoApp) UpdateTask(id int, name string, completed bool) *TodoTask {
	for i, task := range app.Tasks {
		if task.ID == id {
			app.Tasks[i].Title = name
			app.Tasks[i].Completed = completed
			return &app.Tasks[i]
		}
	}

	/* Return nil if the task is not found */
	return nil
}

func (app *TodoApp) DeleteTask(id int) error {
	for i, task := range app.Tasks {
		if task.ID == id {
			app.Tasks = append(app.Tasks[:i], app.Tasks[i+1:]...)
			return nil
		}
	}

	return fmt.Errorf("task with id %d not found", id)
}

func TodoTaskGetOne(c *gin.Context) {
	id := c.Params.ByName("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "ID is required"})
		return
	}

	idInt, err := strconv.Atoi(id)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Id must be an integer"})
		return
	}

	task := globalApp.GetTaskOne(idInt)
	if task == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func TodoTaskGetAll(c *gin.Context) {
	// TODO: Implement the get all tasks endpoint
	c.Status(http.StatusNotImplemented)
}

func TodoTaskCreate(c *gin.Context) {
	type RequestBody struct {
		Name string `json:"name"`
	}
	var body RequestBody
	err := c.ShouldBindBodyWithJSON(&body)

	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTask := globalApp.CreateTask(body.Name)
	c.JSON(http.StatusCreated, newTask)
}

func TodoTaskUpdate(c *gin.Context) {
	// TODO: Implement the update task endpoint
	c.Status(http.StatusNotImplemented)
}

func TodoTaskDelete(c *gin.Context) {
	// TODO: Implement the delete task endpoint
	c.Status(http.StatusNotImplemented)
}

func CorsMiddleware() gin.HandlerFunc {
	return cors.New(cors.Config{
		AllowAllOrigins: true,
		AllowMethods:    []string{"GET", "POST", "PUT", "DELETE"},
		AllowHeaders:    []string{"Origin", "Content-Length", "Content-Type", "Authorization"},
	})
}

func main() {
	/* Create a gin engine instance */
	engine := gin.Default()
	engine.Use(CorsMiddleware())

	engine.POST("/tasks", TodoTaskCreate)
	engine.GET("/tasks/:id", TodoTaskGetOne)

	// TODO: Add the missing endpoints for the update and delete operations

	engine.Run("0.0.0.0:8080")
}
