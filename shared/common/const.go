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
)

type TokenPayload struct {
	UserId int `json:"user_id"`
}

func (t TokenPayload) GetUserId() int {
	return t.UserId
}
