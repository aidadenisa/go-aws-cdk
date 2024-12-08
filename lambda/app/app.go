package app

import (
	"lambda-func/api"
	"lambda-func/database"
)
type App struct {
	ApiHandler api.ApiHandler
}

func NewApp() App {
	//Bootstrap everything

	db := database.NewDynamoDBClient()

	// You inject any struct that respects the UserStore interface 
	// So you can swap out your DB by creating another DB client, eg. Postgres
	// As long as it implements the same interface
	// Follows the dependency inversion principle and the interface segregation principle
	apiHandler := api.NewApiHandler(db)

	return App {
		ApiHandler: apiHandler,
	}
}