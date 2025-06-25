package services

import (
	"log"
	"os"
	"sync"

	"github.com/ctreminiom/go-atlassian/confluence"
	"github.com/pkg/errors"
)

func loadAtlassianCredentials() (host, mail, token string) {
	host = os.Getenv("ATLASSIAN_HOST")
	token = os.Getenv("ATLASSIAN_TOKEN")

	if host == "" || token == "" {
		log.Fatal("ATLASSIAN_HOST, ATLASSIAN_EMAIL, ATLASSIAN_TOKEN are required, please set it in MCP Config")
	}

	return host, mail, token
}

var ConfluenceClient = sync.OnceValue[*confluence.Client](func() *confluence.Client {
	host, _, token := loadAtlassianCredentials()

	instance, err := confluence.New(DefaultHttpClient(), host)
	if err != nil {
		log.Fatal(errors.WithMessage(err, "failed to create confluence client"))
	}

	instance.Auth.SetBearerToken(token)
	return instance
})
