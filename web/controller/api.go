package controller

import (
	"net/http"
	"strconv"

	"github.com/mhsanaei/3x-ui/v2/web/service"
	"github.com/mhsanaei/3x-ui/v2/web/session"

	"github.com/gin-gonic/gin"
)

// APIController handles the main API routes for the 3x-ui panel, including inbounds and server management.
type APIController struct {
	BaseController
	inboundController *InboundController
	serverController  *ServerController
	Tgbot             service.Tgbot
	apiKeyService     service.ApiKeyService
}

// NewAPIController creates a new APIController instance and initializes its routes.
func NewAPIController(g *gin.RouterGroup) *APIController {
	a := &APIController{}
	a.initRouter(g)
	return a
}

// checkAPIAuth is a middleware that returns 404 for unauthenticated API requests
// to hide the existence of API endpoints from unauthorized users
func (a *APIController) checkAPIAuth(c *gin.Context) {
	if !session.IsLogin(c) {
		c.AbortWithStatus(http.StatusNotFound)
		return
	}
	c.Next()
}

// initRouter sets up the API routes for inbounds, server, and other endpoints.
func (a *APIController) initRouter(g *gin.RouterGroup) {
	// Main API group
	api := g.Group("/panel/api")
	api.Use(a.checkAPIAuth)

	// Inbounds API
	inbounds := api.Group("/inbounds")
	a.inboundController = NewInboundController(inbounds)

	// Server API
	server := api.Group("/server")
	a.serverController = NewServerController(server)

	// API Key management
	apikeys := api.Group("/apikeys")
	apikeys.POST("/create", a.createApiKey)
	apikeys.GET("/list", a.listApiKeys)
	apikeys.POST("/delete/:id", a.deleteApiKey)

	// Extra routes
	api.GET("/backuptotgbot", a.BackuptoTgbot)
}

// BackuptoTgbot sends a backup of the panel data to Telegram bot admins.
func (a *APIController) BackuptoTgbot(c *gin.Context) {
	a.Tgbot.SendBackupToAdmins()
}

func (a *APIController) createApiKey(c *gin.Context) {
	user := session.GetLoginUser(c)
	type createReq struct {
		Name      string `json:"name" form:"name"`
		ExpiresAt int64  `json:"expiresAt" form:"expiresAt"`
	}
	var req createReq
	if err := c.ShouldBind(&req); err != nil {
		jsonMsg(c, "invalid request", err)
		return
	}
	apiKey, rawKey, err := a.apiKeyService.Create(user.Id, req.Name, req.ExpiresAt)
	if err != nil {
		jsonMsg(c, "create api key", err)
		return
	}
	jsonObj(c, gin.H{
		"id":        apiKey.Id,
		"name":      apiKey.Name,
		"prefix":    apiKey.Prefix,
		"key":       rawKey,
		"createdAt": apiKey.CreatedAt,
		"expiresAt": apiKey.ExpiresAt,
	}, nil)
}

func (a *APIController) listApiKeys(c *gin.Context) {
	user := session.GetLoginUser(c)
	keys, err := a.apiKeyService.List(user.Id)
	if err != nil {
		jsonMsg(c, "list api keys", err)
		return
	}
	jsonObj(c, keys, nil)
}

func (a *APIController) deleteApiKey(c *gin.Context) {
	user := session.GetLoginUser(c)
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "invalid id", err)
		return
	}
	err = a.apiKeyService.Delete(user.Id, id)
	jsonMsg(c, "delete api key", err)
}
