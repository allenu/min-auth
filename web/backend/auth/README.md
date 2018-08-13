
# Purpose

auth is used for authenticating with OAuth (Twitter only for now) and then
- getting or storing the result in the UserStore database (storing if it's a brand new user)
- storing the userId and username in the session cookie

It is also queried to get the currently signed in user via auth.GetUserInfo()

# How to use

import "blahblah/backend/auth"
import "blahblah/backend/user"

- 
    usp := user.NewUserStoreProvider()
    auth.Setup(usp)

- 

  func PageHandler(w http.ResponseWriter, r *http.Request) {
    user = auth.GetUserInfo(r)

