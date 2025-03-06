package api

import (
	"encoding/json"
	"fmt"
	"lambda-func/types"
	"net/http"

	"github.com/aws/aws-lambda-go/events"
)

type UserStore interface {
	DoesUserExist(username string) (bool, error)
	InsertUser(user types.User) error
	GetUser(username string) (types.User, error)
}

type ApiHandler struct {
	dbStrore UserStore
}

func NewApiHandler(dbStore UserStore) ApiHandler {
	return ApiHandler{
		dbStrore: dbStore,
	}
}

// When we have only 1 normal handler, without proxy requests
// func (api ApiHandler) RegisterUserHandler(event types.RegisterUser) error {
// 	if event.Username == "" || event.Password == "" {
// 		return fmt.Errorf("request has empty params")
// 	}

// 	userExists, err := api.dbStrore.DoesUserExist(event.Username) 
// 	if err != nil {
// 		return fmt.Errorf("there was an error while checking if the user exists: %w", err)
// 	}
// 	if userExists {
// 		return fmt.Errorf("a user with this username already exists")
// 	}

// 	err = api.dbStrore.InsertUser(event)
// 	if err != nil {
// 		return fmt.Errorf("there was an error inserting if the user exists: %w", err)
// 	}

// 	return nil
// }
// With proxy requests
func (api ApiHandler) RegisterUserHandler(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	var registerUser types.RegisterUser

	// This uses the json tags in the struct so the JSON library knows how to unmarshall it
	err := json.Unmarshal([]byte(request.Body), &registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	if registerUser.Username == "" || registerUser.Password == "" {
		return events.APIGatewayProxyResponse{
			Body: "Invalid Request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	userExists, err := api.dbStrore.DoesUserExist(registerUser.Username) 
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}
	if userExists {
		return events.APIGatewayProxyResponse{
			Body: "User Already Exists",
			StatusCode: http.StatusConflict,
		}, nil
	}

	user, err := types.NewUser(registerUser)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal server error", 
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("Could not create new user: %w", err)
	}

	err = api.dbStrore.InsertUser(user)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal server error",
			StatusCode: http.StatusInternalServerError,
		}, fmt.Errorf("error inserting user: %w", err)
	}

	return events.APIGatewayProxyResponse{
		Body: "Successfully registered user", 
		StatusCode: http.StatusOK,
	}, nil
}

func (api ApiHandler) LoginUser(request events.APIGatewayProxyRequest) (events.APIGatewayProxyResponse, error) {
	type LoginRequest struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}

	var loginRequest LoginRequest

	err := json.Unmarshal([]byte(request.Body), &loginRequest)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Invalid request",
			StatusCode: http.StatusBadRequest,
		}, err
	}

	user, err := api.dbStrore.GetUser(loginRequest.Username)
	if err != nil {
		return events.APIGatewayProxyResponse{
			Body: "Internal Server Error",
			StatusCode: http.StatusInternalServerError,
		}, err
	}

	// Validate pass
	if !types.ValidatePassword(user.PasswordHash, loginRequest.Password) {
		return events.APIGatewayProxyResponse{
			Body: "Invalid user Credentials",
			StatusCode: http.StatusUnauthorized,
		}, nil
	}

	accessToken := types.CreateToken(user)
	successMsg := fmt.Sprintf(`{"access_token": %s}`, accessToken)

	return events.APIGatewayProxyResponse{
		Body: successMsg,
		StatusCode: http.StatusOK,
	}, nil
}

