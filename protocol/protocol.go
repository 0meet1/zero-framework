package protocol

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

type ZeroV1MessageOperator interface {
	Operation(*ZeroV1Message) (bool, error)
}
