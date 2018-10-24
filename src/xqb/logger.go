package xqb

import (
	"github.com/jonnywang/go-kits/redis"
	"log"
)

var Logger = redis.Logger

func init() {
	redis.Logger.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}