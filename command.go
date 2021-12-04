package stompserver

// STOMP frame commands. Used upper case naming
// convention to avoid clashing with STOMP header names.
const (
	// Client commands.
	CommandAbort       = "ABORT"
	CommandACK         = "ACK"
	CommandBegin       = "BEGIN"
	CommandCommit      = "COMMIT"
	CommandConnect     = "CONNECT"
	CommandDisconnect  = "DISCONNECT"
	CommandNack        = "NACK"
	CommandSend        = "SEND"
	CommandStomp       = "STOMP"
	CommandSubscribe   = "SUBSCRIBE"
	CommandUnsubscribe = "UNSUBSCRIBE"

	// Server commands.
	CommandConnected = "CONNECTED"
	CommandError     = "ERROR"
	CommandMessage   = "MESSAGE"
	CommandHeartbeat = "RECEIPT"
)

type command struct {
	// Client commands.
	Abort       string
	ACK         string
	Begin       string
	Commit      string
	Connect     string
	Disconnect  string
	Nack        string
	Send        string
	Stomp       string
	Subscribe   string
	Unsubscribe string

	// Server commands.
	Connected string
	Error     string
	Message   string
	Heartbeat string
}

var StompCommand = &command{
	// Client commands.
	Abort:       "ABORT",
	ACK:         "ACK",
	Begin:       "BEGIN",
	Commit:      "COMMIT",
	Connect:     "CONNECT",
	Disconnect:  "DISCONNECT",
	Nack:        "NACK",
	Send:        "SEND",
	Stomp:       "STOMP",
	Subscribe:   "SUBSCRIBE",
	Unsubscribe: "UNSUBSCRIBE",

	// Server commands.
	Connected: "CONNECTED",
	Error:     "ERROR",
	Message:   "MESSAGE",
	Heartbeat: "RECEIPT",
}

// IsClientCommand  验证是否是客户端的命令
//  参数:
//  command	string	包命令
//  return	bool	true 客户端的包 false 异常包(应该发送一个Error包,然后关闭连接)
//        	error	nil / 无效包错误
func IsClientCommand(command string) (bool, error) {
	switch command {
	case CommandAbort, CommandACK, CommandBegin, CommandCommit, CommandConnect, CommandDisconnect, CommandNack, CommandSend, CommandStomp, CommandSubscribe, CommandUnsubscribe:
		return true, nil
	default:
		return false, ErrInvalidCommand
	}
}

// IsServerCommand  是否是server的包
//  参数:
//  command	string	包命令
//  return	bool	true 服务端的包 false 异常包(应该发送一个Error包,然后关闭连接)
//        	error	nil / 无效包错误
func IsServerCommand(command string) (bool, error) {
	switch command {
	case CommandConnected, CommandError, CommandMessage, CommandHeartbeat:
		{
			return true, nil
		}
	default:
		{
			return false, ErrInvalidCommand
		}
	}
}
