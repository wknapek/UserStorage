package dbhandler

import (
	"UserStorage/models"
	"context"
)

type DBHandler interface {
	GetUsers(ctx context.Context) ([]models.User, error)
	GetUser(ctx context.Context, id string) (models.User, error)
	CreateUser(ctx context.Context, usr models.User) error
	UpdateUser(ctx context.Context, usr models.User) error
	DeleteUser(ctx context.Context, id string) error
	AddFileToUser(ctx context.Context, id string, file models.File) error
	DeleteFilesFromUser(ctx context.Context, id string) error
	GetUserFiles(ctx context.Context, id string) ([]models.File, error)
}
