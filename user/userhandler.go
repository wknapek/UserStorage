package user

import (
	"UserStorage/dbhandler"
	"UserStorage/models"
	"UserStorage/queueHandler"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"net/http"
)

type UserHandler struct {
	logger *logrus.Logger
	dbHan  dbhandler.DBHandler
	rabbit queueHandler.QueueHandler
}

func NewUserHandler(logger *logrus.Logger, client dbhandler.DBHandler, han *queueHandler.RabbitHandler) *UserHandler {
	return &UserHandler{logger, client, han}
}

func (uh *UserHandler) CreateUser(c *gin.Context) {

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		uh.logger.Error(err)
		return
	}
	if input.Age < 18 {
		uh.logger.Error("User age is less than 18")
		c.JSON(http.StatusBadRequest, gin.H{"error": "User age is less than 18"})
		return
	}
	err := uh.dbHan.CreateUser(c.Request.Context(), input)
	if err != nil {
		uh.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	uh.rabbit.Publish(models.Event{EventType: "UserCreated", UserID: input.Email, Age: input.Age, NoFiles: len(input.Files)})
	c.JSON(http.StatusCreated, input)
}

func (uh *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	usr, err := uh.dbHan.GetUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		uh.logger.Error(err)
		return
	}

	c.JSON(http.StatusOK, usr)
}

func (uh *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	updUsr := models.User{}
	err := c.ShouldBindJSON(&updUsr)
	if err != nil {
		uh.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	err = uh.dbHan.UpdateUser(c.Request.Context(), updUsr)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		uh.logger.Error(err)
		return
	}
	uh.rabbit.Publish(models.Event{EventType: "UserUpdated", UserID: id, Age: updUsr.Age, NoFiles: len(updUsr.Files)})
	c.JSON(http.StatusOK, updUsr)
}

func (uh *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	err := uh.dbHan.DeleteUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	uh.rabbit.Publish(models.Event{EventType: "UserDeleted", UserID: id})
	c.JSON(http.StatusOK, gin.H{"message": "User deleted"})
}

func (uh *UserHandler) GetAllUsers(c *gin.Context) {
	err, users := uh.dbHan.GetUsers(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err})
		return
	}
	c.JSON(http.StatusOK, users)
}

func (uh *UserHandler) AddFileToUser(c *gin.Context) {
	id := c.Param("id")
	file := models.File{}
	err := c.ShouldBindJSON(&file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		uh.logger.Error(err)
		return
	}
	err = uh.dbHan.AddFileToUser(c.Request.Context(), id, file)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		uh.logger.Error(err)
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "file added"})
}

func (uh *UserHandler) DeleteFilesFromUser(c *gin.Context) {
	id := c.Param("id")
	err := uh.dbHan.DeleteFilesFromUser(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		uh.logger.Error(err)
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "files deleted"})
}

func (uh *UserHandler) GetUserFiles(c *gin.Context) {
	id := c.Param("id")
	files, err := uh.dbHan.GetUserFiles(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		uh.logger.Error(err)
	}
	c.JSON(http.StatusOK, files)
}
