package auth

import (
	"encoding/base64"
	"golang.org/x/oauth2"
	"strings"
	"time"
)

const (
	BasicAuthenticationType = "Basic"
	BearerAuthenticationTyp = "Bearer"
)
const passwordSeparator = ":"

type Token interface {
	oauth2.TokenSource
	Type() string
	String() string
}
type BasicToken struct {
	User     string
	Password string
}

func (b BasicToken) Token() (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken:  base64.URLEncoding.EncodeToString([]byte(b.String())),
		TokenType:    "basic",
		RefreshToken: "",
		Expiry:       time.Now().Add(time.Hour * 24),
	}, nil
}
func (b BasicToken) Type() string {
	return BasicAuthenticationType
}
func (b BasicToken) String() string {
	return b.User + ":" + b.Password
}

func NewBasicToken(user, pass string) Token {
	return BasicToken{
		User:     user,
		Password: pass,
	}
}

type BearerToken struct {
	ID string
}

func (b BearerToken) Token() (*oauth2.Token, error) {
	return &oauth2.Token{
		AccessToken:  b.ID,
		TokenType:    "basic",
		RefreshToken: "",
		Expiry:       time.Now().Add(time.Hour * 24),
	}, nil
}
func (b BearerToken) Type() string {
	return BearerAuthenticationTyp
}
func (b BearerToken) String() string {
	return b.ID
}
func NewBearerToken(id string) Token {
	return BearerToken{
		ID: id,
	}
}

func FromHeaderString(authStr string) (Token, error) {
	if strings.HasPrefix(authStr, BasicAuthenticationType) {
		authStr = strings.TrimLeft(authStr[len(BasicAuthenticationType):], " ")
		userPass, err := base64.URLEncoding.DecodeString(authStr)
		if err != nil {
			return BasicToken{}, AuthenticationUnknown
		}
		strs := strings.Split(string(userPass), passwordSeparator)
		if len(strs) != 2 {
			return BasicToken{}, AuthenticationUnknown
		}
		return BasicToken{
			User:     strs[0],
			Password: strs[1],
		}, nil
	}
	if strings.HasPrefix(authStr, BearerAuthenticationTyp) {
		authStr = strings.TrimLeft(authStr[len(BearerAuthenticationTyp):], " ")
		if len(authStr) == 0 {
			return BasicToken{}, AuthenticationUnknown
		}
		return NewBearerToken(authStr), nil
	}
	return BasicToken{}, AuthenticationUnknown
}
