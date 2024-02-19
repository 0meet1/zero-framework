package protocol

import "github.com/0meet1/zero-framework/server"

const (
	ZEROV1SERV_KEEPER = "ZEROV1SERV_KEEPER"
	ZEROV1SERV_CLIENT = "ZEROV1SERV_CLIENT"
)

var (
	xZERO_MESSAGE_HEAD = []byte{'z', 'e', 'r', 'o'}
	xZERO_MESSAGE_END  = []byte{'Z', 'E', 'R', 'O'}
)

type ZeroV1ServKeeper interface {
	ExecMessage(string, *ZeroV1Message, int) (*ZeroV1Message, error)
	PushMessage(string, *ZeroV1Message) error
}

type ZeroV1Client interface {
	Active() bool
	ExecMessage(*ZeroV1Message, int) (*ZeroV1Message, error)
	PushMessage(*ZeroV1Message) error
}

type ZeroV1MessageOperator interface {
	Operation(server.ZeroConnect, *ZeroV1Message) (bool, error)
}
