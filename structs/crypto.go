package structs

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"
)

func Md5Bytes(b []byte) string {
	w := md5.New()
	w.Write(b)
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	return md5str
}

func HmacSha256(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}
