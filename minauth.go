package minauth

import (
    "github.com/allenu/minauth/user"
    "github.com/allenu/minauth/auth"
    "github.com/allenu/minauth/api"

    "net/http"
)

func init() {
    usp := user.NewUserStoreProvider()
    auth.Setup(usp)
    api.Setup(usp)
}

func GetUserInfo(r *http.Request) auth.UserInfo {
    return auth.GetUserInfo(r)
}

