package dbhandler

import (
	"UserStorage/models"
	"context"
	"go.mongodb.org/mongo-driver/v2/bson"
	"go.mongodb.org/mongo-driver/v2/mongo"
	"go.mongodb.org/mongo-driver/v2/mongo/options"
)

type MongoHandler struct {
	client *mongo.Client
}

func NewMongoHandler(mongoURI string) *MongoHandler {
	client, err := mongo.Connect(options.Client().ApplyURI(mongoURI))
	if err != nil {
		panic(err)
	}
	return &MongoHandler{client}
}

func (m MongoHandler) GetUsers(ctx context.Context) ([]models.User, error) {
	var users []models.User
	cursor, err := m.client.Database("users").Collection("users").Find(ctx, bson.D{})
	if err != nil {
		return nil, err
	}
	if err = cursor.All(ctx, &users); err != nil {
		return nil, err
	}
	return users, nil
}

func (m MongoHandler) GetUser(ctx context.Context, id string) (models.User, error) {
	coll := m.client.Database("users").Collection("users")
	var user models.User
	if err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return models.User{}, err
	}
	return user, nil
}

func (m MongoHandler) CreateUser(ctx context.Context, usr models.User) error {

	_, err := m.client.Database("users").Collection("users").InsertOne(ctx, usr)
	if err != nil {
		return err
	}
	return nil
}

func (m MongoHandler) UpdateUser(ctx context.Context, usr models.User) error {
	_, err := m.client.Database("users").Collection("users").UpdateOne(ctx, bson.M{"_id": usr.Email}, bson.M{"$set": usr})
	if err != nil {
		return err
	}
	return nil
}

func (m MongoHandler) DeleteUser(ctx context.Context, id string) error {
	_, err := m.client.Database("users").Collection("users").DeleteOne(ctx, bson.M{"_id": id})
	if err != nil {
		return err
	}
	return nil
}

func (m MongoHandler) AddFileToUser(ctx context.Context, id string, file models.File) error {
	coll := m.client.Database("users").Collection("users")
	var user models.User
	if err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return err
	}
	user.Files = append(user.Files, file)
	_, err := coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": user})
	if err != nil {
		return err
	}
	return nil
}

func (m MongoHandler) DeleteFilesFromUser(ctx context.Context, id string) error {
	coll := m.client.Database("users").Collection("users")
	var user models.User
	if err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return err
	}
	user.Files = []models.File{}
	_, err := coll.UpdateOne(ctx, bson.M{"_id": id}, bson.M{"$set": user})
	if err != nil {
		return err
	}
	return nil
}

func (m MongoHandler) GetUserFiles(ctx context.Context, id string) ([]models.File, error) {
	coll := m.client.Database("users").Collection("users")
	var user models.User
	if err := coll.FindOne(ctx, bson.M{"_id": id}).Decode(&user); err != nil {
		return nil, err
	}
	return user.Files, nil
}
