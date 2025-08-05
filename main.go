package main

import "chat_app_backend/application"

func main() {
	app := application.Create()
	app.Configure()

	app.Serve()
}
