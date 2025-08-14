package main

import (
	"context"
	"flag"
	"log"
	"time"

	"backend/test/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// hello_client

const (
	defaultName = "风飘儿"
)

var (
	addr = flag.String("addr", "0.0.0.0:8001", "the address to connect to")
	name = flag.String("name", defaultName, "Name to greet")
)

func main() {
	flag.Parse()
	// 连接到server端，此处禁用安全传输
	conn, err := grpc.NewClient(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// 执行RPC调用并打印收到的响应数据
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UploadFile(ctx, &pb.UploadRequest{Url: "https://lingshu-files.oss-cn-hangzhou.aliyuncs.com/chat-files/7100a9f1-f35f-435f-82bc-5787449725f6/测试.txt"})
	if err != nil {
		log.Fatalf("could not greet: %v", err)
	}
	log.Printf("Greeting: %s", r.GetContent())
}
