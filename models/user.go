package models

type User struct {
	Email    string `json:"email" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Age      int    `json:"age" bson:"age"`
	Files    []File `json:"files" bson:"files"`
}

type File struct {
	Name string `bson:"name"`
}
