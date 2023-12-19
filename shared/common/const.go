package common

import (
	"github.com/namsral/flag"
)

var ShowLog = false

func init() {
	flag.BoolVar(&ShowLog, "show-log", false, "show log")
	flag.Parse()
}

const (
	KeyCompFiber = "fiber"
	KeyCompGorm  = "postgres"
	KeyCompRedis = "redis"
	KeyJwt       = "jwt"
	KeyDex       = "dex"
)

type TokenPayload struct {
	UserId   int    `json:"user_id"`
	SubToken string `json:"sub_token"`
}

func (t TokenPayload) GetUserId() int {
	return t.UserId
}

func (t TokenPayload) GetSubToken() string {
	return t.SubToken
}
