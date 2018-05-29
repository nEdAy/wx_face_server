package faceserver

import (
	"encoding/json"
	"fmt"
	"log"
	"os"
	"time"

	"github.com/google/logger"
	"github.com/labstack/echo"
	"github.com/labstack/echo/middleware"
	"github.com/nEdAy/face-login/internal/common"
	"github.com/nEdAy/face-login/internal/config"
	"github.com/nEdAy/face-login/internal/context"
	"github.com/nEdAy/face-login/internal/db"
)

// FaceServer 程序服务对象
type FaceServer struct {
	e   *echo.Echo
	cfg *config.Config
}

// New 创建服务端对象
func New() *FaceServer {
	// 系统日志显示文件和行号
	log.SetFlags(log.Lshortfile | log.LstdFlags)
	// 初始化配置文件
	cfg, err := config.NewConfig("")
	if err != nil {
		logger.Fatalln("配置文件读取失败:", err)
	}
	js, _ := json.Marshal(cfg)
	log.Println(string(js))
	// 初始化google/logger输出到文件
	err = initLogger("faceserver", cfg.Debug)
	if err != nil {
		logger.Fatalln("日志初始化失败:", err)
	}
	// echo对象
	e := echo.New()
	e.Use(context.InitZContext())
	// 注册中间件
	e.Use(middleware.Logger()) // 根据配置将日志输出到哪里
	e.Use(middleware.Recover())
	e.Use(middleware.Gzip())
	e.Use(middleware.CORSWithConfig(middleware.CORSConfig{
		AllowOrigins: []string{"*"},
		AllowMethods: []string{echo.GET, echo.HEAD, echo.PUT, echo.PATCH, echo.POST, echo.DELETE},
	}))
	e.Static("/", common.GetRootDir()+"/public")

	// 初始化mysql
	err = db.InitDB(cfg.MySQL)
	if err != nil {
		logger.Fatalln("mysql连接错误:", err)
	}

	return &FaceServer{
		e:   e,
		cfg: cfg,
	}
}

// Run 启动服务
func (fs *FaceServer) Run() {
	// 路由
	fs.Route()
	// 启动服务
	address := fmt.Sprintf("%s:%d", fs.cfg.Http.Address, fs.cfg.Http.Port)
	// err := fs.e.Start(address)
	err := fs.e.StartTLS(address, "ssl/www.neday.cn_bundle.cer", "ssl/www.neday.cn.key")
	if err != nil {
		fs.e.Logger.Fatal(err)
	}
	fs.e.Logger.Info()

}

// Stop 停止服务
func (fs *FaceServer) Stop() {

}

// 获取log文件对象
func initLogger(name string, verbose bool) error {
	rootPath := common.GetRootDir()
	if rootPath != "" {
		rootPath = fmt.Sprintf("%s%s", rootPath, string(os.PathSeparator))
	}
	logPath := fmt.Sprintf("%slogs%s%s_%d.log", rootPath, string(os.PathSeparator), name, time.Now().Unix())

	lf, err := os.OpenFile(logPath, os.O_CREATE|os.O_WRONLY|os.O_APPEND, 0660)
	if err != nil {
		return err
	}
	logger.Init(name, verbose, false, lf)

	return nil
}
