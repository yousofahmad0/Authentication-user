package services

import (
	"authentication-user/config"
	"authentication-user/entity"
	"authentication-user/requests"
	"authentication-user/response"
	"errors"
	"fmt"
	"github.com/dmitrymomot/go-jwt/blacklist"
	"github.com/golang-jwt/jwt"
	"gorm.io/gorm"
	"os"
	"strconv"
	"time"
)

func SignUp(request requests.SighUpRequest) {
	code := GenerateCode()
	user := entity.User{Firstname: request.Firstname, Lastname: request.Lastname, Username: request.Username, Email: request.Email, Verification_code: strconv.FormatInt(int64(code), 10),
		Phone_number: request.Phone_number, Nationality: request.Nationality}

	// validation
	if validEmail(request.Email) == false {
		return
	}

	var foundUser entity.User
	err := config.Db.Where("email = ?", request.Email).First(&foundUser).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		config.Db.Create(&user)
		Send(code, request.Email)
		return
	} else { //the user is in the database!!
		// if the user requested to verify his email but he is already in the database , techincal problems for example
		if foundUser.Verified == true { //No need 2a3mello she , za3abne
			return
		} else {
			Send(code, request.Email)
			foundUser.Verification_code = strconv.FormatInt(int64(code), 10)
			config.Db.Save(&foundUser)
			return
		}
	}
}

func VerifyEmail(request requests.VerifyEmailRequest) response.BaseResponse {

	var user entity.User
	if validEmail(request.Email) == false {
		return response.BaseResponse{Code: 300, Message: "Invalid Email Address"}
	}
	err := config.Db.Where("email = ?", request.Email).First(&user).Error
	if errors.Is(err, gorm.ErrRecordNotFound) {
		return response.BaseResponse{Code: 301, Message: "Email not found !! If you don't have an account , please create one"}
	}
	config.Db.Where("email = ?", request.Email).First(&user)
	fmt.Println("found user")
	fmt.Println(user)
	if user.Verified == true {
		return response.BaseResponse{Code: 310, Message: "Account is already verified"}
	}
	if user.Verification_code == request.Verification_code {
		user.Verified = true
	} else {
		return response.BaseResponse{Code: 311, Message: "wrong verification code"}
	}
	//	p, _ := GenerateFromPassword(user.Password)
	//pass 8 didgitd
	user.Password = HashPassword(request.Password)
	config.Db.Save(&user)
	//fmt.Println(user)
	return response.BaseResponse{Code: 312, Message: "Account verified"}

}

func Login(request requests.LoginRequest) (string, response.BaseResponse) {
	var user *entity.User
	//validate request
	if validEmail(request.Email) == false {
		return "", response.BaseResponse{Code: 300, Message: "Invalid Email Address"}
	}
	if validPassword(request.Password) == false {
		return "", response.BaseResponse{Code: 302, Message: "Invalid Password"}
	}
	//find the user
	config.Db.Where("email = ? OR username= ?", request.Email, request.Email).Find(&user)
	//fmt.Println(user)
	if user.ID == 0 {
		return "", response.BaseResponse{Code: 303, Message: "User not found"}
	}
	if VerifyPassword(request.Password, user.Password) == true {
		fmt.Println("Logged in ")

	} else {
		return "", response.BaseResponse{Code: 304, Message: "your credentials do not match ! "}
	}
	//generate token :
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"sub": user.ID,
		"exp": time.Now().Add(time.Minute * 2).Unix(),
	})
	//fmt.Println(token)

	tokenString, err := token.SignedString([]byte(os.Getenv("SECRET_KEY")))
	if err != nil {
		return "", response.BaseResponse{Code: 305, Message: "failed to create token"}
	}
	req := requests.TokenSaveRequest{Token: tokenString, UserId: user.ID, Expires_at: time.Now().Add(time.Minute * 2).Unix()}
	SaveToken(req)
	return tokenString, response.BaseResponse{Code: 306, Message: "Token created successfully"}
}

func Logout(id uint) error {
	var token entity.Token
	error2 := config.Db.Where("user_id = ?", id).First(&token)
	if error2 != nil {
		fmt.Println("error in db")
		return error2.Error
	}
	var list blacklist.Blacklist
	err := list.Add(token.Token)
	if err != nil {
		return err
	}
	return nil
}
