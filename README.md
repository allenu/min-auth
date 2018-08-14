# minauth

This is the minimum code required to get auth (via Twitter) and user accounts working with Go and Google Cloud (App Engine).

This is the base user/auth code I wrote for [This Way or That (go-lang edition)](https://this-way-or-that.appspot.com/).

As always, this is a work in progress. Follow me on twitter: [@ussherpress](http://twitter.com/ussherpress)

# How to test it out

git clone into your $GOPATH/pkg/src/minauth.

    cd minauth/example
    cp rename_app.yaml app.yaml

Edit app.yaml and enter your own session secret and twitter consumer key and consumer secret. (These bits come from https://apps.twitter.com/)

Test it out with:

    dev_appserver.py app.yaml

# How to use it in your project

Simply provide a link to /auth/signin to sign in and /auth/signout to sign out. (See REST API section below for more info.)

Sample code or getting the currently signed in user:

    import "github.com/allenu/minauth"

    func myHandler(w http.ResponseWriter, r *http.Request) {
        userInfo := minauth.GetUserInfo(r)

        // userInfo.Username is user's twitter handle or "anonymous" if not signed in
        // userInfo.UserId is the user's unique ID (a UUID generated on first sign-in) or "anonymous" if not signed in
    }

# Parts to this system

## User Storage

User info is stored in Google App Engine datastore (see web/backend/user).

When a user signs in for the first time, their info is added to the datastore. The user info consists of

* UserId, which is a UUID generated for each user
* Username, which is currently just the twitter handle (in the future this could be something of the user's choosing)
* Servicename, which is just the string "twitter" for now (if other auth providers are added, this would be "google" or "facebook")
* ServiceId, which is the unique identifier provided by twitter, prefixed by "twitter-"

When a user signs in again, this info is looked up via the Servicename and ServiceId.

## REST API

To access info on the user or to sign in and out:

    /api/user - this returns the current user info as json ({UserId:, Username:}).

    /auth/signin - initiates the sign-in for twitter oauth
    /auth/signout - signs the current user out
    /auth/callback - handler for Twitter oauth callback

