package mailbox

import (
	"encoding/base64"
	"golang.org/x/oauth2"
	"strings"
	"time"
)

const basicAuthenticationPrefix = "Basic "
const passwordSeparator = ":"

type BasicToken struct {
	User, Password string
}

func (b BasicToken) Token() (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken:  base64.URLEncoding.EncodeToString([]byte(b.User + ":" + b.Password)),
		TokenType:    "basic",
		RefreshToken: "",
		Expiry:       time.Now().Add(time.Hour * 24),
	}, nil
}

func NewBasicToken(user, pass string) oauth2.TokenSource {
	return BasicToken{
		User:     user,
		Password: pass,
	}
}

func FromHeaderString(auth string) (BasicToken, error) {
	if len(auth) <= len(basicAuthenticationPrefix) {
		return BasicToken{}, AuthenticationUnknown
	}
	auth = auth[len(basicAuthenticationPrefix):]
	userPass := base64.URLEncoding.EncodeToString([]byte(auth))
	strs := strings.Split(userPass, passwordSeparator)
	if len(strs) != 2 {
		return BasicToken{}, AuthenticationUnknown
	}
	return BasicToken{
		User:     strs[0],
		Password: strs[1],
	}, nil
}
