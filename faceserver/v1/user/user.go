package user

import (
	"fmt"
	"log"
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
	// 如果未添加到数据库，则删除图片
	defer func() {
		log.Println(user.Id)
		// if user.FaceToken == "" || user.Id == 0 {
		// 	os.Remove(picPath)
		// }
	}()

	faceCount, faceToken, err := face.GetFaceCount(prefixCosUrl, fileName)

	if err != nil {
		logger.Errorln(err)
		return c.JSON(http.StatusBadRequest, err.Error())
	}
	if faceCount == 0 {
		return c.JSON(http.StatusBadRequest, "未检测到人脸信息")
	}
	if faceCount > 1 {
		return c.JSON(http.StatusBadRequest, "请保证人脸照片中只包含一个人脸")
	}

	user.FaceToken = faceToken
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
	pathSeparator := string(os.PathSeparator)
	picPath := fmt.Sprintf("%s%spublic%s", common.GetRootDir(), pathSeparator, user.FaceUrl)
	os.Remove(picPath)

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
