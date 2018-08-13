# min-auth

This is the minimum code required to get auth (via Twitter) and user accounts working with Go and Google Cloud (App Engine).

This is the base user/auth code I wrote for [This Way or That (go-lang edition)](https://this-way-or-that.appspot.com/).

As always, this is a work in progress. Follow me on twitter: [@ussherpress](http://twitter.com/ussherpress)

# How to use

git clone into your $GOPATH/pkg/src/min-auth.

    cp rename_app.yaml app.yaml

Edit app.yaml and enter your own session secret and twitter consumer key and consumer secret. (These bits come from https://apps.twitter.com/)

Test it out with:

    dev_appserver.py app.yaml

