package main




func main() {
	httpServer := NewHttpServer(":8080")
	httpServer.Run()
}

