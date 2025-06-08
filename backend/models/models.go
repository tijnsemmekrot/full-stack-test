package models

type FirstNameRequest struct {
	Name string `json:"first_name"`
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

type IdDeleteRequest struct {
	Id string `json:"id"`
}

type DeleteResponse struct {
	Message string `json:"message"`
	Id      string `json:"id"`
}
