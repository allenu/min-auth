package minauth

import (
    "github.com/allenu/minauth/user"
    "github.com/allenu/minauth/auth"
    "github.com/allenu/minauth/api"
)

func init() {
    usp := user.NewUserStoreProvider()
    auth.Setup(usp)
    api.Setup(usp)
}

