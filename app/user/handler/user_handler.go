package handler

import (
	"errors"
	"learn_oauth/app/user/usecase"
	"learn_oauth/domain"
	"learn_oauth/util"
	"log"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
	"golang.org/x/oauth2/google"
)

type UserHandler struct {
	Userusecase usecase.IUserUsecase
}

var googleOauthConfig *oauth2.Config

func NewUserHandler(userUsecase usecase.IUserUsecase) *UserHandler {
	return &UserHandler{
		Userusecase: userUsecase,
	}
}

func NewGoogleOauthConfig() {
	googleOauthConfig = &oauth2.Config{
		RedirectURL:  "http://localhost:8080/auth/google/callback",
		ClientID:     os.Getenv("GOOGLE_CLIENT_ID"),
		ClientSecret: os.Getenv("GOOGLE_CLIENT_SECRET"),
		Scopes: []string{
			"https://www.googleapis.com/auth/userinfo.email",
			"https://www.googleapis.com/auth/userinfo.profile",
		},
		Endpoint: google.Endpoint,
	}
}

func (handler *UserHandler) OauthGoogleLoginHandler(c *gin.Context) {
	log.Println(googleOauthConfig)
	oauthState := handler.Userusecase.GenerateStateOauthCookie(c)
	u := googleOauthConfig.AuthCodeURL(oauthState)
	http.Redirect(c.Writer, c.Request, u, http.StatusTemporaryRedirect)
}

func (handler *UserHandler) OauthGoogleCallbackHandler(c *gin.Context) {
	oauthState, _ := c.Cookie("oauthstate")
	if c.Request.FormValue("state") != oauthState {
		util.FailedResponse(c, http.StatusUnauthorized, "invalid oauth google state", errors.New("your state is invalid"))
		return
	}

	data, err := handler.Userusecase.GetUserDataFromGoogle(c.Request.FormValue("code"), googleOauthConfig)
	if err != nil {
		util.FailedResponse(c, http.StatusUnauthorized, "failed to get user data from google", err)
		return
	}

	var user domain.User
	user.Name = data["name"].(string)
	user.Email = data["email"].(string)

	//create jwt token
	errObject := handler.Userusecase.GenerateAndSetJWT(c, user)
	if errObject != nil {
		errObject := errObject.(util.ErrorObject)
		util.FailedResponse(c, errObject.Code, errObject.Message, errObject.Err)
		return
	}

	util.SuccessedResponse(c, http.StatusOK, "login success", nil)
}

func (handler *UserHandler) GetLoginUser(c *gin.Context) {
	user, ok := c.Get("user")
	if !ok {
		util.FailedResponse(c, http.StatusInternalServerError, "failed to get login user", errors.New(""))
		return
	}

	util.SuccessedResponse(c, http.StatusOK, "successfully get login user", user)
}
