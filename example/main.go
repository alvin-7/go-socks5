package main

import (
	"log"
	"os"

	socks5 "github.com/alvin-7/go-socks5"
)

func main() {
	// creds := socks5.StaticCredentials{
	// 	"username": "password",
	// }
	conf := &socks5.Config{
		// AuthMethods: []socks5.Authenticator{&socks5.UserPassAuthenticator{
		// 	Credentials: creds,
		// }},
		Resolver: NewDNSResolver(),
		Logger:   log.New(os.Stdout, "[socks5]", log.LstdFlags),
	}
	server, err := socks5.New(conf)
	if err != nil {
		panic(err)
	}

	// Create SOCKS5 proxy on localhost port 8000
	if err := server.ListenAndServe("tcp", "0.0.0.0:1080"); err != nil {
		panic(err)
	}
}
