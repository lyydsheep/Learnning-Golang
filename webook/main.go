package main

func main() {
	server := InitWebServer() //获取虚拟服务器
	server.Run(":8080")
}
