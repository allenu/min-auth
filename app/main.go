package main

import (
    "net/http"
    "html/template"
    "path"

    _ "min-auth/web/backend"
    "min-auth/web/backend/auth"
)

func init() {
    http.HandleFunc("/", indexHandler)
}

func indexHandler(w http.ResponseWriter, r *http.Request) {
    userInfo := auth.GetUserInfo(r)

    fp := path.Join("templates", "index.html")
    if tmpl, err := template.ParseFiles(fp); err == nil {
        if err == nil {
            tmpl.Execute(w, userInfo)
        } else {
            http.Error(w, err.Error(), http.StatusInternalServerError)
        }
    } else {
        http.Error(w, err.Error(), http.StatusInternalServerError)
    }
}

