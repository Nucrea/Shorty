package common

import (
	"crypto/sha512"
	"encoding/hex"
	"fmt"
	"math/rand"
	"regexp"
	"strings"
	"time"

	"github.com/asaskevich/govalidator"
)

var shortIdRegexp = regexp.MustCompile(`^\w+$`)

const shortIdCharset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

func NewShortId(size int) string {
	sb := strings.Builder{}

	src := rand.NewSource(time.Now().UnixMicro())
	r := rand.New(src)

	for range size {
		rNum := r.Intn(len(shortIdCharset))
		randChar := shortIdCharset[rNum]
		sb.WriteByte(randChar)
	}

	return sb.String()
}

func NewAssetHash(asset []byte) string {
	hash := sha512.Sum512(asset)
	return hex.EncodeToString(hash[:])
}

func ValidateShortId(value string) bool {
	return shortIdRegexp.MatchString(value)
}

func ValidateUrl(url string) string {
	if len(url) > 2000 {
		return ""
	}

	url = strings.TrimSpace(url)
	if !govalidator.IsURL(url) {
		return ""
	}

	if !strings.HasSuffix(url, "http") || !strings.HasSuffix(url, "https") {
		url = fmt.Sprintf("https://%s", url)
	}

	return url
}

func NewDigitsString(size int) string {
	sb := strings.Builder{}

	src := rand.NewSource(time.Now().UnixMicro())
	r := rand.New(src)

	for range size {
		randNum := byte(r.Intn(10))
		randByte := byte('0') + randNum
		sb.WriteByte(randByte)
	}

	return sb.String()
}

func MaskSecret(secret string) string {
	hLen := len(secret) / 2

	sb := strings.Builder{}
	sb.WriteString(secret[0:hLen])
	for range hLen {
		sb.WriteRune('*')
	}

	return sb.String()
}
