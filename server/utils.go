package server

import (
	"fmt"
	"shorty/src/common"
	"time"
)

var resourceTokenSecret = common.NewShortId(10)

type ExpiringToken struct {
	Value   string
	Exipres int64
}

// func (e ExpiringToken) ToQuery() string {
// 	return fmt.Sprintf("token=%s&expires=%d", url.QueryEscape(e.Value), e.Exipres)
// }

func NewResourceToken(resource string, expiresAt time.Time) ExpiringToken {
	expires := expiresAt.Unix()
	raw := fmt.Sprintf("%s%d%s", resource, expires, resourceTokenSecret)
	return ExpiringToken{
		Value:   common.HashsumSHA1(raw),
		Exipres: expires,
	}
}

func CheckResourceToken(resource string, expiresAt int64, token string) bool {
	raw := fmt.Sprintf("%s%d%s", resource, expiresAt, resourceTokenSecret)
	return common.HashsumSHA1(raw) == token
}
