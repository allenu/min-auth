
package user

import (
    "net/http"
)

type UserStoreProvider interface {
    NewUserStore(r *http.Request) UserStore
}

type UserStore interface {
    FetchUserByServiceId(serviceId string) (User, error)
    FetchUserByUserId(userId string) (User, error)
    FetchUsers() ([]User, error)
    InsertUser(serviceName string, serviceId string, username string) (User, error)
}

