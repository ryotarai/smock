package server

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/ryotarai/smock/pkg/cli"
	"github.com/slack-go/slack"
)

type Server struct {
	CLI *cli.CLI
}

func New() *Server {
	return &Server{}
}

func (s *Server) Engine() *gin.Engine {
	engine := gin.Default()
	engine.POST("/api/auth.test", s.handleAuthTest)
	engine.POST("/api/chat.postMessage", s.handleChatPostMessage)
	engine.POST("/a/response", s.handleResponse)
	return engine
}

func (s *Server) Run(addr ...string) error {
	return s.Engine().Run(addr...)
}

type authTestResponse struct {
	slack.AuthTestResponse
	Ok    bool   `json:"ok"`
	BotID string `json:"bot_id"`
}

func (s *Server) handleAuthTest(c *gin.Context) {
	r := authTestResponse{
		Ok:    true,
		BotID: "BOTID",
		AuthTestResponse: slack.AuthTestResponse{
			URL:    "",
			Team:   "TEAM",
			User:   "BOTUSER",
			TeamID: "TEAMID",
			UserID: "BOTUSERID",
		},
	}
	c.JSON(200, r)
}

type chatPostMessageRequest struct {
	Token   string `form:"token"`
	Text    string `form:"text"`
	Channel string `form:"channel"`
}

type chatPostMessageResponse struct {
	OK      bool   `json:"ok"`
	Channel string `json:"channel"`
}

func (s *Server) handleChatPostMessage(c *gin.Context) {
	r := &chatPostMessageRequest{}
	if err := c.ShouldBind(r); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}

	s.CLI.OnMessage(&slack.Msg{
		Text: r.Text,
	})

	c.JSON(http.StatusOK, chatPostMessageResponse{
		OK:      true,
		Channel: "CHANNELID",
	})
}

func (s *Server) handleResponse(c *gin.Context) {
	msg := &slack.Msg{}
	if err := c.ShouldBindJSON(msg); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{})
	}

	s.CLI.OnMessage(msg)
}
