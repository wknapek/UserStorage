package main

import (
	"UserStorage/queueHandler"
	"UserStorage/user"
	"context"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
	"os"
	"time"
)

var logger = logrus.New()

func main() {

	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	rabbitHandl := queueHandler.NewRabbitHandler(os.Getenv(""), logger)

	mongoURI := os.Getenv("MONGO_URI")
	if mongoURI == "" {
		logger.Error("MONGO_URI env variable is not set")
		return
	}
	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		logger.Error(err)
		panic(err)
	}
	defer func() {
		if err = client.Disconnect(ctx); err != nil {
			panic(err)
		}
	}()

	r := gin.Default()
	usrHandler := user.NewUserHandler(logger, client, rabbitHandl)

	usersGroup := r.Group("/users")
	{
		usersGroup.GET("", usrHandler.GetAllUsers)
		usersGroup.POST("", usrHandler.CreateUser)
		usersGroup.GET("/:id", usrHandler.GetUser)
		usersGroup.PUT("/:id", usrHandler.UpdateUser)
		usersGroup.DELETE("/:id", usrHandler.DeleteUser)

		usersGroup.GET("/:id/files", usrHandler.GetUserFiles)
		usersGroup.POST("/:id/files", usrHandler.AddFileToUser)
		usersGroup.DELETE("/:id/files", usrHandler.DeleteFilesFromUser)
	}

	err = r.Run(":8080")
	if err != nil {
		return
	}
}

func failOnError(err error, msg string) {
	if err != nil {
		logger.Panicf("%s: %s", msg, err)
	}
}
