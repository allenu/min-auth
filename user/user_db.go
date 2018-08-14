package user

import (
    "errors"
	"net/http"
    "log"

    "golang.org/x/net/context"
	"google.golang.org/appengine"
	"google.golang.org/appengine/datastore"
    "github.com/google/uuid"
)

type DatabaseState struct {
    appEngineContext context.Context
}

const userEntityName = "User"

// API

func NewDatabaseState(r *http.Request) DatabaseState {
	c := appengine.NewContext(r)

    return DatabaseState{appEngineContext: c}
}

func FetchUserByUserId(dbState DatabaseState, userId string) (User, error) {
    key := datastore.NewKey(dbState.appEngineContext, userEntityName, userId, 0, nil)
    var user DatabaseUser

    log.Printf("Fetching user by userId %v", userId)

    err := datastore.Get(dbState.appEngineContext, key, &user)
    log.Printf("err: %v", err)

    return user, err
}

func FetchUserByServiceId(dbState DatabaseState, serviceId string) (User, error) {
    q := datastore.NewQuery(userEntityName).Filter("ServiceId =", serviceId).Limit(1)
    users := make([]DatabaseUser, 0, 1)
    var user User

    _, err := q.GetAll(dbState.appEngineContext, &users)
    if err == nil {
        if len(users) > 0 {
            user = users[0]
        } else {
            err = errors.New("User not found")
        }
    }

    return user, err
}

func InsertUser(dbState DatabaseState, serviceName string, serviceId, username string) (User, error) {
    userId := uuid.New().String()
    user := DatabaseUser{ServiceName: serviceName, ServiceId: serviceId, UserId: userId, Username: username}

    key := datastore.NewKey(dbState.appEngineContext, userEntityName, userId, 0, nil)
    _, err := datastore.Put(dbState.appEngineContext, key, &user)

    return user, err
}

func FetchUsers(dbState DatabaseState) ([]User, error) {
    const maxUsers = 50
    databaseUsers := make([]DatabaseUser, 0, maxUsers)

    q := datastore.NewQuery(userEntityName).Order("Username").Limit(maxUsers)
    if _, err := q.GetAll(dbState.appEngineContext, &databaseUsers); err == nil {
        users := make([]User, 0, len(databaseUsers))
        for _, u := range databaseUsers {
            users = append(users, u)
        }
        return users, nil
    } else {
        return nil, err
    }
}

