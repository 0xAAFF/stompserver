package stompserver

import (
	"strconv"
	"strings"
	"sync"

	"github.com/0xAAFF/stompserver/tools"
)

type StompManager struct {
	// 地址配置
	rootMux                           *sync.Mutex                     // (广播/单播/用户组播)根地址的锁
	rootDestinationSliceToAll         []string                        // 广播根地址,切片
	rootDestinationSliceToApplication []string                        // 单播根地址,切片
	rootDestinationSliceToUser        []string                        // 用户组播根地址,切片
	guid                              string                          // MessageID=guid+"-"+id
	giudMux                           *sync.Mutex                     // MessageID giudMux
	messageTailId                     uint64                          // MessageID=guid+"-"+id
	subscribeSupportTopicSlice        []string                        //支持的订阅地址列表
	stompSubscribeDictionary          map[string]map[string]StompUnit //订阅地址和相应的连接单元 : <订阅地址,<sessionid,连接>>
	stompSubscribeDictionaryMux       *sync.Mutex
}

const (
	Unknow = -1
	All    = iota
	Application
	User
)

var InstancesStompManager = &StompManager{
	rootMux:                           new(sync.Mutex),
	rootDestinationSliceToAll:         []string{},
	rootDestinationSliceToApplication: []string{},
	rootDestinationSliceToUser:        []string{},
	guid:                              tools.TimeUUID().String(),
	giudMux:                           new(sync.Mutex),
	messageTailId:                     uint64(0),
	subscribeSupportTopicSlice:        []string{},
	stompSubscribeDictionary:          make(map[string]map[string]StompUnit),
	stompSubscribeDictionaryMux:       new(sync.Mutex),
}

func (manager *StompManager) NewMessageId() string {
	defer func() { manager.messageTailId++ }()
	if manager.messageTailId > 0xFFFFFFFF {
		manager.messageTailId = 0
	}
	return manager.guid + "-" + strconv.FormatUint(manager.messageTailId, 10)
}

// 群发路径注册
func (manager *StompManager) RegistDestinationToBroadcastAll(destinationRoots ...string) {
	manager.rootMux.Lock()
	defer manager.rootMux.Unlock()

	for _, d := range destinationRoots {
		if !strings.HasSuffix(d, "/") {
			d += "/"
		}
		if !tools.SliceContains(manager.rootDestinationSliceToAll, d) {
			manager.rootDestinationSliceToAll = append(manager.rootDestinationSliceToAll, d)
		}
	}
}

func (manager *StompManager) RegistDestinationToApplication(destinationRoots ...string) {
	manager.rootMux.Lock()
	defer manager.rootMux.Unlock()

	for _, d := range destinationRoots {
		if !strings.HasSuffix(d, "/") {
			d += "/"
		}
		if !tools.SliceContains(manager.rootDestinationSliceToApplication, d) {
			manager.rootDestinationSliceToApplication = append(manager.rootDestinationSliceToApplication, d)
		}
	}
}

func (manager *StompManager) RegistDestinationToUser(destinationRoots ...string) {
	manager.rootMux.Lock()
	defer manager.rootMux.Unlock()

	for _, d := range destinationRoots {
		if !strings.HasSuffix(d, "/") {
			d += "/"
		}
		if !tools.SliceContains(manager.rootDestinationSliceToUser, d) {
			manager.rootDestinationSliceToUser = append(manager.rootDestinationSliceToUser, d)
		}
	}
}

func (manager *StompManager) SubscribeFlowDirction(destination string) int {
	index := strings.Index(destination[1:], "/")
	if index > -1 {
		destination = destination[:index+2]
	}
	if tools.SliceContains(manager.rootDestinationSliceToAll, destination) {
		return All
	} else if tools.SliceContains(manager.rootDestinationSliceToApplication, destination) {
		return Application
	} else if tools.SliceContains(manager.rootDestinationSliceToUser, destination) {
		return User
	} else {
		return Unknow
	}
}

func (manager *StompManager) AddSubscribeDestination(destinations ...string) {
	manager.stompSubscribeDictionaryMux.Lock()
	defer manager.stompSubscribeDictionaryMux.Unlock()

	for _, d := range destinations {
		if !tools.SliceContains(manager.subscribeSupportTopicSlice, d) {
			manager.subscribeSupportTopicSlice = append(manager.subscribeSupportTopicSlice, d)
		}
	}
}

//  SubScribe  添加一个新的订阅地址
//  参数:
//  	subScribeFrame	客户端的订阅数据包
//  	stompUnit		客户端连接对象
//  return	nil:Success, not NULL:ERROR
func (manager *StompManager) SubScribe(subScribeFrame *Frame, stompUnit *StompUnit) *Frame {
	manager.stompSubscribeDictionaryMux.Lock()
	defer manager.stompSubscribeDictionaryMux.Unlock()
	desinationValue := subScribeFrame.Header.Get(StompHeaders.Destination) // Frame 中destination的值
	subidValue := subScribeFrame.Header.Get(StompHeaders.Id)               // Frame 中sub-id的值

	// If the server cannot successfully create the subscription,
	// the server MUST send the client an ERROR frame and disconnect the client.
	if !subScribeFrame.Header.ContainsKey(StompHeaders.Destination) {
		return NewErrorFrame("Headers missing:'destination'", "SUBSCRIBE 包必携带'destination'头,而且值不能为空.\nSUBSCRIBE Frame's headers MUST Contain 'destination',and the value MUST NOT null/empty", subScribeFrame)
	} else if desinationValue == "" {
		return NewErrorFrame("Header['destination'] MUST NOT be null or empty", "SUBSCRIBE 包规定header['destination']的值不能为空.\nSUBSCRIBE Frame's header['destination']: MUST NOT be null/empty", subScribeFrame)
	} else if !tools.SliceContains(manager.subscribeSupportTopicSlice, desinationValue) {
		return NewErrorFrame("Stomp Server not support subscribe this destination", "Stomp Server not support subscribe this destination : '"+desinationValue+"\n", subScribeFrame)
	} else if !subScribeFrame.Header.ContainsKey(StompHeaders.Id) {
		return NewErrorFrame("Headers missing: 'id'", "SUBSCRIBE 包必须携带'id'头,而且值不能为空.\nSUBSCRIBE Frame's headers MUST Contain 'id',and the value MUST uniquely"+"\nSince a single connection can have multiple open subscriptions with a server, an id header MUST be included in the frame to uniquely identify the subscription.", subScribeFrame)
	} else if subidValue == "" {
		return NewErrorFrame("Header['id'] MUST NOT be null or empty,and MUST uniquely.", "SUBSCRIBE 包必须携带'id'头,而且值必须唯一.\nSUBSCRIBE Frame's header['id'] MUST NOT null/empty,and MUST uniquely.", subScribeFrame)
	}

	// Check sub-id
	if tools.MapContainsValue(stompUnit.TopicSubidDictionary, subidValue) { // stomp单元已经使用该sub-id订阅了toppic地址
		if tools.MapContainsKey(stompUnit.TopicSubidDictionary, desinationValue) && stompUnit.TopicSubidDictionary[desinationValue] == subidValue {
			return nil
		} else {
			destinationOfsubid := tools.MapGetFirstKeyByValue(stompUnit.TopicSubidDictionary, desinationValue)
			return NewErrorFrame("SUBSCRIBE: ['id'] MUST uniquely", "this sub-id is alread been used with this destination:"+destinationOfsubid+".\n sub-id MUST uniquely,please use another id", subScribeFrame)
		}
	}

	if tools.MapContainsKey(stompUnit.TopicSubidDictionary, desinationValue) {
		if stompUnit.TopicSubidDictionary[desinationValue] == subidValue { // 已经订阅过,id都匹配
			return nil
		} else { //换了一个sub-id
			// destination已经匹配了一个sub-id
			//return new ErrorFrame("SUBSCRIBE: ['destination'] already linked a sub-id", "this destination is alread linked a sub-id:" + stompBehavior.topicPathSubidDictionary[subScribeFrame.Headers[StompHeaders.Destination]] + ".\n Stomp Server Don't want client Subscribe one destination by many times.", subScribeFrame);
			stompUnit.TopicSubidDictionary[desinationValue] = subidValue // update sub-id
		}
	} else {
		stompUnit.TopicSubidDictionary[desinationValue] = subidValue
	}

	//
	// ack ... Not support!
	//
	//
	// Stomp1.0
	// The body of the SUBSCRIBE command is ignored.
	//

	if stompUnitMapContainsKey(manager.stompSubscribeDictionary, desinationValue) {
		idUnitKV := manager.stompSubscribeDictionary[desinationValue]
		if !sessionIdStompUnitMapContainsKey(idUnitKV, stompUnit.SessionId) {
			manager.stompSubscribeDictionary[desinationValue][stompUnit.SessionId] = *stompUnit
		}
	} else {
		idUnitMap := make(map[string]StompUnit)
		idUnitMap[stompUnit.SessionId] = *stompUnit
		manager.stompSubscribeDictionary[desinationValue] = idUnitMap
	}
	return nil
}

func stompUnitMapContainsKey(imap map[string]map[string]StompUnit, key string) bool {
	for k := range imap {
		if k == key {
			return true
		}
	}
	return false
}

func sessionIdStompUnitMapContainsKey(imap map[string]StompUnit, key string) bool {
	for k := range imap {
		if k == key {
			return true
		}
	}
	return false
}

/// 客户端取消订阅
func (manager *StompManager) UnSubScribe(unSubScribeFrame *Frame, stompUnit *StompUnit) {
	manager.stompSubscribeDictionaryMux.Lock()
	defer manager.stompSubscribeDictionaryMux.Unlock()
	stompUnit.TopicSubidDictionaryMutex.Lock()
	{
		if unSubScribeFrame.Header.ContainsKey(StompHeaders.Id) && unSubScribeFrame.Header.Get(StompHeaders.Id) != "" {
			destinationOfSubid := tools.MapGetFirstKeyByValue(stompUnit.TopicSubidDictionary, unSubScribeFrame.Header.Get(StompHeaders.Id))

			if stompUnitMapContainsKey(manager.stompSubscribeDictionary, destinationOfSubid) {
				if sessionIdStompUnitMapContainsKey(manager.stompSubscribeDictionary[destinationOfSubid], stompUnit.SessionId) {
					delete(manager.stompSubscribeDictionary[destinationOfSubid], stompUnit.SessionId)
				}
			}
		}
	}
	stompUnit.TopicSubidDictionaryMutex.Unlock()
}

// when stompBehavior closing,need remove subscribe
func (manager *StompManager) StompUnitOnClose(stompUnit *StompUnit) {
	manager.stompSubscribeDictionaryMux.Lock()
	defer manager.stompSubscribeDictionaryMux.Unlock()

	stompUnit.TopicSubidDictionaryMutex.Lock()
	{
		for destination := range stompUnit.TopicSubidDictionary {
			if stompUnitMapContainsKey(manager.stompSubscribeDictionary, destination) {
				if sessionIdStompUnitMapContainsKey(manager.stompSubscribeDictionary[destination], stompUnit.SessionId) {
					delete(manager.stompSubscribeDictionary[destination], stompUnit.SessionId)
				}

				if len(manager.stompSubscribeDictionary[destination]) == 0 {
					delete(manager.stompSubscribeDictionary, destination)
				}
			}
		}
	}
	stompUnit.TopicSubidDictionaryMutex.Unlock()

	// for k := range manager.stompSubscribeDictionary {
	// 	Ilog.D(DebugLevel.Debug, fmt.Sprint("SubScribe Dictionary:k", k, "   Len = ", len(manager.stompSubscribeDictionary[k])))
	// }
}

//  指定StompUnit.SessionId查找StompUnit发送stompFrame
// LetStompUnitSend   指定StompUnit.SessionId查找StompUnit发送stompFrame
//  参数:
// 	stompMessageFrame	*Frame	指定的Stomp数据包
// 	unitSessionId		string	指定的StompUnit的SessionId
func (manager *StompManager) LetStompUnitSend(stompMessageFrame *Frame, unitSessionId string) {
	if stompMessageFrame.Command == StompCommand.Message && stompMessageFrame.Header.ContainsKey(StompHeaders.Destination) && stompMessageFrame.Header.Get(StompHeaders.Destination) != "" && stompMessageFrame.Header.ContainsKey(StompHeaders.MessageId) && stompMessageFrame.Header.Get(StompHeaders.MessageId) != "" {
		// 根据unitSessionId 找到对应的unit
		if sessionIdStompUnitMapContainsKey(manager.stompSubscribeDictionary[stompMessageFrame.Header.Get(StompHeaders.Destination)], unitSessionId) {
			stompUnit := manager.stompSubscribeDictionary[stompMessageFrame.Header.Get(StompHeaders.Destination)][unitSessionId]
			stompUnit.SendStompMessage(stompMessageFrame)
		}
	}
}

// 用来发布订阅的内容.
func (manager *StompManager) Publish(stompMessageFrame *Frame) {
	if stompMessageFrame.Command == StompCommand.Message && stompMessageFrame.Header.ContainsKey(StompHeaders.Destination) && stompMessageFrame.Header.Get(StompHeaders.Destination) != "" && stompMessageFrame.Header.ContainsKey(StompHeaders.MessageId) && stompMessageFrame.Header.Get(StompHeaders.MessageId) != "" {
		destinationValue := stompMessageFrame.Header.Get(StompHeaders.Destination)
		flowDirction := manager.SubscribeFlowDirction(destinationValue)
		switch flowDirction {
		case All:
			{
				for _, unit := range manager.stompSubscribeDictionary[destinationValue] {
					unit.SendStompMessage(stompMessageFrame)
				}
				break
			}
		case User:
			{
				for _, v := range manager.stompSubscribeDictionary[destinationValue] {
					if v.IsUser() {
						v.SendStompMessage(stompMessageFrame)
					}
				}
				break
			}
		default:
			{
				return
			}
		}
	}
}
