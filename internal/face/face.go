package face

import (
	"log"
	"time"
	"golang.org/x/net/context"
	pb "github.com/nEdAy/face-login/internal/face/face_recognition"
	"google.golang.org/grpc"
)

const (
	address = "localhost:50052"
)

func GetFaceCount(prefixCosUrl string, fileName string) (int32, string, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewFaceRecognitionClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.GetFaceCount(ctx, &pb.GetFaceCountRequest{PrefixCosUrl: prefixCosUrl, FileName: fileName})
	if err != nil {
		log.Fatalf("could not Count: %v", err)
	}
	log.Printf("Count: %d", r.Count)
	return r.Count, r.UnknownFaceEncodings, err
}

func IsMatchFace(prefixCosUrl string, fileName string, knownFaceEncoding string) (bool, error) {
	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	c := pb.NewFaceRecognitionClient(conn)

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	r, err := c.IsMatchFace(ctx, &pb.IsMatchFaceRequest{PrefixCosUrl: prefixCosUrl, FileName: fileName, KnownFaceEncoding: knownFaceEncoding})
	if err != nil {
		log.Fatalf("could not IsMatchFace: %v", err)
	}
	log.Printf("IsMatchFace: %t", r.IsMatchFace)
	return r.IsMatchFace, err
}
