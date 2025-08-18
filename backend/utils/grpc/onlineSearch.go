package grpc

import (
	"backend/utils/grpc/pb/onlineSearch"
	"context"
	"fmt"
	"os"
	"time"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

// grpc调用
func OnlineSearch(keyWord string) (string, error) {
	// 读取 gRPC 地址，优先使用环境变量 GRPC_ADDR，其次使用默认值
	addr := os.Getenv("GRPC_ONLINESEARCH")
	if addr == "" {
		addr = "0.0.0.0:8003"
	}
	// 连接到server端，此处禁用安全传输
	conn, err := grpc.NewClient(addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		fmt.Println("OnlineSearch连接失败", err)
		return "", err
	}
	defer conn.Close()
	c := onlineSearch.NewGreeterClient(conn)

	// 执行RPC调用并打印收到的响应数据
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.OnlineSearch(ctx, &onlineSearch.SearchRequest{KeyWord: keyWord})
	if err != nil {
		fmt.Println("OnlineSearch解析失败", err)
		return "", err
	}
	return r.GetReply(), err
}
