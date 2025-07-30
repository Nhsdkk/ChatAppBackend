package main

import "chat_app_backend/application"

func main() {
	config := application.Config{
		Url: "localhost:8080",
	}
	app := application.Create(&config)
	app.Configure()

	app.Serve()
}
