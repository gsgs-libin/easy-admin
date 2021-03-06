package oauth

import (
	"net/http"
	"easy-api-template/conf"
	"easy-api-template/model/auth"
	"easy-api-template/serializer"
	"easy-api-template/util"
)

type LoginParams struct {
	UserName string `form:"user_name" json:"user_name" binding:"required"` // 用户名
	Password string `form:"password" json:"password" binding:"required"`   // 密码
}

type RegisterParams struct {
	UserName string `form:"user_name" json:"user_name" binding:"required"` // 用户名
	Password string `form:"password" json:"password" binding:"required"`   // 密码
	Email    string `form:"email" json:"email" binding:"required"`         // 邮箱
}

func (p *LoginParams) Login() serializer.SsopaResponse {
	var SsoPaUserModel auth.SsoPaUsers
	err := conf.Orm.Where("user_name = ?", p.UserName).Find(&SsoPaUserModel).RowsAffected
	if err == 0 {
		return serializer.SsopaResponse{
			Response: serializer.Response{
				Code: http.StatusOK,
				Data: err,
				Msg:  "登录失败，用户未找到，请注册！",
			},
			ResCode: serializer.USERNOTEXISTS,
		}
	}
	validatePassStatus := util.ValidatePasswords(SsoPaUserModel.Password, []byte(p.Password))
	if !validatePassStatus {
		return serializer.SsopaResponse{
			Response: serializer.Response{
				Code: http.StatusOK,
				Data: nil,
				Msg:  "登录失败，请检查用户密码！",
			},
			ResCode: serializer.PASSWORDERROR,
		}
	}
	return serializer.SsopaResponse{
		Response: serializer.Response{
			Code: http.StatusOK,
			Data: SsoPaUserModel.Email,
			Msg:  "登录成功",
		},
		ResCode: serializer.LOGINSUCCESS,
	}
}

func (p *RegisterParams) Register() serializer.SsopaResponse {
	var SsoPaUserModel auth.SsoPaUsers
	err := conf.Orm.Where("user_name = ?", p.UserName).Find(&SsoPaUserModel).RowsAffected
	if err >= 1 {
		return serializer.SsopaResponse{
			Response: serializer.Response{
				Code: http.StatusOK,
				Data: err,
				Msg:  "用户已经存在，请直接登录！",
			},
			ResCode: serializer.USEREXISTS,
		}
	}
	SsoPaUserModel.UserName = p.UserName
	SsoPaUserModel.Password = p.Password
	SsoPaUserModel.Email = p.Email
	createErr := conf.Orm.Create(&SsoPaUserModel).Error
	if createErr != nil {
		return serializer.SsopaResponse{
			Response: serializer.Response{
				Code: http.StatusOK,
				Data: createErr,
				Msg:  "创建用户失败！",
			},
			ResCode: serializer.USERCREATEERROR,
		}
	}
	return serializer.SsopaResponse{
		Response: serializer.Response{
			Code: http.StatusOK,
			Data: nil,
			Msg:  "注册成功，请登录",
		},
		ResCode: serializer.CREATEUSERSUCCESS,
	}
}
