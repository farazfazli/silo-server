// TODO: Find a decent way to benchmark Android client connections
package main

import (
	"crypto/tls"
	"fmt"
	"log"
	"net"
	"net/http"
	"time"
	"strings"
)

const BASE_URL = "https://silo.live"

func main() {
	Clients = make(map[string]net.Conn)
	// TODO: get domain & real certificate
	cert, err := tls.LoadX509KeyPair("server.pem", "server.key")
	if err != nil {
		log.Fatal(err)
	}
	config := tls.Config{Certificates: []tls.Certificate{cert}}
	ln, err := net.Listen("tcp", ":1337") // External port for socket client connection
	if err != nil {
		fmt.Println(err)
	}
	go ListenForClients(ln)

	fileServer := http.FileServer(http.Dir("static"))
	mux := http.NewServeMux()
	mux.Handle("/register", CredentialsMiddleware(AuthMiddleware(SecurityMiddleware(http.HandlerFunc(Register)))))
	mux.Handle("/upload", FileUploadMiddleware(AuthMiddleware(SecurityMiddleware(http.HandlerFunc(Upload)))))
	mux.Handle("/login", CredentialsMiddleware(AuthMiddleware(SecurityMiddleware(http.HandlerFunc(Login)))))
	mux.Handle("/query", SecurityMiddleware(http.HandlerFunc(Query)))
	mux.Handle("/logout", AuthMiddleware(SecurityMiddleware(http.HandlerFunc(Logout))))
	
	mux.Handle("/", SecurityMiddleware(fileServer))

	httpServer := &http.Server{
		Addr:         ":80", // 80
		Handler:      http.HandlerFunc(Redirect),
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
	}

	httpsServer := &http.Server{
		Addr:         ":443", // 443
		Handler:      mux,
		ReadTimeout:  5 * time.Second,
		WriteTimeout: 5 * time.Second,
		TLSConfig:    &config,
	}

	fmt.Printf("Listening on %s and %s\n", httpServer.Addr, httpsServer.Addr)
	go func() {
		httpServer.ListenAndServe() // HandleFunc here reroutes all
	}()
	fmt.Println(httpsServer.ListenAndServeTLS("", ""))
}

func IsNotApk(filename string) bool {
	return strings.HasSuffix(filename, "apk") == false
}