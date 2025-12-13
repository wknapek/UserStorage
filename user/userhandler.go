package user

import (
	"UserStorage/models"
	"UserStorage/queueHandler"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"net/http"
)

type UserHandler struct {
	logger *logrus.Logger
	client *mongo.Client
	rabbit *queueHandler.RabbitHandler
}

func NewUserHandler(logger *logrus.Logger, client *mongo.Client, han *queueHandler.RabbitHandler) *UserHandler {
	return &UserHandler{logger, client, han}
}

func (uh *UserHandler) CreateUser(c *gin.Context) {

	var input models.User
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		uh.logger.Error(err)
		return
	}
	one, err := uh.client.Database("users").Collection("users").InsertOne(c.Request.Context(), input)
	if err != nil {
		uh.logger.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	uh.rabbit.Publish(models.Event{EventType: "UserCreated", User: input.Email})
	c.JSON(http.StatusCreated, one)
}

func (uh *UserHandler) GetUser(c *gin.Context) {
	id := c.Param("id")

	coll := uh.client.Database("users").Collection("users")
	var user models.User
	if err := coll.FindOne(c.Request.Context(), bson.M{"_id": id}).Decode(&user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		uh.logger.Error(err)
	}

	c.JSON(http.StatusOK, user)
}

func (uh *UserHandler) UpdateUser(c *gin.Context) {
	id := c.Param("id")
	res, err := uh.client.Database("users").Collection("users").UpdateOne(c.Request.Context(), bson.M{"_id": id}, bson.M{"$set": c.Request.Body})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		uh.logger.Error(err)
	}
	uh.rabbit.Publish(models.Event{EventType: "UserUpdated", User: id})
	c.JSON(http.StatusOK, res)
}

func (uh *UserHandler) DeleteUser(c *gin.Context) {
	id := c.Param("id")
	res, err := uh.client.Database("users").Collection("users").DeleteOne(c.Request.Context(), bson.M{"_id": id})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	uh.rabbit.Publish(models.Event{EventType: "UserDeleted", User: id})
	c.JSON(http.StatusOK, res)
}

func (uh *UserHandler) GetAllUsers(c *gin.Context) {
	var users []models.User
	cursor, err := uh.client.Database("users").Collection("users").Find(c.Request.Context(), bson.D{})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	if err = cursor.All(c.Request.Context(), &users); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
	}
	c.JSON(http.StatusOK, users)
}

func (uh *UserHandler) AddFileToUser(c *gin.Context) {
	id := c.Param("id")
	coll := uh.client.Database("users").Collection("users")
	var user models.User
	if err := coll.FindOne(c.Request.Context(), bson.M{"_id": id}).Decode(&user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		uh.logger.Error(err)
	}
	file := models.File{}
	err := c.ShouldBindJSON(&file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		uh.logger.Error(err)
		return
	}
	user.Files = append(user.Files, file)
	_, err = coll.UpdateOne(c.Request.Context(), bson.M{"_id": id}, bson.M{"$set": user})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		uh.logger.Error(err)
	}
	c.JSON(http.StatusOK, user)
}

func (uh *UserHandler) DeleteFilesFromUser(c *gin.Context) {
	id := c.Param("id")
	coll := uh.client.Database("users").Collection("users")
	var user models.User
	if err := coll.FindOne(c.Request.Context(), bson.M{"_id": id}).Decode(&user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		uh.logger.Error(err)
	}
	file := models.File{}
	err := c.ShouldBindJSON(&file)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		uh.logger.Error(err)
		return
	}
	user.Files = []models.File{}
	_, err = coll.UpdateOne(c.Request.Context(), bson.M{"_id": id}, bson.M{"$set": user})
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		uh.logger.Error(err)
	}
	c.JSON(http.StatusOK, user)
}

func (uh *UserHandler) GetUserFiles(c *gin.Context) {
	id := c.Param("id")
	coll := uh.client.Database("users").Collection("users")
	var user models.User
	if err := coll.FindOne(c.Request.Context(), bson.M{"_id": id}).Decode(&user); err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		uh.logger.Error(err)
	}
	c.JSON(http.StatusOK, user.Files)
}
