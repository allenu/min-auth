package auth

import (
    "os"
    "encoding/json"
    "log"
    "net/http"

    userpkg "min-auth/web/backend/user"

	"github.com/mrjones/oauth"
    "github.com/gorilla/sessions"
    "google.golang.org/appengine"
    "google.golang.org/appengine/urlfetch"
)

type TwitterUserInfo struct {
    Id_str string
    Screen_name string
    Location string
    Description string
}

type User interface {
    GetUserId() string
    GetUsername() string
}

type UserInfo struct {
    UserId string
    Username string
    IsAdmin bool
    IsSignedIn bool
}

func (u UserInfo) GetUserId() string {
    return u.UserId
}

func (u UserInfo) GetIsAdmin() bool {
    return u.IsAdmin
}

var c *oauth.Consumer
var requestToken *oauth.RequestToken
var consumerKey string
var consumerSecret string

const sessionName = "ThisWayOrThat-Session"

var store *sessions.CookieStore

var userStoreProvider userpkg.UserStoreProvider

func Setup(usp userpkg.UserStoreProvider) {
    userStoreProvider = usp

    cookieStoreSecret := os.Getenv("sessions_secret")
    store = sessions.NewCookieStore([]byte(cookieStoreSecret))
    consumerKey = os.Getenv("twitter_consumer_key")
    consumerSecret = os.Getenv("twitter_consumer_secret")

    http.HandleFunc("/auth/signin", authSignInHandler)
    http.HandleFunc("/auth/signout", authSignOutHandler)
    http.HandleFunc("/auth/callback", authCallbackHandler)
}

func authSignInHandler(w http.ResponseWriter, r *http.Request) {
    var u string
    var err error

    ctx := appengine.NewContext(r)
    httpClient := urlfetch.Client(ctx)

	c = oauth.NewCustomHttpClientConsumer(
		consumerKey,
		consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		},
        httpClient)

    c.Debug(true)

    url_callback := "http://" + r.Host + "/auth/callback"
    requestToken, u, err = c.GetRequestTokenAndUrl(url_callback)
	if err != nil {
		log.Fatal(err)
	}

	http.Redirect(w, r, u, http.StatusTemporaryRedirect)
}

func authSignOutHandler(w http.ResponseWriter, r *http.Request) {
    session, err := store.Get(r, sessionName)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }
    delete(session.Values, "username")
    session.Save(r, w)

    http.Redirect(w, r, "http://" + r.Host, http.StatusTemporaryRedirect)
}

func authCallbackHandler(w http.ResponseWriter, r *http.Request) {
	values := r.URL.Query()
    verificationCode := values.Get("oauth_verifier")

    // oauth_token := values.Get("oauth_token")
    //fmt.Printf("oauth_token: %v\noauth_verifier: %v\n", oauth_token, verificationCode)

    ctx := appengine.NewContext(r)
    httpClient := urlfetch.Client(ctx)

	c = oauth.NewCustomHttpClientConsumer(
		consumerKey,
		consumerSecret,
		oauth.ServiceProvider{
			RequestTokenUrl:   "https://api.twitter.com/oauth/request_token",
			AuthorizeTokenUrl: "https://api.twitter.com/oauth/authorize",
			AccessTokenUrl:    "https://api.twitter.com/oauth/access_token",
		},
        httpClient)

    if c == nil {
        log.Fatal("bad client")
    }

    c.Debug(true)

    accessToken, err := c.AuthorizeToken(requestToken, verificationCode)
    if err != nil {
		log.Fatal(err)
    }

	verifyClient, err := c.MakeHttpClient(accessToken)
	if err != nil {
		log.Fatal(err)
    }

	response, err := verifyClient.Get("https://api.twitter.com/1.1/account/verify_credentials.json")
	if err != nil {
		log.Fatal(err)
	}
    defer response.Body.Close()

    // from https://stackoverflow.com/questions/17156371/how-to-get-json-response-in-golang
    target := new(TwitterUserInfo)
    err = json.NewDecoder(response.Body).Decode(target)
    if err != nil {
        log.Fatal(err)
    }

	/*
    log.Printf("Response to verify_credentials is %v\n", target)

    log.Printf("user_id: %v\n", target.Id_str)
    log.Printf("username: %v\n", target.Screen_name)
    log.Printf("location: %v\n", target.Location)
    log.Printf("description: %v\n", target.Description)
    */

    session, err := store.Get(r, sessionName)
    if err != nil {
        http.Error(w, err.Error(), http.StatusInternalServerError)
        return
    }

    // To future-proof us for other auth, prefix with "twitter-"
    serviceId := "twitter-" + target.Id_str
    serviceName := "twitter"

    // Look up user by serviceId. If we can't find it, create a new entry.
    userStore := userStoreProvider.NewUserStore(r)
    var user User
    user, err = userStore.FetchUserByServiceId(serviceId)

    var username string
    if err == nil {
        // log.Printf("Loaded user: %v", user)
        username = user.GetUsername()
    } else {
        // create a new user for him
        username = target.Screen_name
        user, err = userStore.InsertUser(serviceName, serviceId, username)
        // log.Printf("Couldn't find user, created: %v with err: %v", user, err)
    }
    session.Values["username"] = username
    session.Values["userid"] = user.GetUserId()
    session.Save(r, w)

    http.Redirect(w, r, "http://" + r.Host, http.StatusTemporaryRedirect)
}

func GetUserInfo(r *http.Request) UserInfo {
    session, err := store.Get(r, sessionName)

    var username string = "anonymous"
    var userId string = "anonymous"
    var isAdmin = false
    var isSignedIn = false

    if err == nil {
        maybeUsername, gotUsername := session.Values["username"]
        maybeUserId, gotUserId := session.Values["userid"]
        if gotUsername && gotUserId {
            username = maybeUsername.(string)
            userId = maybeUserId.(string)
            // log.Print("Found user ", username)
            // log.Print("Values is ", session.Values)
            isSignedIn = true
        } else {
            // log.Print("Couldn't find user - sessionName: ", sessionName, r.URL.Path)
            // log.Print("Values is ", session.Values)
        }
    }

    userInfo := UserInfo{UserId: userId, Username: username, IsAdmin: isAdmin, IsSignedIn: isSignedIn}
    return userInfo
}

