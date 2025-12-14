package models

type User struct {
	Email    string `json:"email" bson:"_id"`
	Username string `json:"username" bson:"username"`
	Age      int    `json:"age" bson:"age"`
	Password string `json:"-" bson:"password"`
	Files    []File `json:"files" bson:"files"`
}

type File struct {
	Name string `json:"name" bson:"name"`
}
