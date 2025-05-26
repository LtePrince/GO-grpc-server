package service

import (
	"io"
	"mime"
	"os"
	"path/filepath"
	"strings"

	pb "github.com/LtePrince/GO-grpc-server/pkg/api"
)

type SystemServiceServer struct {
	pb.UnimplementedSystemServiceServer
	DataDir string // 你的 data 文件夹路径
}

func NewSystemServiceServer(dataDir string) *SystemServiceServer {
	return &SystemServiceServer{DataDir: dataDir}
}

func (s *SystemServiceServer) SendFile(req *pb.SendFileRequest, stream pb.SystemService_SendFileServer) error {
	// 拼接文件路径，防止目录穿越
	cleanPath := filepath.Clean(req.FilePath)
	filePath := filepath.Join(s.DataDir, cleanPath)

	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()

	// 获取文件信息
	stat, err := file.Stat()
	if err != nil {
		return err
	}

	// 发送文件元数据
	mimeType := mime.TypeByExtension(strings.ToLower(filepath.Ext(filePath)))
	if mimeType == "" {
		mimeType = "application/octet-stream"
	}
	meta := &pb.SendFileResponse{
		Data: &pb.SendFileResponse_Metadata{
			Metadata: &pb.FileMetadata{
				FileName: stat.Name(),
				MimeType: mimeType,
				FileSize: uint64(stat.Size()),
			},
		},
	}
	if err := stream.Send(meta); err != nil {
		return err
	}

	// 发送文件内容
	buf := make([]byte, 32*1024)
	for {
		n, err := file.Read(buf)
		if n > 0 {
			chunk := &pb.SendFileResponse{
				Data: &pb.SendFileResponse_Chunk{
					Chunk: buf[:n],
				},
			}
			if err := stream.Send(chunk); err != nil {
				return err
			}
		}
		if err == io.EOF {
			break
		}
		if err != nil {
			return err
		}
	}
	return nil
}
