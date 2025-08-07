package protocol

import "github.com/0meet1/zero-framework/server"

const (
	ZEROKMSG_SERVER = "ZEROKMSG_SERVER"
	ZEROKMSG_CLIENT = "ZEROKMSG_CLIENT"
)

var (
	kZERO_MESSAGE_HEAD = []byte{'z', 'e', 'r', 'o'}
	kZERO_MESSAGE_END  = []byte{'Z', 'E', 'R', 'O'}
)

type ZeroKMessageServer interface {
	ExecMessage(string, *ZeroKMessage, int) (*ZeroKMessage, error)
	PushMessage(string, *ZeroKMessage) error
}

type ZeroKMessageClient interface {
	Active() bool
	ExecMessage(*ZeroKMessage, int) (*ZeroKMessage, error)
	PushMessage(*ZeroKMessage) error
}

type ZeroKMessageOperator interface {
	Operation(server.ZeroConnect, *ZeroKMessage) (bool, error)
}
