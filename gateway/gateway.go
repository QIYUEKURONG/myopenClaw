package gateway

import (
	"context"
	"errors"
	"fmt"
	"log"
	"myopenclaw/agent"
	"myopenclaw/storage"
	"myopenclaw/types"
	"strings"
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
	//load session
	sessionMap, err := storage.LoadSessionIndex()
	if err != nil {
		fmt.Errorf("Error loading session index: %v", err)
		panic(err)
	}

	result := make(map[string]*types.Session)
	for sessionKey, sessionId := range sessionMap {
		result[sessionKey] = &types.Session{
			ID: sessionId,
		}
	}

	return &Gateway{
		Session: result,
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

	index, _ := storage.LoadSessionIndex()
	index[sessionKey] = newSession.ID

	err := storage.SaveSessionIndex(index)
	if err != nil {
		log.Fatalf("Error saving session: %v", err)
		return nil, fmt.Errorf("Error saving session: %v", err)
	}

	return &newSession, nil
}

func (g *Gateway) handleCommand(ctx context.Context, msg *types.Message) (*types.Response, error) {
	switch msg.Content {
	case "/new":
		//从内存删了
		sessionKey := msg.UserID + "_" + msg.Channel
		if _, ok := g.Session[sessionKey]; ok {
			delete(g.Session, sessionKey)
		}
		return &types.Response{ID: uuid.NewString(), Content: "已创建新会话", SessionID: msg.SessionID, CreatedTime: time.Now()}, nil
	case "/reset":
		sessionKey := msg.UserID + "_" + msg.Channel
		if _, ok := g.Session[sessionKey]; !ok {
			return &types.Response{ID: uuid.NewString(), Content: "未找到对应的用户", SessionID: msg.SessionID, CreatedTime: time.Now()},
				fmt.Errorf("未找到对应的用户信息")
		}
		//storage.SaveSessionIndex()

	case "history":

	default:
		return &types.Response{Content: "未知命令: " + msg.Content}, nil
	}

}

// HandleMessage 处理消息（核心方法）
func (g *Gateway) HandleMessage(ctx context.Context, msg *types.Message) (*types.Response, error) {
	//对ai模式进行
	if strings.HasPrefix(msg.Channel, "/") {
		return g.handleCommand(ctx, msg)
	}

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
