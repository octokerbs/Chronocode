package main

func main() {
	server, err := NewServer(":8080")
	if err != nil {
		panic(err)
	}

	server.Run()
}
