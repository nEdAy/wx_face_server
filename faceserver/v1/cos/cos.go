package cos

import (
	"net/http"
	"github.com/labstack/echo"
)

// CosController Cos可访问接口
type CosController struct {
}

// NewAuthorization 生产鉴权签名
func (pc *CosController) NewAuthorization(c echo.Context) error {

	secretID := "AKIDLbhdR6zt7A2jEZJjZEj3CD4kFhQz7acT"
	secretKey := "ufLr12wfXFm8Tku8HI8wd9WijZEk5FF4"
	host := "face-recognition-1253284991.cos.ap-beijing.myqcloud.com"
	uri := "http://face-recognition-1253284991.cos.ap-beijing.myqcloud.com/faces"

	req, _ := http.NewRequest("PUT", uri, nil)
	req.Header.Add("Host", host)
	req.Header.Add("x-cos-content-sha1", "db8ac1c259eb89d4a131b253bacfca5f319d54f2")
	req.Header.Add("x-cos-stroage-class", "nearline")

	authTime := NewAuthTime(defaultAuthExpire)

	return c.JSON(http.StatusOK, newAuthorization(secretID, secretKey, req, authTime))
}
