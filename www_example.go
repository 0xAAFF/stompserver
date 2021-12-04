package stompserver

import (
	"fmt"
	"net/http"

	"golang.org/x/net/websocket"
)

/*
steps
*/

// those funcs not just in one file

func WebServerStart() {
	go initServer()
}

//
// start stompserver
// stompserver based on "golang.org/x/net/websocket"
//
func initServer() {

	port := 80

	//http.HandleFunc("/", HttpVueServer)                                        // Web服务
	//fmt.Println("Web Server  : http://127.0.0.1:", 80, "/")

	http.Handle("/stomp", websocket.Handler(StompServerInstance.NewStompUnit)) // Stomp服务
	fmt.Println("Stomp Server: ws://localhost:", 80, "/stomp")

	err := http.ListenAndServe(fmt.Sprintf("0.0.0.0:%d", port), nil) //only local  so client origin:127.0.0.1

	if err != nil {
		fmt.Println("init Serve: " + err.Error())
		panic("initServer: " + err.Error())
	}
}

// 1 注册Stomp 关联的地址
func Regist() {
	/*
		like java config broker
		@Override
		public void configureMessageBroker(MessageBrokerRegistry config) {
			config.enableSimpleBroker("/queue", "/topic");
			config.setApplicationDestinationPrefixes("/app");
			config.setUserDestinationPrefix("/user");
		}
		这里需要优先规范stop的地址的开头的表示信息
		RegistDestinationToBroadcastAll 表示以设定地址开头的地址 都将会被群发 例如设定了 "/broadcast",那么 "/broadcast/path","/broadcast/path/a","/broadcast/"的目的地址上的消息都会被群发
		RegistDestinationToApplication  表示单发,哪个客户端发送的数据,则返回数据就发送给谁
		RegistDestinationToUser			表示组发,例如登录成功的客户端才能收到数据,未登录的客户端则不会收到数据
	*/
	InstancesStompManager.RegistDestinationToBroadcastAll("/queue", "/topic", "/broadcast") //  message will to All client when (destination) start with "/queue" or "/topic" or "/broadcast"
	InstancesStompManager.RegistDestinationToApplication("/application", "/applications/")  //  message will to the sender client when (destination) start with "/application" or "/applications/"
	InstancesStompManager.RegistDestinationToUser("/user")                                  //  message will to the group client when (destination) start with "/user"
}

// 2 提供订阅地址
func AddSubscribeInterface() {
	// Stomp的订阅地址
	vlist := make([]string, 0)
	for _, v := range SubscribeMap {
		vlist = append(vlist, v)
	}
	/*
	* client will subscribe the paths/destination
	 */
	InstancesStompManager.AddSubscribeDestination(vlist...)
}

// 客户端信息返回的订阅地址  stomp_destinations.go
var SubscribeMap map[string]string = map[string]string{
	"Whatamidoing": "/broadcast/whatamidoing",
	"Whoami":       "/broadcast/whoami",
	"Whereami":     "/broadcast/whereami",
}

// 3 websocket

// 新建立的stomp连接
// #region 区域折叠
type StompServer struct {
	IStompManager StompManager // 所有管理相关都放入此处
	ReflexHandle  func(sourceStompMessage *Frame, unit *StompUnit)
	// or more ...
}

// 新建立的stomp连接
func (stompServer *StompServer) NewStompUnit(ws *websocket.Conn) {
	stompUnit := NewStompUnit(ws, &StompServerInstance.IStompManager, stompServer.ReflexHandle)
	stompUnit.Run()
}

var StompServerInstance = &StompServer{
	IStompManager: *InstancesStompManager,
}

// 4 关联响应接口的函数
func SetReflex(reflex func(sourceStompMessage *Frame, unit *StompUnit)) {
	var surfaceServer *StompServer = StompServerInstance // 基于前端的服务,接受和发送给前端处理的服务
	surfaceServer.ReflexHandle = reflex
}

// 5 ReflexHandle reflex_stomp.go
/*
针对前端的Stomp的反射函数将在这里实现
Stomp是根据Destination的地址来进行接口的数据解析的
这里着重实现func Reflex()函数
*/

// 客户端针对不同的地址,做出对应的反应 reflex_stomp.go
func Reflex(sourceStompMessage *Frame, unit *StompUnit) {
	/*
	* 如下 是示例代码 供测试(需要配合相应的Web项目)
	 */

	// 此行代码勿删,防止连接未建立,客户端非法提交Send数据包
	if !unit.IsConnected {
		return
	}
	switch sourceStompMessage.Header.Get(StompHeaders.Destination) {
	case AcceptInterface["AskWhatamidoing"]:
		{
			go Whatamidoing(sourceStompMessage, unit)
			break
		}
	case AcceptInterface["AskWhereami"]:
		{
			// go Whereami(sourceStompMessage, unit)
			break
		}
	case AcceptInterface["AskWhoami"]:
		{
			//go Whoami(sourceStompMessage, unit)
			break
		}
		// ... or more
	}

}

// 提供给客户端访问的地址 stomp_destinations.go
var AcceptInterface map[string]string = map[string]string{
	"AskWhatamidoing": "/whatamidoing",
	"AskWhoami":       "/whoami",
	"AskWhereami":     "/whereami",
	// or more ...
}

// reflex_stomp.go
func Whatamidoing(sourceStompMessage *Frame, unit *StompUnit) {

	var subid string
	unit.TopicSubidDictionaryMutex.Lock()
	{
		subid = unit.TopicSubidDictionary[SubscribeMap["Whatamidoing"]]
	}
	unit.TopicSubidDictionaryMutex.Unlock()

	messageFrame, errtxt := NewMessageFrame(SubscribeMap["Whatamidoing"], StompServerInstance.IStompManager.NewMessageId(), subid)
	if errtxt != "" {
		messageFrame = NewErrorFrame("AskWhatamidoing Error", errtxt, sourceStompMessage)
	} else {
		var jsonText = `{"name:"wstomp"}` // Whatamidoing.Serialize()
		messageFrame.SetBody(jsonText)
	}
	unit.SendStompMessage(messageFrame)
}
