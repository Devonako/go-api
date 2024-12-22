package main

import (
    "net/http"
    "strconv"

    "github.com/gin-gonic/gin"
    "todo-api/database"
    "todo-api/models"
)

func main() {
    db, err := database.InitDB()
    if err != nil {
        panic("failed to connect database")
    }

    router := gin.Default()

    v1 := router.Group("/api/v1")
    {
        todos := v1.Group("/todos")
        {
            todos.POST("/", createTodo(db))
            todos.GET("/", fetchAllTodos(db))
            todos.GET("/:id", fetchSingleTodo(db))
            todos.PUT("/:id", updateTodo(db))
            todos.DELETE("/:id", deleteTodo(db))
        }
    }

    router.Run()
}

func createTodo(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var input models.Todo
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        todo := models.Todo{Task: input.Task, Completed: input.Completed}
        db.Create(&todo)

        c.JSON(http.StatusCreated, todo)
    }
}

func fetchAllTodos(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var todos []models.Todo
        db.Find(&todos)
        c.JSON(http.StatusOK, todos)
    }
}

func fetchSingleTodo(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var todo models.Todo
        todoID := c.Param("id")

        // Convert the ID to an integer
        id, err := strconv.Atoi(todoID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
            return
        }

        if err := db.First(&todo, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
            return
        }

        c.JSON(http.StatusOK, todo)
    }
}

func updateTodo(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var todo models.Todo
        todoID := c.Param("id")

        // Convert the ID to an integer
        id, err := strconv.Atoi(todoID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
            return
        }

        if err := db.First(&todo, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
            return
        }

        var input models.Todo
        if err := c.ShouldBindJSON(&input); err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
            return
        }

        db.Model(&todo).Updates(input) // Update the fields

        c.JSON(http.StatusOK, todo)
    }
}

func deleteTodo(db *gorm.DB) gin.HandlerFunc {
    return func(c *gin.Context) {
        var todo models.Todo
        todoID := c.Param("id")

        // Convert the ID to an integer
        id, err := strconv.Atoi(todoID)
        if err != nil {
            c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
            return
        }

        if err := db.First(&todo, id).Error; err != nil {
            c.JSON(http.StatusNotFound, gin.H{"error": "Todo not found"})
            return
        }

        db.Delete(&todo)

        c.JSON(http.StatusOK, gin.H{"data": true})
    }
}