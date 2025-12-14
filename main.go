package main

import (
	"UserStorage/dbhandler"
	"UserStorage/queueHandler"
	"UserStorage/secutiry"
	"UserStorage/user"
	"flag"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"os"
)

var logger = logrus.New()

func main() {

	mongoURI := flag.String("mongo-uri", "", "mongo uri")
	rabbitURI := flag.String("rabbit-uri", "", "rabbit uri")
	secret := flag.String("secret", "", "secret for jwt")
	flag.Parse()
	if *mongoURI == "" {
		logger.Error("mongo-uri is not set")
		return
	}

	if *rabbitURI == "" {
		logger.Error("rabbit-uri is not set")
		return
	}

	logger.SetOutput(os.Stdout)
	logger.SetFormatter(&logrus.JSONFormatter{})

	rabbitHandl := queueHandler.NewRabbitHandler(*rabbitURI, logger)
	auth := secutiry.NewAuthObj([]byte(*secret))

	r := gin.Default()
	dbHan := dbhandler.NewMongoHandler(*mongoURI)
	usrHandler := user.NewUserHandler(logger, dbHan, rabbitHandl, auth)

	usersGroup := r.Group("/users")
	usersGroup.Use(auth.Auth())
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

	err := r.Run(":8080")
	if err != nil {
		return
	}
}
