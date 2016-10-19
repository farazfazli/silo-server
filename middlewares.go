package main

import (
	"fmt"
	"net"
	"net/http"
	"strings"
)

// TODO: Enhance
func SecurityMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubdomains; preload")
		w.Header().Set("X-Frame-Options", "DENY")
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			fmt.Println("Error getting IP:", err)
		}
		fmt.Println("New connection from IP:", ip)
		h.ServeHTTP(w, r)
	})
}

func AuthMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubdomains; preload")
		w.Header().Set("X-Frame-Options", "DENY")
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			fmt.Println("Error getting IP:", err)
		}
		fmt.Println("New connection from IP:", ip)
		token := r.Header.Get("Authorization")
		fmt.Println("Authorization:", token)
		if token != "" && strings.Contains(token, "Bearer ") {
			token = strings.Split(token, " ")[1]
			tokens, ok := VerifyToken(token)
			if ok {
				fmt.Fprintf(w, "%s\n", tokens["aud"])
				r.Header.Set("Authorization", token)
			h.ServeHTTP(w, r)
			} else {
			w.WriteHeader(401)
			return
			}
		} else {
			// Go to Login middleware
			h.ServeHTTP(w, r)
		}
	})
}

func CredentialsMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Strict-Transport-Security", "max-age=31536000; includeSubdomains; preload")
		w.Header().Set("X-Frame-Options", "DENY")
		ip, _, err := net.SplitHostPort(r.RemoteAddr)
		if err != nil {
			fmt.Println("Error getting IP:", err)
		}
		fmt.Println("New connection from IP:", ip)
			email := Validate(r.FormValue("email"))
			password := r.FormValue("password")
		if len(email) < 5 && len(password) > 0 {
			fmt.Fprintf(w, "Error registering, invalid length")
			w.WriteHeader(401)
			return
		}
			h.ServeHTTP(w, r)
	})
}

func FileUploadMiddleware(h http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		r.Body = http.MaxBytesReader(w, r.Body, 10 << 20)
		err := r.ParseForm()
		if err != nil || r.FormValue("name") != "" && r.FormValue("hash") != "" {
			w.WriteHeader(400)
		}
		h.ServeHTTP(w, r)
	})
}