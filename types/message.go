package types

import "time"

//此消息是为了承接来自不同channel的消息
//不同channel的消息是不一样的
//这里的message是为了统一格式
type Message struct {
	//发送的消息
	ID string `json:"id"`
	//消息
	Content string `json:"content"`
	//每个会话ID
	SessionID string `json:"sessionID"`
	//TODO 最后搞清楚为什么这个id要存在这里
	UserID string `json:"userId"`
	//来源
	Channel string `json:"channel"`

	//创建时间
	CreatedTime time.Time `json:"createdTime"`
}

//此消息是agent处理完成后 返送给gataway的消息
type Response struct {
	//回复的唯一ID
	ID string `json:"id"`
	//回复的消息
	Content string `json:"content"`
	//这里的to其实是告诉应该发送给那个
	SessionID string `json:"sessionID"`
	//生产的时间
	CreatedTime time.Time `json:"createdTime"`
}

type Session struct {
	//session 对应的id 这个是随机产生的
	ID string `json:"id"`

	//userID 其实就是绑定信息 比如 绑定了微信这个时候微信和本地的服务有一个协商的唯一id
	UserID string `json:"userId"`

	//本次session建立的时候绑定来源
	Channel string `json:"channel"`

	//config 是本次会话绑定的配置 这个配置里面包含了对应的agent的使用等等
	//Config string `json:"config"`
	//
	////最后一条信息的id
	//LastMessageID string `json:"lastMessageID"`
	//
	//MessageCount int `json:"messageCount"`

	LastActivityTime time.Time `json:"lastActivityTime"`
	//产生session 的时间
	CreatedTime time.Time `json:"createdTime"`
}
