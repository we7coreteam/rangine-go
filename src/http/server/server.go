package server

import (
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"github.com/we7coreteam/w7-rangine-go/src/core/server"
	"github.com/we7coreteam/w7-rangine-go/src/http/session"
)

var GHttpServer *Server

func GetServer() *Server {
	return GHttpServer
}

func Use(middleware ...gin.HandlerFunc) gin.IRouter {
	GHttpServer.Engine.Use(middleware...)

	return GHttpServer.Engine
}

func RegisterRouters(register func(engine *gin.Engine)) *Server {
	register(GHttpServer.Engine)
	return GHttpServer
}

func GetSession() *session.Session {
	return GHttpServer.Session
}

type Server struct {
	server.Interface
	config *viper.Viper

	Engine  *gin.Engine
	Session *session.Session
}

func NewHttpDefaultServer(config *viper.Viper) *Server {
	var sessionConfig session.SessionConf
	var cookieConfig session.Cookie
	err := config.UnmarshalKey("session", &sessionConfig)
	if err != nil {
		panic(err)
	}
	err = config.UnmarshalKey("cookie", &cookieConfig)
	if err != nil {
		panic(err)
	}

	server := NewServer(config)
	server.Session = session.NewSession(sessionConfig, cookieConfig)

	return server
}

func NewServer(config *viper.Viper) *Server {
	server := &Server{
		config: config,
	}
	server.initGinEngine()
	GHttpServer = server

	return server
}

func (server *Server) initGinEngine() {
	gin.SetMode("release")
	server.Engine = gin.New()
	server.Engine.Routes()
}

func (server Server) GetServerName() string {
	return "http"
}

func (server Server) GetOptions() map[string]string {
	return map[string]string{
		"Host": server.config.GetString("server.http.host"),
		"Port": server.config.GetString("server.http.port"),
	}
}

func (server Server) Start() {
	err := server.Engine.Run(server.config.GetString("server.http.host") + ":" + server.config.GetString("server.http.port"))
	if err != nil {
		panic(err)
	}
}
