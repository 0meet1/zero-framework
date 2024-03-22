package ossminiv2

import "github.com/0meet1/zero-framework/global"

type OssminiV2Keeper interface {
	Staging(string, []byte) (string, error)
	Submit(string, []byte) (string, error)
	Exchange(string) (string, error)
	Complete(string) (string, error)
	Fetch(string) ([]byte, error)
}

type xOssminiV2Keeper struct {
	serverAddr string
	appId      string
	appSecret  string

	stagingAppId     string
	stagingAppSecret string

	bucketName string
	useSSL     bool
}

func NewKeeper(serverAddr string, bucketName ...string) *xOssminiV2Keeper {
	if len(bucketName) > 0 {
		return &xOssminiV2Keeper{
			serverAddr: serverAddr,
			bucketName: bucketName[0],
		}
	} else {
		return &xOssminiV2Keeper{
			serverAddr: serverAddr,
		}
	}
}

func (oKeeper *xOssminiV2Keeper) Bucket(bucketName string) *xOssminiV2Keeper {
	oKeeper.bucketName = bucketName
	return oKeeper
}

func (oKeeper *xOssminiV2Keeper) StagingSecret(appId, appSecret string) *xOssminiV2Keeper {
	oKeeper.stagingAppId = appId
	oKeeper.stagingAppSecret = appSecret
	return oKeeper
}

func (oKeeper *xOssminiV2Keeper) Secret(appId, appSecret string) *xOssminiV2Keeper {
	oKeeper.appId = appId
	oKeeper.appSecret = appSecret
	return oKeeper
}

func (oKeeper *xOssminiV2Keeper) UseSSL() *xOssminiV2Keeper {
	oKeeper.useSSL = true
	return oKeeper
}

func (oKeeper *xOssminiV2Keeper) Mount(registerName string) {
	global.Key(registerName, oKeeper)
}

func (oKeeper *xOssminiV2Keeper) Staging(filename string, imageBytes []byte) (string, error) {
	client := NewClient(oKeeper.serverAddr, oKeeper.bucketName).
		Secret(oKeeper.appId, oKeeper.appSecret).
		StagingSecret(oKeeper.stagingAppId, oKeeper.stagingAppSecret)
	if oKeeper.useSSL {
		client.UseSSL()
	}
	return client.Staging(filename, imageBytes)
}

func (oclient *xOssminiV2Keeper) Submit(filename string, imageBytes []byte) (string, error) {
	client := NewClient(oclient.serverAddr, oclient.bucketName).
		Secret(oclient.appId, oclient.appSecret).
		StagingSecret(oclient.stagingAppId, oclient.stagingAppSecret)
	if oclient.useSSL {
		client.UseSSL()
	}
	return client.Submit(filename, imageBytes)
}

func (oclient *xOssminiV2Keeper) Exchange(ticket string) (string, error) {
	client := NewClient(oclient.serverAddr, oclient.bucketName).
		Secret(oclient.appId, oclient.appSecret).
		StagingSecret(oclient.stagingAppId, oclient.stagingAppSecret)
	if oclient.useSSL {
		client.UseSSL()
	}
	return client.Exchange(ticket)
}

func (oclient *xOssminiV2Keeper) Complete(srcpath string) (string, error) {
	client := NewClient(oclient.serverAddr, oclient.bucketName).
		Secret(oclient.appId, oclient.appSecret).
		StagingSecret(oclient.stagingAppId, oclient.stagingAppSecret)
	if oclient.useSSL {
		client.UseSSL()
	}
	return client.Complete(srcpath)
}

func (oclient *xOssminiV2Keeper) Fetch(srcpath string) ([]byte, error) {
	client := NewClient(oclient.serverAddr, oclient.bucketName).
		Secret(oclient.appId, oclient.appSecret).
		StagingSecret(oclient.stagingAppId, oclient.stagingAppSecret)
	if oclient.useSSL {
		client.UseSSL()
	}
	return client.Fetch(srcpath)
}
