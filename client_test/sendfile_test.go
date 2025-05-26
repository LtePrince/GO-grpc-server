package main

import (
	"bytes"
	"context"
	"io"
	"net"
	"os"
	"testing"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/test/bufconn"

	pb "github.com/LtePrince/GO-grpc-server/pkg/api"
	"github.com/LtePrince/GO-grpc-server/pkg/service"
)

const bufSize = 1024 * 1024

var lis *bufconn.Listener

func bufDialer(context.Context, string) (net.Conn, error) {
	return lis.Dial()
}

func TestSendFile(t *testing.T) {
	// 测试文件路径
	testFile := "../data/100Hz-44.1K-sine_0dB.wav"
	origin, err := os.ReadFile(testFile)
	if err != nil {
		t.Fatalf("failed to read test file: %v", err)
	}

	// 启动内存gRPC服务
	lis = bufconn.Listen(bufSize)
	s := grpc.NewServer()
	srv := service.NewSystemServiceServer("../data")
	pb.RegisterSystemServiceServer(s, srv)
	go s.Serve(lis)
	defer s.Stop()

	conn, err := grpc.NewClient(
		"passthrough://test",
		grpc.WithContextDialer(bufDialer),
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
	if err != nil {
		t.Fatalf("failed to dial bufnet: %v", err)
	}
	defer conn.Close()
	client := pb.NewSystemServiceClient(conn)

	// 调用SendFile
	stream, err := client.SendFile(context.Background(), &pb.SendFileRequest{FilePath: "100Hz-44.1K-sine_0dB.wav"})
	if err != nil {
		t.Fatalf("SendFile failed: %v", err)
	}

	var recv bytes.Buffer
	// 跳过第一个元数据包
	first, err := stream.Recv()
	if err != nil {
		t.Fatalf("failed to receive metadata: %v", err)
	}
	if first.GetMetadata() == nil {
		t.Fatalf("first message should be metadata")
	}

	// 读取所有chunk
	for {
		resp, err := stream.Recv()
		if err == io.EOF {
			break
		}
		if err != nil {
			t.Fatalf("stream.Recv error: %v", err)
		}
		if chunk := resp.GetChunk(); chunk != nil {
			recv.Write(chunk)
		}
	}

	// 断言内容一致
	if !bytes.Equal(origin, recv.Bytes()) {
		t.Errorf("file content mismatch: sent %d bytes, received %d bytes", len(origin), recv.Len())
	}
}
