package main

func main() {

	currentApp := GenerateApp()

	currentApp.Listen("8080")

}
