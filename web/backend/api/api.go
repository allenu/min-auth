package api

import (
    "encoding/json"
    "net/http"
    "strings"

    userpkg "min-auth/web/backend/user"
)

const (
    ResponseSuccess = "OK"
    ResponseNoAccess = "No Access"
    ResponseDatabaseError = "Database Error"
    ResponseMissingTitle = "Missing Title"
    ResponseMalformedRequest = "Malformed Request"
    ResponseInvalidParent = "Invalid Parent"
    ResponseNoEntity = "No Entity Found"
)

type ReadUserRequest struct {
    UserId string
}
type ReadUserResponse struct {
    ResponseCode string
    User User
}
type User struct {
    UserId string
    Username string
}

var userStoreProvider userpkg.UserStoreProvider

func Setup(usp userpkg.UserStoreProvider) {
    userStoreProvider = usp

    http.HandleFunc("/api/user/", apiFetchUser)
}

func apiFetchUser(w http.ResponseWriter, r *http.Request) {
    userId := strings.TrimPrefix(r.URL.Path, "/api/user/")

    userStore := userStoreProvider.NewUserStore(r)
    user, err := userStore.FetchUserByUserId(userId)

    /*
    dbState := userpkg.NewDatabaseState(r)
    user, err := userpkg.FetchUserByUserId(dbState, userId)
    */

    var responseCode string
    var userResponse User
    if err != nil {
        responseCode = ResponseNoEntity
        userResponse = User{UserId: "anonymous", Username: "anonymous"}
    } else {
        responseCode = ResponseSuccess
        userResponse = User{UserId: user.GetUserId(), Username: user.GetUsername()}
    }
    response := ReadUserResponse{ResponseCode: responseCode, User: userResponse}

    js, _ := json.Marshal(response)
    w.Header().Set("Content-Type", "application/json")
    w.Write(js)
}

