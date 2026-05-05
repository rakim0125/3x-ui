package controller

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"

	"github.com/mhsanaei/3x-ui/v2/database/model"
	"github.com/mhsanaei/3x-ui/v2/web/service"

	"github.com/gin-gonic/gin"
)

// OpenAPIController exposes configuration and statistics APIs
// authenticated via API key (header X-API-Key or query param api_key).
type OpenAPIController struct {
	apiKeyService  service.ApiKeyService
	inboundService service.InboundService
	settingService service.SettingService
	serverService  service.ServerService
	xrayService    service.XrayService
	outboundService service.OutboundService
	xraySettingService service.XraySettingService
}

func NewOpenAPIController(g *gin.RouterGroup) *OpenAPIController {
	a := &OpenAPIController{}
	a.initRouter(g)
	return a
}

// apiKeyAuth validates the API key from header or query parameter.
func (a *OpenAPIController) apiKeyAuth(c *gin.Context) {
	key := c.GetHeader("X-API-Key")
	if key == "" {
		key = c.Query("api_key")
	}
	if key == "" {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"msg":     "missing api key: provide X-API-Key header or api_key query parameter",
		})
		return
	}

	apiKey, err := a.apiKeyService.Validate(key)
	if err != nil {
		c.AbortWithStatusJSON(http.StatusUnauthorized, gin.H{
			"success": false,
			"msg":     err.Error(),
		})
		return
	}

	c.Set("api_user_id", apiKey.UserId)
	c.Next()
}

func (a *OpenAPIController) getUserId(c *gin.Context) int {
	id, _ := c.Get("api_user_id")
	return id.(int)
}

func (a *OpenAPIController) initRouter(g *gin.RouterGroup) {
	api := g.Group("/openapi")
	api.Use(a.apiKeyAuth)

	// ---------- Inbounds (configuration) ----------
	inbounds := api.Group("/inbounds")
	inbounds.GET("/list", a.getInbounds)
	inbounds.GET("/get/:id", a.getInbound)
	inbounds.POST("/add", a.addInbound)
	inbounds.POST("/del/:id", a.delInbound)
	inbounds.POST("/update/:id", a.updateInbound)
	inbounds.POST("/addClient", a.addInboundClient)
	inbounds.POST("/:id/delClient/:clientId", a.delInboundClient)
	inbounds.POST("/updateClient/:clientId", a.updateInboundClient)
	inbounds.POST("/import", a.importInbound)

	// ---------- Inbounds (statistics / traffic) ----------
	inbounds.GET("/getClientTraffics/:email", a.getClientTraffics)
	inbounds.GET("/getClientTrafficsById/:id", a.getClientTrafficsById)
	inbounds.POST("/clientIps/:email", a.getClientIps)
	inbounds.POST("/clearClientIps/:email", a.clearClientIps)
	inbounds.POST("/:id/resetClientTraffic/:email", a.resetClientTraffic)
	inbounds.POST("/resetAllTraffics", a.resetAllTraffics)
	inbounds.POST("/resetAllClientTraffics/:id", a.resetAllClientTraffics)
	inbounds.POST("/delDepletedClients/:id", a.delDepletedClients)
	inbounds.POST("/onlines", a.onlines)
	inbounds.POST("/lastOnline", a.lastOnline)
	inbounds.POST("/updateClientTraffic/:email", a.updateClientTraffic)
	inbounds.POST("/:id/delClientByEmail/:email", a.delInboundClientByEmail)

	// ---------- Xray (configuration) ----------
	xrayGroup := api.Group("/xray")
	xrayGroup.GET("/getDefaultJsonConfig", a.getDefaultXrayConfig)
	xrayGroup.GET("/getOutboundsTraffic", a.getOutboundsTraffic)
	xrayGroup.GET("/getXrayResult", a.getXrayResult)
	xrayGroup.POST("/", a.getXraySetting)
	xrayGroup.POST("/update", a.updateXraySetting)
	xrayGroup.POST("/resetOutboundsTraffic", a.resetOutboundsTraffic)

	// ---------- Server (statistics) ----------
	server := api.Group("/server")
	server.GET("/status", a.status)
	server.GET("/cpuHistory/:bucket", a.getCpuHistoryBucket)
	server.GET("/getConfigJson", a.getConfigJson)
}

// ==================== Inbound handlers ====================

func (a *OpenAPIController) getInbounds(c *gin.Context) {
	userId := a.getUserId(c)
	inbounds, err := a.inboundService.GetInbounds(userId)
	if err != nil {
		jsonMsg(c, "get inbounds", err)
		return
	}
	jsonObj(c, inbounds, nil)
}

func (a *OpenAPIController) getInbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "invalid id", err)
		return
	}
	inbound, err := a.inboundService.GetInbound(id)
	if err != nil {
		jsonMsg(c, "get inbound", err)
		return
	}
	jsonObj(c, inbound, nil)
}

func (a *OpenAPIController) addInbound(c *gin.Context) {
	inbound := &model.Inbound{}
	if err := c.ShouldBind(inbound); err != nil {
		jsonMsg(c, "invalid request", err)
		return
	}
	inbound.UserId = a.getUserId(c)
	if inbound.Listen == "" || inbound.Listen == "0.0.0.0" || inbound.Listen == "::" || inbound.Listen == "::0" {
		inbound.Tag = fmt.Sprintf("inbound-%v", inbound.Port)
	} else {
		inbound.Tag = fmt.Sprintf("inbound-%v:%v", inbound.Listen, inbound.Port)
	}
	inbound, needRestart, err := a.inboundService.AddInbound(inbound)
	if err != nil {
		jsonMsg(c, "add inbound", err)
		return
	}
	jsonMsgObj(c, "inbound added", inbound, nil)
	if needRestart {
		a.xrayService.SetToNeedRestart()
	}
}

func (a *OpenAPIController) delInbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "invalid id", err)
		return
	}
	needRestart, err := a.inboundService.DelInbound(id)
	if err != nil {
		jsonMsg(c, "delete inbound", err)
		return
	}
	jsonMsgObj(c, "inbound deleted", id, nil)
	if needRestart {
		a.xrayService.SetToNeedRestart()
	}
}

func (a *OpenAPIController) updateInbound(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "invalid id", err)
		return
	}
	inbound := &model.Inbound{Id: id}
	if err := c.ShouldBind(inbound); err != nil {
		jsonMsg(c, "invalid request", err)
		return
	}
	inbound, needRestart, err := a.inboundService.UpdateInbound(inbound)
	if err != nil {
		jsonMsg(c, "update inbound", err)
		return
	}
	jsonMsgObj(c, "inbound updated", inbound, nil)
	if needRestart {
		a.xrayService.SetToNeedRestart()
	}
}

func (a *OpenAPIController) addInboundClient(c *gin.Context) {
	data := &model.Inbound{}
	if err := c.ShouldBind(data); err != nil {
		jsonMsg(c, "invalid request", err)
		return
	}
	needRestart, err := a.inboundService.AddInboundClient(data)
	if err != nil {
		jsonMsg(c, "add client", err)
		return
	}
	jsonMsg(c, "client added", nil)
	if needRestart {
		a.xrayService.SetToNeedRestart()
	}
}

func (a *OpenAPIController) delInboundClient(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "invalid id", err)
		return
	}
	clientId := c.Param("clientId")
	needRestart, err := a.inboundService.DelInboundClient(id, clientId)
	if err != nil {
		jsonMsg(c, "delete client", err)
		return
	}
	jsonMsg(c, "client deleted", nil)
	if needRestart {
		a.xrayService.SetToNeedRestart()
	}
}

func (a *OpenAPIController) updateInboundClient(c *gin.Context) {
	clientId := c.Param("clientId")
	inbound := &model.Inbound{}
	if err := c.ShouldBind(inbound); err != nil {
		jsonMsg(c, "invalid request", err)
		return
	}
	needRestart, err := a.inboundService.UpdateInboundClient(inbound, clientId)
	if err != nil {
		jsonMsg(c, "update client", err)
		return
	}
	jsonMsg(c, "client updated", nil)
	if needRestart {
		a.xrayService.SetToNeedRestart()
	}
}

func (a *OpenAPIController) importInbound(c *gin.Context) {
	inbound := &model.Inbound{}
	if err := json.Unmarshal([]byte(c.PostForm("data")), inbound); err != nil {
		jsonMsg(c, "invalid data", err)
		return
	}
	inbound.Id = 0
	inbound.UserId = a.getUserId(c)
	if inbound.Listen == "" || inbound.Listen == "0.0.0.0" || inbound.Listen == "::" || inbound.Listen == "::0" {
		inbound.Tag = fmt.Sprintf("inbound-%v", inbound.Port)
	} else {
		inbound.Tag = fmt.Sprintf("inbound-%v:%v", inbound.Listen, inbound.Port)
	}
	for index := range inbound.ClientStats {
		inbound.ClientStats[index].Id = 0
		inbound.ClientStats[index].Enable = true
	}
	inbound, needRestart, err := a.inboundService.AddInbound(inbound)
	jsonMsgObj(c, "inbound imported", inbound, err)
	if err == nil && needRestart {
		a.xrayService.SetToNeedRestart()
	}
}

// ==================== Traffic / statistics handlers ====================

func (a *OpenAPIController) getClientTraffics(c *gin.Context) {
	email := c.Param("email")
	traffics, err := a.inboundService.GetClientTrafficByEmail(email)
	if err != nil {
		jsonMsg(c, "get client traffics", err)
		return
	}
	jsonObj(c, traffics, nil)
}

func (a *OpenAPIController) getClientTrafficsById(c *gin.Context) {
	id := c.Param("id")
	traffics, err := a.inboundService.GetClientTrafficByID(id)
	if err != nil {
		jsonMsg(c, "get client traffics", err)
		return
	}
	jsonObj(c, traffics, nil)
}

func (a *OpenAPIController) getClientIps(c *gin.Context) {
	email := c.Param("email")
	ips, err := a.inboundService.GetInboundClientIps(email)
	if err != nil || ips == "" {
		jsonObj(c, "No IP Record", nil)
		return
	}
	jsonObj(c, ips, nil)
}

func (a *OpenAPIController) clearClientIps(c *gin.Context) {
	email := c.Param("email")
	err := a.inboundService.ClearClientIps(email)
	jsonMsg(c, "clear client ips", err)
}

func (a *OpenAPIController) resetClientTraffic(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "invalid id", err)
		return
	}
	email := c.Param("email")
	needRestart, err := a.inboundService.ResetClientTraffic(id, email)
	if err != nil {
		jsonMsg(c, "reset client traffic", err)
		return
	}
	jsonMsg(c, "client traffic reset", nil)
	if needRestart {
		a.xrayService.SetToNeedRestart()
	}
}

func (a *OpenAPIController) resetAllTraffics(c *gin.Context) {
	err := a.inboundService.ResetAllTraffics()
	if err != nil {
		jsonMsg(c, "reset all traffics", err)
		return
	}
	a.xrayService.SetToNeedRestart()
	jsonMsg(c, "all traffics reset", nil)
}

func (a *OpenAPIController) resetAllClientTraffics(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "invalid id", err)
		return
	}
	err = a.inboundService.ResetAllClientTraffics(id)
	if err != nil {
		jsonMsg(c, "reset all client traffics", err)
		return
	}
	a.xrayService.SetToNeedRestart()
	jsonMsg(c, "all client traffics reset", nil)
}

func (a *OpenAPIController) delDepletedClients(c *gin.Context) {
	id, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "invalid id", err)
		return
	}
	err = a.inboundService.DelDepletedClients(id)
	jsonMsg(c, "delete depleted clients", err)
}

func (a *OpenAPIController) onlines(c *gin.Context) {
	jsonObj(c, a.inboundService.GetOnlineClients(), nil)
}

func (a *OpenAPIController) lastOnline(c *gin.Context) {
	data, err := a.inboundService.GetClientsLastOnline()
	jsonObj(c, data, err)
}

func (a *OpenAPIController) updateClientTraffic(c *gin.Context) {
	email := c.Param("email")
	type TrafficUpdateRequest struct {
		Upload   int64 `json:"upload"`
		Download int64 `json:"download"`
	}
	var req TrafficUpdateRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		jsonMsg(c, "invalid request", err)
		return
	}
	err := a.inboundService.UpdateClientTrafficByEmail(email, req.Upload, req.Download)
	jsonMsg(c, "client traffic updated", err)
}

func (a *OpenAPIController) delInboundClientByEmail(c *gin.Context) {
	inboundId, err := strconv.Atoi(c.Param("id"))
	if err != nil {
		jsonMsg(c, "invalid inbound id", err)
		return
	}
	email := c.Param("email")
	needRestart, err := a.inboundService.DelInboundClientByEmail(inboundId, email)
	if err != nil {
		jsonMsg(c, "delete client by email", err)
		return
	}
	jsonMsg(c, "client deleted", nil)
	if needRestart {
		a.xrayService.SetToNeedRestart()
	}
}

// ==================== Xray handlers ====================

func (a *OpenAPIController) getDefaultXrayConfig(c *gin.Context) {
	config, err := a.settingService.GetDefaultXrayConfig()
	if err != nil {
		jsonMsg(c, "get xray config", err)
		return
	}
	jsonObj(c, config, nil)
}

func (a *OpenAPIController) getOutboundsTraffic(c *gin.Context) {
	traffic, err := a.outboundService.GetOutboundsTraffic()
	if err != nil {
		jsonMsg(c, "get outbounds traffic", err)
		return
	}
	jsonObj(c, traffic, nil)
}

func (a *OpenAPIController) getXrayResult(c *gin.Context) {
	jsonObj(c, a.xrayService.GetXrayResult(), nil)
}

func (a *OpenAPIController) getXraySetting(c *gin.Context) {
	xraySetting, err := a.settingService.GetXrayConfigTemplate()
	if err != nil {
		jsonMsg(c, "get xray setting", err)
		return
	}
	inboundTags, err := a.inboundService.GetInboundTags()
	if err != nil {
		jsonMsg(c, "get inbound tags", err)
		return
	}
	outboundTestUrl, _ := a.settingService.GetXrayOutboundTestUrl()
	if outboundTestUrl == "" {
		outboundTestUrl = "https://www.google.com/generate_204"
	}
	resp := map[string]interface{}{
		"xraySetting":     json.RawMessage(xraySetting),
		"inboundTags":     json.RawMessage(inboundTags),
		"outboundTestUrl": outboundTestUrl,
	}
	result, err := json.Marshal(resp)
	if err != nil {
		jsonMsg(c, "marshal xray setting", err)
		return
	}
	jsonObj(c, string(result), nil)
}

func (a *OpenAPIController) updateXraySetting(c *gin.Context) {
	xraySetting := c.PostForm("xraySetting")
	if err := a.xraySettingService.SaveXraySetting(xraySetting); err != nil {
		jsonMsg(c, "update xray setting", err)
		return
	}
	outboundTestUrl := c.PostForm("outboundTestUrl")
	if outboundTestUrl == "" {
		outboundTestUrl = "https://www.google.com/generate_204"
	}
	_ = a.settingService.SetXrayOutboundTestUrl(outboundTestUrl)
	jsonMsg(c, "xray setting updated", nil)
}

func (a *OpenAPIController) resetOutboundsTraffic(c *gin.Context) {
	tag := c.PostForm("tag")
	err := a.outboundService.ResetOutboundTraffic(tag)
	jsonMsg(c, "reset outbound traffic", err)
}

// ==================== Server status handlers ====================

func (a *OpenAPIController) status(c *gin.Context) {
	s := a.serverService.GetStatus(nil)
	jsonObj(c, s, nil)
}

func (a *OpenAPIController) getCpuHistoryBucket(c *gin.Context) {
	bucketStr := c.Param("bucket")
	bucket, err := strconv.Atoi(bucketStr)
	if err != nil || bucket <= 0 {
		jsonMsg(c, "invalid bucket", fmt.Errorf("bad bucket"))
		return
	}
	allowed := map[int]bool{2: true, 30: true, 60: true, 120: true, 180: true, 300: true}
	if !allowed[bucket] {
		jsonMsg(c, "invalid bucket", fmt.Errorf("unsupported bucket"))
		return
	}
	points := a.serverService.AggregateCpuHistory(bucket, 60)
	jsonObj(c, points, nil)
}

func (a *OpenAPIController) getConfigJson(c *gin.Context) {
	configJson, err := a.serverService.GetConfigJson()
	if err != nil {
		jsonMsg(c, "get config json", err)
		return
	}
	jsonObj(c, configJson, nil)
}
