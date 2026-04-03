package gateway

import (
	"context"
	"errors"
	"fmt"
	"myopenclaw/agent"
	"myopenclaw/types"
	"sync"
	"time"

	"github.com/google/uuid"
)

type Gateway struct {
	//birge 是外部不同协议和内部消息的桥梁 只要实现get和push就可以
	//Birge []Birge `json:"birge"`
	//gateway管理所有的session
	Session map[string]*types.Session `json:"session"`
	//MessageRoute map[types.Session]agent.AgentRunime `json:"messageRoute"`

	//之所以定义为map是因为同一个session必须排队 但是不同的session可以并发进行
	//muList map[string]sync.RWMutex
	//并发的数量
	//ConcurrenceSession int
	RunTime *agent.Runtime

	GlobalRw sync.RWMutex
}

// Birge 是只要存在get和push消息就行
type Birge struct {
	GetMessage  interface{} `json:"getMessage"`
	PushMessage interface{} `json:"pushMessage"`
}

func NewGateway() *Gateway {
	return &Gateway{
		Session: make(map[string]*types.Session),
	}
}

func (g *Gateway) getOrCreateSession(msg *types.Message) (*types.Session, error) {
	//添加一个RwMutext
	g.GlobalRw.Lock()
	defer g.GlobalRw.Unlock()

	if msg == nil {
		return nil, errors.New("msg == nil")
	}

	sessionKey := msg.UserID + "_" + msg.Channel

	if val, ok := g.Session[sessionKey]; ok {
		return val, nil
	}

	var newSession types.Session
	//这个id对应的本次的唯一的id
	newSession.ID = uuid.New().String()
	newSession.CreatedTime = time.Now()
	newSession.LastActivityTime = time.Now()
	newSession.UserID = msg.UserID
	newSession.Channel = msg.Channel
	newSession.Messages = make([]types.LLMMessage, 0)

	g.Session[sessionKey] = &newSession

	return &newSession, nil
}

// HandleMessage 处理消息（核心方法）
func (g *Gateway) HandleMessage(ctx context.Context, msg *types.Message) (*types.Response, error) {
	// 你来实现：
	// 1. 获取或创建 Session
	// 2. 将 Session.ID 赋值给 msg.SessionID
	// 3. 打印日志（方便调试）
	// 4. 暂时返回一个假的 Response（因为还没有 Agent Runtime）
	session, err := g.getOrCreateSession(msg)
	if err != nil {
		fmt.Errorf("getOrCreateSession find error: +%v", err.Error())
		return nil, err
	}
	msg.SessionID = session.ID
	fmt.Printf("%v : %v\n", msg, session)

	response, err := g.RunTime.ProcessMessage(ctx, msg, session)
	if err != nil {
		fmt.Errorf("RunTime ProcessMessage error: +%v", err.Error())
		return nil, err
	}

	return response, nil
}
