package main

import (
	"deficonnect/defipayapi/handlers"

	"github.com/aws/aws-lambda-go/lambda"
	"github.com/joho/godotenv"
)

func main() {
	err := godotenv.Load()
	if err != nil {
		println("error loading .env file %v", err)
	}
	app, err := handlers.InitSlsApp()
	if err != nil {
		panic(err)
	}
	lambda.Start(app.CheckTransactionStatus)
}
