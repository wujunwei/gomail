package auth

import (
	"context"
	"errors"
	"gomail/pkg/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
	"strings"
)

var (
	AuthenticationNotFound = errors.New("can not found auth information")
	AuthenticationUnknown  = errors.New("auth string is unknown")
	AuthenticationFailed   = errors.New("user not found or wrong password")

	WhiteList = []string{"proto.MailBox/Register", "proto.MailBox/Login"}
)

type Interceptor interface {
	StreamAuth(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error
	UnaryAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error)
}

func InWhiteList(url string) bool {
	for _, s := range WhiteList {
		if s == url {
			return true
		}
	}
	return false
}
func NewAuthInterceptor(storage db.Storage, sess db.Session) Interceptor {
	return &defaultInterceptor{registry: storage, sess: sess}
}

type defaultInterceptor struct {
	registry db.Storage
	sess     db.Session
}

type User struct {
	ID       string `bson:"_id"`
	Name     string `bson:"user"`
	Password string `bson:"password"`
	Weight   int32  `bson:"weight"`
}

func (d *defaultInterceptor) getUser(token Token) *User {
	res := &User{}
	conditions := map[string]interface{}{}
	switch token.Type() {
	case BearerAuthenticationTyp:
		conditions["_id"] = token.String()
	case BasicAuthenticationType:
		authStr := token.String()
		strings.Split(authStr, passwordSeparator)
		if len(authStr) != 2 {
			return nil
		}
		conditions["user"] = authStr[0]
		conditions["password"] = authStr[1]
	default:
		return nil
	}
	err := d.sess.Get(conditions, res)
	if err != nil {
		log.Printf("user:%s cannot found because error : %v", token, err)
		return nil
	}
	return res
}

func (d *defaultInterceptor) check(ctx context.Context, method string) error {
	md, ok := metadata.FromIncomingContext(ctx)
	if !ok || len(md["authorization"]) == 0 || md["authorization"][0] == "" {
		return AuthenticationNotFound
	}
	if !InWhiteList(method) {
		tk, err := FromHeaderString(md["authorization"][0])
		if err != nil {
			return err
		}
		if d.getUser(tk) == nil {
			return AuthenticationFailed
		}
	}
	return nil
}
func (d *defaultInterceptor) StreamAuth(srv interface{}, ss grpc.ServerStream, info *grpc.StreamServerInfo, handler grpc.StreamHandler) error {
	if info.IsClientStream {
		if err := d.check(ss.Context(), info.FullMethod); err != nil {
			return err
		}
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(srv, ss)
}

func (d *defaultInterceptor) UnaryAuth(ctx context.Context, req interface{}, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (interface{}, error) {
	if err := d.check(ctx, info.FullMethod); err != nil {
		return nil, err
	}
	// Continue execution of handler after ensuring a valid token.
	return handler(ctx, req)
}
