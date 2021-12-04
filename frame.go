/*
package stompserver provides functionality for manipulating STOMP frames.
*/
package stompserver

import (
	"fmt"
	"strings"
)

// A Frame represents a STOMP frame. A frame consists of a command
// followed by a collection of header entries, and then an optional
// body.
type Frame struct {
	Command string
	Header  *Header
	body    []byte // SetBody
}

// New creates a new STOMP frame with the specified command and headers.
// The headers should contain an even number of entries. Each even index is
// the header name, and the odd indexes are the assocated header values.
func New(command string, headers ...string) *Frame {
	f := &Frame{Command: command, Header: &Header{}}
	for index := 0; index < len(headers); index += 2 {
		f.Header.Add(headers[index], headers[index+1])
	}
	return f
}

// Clone creates a deep copy of the frame and its header. The cloned
// frame shares the body with the original frame.
func (f *Frame) Clone() *Frame {
	fc := &Frame{Command: f.Command}
	if f.Header != nil {
		fc.Header = f.Header.Clone()
	}
	if f.body != nil {
		fc.body = make([]byte, len(f.body))
		copy(fc.body, f.body)
	}
	return fc
}

func (f *Frame) Serialize() string {
	// 数据帧以 COMMAND 命令开始,以 end-of-line (EOL)结束
	frameText := f.Command + "\n"

	for i := 0; i < f.Header.Len(); i++ {
		k, v := f.Header.GetAt(i)
		frameText += k + ":" + v + "\n"
	}
	// header 头信息完毕后,下面跟一个 EOL,区分头部和消息体 Body.
	frameText += "\n"
	frameText += string(f.body)

	// 消息体 Body 后面跟一个 NULL,表示结束
	frameText += "\x00"
	return frameText
}

func (f *Frame) SetBody(body string) {
	f.body = []byte(body)
	f.Header.Add(StompHeaders.ContentLength, fmt.Sprint(len(f.body)))
}

func (f *Frame) Body() []byte {
	return f.body
}

func Deserialize(stompText string) (*Frame, error) {
	reader := NewReader(strings.NewReader(stompText))
	return reader.Read()
}

func NewErrorFrame(tips string, detail string, originFrame *Frame) *Frame {

	frame := New(CommandError)

	if tips != "" {
		frame.Header.Add(Message, tips)
	}
	frame.Header.Add(ContentType, Plain+";charset=UTF-8")
	errorbody := "Server Cannot processed This Frame\n=================\n"

	if originFrame != nil {
		errorbody += originFrame.Serialize()
		errorbody += "\n=================\n"
		if v, isContains := originFrame.Header.Contains(ReceiptId); isContains && v != "" {
			frame.Header.Add(ReceiptId, string(encodeValue(originFrame.Header.Get(ReceiptId))))
		}
	}
	errorbody += detail
	frame.body = []byte(errorbody) //append(frame.Body,errorbody)
	frame.Header.Add(ContentLength, fmt.Sprint(len(frame.body)))

	return frame
}

/*
	STOMP 1.2 servers MUST set the following headers:

	version : The version of the STOMP protocol the session will be using. See Protocol Negotiation for more details.
	STOMP 1.2 servers MAY set the following headers:

	heart-beat : The Heart-beating settings.

	session : A session identifier that uniquely identifies the session.

	server : A field that contains information about the STOMP server. The field MUST contain a server-name field and MAY be followed by optional comment fields delimited by a space character.

	The server-name field consists of a name token followed by an optional version number token.

	server = name ["/" version] *(comment)

	Example:

	server:Apache/1.3.9
*/
func NewConnectedFrame(connectFrame *Frame, session string) *Frame {
	connectedFrame := New(StompCommand.Connected)

	if connectFrame.Header.Get(StompHeaders.AcceptVersion) != "" {
		maxVersion := -1

		sVerArr := strings.Split(connectFrame.Header.Get(StompHeaders.AcceptVersion), ",")

		for i := len(sVerArr) - 1; i > -1; i-- {
			if sVerArr[i] == "1.2" {
				maxVersion = 2
				break
			} else if sVerArr[i] == "1.1" {
				if maxVersion < 2 {
					maxVersion = 1
				}
			} else if sVerArr[i] == "1.0" {
				if maxVersion < 1 {
					maxVersion = 0
				}
			}
		}

		if maxVersion > -1 {
			connectedFrame.Header.Set(StompHeaders.Version, fmt.Sprint("1.", maxVersion))
		}

		if connectedFrame.Header.Get(StompHeaders.Version) == "" {
			return NewErrorFrame("Version Unsupport", "XStomp Server Supported Protocol Versions Are 1.0,1.1,1.2", connectFrame)
		}
	} else {
		connectedFrame.Header.Set(StompHeaders.Version, "1.0")
	}

	connectedFrame.Header.Set(StompHeaders.Session, session)

	connectedFrame.Header.Set(StompHeaders.Server, "XStompServer/1.0.0")

	connectedFrame.Header.Set(StompHeaders.HeartBeat, "0,0")

	return connectedFrame
}

func NewMessageFrame(destination string, messageID string, subscription string) (*Frame, string) {
	errorText := ""
	if destination == "" {
		errorText = "destination can not be NULL/Empty." //+ "The MESSAGE frame MUST include a destination header indicating the destination the message was sent to. \nIf the message has been sent using STOMP, this destination header SHOULD be identical to the one used in the corresponding SEND frame."
		return nil, errorText
	}
	if messageID == "" {
		errorText = "message-id can not be NULL/Empty." //+ "The MESSAGE frame MUST also contain a message-id header with a unique identifier for that message and a subscription header matching the identifier of the subscription that is receiving the message."
		return nil, errorText
	}
	// 这里忽略 subscription,因为在群发的时候要重新赋值,但是一定要在发送之前添加subscription
	// if subscription == "" {
	// 	errorText = "subscription can not be NULL/Empty." //+ "订阅包里的id头必须是唯一的,sub-id将用在MESSAGE和取消订阅的操作中,否则客户端无法分辨数据来源于哪个订阅id\nSince a single connection can have multiple open subscriptions with a server, an id header MUST be included in the frame to uniquely identify the subscription. The id header allows the client and server to relate subsequent MESSAGE or UNSUBSCRIBE frames to the original subscription."
	// 	return nil, errorText
	// }
	mHeader := NewHeader()
	mHeader.Add(StompHeaders.Destination, destination)
	mHeader.Add(StompHeaders.MessageId, messageID)
	if subscription != "" {
		mHeader.Add(StompHeaders.Subscription, subscription)
	}
	mHeader.Add(StompHeaders.ContentType, Plain+";charset=UTF-8")

	// ack
	/*
	 * SpringBoot WebSocket 不支持 ack,receipts ,仅仅实现了stomp的子集.当前版本也决定暂时不实现此功能
	 * SpringBoot WebSocket Not Support ack,receipts
	 * The simple broker is great for getting started but supports only a subset of STOMP commands (e.g. no acks, receipts, etc.), relies on a simple message sending loop, and is not suitable for clustering. As an alternative, applications can upgrade to using a full-featured message broker.
	 */
	// if ackValue != "" {
	// 	mHeader.Add(StompHeaders.Ack, ackValue)
	// }

	messageFrame := &Frame{
		Command: StompCommand.Message,
		Header:  mHeader,
	}
	return messageFrame, errorText
}
