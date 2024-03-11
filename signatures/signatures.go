package signatures

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"math"
	"net/http"
	"sort"
	"strconv"
	"time"

	"github.com/0meet1/zero-framework/database"
	"github.com/0meet1/zero-framework/global"
	"github.com/0meet1/zero-framework/server"
	"github.com/0meet1/zero-framework/structs"
	"github.com/gofrs/uuid"
)

const (
	ZERO_SIGNATURE_NONCE_PREFIX = "ZERO_NONCE_"
	ZERO_SIGNATURE_NONCE_VALUE  = "cooling"

	HEADER_SIGNATURE_APP       = "ZoXappname"
	HEADER_SIGNATURE_SECRET    = "ZoXsecret"
	HEADER_SIGNATURE_NONCE     = "ZoXnonce"
	HEADER_SIGNATURE_TIMESTAMP = "ZoXtimestamp"
	HEADER_SIGNATURE_SIGN      = "ZoXsignature"

	SIGNATURE_RAW = "ZoXraw"

	SIGNATURE_OVERFLOW_TIME  = 600
	SIGNATURE_VALID_DURATION = 1200
)

type ZeroAppSecretFetcher interface {
	FetchSecret(*ZeroSignature) string
}

type ZeroSignature struct {
	xhttpreq *http.Request
	fetcher  ZeroAppSecretFetcher

	ZoXappname   string
	ZoXsecret    string
	ZoXnonce     string
	ZoXtimestamp int64
	ZoXsignature string
	params       map[string]string
}

func NewSignatureMaker(appname, secret string) *ZeroSignature {
	return &ZeroSignature{
		ZoXappname: appname,
		ZoXsecret:  secret,
		params:     make(map[string]string),
	}
}

func NewSignatureParser(xhttpreq *http.Request, fetchers ...ZeroAppSecretFetcher) (*ZeroSignature, error) {
	timestamp := int64(0)
	if len(xhttpreq.Header.Get(HEADER_SIGNATURE_TIMESTAMP)) > 0 {
		xttimestamp, err := strconv.ParseInt(xhttpreq.Header.Get(HEADER_SIGNATURE_TIMESTAMP), 10, 64)
		if err != nil {
			return nil, err
		}
		timestamp = xttimestamp
	}
	if len(fetchers) > 0 {
		return &ZeroSignature{
			xhttpreq:     xhttpreq,
			fetcher:      fetchers[0],
			ZoXappname:   xhttpreq.Header.Get(HEADER_SIGNATURE_APP),
			ZoXsecret:    xhttpreq.Header.Get(HEADER_SIGNATURE_SECRET),
			ZoXnonce:     xhttpreq.Header.Get(HEADER_SIGNATURE_NONCE),
			ZoXtimestamp: timestamp,
			ZoXsignature: xhttpreq.Header.Get(HEADER_SIGNATURE_SIGN),
			params:       make(map[string]string),
		}, nil
	} else {
		return &ZeroSignature{
			xhttpreq:     xhttpreq,
			fetcher:      nil,
			ZoXappname:   xhttpreq.Header.Get(HEADER_SIGNATURE_APP),
			ZoXsecret:    xhttpreq.Header.Get(HEADER_SIGNATURE_SECRET),
			ZoXnonce:     xhttpreq.Header.Get(HEADER_SIGNATURE_NONCE),
			ZoXtimestamp: timestamp,
			ZoXsignature: xhttpreq.Header.Get(HEADER_SIGNATURE_SIGN),
			params:       make(map[string]string),
		}, nil
	}
}

func (zox *ZeroSignature) AddParams(params map[string]string) *ZeroSignature {
	for k, v := range params {
		zox.params[k] = v
	}
	return zox
}

func (zox *ZeroSignature) AddParam(k string, v string) *ZeroSignature {
	zox.params[k] = v
	return zox
}

func (zox *ZeroSignature) AddRaw(raw []byte) *ZeroSignature {
	zox.params[SIGNATURE_RAW] = structs.Md5Bytes(raw)
	return zox
}

func (zox *ZeroSignature) AddStream(name string, stream []byte) *ZeroSignature {
	zox.params[name] = structs.Md5Bytes(stream)
	return zox
}

func (zox *ZeroSignature) signature() (string, error) {
	if zox.fetcher != nil {
		zox.ZoXsecret = zox.fetcher.FetchSecret(zox)
	}

	if len(zox.ZoXsecret) <= 0 {
		return "", errors.New(fmt.Sprintf("invalid appname: %s", zox.ZoXappname))
	}

	keys := make([]string, 0)
	keys = append(keys, HEADER_SIGNATURE_APP)
	keys = append(keys, HEADER_SIGNATURE_NONCE)
	keys = append(keys, HEADER_SIGNATURE_TIMESTAMP)

	sparams := make(map[string]string)
	sparams[HEADER_SIGNATURE_APP] = zox.ZoXappname
	sparams[HEADER_SIGNATURE_NONCE] = zox.ZoXnonce
	sparams[HEADER_SIGNATURE_TIMESTAMP] = fmt.Sprintf("%d", zox.ZoXtimestamp)

	for k, v := range zox.params {
		keys = append(keys, k)
		sparams[k] = v
	}

	paramstr := ""
	sort.Sort(sort.StringSlice(keys))
	for _, key := range keys {
		if len(paramstr) > 0 {
			paramstr = fmt.Sprintf("%s&%s=%s", paramstr, key, sparams[key])
		} else {
			paramstr = fmt.Sprintf("%s=%s", key, sparams[key])
		}
	}

	return structs.HmacSha256(paramstr, zox.ZoXsecret), nil
}

func (zox *ZeroSignature) Complete() error {
	uid, err := uuid.NewV4()
	if err != nil {
		return err
	}
	zox.ZoXnonce = uid.String()
	zox.ZoXtimestamp = time.Now().Unix()
	signstr, err := zox.signature()
	if err != nil {
		return err
	}
	zox.ZoXsignature = signstr
	return nil
}

func (zox *ZeroSignature) Check() error {
	durationTime := time.Now().Unix() - zox.ZoXtimestamp
	if durationTime > SIGNATURE_VALID_DURATION {
		return errors.New("signature expired")
	}
	if durationTime < 0 && math.Abs(float64(durationTime)) > SIGNATURE_OVERFLOW_TIME {
		return errors.New("invalid timestamp")
	}
	err := zox.checknonce()
	if err != nil {
		return err
	}
	signstr, err := zox.signature()
	if err != nil {
		return err
	}
	if signstr != zox.ZoXsignature {
		return errors.New("invalid signature")
	}
	return nil
}

func (zox *ZeroSignature) checknonce() error {
	redisKeeper := global.Value(database.DATABASE_REDIS).(database.RedisKeeper)
	xnonce, err := redisKeeper.Get(fmt.Sprintf("%s%s", ZERO_SIGNATURE_NONCE_PREFIX, zox.ZoXnonce))
	if err != nil {
		return err
	}
	if xnonce == ZERO_SIGNATURE_NONCE_VALUE {
		return errors.New("signature repeat")
	} else {
		return redisKeeper.SetEx(
			fmt.Sprintf("%s%s", ZERO_SIGNATURE_NONCE_PREFIX, zox.ZoXnonce),
			ZERO_SIGNATURE_NONCE_VALUE,
			SIGNATURE_VALID_DURATION)
	}
}

func NewSignatureRawParser(req *http.Request, fetchers ...ZeroAppSecretFetcher) (*ZeroSignature, *structs.ZeroRequest, error) {
	xParser, err := NewSignatureParser(req, fetchers...)
	if err != nil {
		return nil, nil, err
	}
	body, err := io.ReadAll(req.Body)
	if err != nil {
		return nil, nil, err
	}
	xParser.AddRaw(body)

	err = xParser.Check()
	if err != nil {
		return nil, nil, err
	}

	var request structs.ZeroRequest
	err = json.Unmarshal(body, &request)
	if err != nil {
		return nil, nil, err
	}
	return xParser, &request, nil
}

func NewSignatureStreamParser(req *http.Request, maxmem int64, fetchers ...ZeroAppSecretFetcher) (*ZeroSignature, *server.XhttpFromFile, error) {
	xParser, err := NewSignatureParser(req, fetchers...)
	if err != nil {
		return nil, nil, err
	}

	files, err := server.XhttpFromFileRequest(req, maxmem)
	if err != nil {
		return nil, nil, err
	}

	if len(files) <= 0 {
		return nil, nil, errors.New("no files in request")
	}
	if len(files) > 1 {
		return nil, nil, errors.New("too many files in request")
	}

	xParser.AddStream(files[0].FileName(), files[0].FilesBytes())

	return xParser, files[0], nil
}

func NewSignatureKeyValueParser(req *http.Request, fetchers ...ZeroAppSecretFetcher) (*ZeroSignature, map[string]string, error) {
	xParser, err := NewSignatureParser(req, fetchers...)
	if err != nil {
		return nil, nil, err
	}

	kv := server.XhttpKeyValueRequest(req)
	for k, v := range kv {
		switch k {
		case HEADER_SIGNATURE_APP:
			xParser.ZoXappname = v
		case HEADER_SIGNATURE_NONCE:
			xParser.ZoXnonce = v
		case HEADER_SIGNATURE_TIMESTAMP:
			timestamp, err := strconv.ParseInt(v, 10, 64)
			if err != nil {
				return nil, nil, err
			}
			xParser.ZoXtimestamp = timestamp
		case HEADER_SIGNATURE_SIGN:
			xParser.ZoXsecret = v
		default:
			xParser.AddParam(k, v)
		}
	}

	return xParser, kv, nil
}
