package webscoket

import (
	"fmt"
	"github.com/gorilla/websocket"
	"net/http"
	"z3/controllers"
)


//用户状态
const(
	ON_LINE  	= 1		//在线
	OFF_LINE	= 2		//离线
	INVISIBLE	= 3		//隐身
)

//客户端
type Client struct {
	id 			int					//用户id
	socket 		*websocket.Conn		//用户长链接
	sendMsg		chan interface{}	//用户消息
	Status		int					//用户状态
}

//消息
type Message struct {
	SendId 		int					//发送者
	ReceiveId	int					//接受者
	Content		string				//内容
}

//客户端管理
type ClientManage struct {
	Register	chan *Client
	LogOff		chan *Client
	Clients 	map[int]*Client
	Broadcast	chan interface{}
}

var Manager = &ClientManage{
	Register:	make(chan *Client),
	LogOff:		make(chan *Client),
	Clients:	make(map[int]*Client),
	Broadcast:	make(chan interface{}),
}


//注册链接
func (manager *ClientManage) RegisterSocket(c *Client) {
	Manager.Register <- c
	go c.read()
	go c.write()
}

//广播
func (manager *ClientManage) ListenEvent() {
	for {
		select {
		//注册链接
		case conn := <- manager.Register:
			manager.Clients[conn.id] = conn

		//注销链接
		case conn := <- manager.LogOff:
			delete(manager.Clients,conn.id)

		//向所有在线用户推送广播
		case msg :=  <-manager.Broadcast:
			for _,conn := range manager.Clients {
				if conn.Status != OFF_LINE {
					conn.sendMsg <- msg
				}else{
					continue
				}
			}
		}
	}
}


//读取消息
func (c *Client) read(){
	defer func() {
		Manager.LogOff <- c
		if err := c.socket.Close();err != nil {
			controllers.WriteLog("ws-注销链接失败-defer,Error:"+err.Error())
		}
	}()
	for {
		_,message,err := c.socket.ReadMessage()
		if err != nil {
			Manager.LogOff <- c
			if err := c.socket.Close();err != nil {
				controllers.WriteLog("ws-注销链接失败,Error:"+err.Error())
				break
			}
		}
		msg := &Message{SendId: c.id, Content: string(message)}
		Manager.Broadcast <- msg
	}
}

//发送消息
func (c *Client) write(){
	for  {
		select {
		case msg := <- c.sendMsg:
			if err := c.socket.WriteJSON(msg); err != nil {
				controllers.WriteLog("ws发送消息失败,Error:"+err.Error())
			}
		}
	}

}

var token string

func WsHandler(res http.ResponseWriter, req *http.Request) {
	defer func() {
		if err:= recover();err != nil {
			fmt.Println(err)
		}
	}()

	//进行用户登录权限验证
	if err := AuthcheckUser(req);err != "" {
		res.WriteHeader(403)
		return
	}

	response := &http.Response{
		Header:make(map[string][]string),
	}

	response.Header.Add("Sec-Websocket-Protocol",token)
	//将http升级为websocket
	conn , err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool {
			return true
		}}).Upgrade(res,req,response.Header)
	if err != nil {
		http.NotFound(res,req)
		return
	}

	client := &Client{
		id: 		controllers.MemberId,
		socket: 	conn,
		sendMsg: 	make(chan interface{}),
		Status:		ON_LINE,
	}
	//注册一个新的链接
	Manager.RegisterSocket(client)
}

func AuthcheckUser(req *http.Request) string{
	header := req.Header
	if header == nil || header["Sec-Websocket-Protocol"] == nil || header["Sec-Websocket-Protocol"][0] == "" {
		return  "暂无权限进行此操作"
	}

	token = header["Sec-Websocket-Protocol"][0]
	if err := controllers.ResUserKeyMsg(token);err != "" {
		return err
	}
	return ""
}
