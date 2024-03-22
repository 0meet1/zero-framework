package ossminiv2

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"path"
	"strings"

	"github.com/0meet1/zero-framework/signatures"
	"github.com/0meet1/zero-framework/structs"
)

const (
	xOSSMINI_V2_STAGING  = "/zeroapi/v1/ossminiv2/store/staging/submit"
	xOSSMINI_V2_SUBMIT   = "/zeroapi/v1/ossminiv2/store/::bucketName/submit"
	xOSSMINI_V2_EXCHANGE = "/zeroapi/v1/ossminiv2/store/::bucketName/exchange"
	xOSSMINI_V2_REMOVE   = "/zeroapi/v1/ossminiv2/store/::bucketName/remove"
	xOSSMINI_V2_PATH     = "/zeroapi/v1/ossminiv2/store/::bucketName/::path"
)

type OssminiV2Client struct {
	serverAddr string
	appId      string
	appSecret  string

	stagingAppId     string
	stagingAppSecret string

	bucketName string
	useSSL     bool
}

func NewClient(serverAddr string, bucketName ...string) *OssminiV2Client {
	if len(bucketName) > 0 {
		return &OssminiV2Client{
			serverAddr: serverAddr,
			bucketName: bucketName[0],
		}
	} else {
		return &OssminiV2Client{
			serverAddr: serverAddr,
		}
	}
}

func (oclient *OssminiV2Client) Bucket(bucketName string) *OssminiV2Client {
	oclient.bucketName = bucketName
	return oclient
}

func (oclient *OssminiV2Client) StagingSecret(appId, appSecret string) *OssminiV2Client {
	oclient.stagingAppId = appId
	oclient.stagingAppSecret = appSecret
	return oclient
}

func (oclient *OssminiV2Client) Secret(appId, appSecret string) *OssminiV2Client {
	oclient.appId = appId
	oclient.appSecret = appSecret
	return oclient
}

func (oclient *OssminiV2Client) UseSSL() *OssminiV2Client {
	oclient.useSSL = true
	return oclient
}

func (oclient *OssminiV2Client) makeaddr(xxpaths ...string) string {
	if oclient.useSSL {
		return fmt.Sprintf("https://%s", path.Join(oclient.serverAddr, path.Join(xxpaths...)))
	}
	return fmt.Sprintf("http://%s", path.Join(oclient.serverAddr, path.Join(xxpaths...)))
}

func (oclient *OssminiV2Client) do(req *http.Request) ([]byte, error) {
	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()
	return io.ReadAll(resp.Body)
}

func (oclient *OssminiV2Client) exec(req *http.Request) (*structs.ZeroResponse, error) {
	body, err := oclient.do(req)
	if err != nil {
		return nil, err
	}
	var resp structs.ZeroResponse
	err = json.Unmarshal(body, &resp)
	if err != nil {
		return nil, err
	}

	if resp.Code != 200 {
		return nil, errors.New(resp.Message)
	}
	return &resp, nil
}

func (oclient *OssminiV2Client) pretreat(filename string, imageBytes []byte) (*bytes.Buffer, *multipart.Writer, error) {
	bytesBuffer := &bytes.Buffer{}
	writer := multipart.NewWriter(bytesBuffer)

	field, err := writer.CreateFormFile(filename, filename)
	if err != nil {
		return nil, nil, err
	}

	_, err = io.Copy(field, bytes.NewBuffer(imageBytes))
	if err != nil {
		return nil, nil, err
	}

	err = writer.Close()
	if err != nil {
		return nil, nil, err
	}
	return bytesBuffer, writer, nil
}

func (oclient *OssminiV2Client) xTicket(req *http.Request, filename string) (string, error) {
	resp, err := oclient.exec(req)
	if err != nil {
		return "", err
	}

	if resp.Expands == nil {
		return "", errors.New(" expands not found ")
	}
	ticket, ok := resp.Expands[filename]
	if !ok {
		return "", errors.New(" ticket not found ")
	}
	return ticket.(string), nil
}

func (oclient *OssminiV2Client) xPath(req *http.Request, key string) (string, error) {
	resp, err := oclient.exec(req)
	if err != nil {
		return "", err
	}

	if resp.Expands == nil {
		return "", errors.New(" expands not found ")
	}
	path, ok := resp.Expands[key]
	if !ok {
		return "", errors.New(" path not found ")
	}
	fmt.Println(resp.Expands)
	return path.(string), nil
}

func (oclient *OssminiV2Client) Staging(filename string, imageBytes []byte) (string, error) {
	bytesBuffer, writer, err := oclient.pretreat(filename, imageBytes)
	if err != nil {
		return "", err
	}

	req, err := http.NewRequest("POST", oclient.makeaddr(xOSSMINI_V2_STAGING), bytesBuffer)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	if len(oclient.stagingAppId) > 0 && len(oclient.stagingAppSecret) > 0 {
		xSignature := signatures.NewSignatureMaker(oclient.stagingAppId, oclient.stagingAppSecret)
		xSignature.AddStream(filename, imageBytes)
		err = xSignature.Complete()
		if err != nil {
			return "", err
		}
		req.Header.Set(signatures.HEADER_SIGNATURE_APP, xSignature.ZoXappname)
		req.Header.Set(signatures.HEADER_SIGNATURE_TIMESTAMP, fmt.Sprintf("%d", xSignature.ZoXtimestamp))
		req.Header.Set(signatures.HEADER_SIGNATURE_NONCE, xSignature.ZoXnonce)
		req.Header.Set(signatures.HEADER_SIGNATURE_SIGN, xSignature.ZoXsignature)
	}
	return oclient.xTicket(req, filename)
}

func (oclient *OssminiV2Client) Submit(filename string, imageBytes []byte) (string, error) {
	bytesBuffer, writer, err := oclient.pretreat(filename, imageBytes)
	if err != nil {
		return "", err
	}

	xURI := oclient.makeaddr(strings.ReplaceAll(xOSSMINI_V2_SUBMIT, "::bucketName", oclient.bucketName))
	req, err := http.NewRequest("POST", xURI, bytesBuffer)
	if err != nil {
		return "", err
	}
	req.Header.Set("Content-Type", writer.FormDataContentType())

	if len(oclient.appId) > 0 && len(oclient.appSecret) > 0 {
		xSignature := signatures.NewSignatureMaker(oclient.appId, oclient.appSecret)
		xSignature.AddStream(filename, imageBytes)
		err = xSignature.Complete()
		if err != nil {
			return "", err
		}
		req.Header.Set(signatures.HEADER_SIGNATURE_APP, xSignature.ZoXappname)
		req.Header.Set(signatures.HEADER_SIGNATURE_TIMESTAMP, fmt.Sprintf("%d", xSignature.ZoXtimestamp))
		req.Header.Set(signatures.HEADER_SIGNATURE_NONCE, xSignature.ZoXnonce)
		req.Header.Set(signatures.HEADER_SIGNATURE_SIGN, xSignature.ZoXsignature)
	}

	return oclient.xPath(req, filename)
}

func (oclient *OssminiV2Client) Exchange(ticket string) (string, error) {
	xRequest := &structs.ZeroRequest{
		Expands: map[string]interface{}{
			"ticket": ticket,
		},
	}

	jsonBytes, err := json.Marshal(xRequest)
	if err != nil {
		return "", err
	}

	xURI := oclient.makeaddr(strings.ReplaceAll(xOSSMINI_V2_EXCHANGE, "::bucketName", oclient.bucketName))
	request, err := http.NewRequest("POST", xURI, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return "", err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	return oclient.xPath(request, ticket)
}

func (oclient *OssminiV2Client) Remove(srcpath string) error {
	xRequest := &structs.ZeroRequest{
		Expands: map[string]interface{}{
			"path": srcpath,
		},
	}

	jsonBytes, err := json.Marshal(xRequest)
	if err != nil {
		return err
	}

	xURI := oclient.makeaddr(strings.ReplaceAll(xOSSMINI_V2_REMOVE, "::bucketName", oclient.bucketName))
	request, err := http.NewRequest("POST", xURI, bytes.NewBuffer(jsonBytes))
	if err != nil {
		return err
	}
	request.Header.Set("Content-Type", "application/json; charset=UTF-8")

	if len(oclient.appId) > 0 && len(oclient.appSecret) > 0 {
		xSignature := signatures.NewSignatureMaker(oclient.appId, oclient.appSecret)
		xSignature.AddParam("path", fmt.Sprintf("%s/%s", oclient.bucketName, srcpath))
		err = xSignature.Complete()
		if err != nil {
			return err
		}
		request.Header.Set(signatures.HEADER_SIGNATURE_APP, xSignature.ZoXappname)
		request.Header.Set(signatures.HEADER_SIGNATURE_TIMESTAMP, fmt.Sprintf("%d", xSignature.ZoXtimestamp))
		request.Header.Set(signatures.HEADER_SIGNATURE_NONCE, xSignature.ZoXnonce)
		request.Header.Set(signatures.HEADER_SIGNATURE_SIGN, xSignature.ZoXsignature)
	}

	_, err = oclient.exec(request)
	return err
}

func (oclient *OssminiV2Client) Complete(srcpath string) (string, error) {
	xURI := strings.ReplaceAll(xOSSMINI_V2_PATH, "::bucketName", oclient.bucketName)
	if len(oclient.appId) > 0 && len(oclient.appSecret) > 0 {
		xSignature := signatures.NewSignatureMaker(oclient.appId, oclient.appSecret)
		xSignature.AddParam("path", fmt.Sprintf("%s/%s", oclient.bucketName, srcpath))
		err := xSignature.Complete()
		if err != nil {
			return "", err
		}
		xQueryPath := fmt.Sprintf("%s?%s=%s&%s=%d&%s=%s&%s=%s",
			srcpath,
			signatures.HEADER_SIGNATURE_APP, xSignature.ZoXappname,
			signatures.HEADER_SIGNATURE_TIMESTAMP, xSignature.ZoXtimestamp,
			signatures.HEADER_SIGNATURE_NONCE, xSignature.ZoXnonce,
			signatures.HEADER_SIGNATURE_SIGN, xSignature.ZoXsignature)
		return strings.ReplaceAll(xURI, "::path", xQueryPath), nil
	}
	return strings.ReplaceAll(xURI, "::path", srcpath), nil
}

func (oclient *OssminiV2Client) Fetch(srcpath string) ([]byte, error) {
	xURI, err := oclient.Complete(srcpath)
	if err != nil {
		return nil, err
	}
	request, err := http.NewRequest("GET", oclient.makeaddr(xURI), nil)
	if err != nil {
		return nil, err
	}
	return oclient.do(request)
}
