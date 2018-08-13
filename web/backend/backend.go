package backend

import (
    "min-auth/web/backend/user"
    "min-auth/web/backend/auth"
    "min-auth/web/backend/api"
)

func init() {
    usp := user.NewUserStoreProvider()
    auth.Setup(usp)
    api.Setup(usp)
}

