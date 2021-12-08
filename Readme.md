# stompserver

go stomp server base on "net/http"

base on "net/http" and "golang.org/x/net/websocket"
so use one port, you can be WebServer or StompServer or websocket or both

基于"net/http"和"golang.org/x/net/websocket",所以一个端口,就可以同时实现web服务,websocket,stomp服务.

# example
https://github.com/0xAAFF/stompserver_example


# 使用方式

```
	//http.HandleFunc("/", HttpVueServer)                                        // Web服务
	//fmt.Println("Web Server  : http://127.0.0.1:", 80, "/")

	http.Handle("/stomp", websocket.Handler(StompServerInstance.NewStompUnit)) // Stomp服务
	fmt.Println("Stomp Server: ws://localhost:", 80, "/stomp")
```

see www_example.go


## 模块使用说明
/example/

- Web 网页解析模块
    1. 模块路径
        server_web.go
    2. 填充代码 实现接口 -> onReflex
        server_web.go已经实现了针对Vue2/3/4 build项目的文件支持,
        后期调用接口主要写在func onReflex(responseW http.ResponseWriter, request *http.Request)函数中
        当然,因为项目已经支持了Stomp协议,建议项目中此模块只用于网页资源的解析即可
    3. 测试建议
        - 基础web资源解析(通过浏览器可以访问到web项目中的各种资源,html,js,css,图片,字体等web资源)
        - 针对web路径做安全测试,例如访问js目录中非js文件,或者构造其他路径,使得程序奔溃,或者构造伪装路径获取其他目录资源,甚至构造路径尝试执行代码等

- Stomp 模块
    1. 模块路径
        server_stomp.go
    2. 填充代码,实现接口
        - 配置群发,组发,单发地址 server_stomp.go->init()
            在init()中,需要实现两个:注册群发组发单发的根地址标识 和 注册可被客户端订阅的地址

        - 实现每个订阅地址的实际功能 ./service/kernel/stomp_reflex.go->Reflex(sourceStompMessage *xstomp.Frame, unit *xstomp.StompUnit)
            在Reflex(sourceStompMessage *xstomp.Frame, unit *xstomp.StompUnit)函数中,实现针对每个地址的接口访问
            这样的好处是一处提交,数据同步所有客户端,对于状态改变等非常迅速.主要实现的接口将在此处



## Copyright

```
/*******************************************************************************
 *
 *  Copyright (c) 2021 0xAAFF<littletools@outlook.com>
 *
 *  Licensed under the Apache License, Version 2.0 (the "License");
 *  you may not use this file except in compliance with the License.
 *  You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 *  Unless required by applicable law or agreed to in writing, software
 *  distributed under the License is distributed on an "AS IS" BASIS,
 *  WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 *  See the License for the specific language governing permissions and
 *  limitations under the License.
 * 
 *******************************************************************************/
 ```


 ## AboutMe

 [MyPage](https://www.jianshu.com/u/dbdb14e006b3)