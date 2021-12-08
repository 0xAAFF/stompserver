package stompserver

import (
	"fmt"
	"sync"
	"time"

	"github.com/0xAAFF/stompserver/tools"

	"golang.org/x/net/websocket"
)

type StompUnit struct {
	xstomp                    *websocket.Conn                                  // 连接本身,底层是websocket连接
	address                   string                                           // ip:port 地址
	SessionId                 string                                           // 使用UUID分配id
	TopicSubidDictionary      map[string]string                                // 按照订阅地址关联Stomp客户端连接,key:Destination  Value:Sub-id
	TopicSubidDictionaryMutex *sync.Mutex                                      // 锁
	IsConnected               bool                                             // 是否建立stomp连接,更改前,需要建立connect
	PublishManager            *StompManager                                    // 主服务的管理,类似服务的管家
	isUser                    bool                                             // 是否是用户组成员
	reflexHandle              func(sourceStompMessage *Frame, unit *StompUnit) // 正式数据处理函数
	sendStompQueue            chan Frame                                       // 发送信道
	closeOK                   bool
	ControllerUnitId          string // 登录时使用,用来保持连接的

	//xstompRunnerChannel  chan int                                         // 当前stomp连接状态,如果中断,则通知主服务,移除本连接
	/* not support
	stompVersion string
	heartbeat string
	*/

}

//NestompserverUnit 创建一个Tcp连接单元
func NewStompUnit(connection *websocket.Conn, publishManager *StompManager, reflexHandle func(sourceStompMessage *Frame, unit *StompUnit)) StompUnit {
	return StompUnit{
		xstomp:                    connection,
		address:                   connection.Request().RemoteAddr,
		SessionId:                 tools.TimeUUID().String(), // 分配UUID
		TopicSubidDictionary:      make(map[string]string),
		TopicSubidDictionaryMutex: new(sync.Mutex),
		IsConnected:               false,
		PublishManager:            publishManager,
		isUser:                    false,
		reflexHandle:              reflexHandle,
		sendStompQueue:            make(chan Frame),
		closeOK:                   false,
	}
}

//Run 执行工作
func (stompUnit *StompUnit) Run() {
	defer stompUnit.onClose()

	go stompUnit.send()
	stompUnit.recevice()
}

func (stompUnit *StompUnit) send() {
	for stompMessage := range stompUnit.sendStompQueue {
		// Ilog.Debug(stompMessage.Serialize())
		if stompUnit.closeOK {
			return
		}
		_, err := stompUnit.xstomp.Write([]byte(stompMessage.Serialize()))
		if err != nil {
			return
		}
	}
}

func (stompUnit *StompUnit) recevice() {
	for {
		if stompUnit.closeOK {
			return
		}
		var reply string
		err := websocket.Message.Receive(stompUnit.xstomp, &reply)
		if err != nil {
			return
		}
		go func() {
			frameMessage, err := Deserialize(reply)
			if err != nil {
				// Ilog.Debug("StompUnit Run Read:Error Message (数据格式错误)=>", reply)
				return
			}

			// Ilog.Debug(frameMessage.Serialize())
			stompUnit.onStompProtocol(frameMessage)
		}()
	}

}

func (stompUnit *StompUnit) onStompProtocol(stomp *Frame) {
	switch stomp.Command {
	case StompCommand.Connect:
		{
			/*
			 * 这里需要做一些协议支持处理 http://stomp.github.io/stomp-specification-1.1.html#Protocol_Negotiation
			 *
			 * 可接受的协议 Done.
			 * host验证
			 *
			 * login
			 * passcode
			 */
			/*
			 * TODO secured STOMP server
			 */
			connectedFrame := NewConnectedFrame(stomp, stompUnit.SessionId)
			stompUnit.SendStompMessage(connectedFrame)

			if connectedFrame.Command != StompCommand.Error {
				stompUnit.IsConnected = true
			}
			break
		}
	case StompCommand.Subscribe:
		{
			errorFrame := stompUnit.PublishManager.SubScribe(stomp, stompUnit)
			if errorFrame != nil {
				stompUnit.SendStompMessage(errorFrame)
			}
			stompUnit.reflexHandle(stomp, stompUnit)
			break
		}
	case StompCommand.Unsubscribe:
		{
			stompUnit.PublishManager.UnSubScribe(stomp, stompUnit)
			break
		}
	case StompCommand.Send:
		{
			stompUnit.reflexHandle(stomp, stompUnit)
			break
		}
	default:
		{
			stompFrame := NewErrorFrame("Stomp Server Can not Support This Command", "StompServer Only Support Connect/Subscribe/Unsubscribe/Send", stomp)
			//
			stompUnit.SendStompMessage(stompFrame)
		}
	}
}

func (stompUnit *StompUnit) SendStompMessage(stompMessageFrame *Frame) {
	if ok, error := IsServerCommand(stompMessageFrame.Command); ok {
		stompFrame := *stompMessageFrame
		// 针对MessageFrame的Subscription:subid 做处理
		if stompFrame.Command == CommandMessage && !stompFrame.Header.ContainsKey(StompHeaders.Subscription) {
			frameDestination := stompFrame.Header.Get(StompHeaders.Destination)
			if frameDestination == "" {
				return
			}
			stompUnit.TopicSubidDictionaryMutex.Lock()
			{
				if subid, ok := stompUnit.TopicSubidDictionary[frameDestination]; ok {
					stompFrame.Header.Del(StompHeaders.Subscription)
					stompFrame.Header.Add(StompHeaders.Subscription, subid)
				}
			}
			stompUnit.TopicSubidDictionaryMutex.Unlock()
		}
		// Ilog.Debug("Is Server Command")
		stompUnit.sendStompQueue <- stompFrame
	} else {
		fmt.Println("Unit SendStompMessage :" + error.Error())
		// Ilog.D(DebugLevel.Error, "Unit SendStompMessage :"+error.Error())
	}
	// If the server cannot successfully process the SEND frame for any reason, the server MUST send the client an ERROR frame and then close the connection.
	if stompMessageFrame.Command == StompCommand.Error {
		//Ilog.DD(DebugLevel.Error, fmt.Sprint("Stomp Send ERROR Frame 后关闭:", stompMessageFrame.Serialize()))
		stompUnit.onClose()
	}
}

func (stompUnit *StompUnit) onClose() {
	if stompUnit.closeOK {
		return
	}
	stompUnit.closeOK = true
	// 连接移除
	stompUnit.PublishManager.StompUnitOnClose(stompUnit)
	// 休眠300毫秒
	// server will likely only allow closed connections to linger for short time before the connection is reset.
	time.Sleep(time.Duration(300) * time.Millisecond)
	// 关闭连接
	stompUnit.xstomp.Close()
}

func (stompUnit *StompUnit) IsUser() bool {
	return stompUnit.isUser
}
