package session

import (
	"log"
	"test/config"

	"github.com/gorilla/sessions"
)

var Store = sessions.NewCookieStore([]byte(getSessionSecret()))

func getSessionSecret() string {
	config.LoadEnv()
	secret := config.SessionSecret
	if secret == "" {
		log.Fatal("SESSION_SECRET environment variable is not set")
	}

	return secret

}
