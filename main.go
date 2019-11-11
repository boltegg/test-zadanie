package main

func main() {

	err := RunHttpServer(":80")
	if err != nil {
		panic(err)
	}
}