package mailbox

import (
	"context"
	"errors"
	"gomail/pkg/db"
	"google.golang.org/grpc"
	"google.golang.org/grpc/metadata"
	"log"
)

var (
	AuthenticationNotFound = errors.New("can not found auth information")
	AuthenticationUnknown  = errors.New("auth string is unknown")
	AuthenticationFailed   = errors.New("user not found or wrong password")

	WhiteList = []string{"proto.MailBox/Register"}
)

type AuthInterceptor interface {
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
func NewAuthInterceptor(storage db.Storage) AuthInterceptor {
	return &defaultInterceptor{registry: storage}
}

type defaultInterceptor struct {
	registry db.Storage
}

type User struct {
	Name     string `bson:"user"`
	Password string `bson:"password"`
}

func (d *defaultInterceptor) valid(token BasicToken) bool {
	res := &User{}
	err := d.registry.Get(map[string]interface{}{"name": token.User, "password": token.Password}, res)
	if err != nil {
		log.Printf("user:%s cannot found because error : %v", token.User, err)
		return false
	}
	return token.User == res.Name && token.Password == res.Password
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
		if !d.valid(tk) {
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
