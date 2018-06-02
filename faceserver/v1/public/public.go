package public

import (
	"net/http"
	"github.com/google/logger"
	"github.com/labstack/echo"
	"github.com/nEdAy/face-login/internal/db"
	"github.com/nEdAy/face-login/model"
	"log"
	"github.com/nEdAy/face-login/internal/face"
)

// PublicController 公共可访问接口
type PublicController struct {
}

type LoginUserModel struct {
	Username     string `json:"username"`
	PrefixCosUrl string `json:"prefixCosUrl"`
	FileName     string `json:"fileName"`
}

// Login 登录
func (pc *PublicController) Login(c echo.Context) error {
	loginUserModel := new(LoginUserModel)
	if err := c.Bind(loginUserModel); err != nil {
		return c.JSON(http.StatusBadRequest, "参数格式错误")
	}
	user := new(model.UserModel)
	user.Username = loginUserModel.Username
	if user.Username == "" {
		return c.JSON(http.StatusBadRequest, "用户名不能为空")
	}
	prefixCosUrl := loginUserModel.PrefixCosUrl
	if prefixCosUrl == "" {
		return c.JSON(http.StatusBadRequest, "图片地址不能为空")
	}
	fileName := loginUserModel.FileName
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

	// 查询用户信息
	findUser := new(model.UserModel)
	err := db.DB.Where("username = ?", user.Username).Find(findUser).Limit(1).Error
	if err != nil {
		logger.Errorln(err)
		if err.Error() == "record not found" {
			return c.JSON(http.StatusBadRequest, "用户<"+user.Username+">不存在")
		}
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	isMatchFace, err := face.IsMatchFace(prefixCosUrl, fileName, findUser.FaceToken)
	if err != nil {
		logger.Errorln(err)
		/*		if faceCount == 0 {
					return c.JSON(http.StatusBadRequest, "未检测到人脸信息")
				}
				if faceCount > 1 {
					return c.JSON(http.StatusBadRequest, "请保证人脸照片中只包含一个人脸")
				}*/
		return c.JSON(http.StatusBadRequest, err.Error())
	}

	if !isMatchFace {
		return c.JSON(http.StatusBadRequest, "拍摄照片中未检测到该用户人脸")
	}
	user.Password = ""
	user.FaceToken = ""
	return c.JSON(http.StatusOK, findUser)
}
