package structs

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/hex"
	"fmt"

	"github.com/teris-io/shortid"
)

var Md5Bytes = func(b []byte) string {
	w := md5.New()
	w.Write(b)
	md5str := fmt.Sprintf("%x", w.Sum(nil))
	return md5str
}

var Md5 = func(s string) string {
	return Md5Bytes([]byte(s))
}

var HmacSha256 = func(data string, secret string) string {
	h := hmac.New(sha256.New, []byte(secret))
	h.Write([]byte(data))
	return hex.EncodeToString(h.Sum(nil))
}

const ALPHABET_X1 = "0123456789abcdefghijklmnopqrstuvwxyz@|ABCDEFGHIJKLMNOPQRSTUVWXYZ"
const ALPHABET_X2 = "abcdefghijklmnopqrstuvwxyz@|ABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"

var GenerateSUID = func() (string, error) {
	generaterf, err := shortid.New(1, ALPHABET_X1, 2508)
	if err != nil {
		return "", err
	}

	siuf, err := generaterf.Generate()
	if err != nil {
		return "", err
	}

	generaters, err := shortid.New(2, ALPHABET_X2, 818)
	if err != nil {
		return "", err
	}

	sius, err := generaters.Generate()
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%s-%s", siuf, sius), nil
}
