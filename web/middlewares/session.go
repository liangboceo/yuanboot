package middlewares

import (
	"github.com/liangboceo/yuanboot/abstractions"
	"github.com/liangboceo/yuanboot/web/context"
	"github.com/liangboceo/yuanboot/web/session"
	"github.com/liangboceo/yuanboot/web/session/identity"
	"github.com/liangboceo/yuanboot/web/session/store"
	"sync"
)

type SessionMiddleware struct {
	*BaseMiddleware
	sessionMgr   context.ISessionManager
	sessionStore store.ISessionStore
	identity     identity.IProvider
	sessionName  string
	mMaxLifeTime int64
	lock         sync.Mutex
}

type SessionConfig struct {
	Name    string `mapstructure:"name" config:"name"`
	TimeOut int64  `mapstructure:"timeout" config:"timeout"`
}

func NewSessionWith(provider identity.IProvider, store store.ISessionStore, config abstractions.IConfiguration) *SessionMiddleware {
	var sessionConfig *SessionConfig
	config.GetSection("yuanboot.application.server.session").Unmarshal(&sessionConfig)
	if sessionConfig.TimeOut == 0 {
		sessionConfig.TimeOut = 3600
	}
	if sessionConfig.Name != "" {
		provider.SetName(sessionConfig.Name)
	}
	store.SetMaxLifeTime(sessionConfig.TimeOut)
	mgr := session.NewSessionWithStore(store)
	return &SessionMiddleware{BaseMiddleware: &BaseMiddleware{}, sessionMgr: mgr, identity: provider, sessionStore: store}
}

func (sessionMid *SessionMiddleware) Inovke(ctx *context.HttpContext, next func(ctx *context.HttpContext)) {
	sessionMid.lock.Lock()
	sessionMid.identity.SetContext(ctx)
	sessionId := sessionMid.sessionMgr.Load(sessionMid.identity)
	sessionMid.lock.Unlock()

	sessionId = sessionMid.sessionMgr.NewSession(sessionId)
	if sessionId != "" {
		ctx.SetItem("sessionId", sessionId)
		ctx.SetItem("sessionMgr", sessionMid.sessionMgr)
	}
	next(ctx)
}
