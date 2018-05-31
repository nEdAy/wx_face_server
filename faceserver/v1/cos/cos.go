package cos

import (
	"log"
	"net/http"
	"github.com/labstack/echo"
	"google.golang.org/grpc"
	"time"
	"golang.org/x/net/context"
	pb "github.com/nEdAy/face-login/faceserver/v1/cos/wx_cos_auth"
)

const (
	address = "localhost:50051"
)

// CosController Cos可访问接口
type CosController struct {
}

// NewAuthorization 生产鉴权签名
func (pc *CosController) NewAuthorization(c echo.Context) error {

	method := c.FormValue("method")
	if method == "" {
		return c.JSON(http.StatusBadRequest, "method不能为空")
	}

	pathname := c.FormValue("pathname")
	if pathname == "" {
		return c.JSON(http.StatusBadRequest, "pathname不能为空")
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
		return c.JSON(http.StatusBadRequest, "did not connect: " + err.Error())
	}
	defer conn.Close()
	client := pb.NewWXCosAuthClient(conn)

	// Contact the server and print out its response.
	ctx, cancel := context.WithTimeout(context.Background(), time.Second)
	defer cancel()
	r, err := client.GetAuthData(ctx, &pb.GetAuthDataRequest{Method: method,Pathname: pathname})
	if err != nil {
		log.Fatalf("could not AuthData: %v", err)
		return c.JSON(http.StatusBadRequest, "could not AuthData: "+ err.Error())
	}
	log.Printf("AuthData: %s", r.AuthData)

	return c.String(http.StatusOK, r.AuthData)
}
