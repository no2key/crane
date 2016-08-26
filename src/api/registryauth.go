package api

import (
	"github.com/Dataman-Cloud/go-component/auth"
	rauth "github.com/Dataman-Cloud/go-component/registryauth"
	"github.com/Dataman-Cloud/go-component/utils/dmerror"
	"github.com/Dataman-Cloud/go-component/utils/dmgin"

	log "github.com/Sirupsen/logrus"
	"github.com/gin-gonic/gin"
)

const (
	CodeCreateRegistryAuthParamError = "400-17001"
	CodeGetRegistryAuthExistError    = "400-17002"
	CodeDeleteRegistryAuthParamError = "400-17003"
	CodeRegistryAuthInvalidUserError = "401-17004"
)

func (api *Api) Create(ctx *gin.Context) {
	var registryAuth rauth.RegistryAuth

	if err := ctx.BindJSON(&registryAuth); err != nil {
		log.Errorf("create registryAuth param error: %v", err)
		dmgin.HttpErrorResponse(ctx, dmerror.NewError(CodeCreateRegistryAuthParamError, err.Error()))
		return
	}

	rs, err := rauth.GetHubApi().List(&rauth.RegistryAuth{Name: registryAuth.Name, AccountId: registryAuth.AccountId})
	if err != nil {
		dmgin.HttpErrorResponse(ctx, err)
		return
	}

	if len(rs) > 0 {
		dmgin.HttpErrorResponse(ctx, dmerror.NewError(CodeGetRegistryAuthExistError, "registryAuth exists"))
		return
	}

	if err := rauth.GetHubApi().Create(&registryAuth); err != nil {
		log.Errorf("create registryAuth operation error: %v", err)
		dmgin.HttpErrorResponse(ctx, err)
		return
	}

	dmgin.HttpOkResponse(ctx, "create success")
}

func (api *Api) List(ctx *gin.Context) {
	account, ok := ctx.Get("account")
	if !ok {
		log.Error("get registryAuths invalid user")
		dmgin.HttpErrorResponse(ctx, dmerror.NewError(CodeRegistryAuthInvalidUserError, "invalid user"))
		return
	}

	registryAuth, err := rauth.GetHubApi().List(&rauth.RegistryAuth{AccountId: account.(auth.Account).ID})
	if err != nil {
		log.Errorf("get registryAuth by name error: %v", err)
		dmgin.HttpErrorResponse(ctx, err)
		return
	}

	dmgin.HttpOkResponse(ctx, registryAuth)
}

func (api *Api) Delete(ctx *gin.Context) {
	account, ok := ctx.Get("account")
	if !ok {
		log.Error("delete registryAuth invalid user")
		dmgin.HttpErrorResponse(ctx, dmerror.NewError(CodeRegistryAuthInvalidUserError, "invalid user"))
		return
	}

	name := ctx.Param("rauth_name")
	if name == "" {
		log.Errorf("get registryAuth name invalid")
		dmgin.HttpErrorResponse(ctx, dmerror.NewError(CodeDeleteRegistryAuthParamError, "registryAuth name invalid"))
		return
	}

	if err := rauth.GetHubApi().Delete(&rauth.RegistryAuth{Name: name, AccountId: account.(auth.Account).ID}); err != nil {
		log.Errorf("delete registryAuth error: %v", err)
		dmgin.HttpErrorResponse(ctx, err)
		return
	}

	dmgin.HttpOkResponse(ctx, "delete success")
}