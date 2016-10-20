package main

import (
	"fmt"
	 "net/http"
	 "os"
	 "io"
	 "encoding/json"
	 "strings"
)

type Response struct {
	PackageName string
	DownloadLink string
}

func Register(w http.ResponseWriter, r *http.Request) {
	response := RegisterUser(r.FormValue("email"), r.FormValue("password"))
	if response == "" {
		w.WriteHeader(504)
	}
	w.Header().Set("Authorization", response)
}

func Login(w http.ResponseWriter, r *http.Request) {
	response := LoginUser(r.FormValue("email"), r.FormValue("password"))
	if response == "" {
		w.WriteHeader(504)
	}
	w.Header().Set("Authorization", response)
}

func Query(w http.ResponseWriter, r *http.Request) {
	query := GetInfo(strings.TrimSpace(r.FormValue("hash")))
	if query == nil {
		fmt.Println(query)
		w.WriteHeader(504)
		return;
	}
	responseStruct := Response{query[0], query[1]}
	response, err := json.Marshal(responseStruct)
	if err != nil {
		w.WriteHeader(504)
		return;
	}
	reply := string(response)
	reply += "\n"
	fmt.Println("Reply" + reply)
	fmt.Fprintf(w, reply)
}

func Logout(w http.ResponseWriter, r *http.Request) {
	token :=  r.Header.Get("Authorization")
	claims, _ := VerifyToken(token)
	fmt.Println("logging out:",claims["aud"].(string))
	DelToken(claims["aud"].(string), token)
	w.WriteHeader(200)
}

func Upload(w http.ResponseWriter, r *http.Request) {
	r.ParseMultipartForm(10 << 20)
	file, handler, err := r.FormFile("apk")
	if err != nil || IsNotApk(handler.Filename) {
		w.WriteHeader(400)
		return
	}
	defer file.Close()
	fileName := RandomString() + ".apk"
	// TODO: look into exact meaning of os flags in line below
	f, err := os.OpenFile("./static/apks/" + fileName, os.O_WRONLY | os.O_CREATE, 0666)
        if err != nil {
            fmt.Println(err)
            w.WriteHeader(400)
            return
        }
        defer f.Close()
    io.Copy(f, file)
    signHash := r.FormValue("hash")
    packageName := r.FormValue("name")
	hostPackageName := r.FormValue("package")

    SetInfo(signHash, packageName, fileName)
    BroadcastUpdate(hostPackageName)
    w.WriteHeader(200)
    return
}

func Redirect(w http.ResponseWriter, r *http.Request) {
	http.Redirect(w, r, "https://"+r.Host+":443"+r.RequestURI, http.StatusMovedPermanently)
}

func errorHandler(w http.ResponseWriter, r *http.Request, responseCode int) {
	w.WriteHeader(responseCode)
	switch responseCode {
	case http.StatusNotFound:
		fmt.Fprintf(w, "404")
	}
}