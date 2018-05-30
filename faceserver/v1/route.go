package v1

import (
	"github.com/labstack/echo"
	"github.com/nEdAy/face-login/faceserver/v1/public"
	"github.com/nEdAy/face-login/faceserver/v1/user"
	"github.com/nEdAy/face-login/faceserver/v1/cos"
)

// Route v1版本的路由文件
func Route(g *echo.Group) {
	/* 公共访问接口 */
	publicC := new(public.PublicController)
	g.POST("/login", publicC.Login)

	/* Cos访问接口 */
	CosC := new(cos.CosController)
	g.GET("/newAuthorization", CosC.NewAuthorization)

	/* 用户管理 */
	userC := new(user.UserController)
	userR := g.Group("/user")
	userR.POST("", userC.AddUser)      // 添加用户
	userR.GET("", userC.UserList)      // 用户列表
	userR.DELETE("", userC.DelUser)    // 删除用户
	userR.DELETE("/all", userC.DelAll) // 清空用户

}
