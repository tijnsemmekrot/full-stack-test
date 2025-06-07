package models

type FirstNameRequest struct {
	FirstName string `json:"first_name"`
}

type Person struct {
	Name string `bson:"name"`
}

type PersonResponse struct {
	ID   string `json:"_id"`
	Name string `json:"name"`
}

type Response struct {
	Message string `json:"message"`
	ID      string `json:"id"`
}
