package usecase

import (
	"context"
	"crypto/rand"
	"encoding/base64"
	"encoding/json"
	"errors"
	"learn_oauth/app/user/repository"
	"learn_oauth/domain"
	"learn_oauth/util"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/oauth2"
)

type IUserUsecase interface {
	GenerateStateOauthCookie(c *gin.Context) string
	GetUserDataFromGoogle(code string, googleOauthConfig *oauth2.Config) (map[string]any, error)
	GenerateAndSetJWT(c *gin.Context, user domain.User) interface{}
}

type Userusecase struct {
	UserRepository repository.IUserRepository
}

func NewUserUsecase(userRepository repository.IUserRepository) IUserUsecase {
	return &Userusecase{
		UserRepository: userRepository,
	}
}

func (useCase *Userusecase) GenerateStateOauthCookie(c *gin.Context) string {
	b := make([]byte, 16)
	rand.Read(b)
	state := base64.URLEncoding.EncodeToString(b)
	c.SetCookie("oauthstate", state, 3600*24, "", "", false, true)

	return state
}

func (usecase *Userusecase) GetUserDataFromGoogle(code string, googleOauthConfig *oauth2.Config) (map[string]any, error) {
	// Use code to get token and get user info from Google.
	const oauthGoogleUrlAPI = "https://www.googleapis.com/oauth2/v2/userinfo?access_token="

	token, err := googleOauthConfig.Exchange(context.Background(), code)
	if err != nil {
		return nil, errors.New("code exchange wrong: " + err.Error())
	}

	response, err := http.Get(oauthGoogleUrlAPI + token.AccessToken)
	if err != nil {
		return nil, errors.New("failed getting user info: " + err.Error())
	}
	defer response.Body.Close()

	var contents map[string]any
	err = json.NewDecoder(response.Body).Decode(&contents)
	if err != nil {
		return nil, errors.New("failed read response: " + err.Error())
	}

	return contents, nil
}

func (usecase *Userusecase) GenerateAndSetJWT(c *gin.Context, user domain.User) interface{} {
	user.ID = uuid.New()
	rowsAffected, err := usecase.UserRepository.UpdateUserReturnAffectedRow(&user)
	if err != nil {
		return util.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "failed to update user",
			Err:     err,
		}
	}

	if rowsAffected == 0 {
		err := usecase.UserRepository.CreateUser(&user)
		if err != nil {
			return util.ErrorObject{
				Code:    http.StatusInternalServerError,
				Message: "failed to create user",
				Err:     err,
			}
		}
	}

	err = util.GenerateJWT(c, user.ID)
	if err != nil {
		return util.ErrorObject{
			Code:    http.StatusInternalServerError,
			Message: "failed to generate jwt",
			Err:     err,
		}
	}

	return nil
}
