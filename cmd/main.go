package main

import (
	"fmt"

	"example/pkg/zkstore"
)

func main() {
	servers := []string{"0.0.0.0:2181"}
	client, err := zkstore.NewZkClient(servers, "/api", 10)
	if err != nil {
		panic(err)
	}
	defer client.Close()
	node1 := &zkstore.ServerNode{"user", "127.0.0.1", 4000}
	if err := client.Register(node1); err != nil {
		panic(err)
	}
	nodes, err := client.GetNodes("user")
	for _, v := range nodes {
		fmt.Println(v.Name, v.Host, v.Port)
	}
	fmt.Println(err)
}
