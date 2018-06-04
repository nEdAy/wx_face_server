package user

import (
	"fmt"
	"net/http"
	"os"
	"time"

	"github.com/labstack/echo"
	"github.com/nEdAy/face-login/internal/common"
	"github.com/nEdAy/face-login/internal/db"
	"github.com/nEdAy/face-login/model"
	"github.com/google/logger"
	"github.com/nEdAy/face-login/internal/face"
)

// UserController 用户管理，即FaceSet管理
type UserController struct {
}

type AddUserModel struct {
	Username     string `json:"username"`
	Password     string `json:"password"`
	PrefixCosUrl string `json:"prefixCosUrl"`
	FileName     string `json:"fileName"`
}

// AddUser 添加用户
func (uc *UserController) AddUser(c echo.Context) error {
	addUserModel := new(AddUserModel)
	if err := c.Bind(addUserModel); err != nil {
		return c.JSON(http.StatusBadRequest, "参数格式错误")
	}
	user := new(model.UserModel)
	user.Username = addUserModel.Username
	if user.Username == "" {
		return c.JSON(http.StatusBadRequest, "用户名不能为空")
	}
	user.Password = common.UserPwdEncrypt(addUserModel.Password)
	prefixCosUrl := addUserModel.PrefixCosUrl
	if prefixCosUrl == "" {
		return c.JSON(http.StatusBadRequest, "图片地址不能为空")
	}
	fileName := addUserModel.FileName
	if fileName == "" {
		return c.JSON(http.StatusBadRequest, "图片地址不能为空")
	}
	user.FaceUrl = prefixCosUrl + fileName

	user.FaceToken = common.UserPwdEncrypt(user.Username)
	faceCount, err := face.GetFaceCount(prefixCosUrl, fileName, user.FaceToken)

	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if faceCount == -1 {
		return c.JSON(http.StatusBadRequest, "已存在该用户名的人脸信息")
	}
	if faceCount == 0 {
		return c.JSON(http.StatusBadRequest, "未检测到人脸信息")
	}
	if faceCount > 1 {
		return c.JSON(http.StatusBadRequest, "请保证人脸照片中只包含一个人脸")
	}

	user.CreateTime = time.Now().Unix()

	// 查询用户信息
	err = db.DB.Where("username = ?", user.Username).Find(new(model.UserModel)).Limit(1).Error
	if err != nil {
		logger.Errorln(err)
		if err.Error() == "record not found" {
			err = db.DB.Create(user).Error
			if err != nil {
				logger.Errorln(err)
				return c.JSON(http.StatusBadRequest, err.Error())
			}
			return c.JSON(http.StatusOK, user)
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	} else {
		return c.JSON(http.StatusBadRequest, "用户<"+user.Username+">已注册")
	}
}

// UserList 用户列表
func (uc *UserController) UserList(c echo.Context) error {
	list := make([]*model.UserModel, 0)
	err := db.DB.Model(&model.UserModel{}).Order("id desc").Scan(&list).Error
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, list)
}

// DelUser 删除用户
func (uc *UserController) DelUser(c echo.Context) error {
	id := c.FormValue("id")
	// 查询用户信息
	user := new(model.UserModel)
	err := db.DB.Where("id = ?", id).Find(user).Limit(1).Error
	if err != nil {
		logger.Errorln(err)
		if err.Error() == "record not found" {
			return c.JSON(http.StatusBadRequest, "用户信息不存在")
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if user.Id == 0 {
		return c.JSON(http.StatusBadRequest, "用户信息不存在")
	}

	// 删除用户
	err = db.DB.Where("id = ?", id).Delete(model.UserModel{}).Error
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	//TODO: 删除cos上的图片

	// 删除服务器上的图片
	pathSeparator := string(os.PathSeparator)
	picPath := fmt.Sprintf("%s%s%s", "cache/faces", pathSeparator, user.FaceToken)
	os.RemoveAll(picPath)

	return c.JSON(http.StatusOK, "ok")
}

// DelAll 删除全部用户
func (uc *UserController) DelAll(c echo.Context) error {
	// 删除用户
	err := db.DB.Delete(model.UserModel{}).Error
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	// 删除所有图片
	pathSeparator := string(os.PathSeparator)
	picPath := fmt.Sprintf("%s%spublic%sfaces", common.GetRootDir(), pathSeparator, pathSeparator)
	err = os.RemoveAll(picPath)
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	err = os.MkdirAll(picPath, os.ModePerm)
	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	return c.JSON(http.StatusOK, "ok")
}
