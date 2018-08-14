package user

import (
	"net/http"
	"google.golang.org/appengine"
)

func NewUserStoreProvider() GoogleUserStoreProvider {
    return GoogleUserStoreProvider{}
}

type GoogleUserStoreProvider struct {
}

func (g GoogleUserStoreProvider) NewUserStore(r *http.Request) UserStore {
    ctx := appengine.NewContext(r)

    return GoogleUserStore{dbState: DatabaseState{appEngineContext: ctx}}
}

type GoogleUserStore struct {
    dbState DatabaseState
}

func (g GoogleUserStore) FetchUserByServiceId(serviceId string) (User, error) {
    return FetchUserByServiceId(g.dbState, serviceId)
}

func (g GoogleUserStore) FetchUserByUserId(userId string) (User, error) {
    return FetchUserByUserId(g.dbState, userId)
}

func (g GoogleUserStore) FetchUsers() ([]User, error) {
    return FetchUsers(g.dbState)
}

func (g GoogleUserStore) InsertUser(serviceName string, serviceId string, username string) (User, error) {
    return InsertUser(g.dbState, serviceName, serviceId, username)
}

