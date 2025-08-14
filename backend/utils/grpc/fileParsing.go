package grpc

import (
	"context"
	"os"
	"time"

	"backend/test/pb"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// grpc调用
func FileParsing(file_url string) (string, error) {
	// 读取 gRPC 地址，优先使用环境变量 GRPC_ADDR，其次使用默认值
	addr := os.Getenv("GRPC_FILEPARSING")
	if addr == "" {
		addr = "0.0.0.0:8001"
	}
	// 连接到server端，此处禁用安全传输
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		return "", err
	}
	defer conn.Close()
	c := pb.NewGreeterClient(conn)

	// 执行RPC调用并打印收到的响应数据
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := c.UploadFile(ctx, &pb.UploadRequest{Url: file_url})
	if err != nil {
		return "", err
	}
	return r.GetContent(), err
}
