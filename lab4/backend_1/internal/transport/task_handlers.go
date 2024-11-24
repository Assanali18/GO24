package transport

import (
	"backend/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

var tasks = make(map[uint]models.Task)
var nextID uint = 1


func GetTasks(c *gin.Context) {
	tasksList := make([]models.Task, 0, len(tasks))
	for _, task := range tasks {
		tasksList = append(tasksList, task)
	}
	c.JSON(http.StatusOK, tasksList)
}


func CreateTask(c *gin.Context) {
	var newTask models.Task
	if err := c.ShouldBindJSON(&newTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	newTask.ID = nextID
	tasks[nextID] = newTask
	nextID++
	c.JSON(http.StatusCreated, newTask)
}


func UpdateTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	var updatedTask models.Task
	if err := c.ShouldBindJSON(&updatedTask); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	task, exists := tasks[uint(id)]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	task.Title = updatedTask.Title
	task.Description = updatedTask.Description
	task.Status = updatedTask.Status
	task.Category = updatedTask.Category
	tasks[uint(id)] = task

	c.JSON(http.StatusOK, task)
}

func DeleteTask(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	if _, exists := tasks[uint(id)]; !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	delete(tasks, uint(id))
	c.JSON(http.StatusNoContent, gin.H{"message": "Task deleted"})
}

func GetTaskDetail(c *gin.Context) {
	id, err := strconv.ParseUint(c.Param("id"), 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	task, exists := tasks[uint(id)]
	if !exists {
		c.JSON(http.StatusNotFound, gin.H{"error": "Task not found"})
		return
	}

	c.JSON(http.StatusOK, task)
}

func GetStats(c *gin.Context) {
	taskCount := uint(len(tasks)) 
	c.JSON(http.StatusOK, gin.H{
		"user_count": userCount,
		"task_count": taskCount,
	})
}
