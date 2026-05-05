package job

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/mhsanaei/3x-ui/v2/config"
	"github.com/mhsanaei/3x-ui/v2/logger"
	"github.com/mhsanaei/3x-ui/v2/web/service"
)

// Public IP detection services, tried in order until one succeeds.
var publicIPEndpoints = []string{
	"https://api.ipify.org",
	"https://ifconfig.me/ip",
	"https://icanhazip.com",
	"https://ipinfo.io/ip",
}

// NodeRegisterJob periodically registers this node with the registry centers
// as a heartbeat. Both registry nodes receive the registration (primary + backup).
//
// The node's public IP is auto-detected on startup via external services.
// Registry URLs come from config.GetRegistryNodes() (see config package):
//   - XUI_REGISTRY_NODES / REGISTRY_NODES — comma-separated override
//   - XUI_REGISTRY_NODES_FILE — path to a line-based config file
//   - {XUI_DB_FOLDER or default}/registry_nodes — line-based file
//   - embedded default in config/registry_nodes.default
// Optional env:
//   - NODE_IP: override auto-detected public IP
//   - NODE_PORT: override the panel port for registration
//   - NODE_DESCRIPTION: node description (defaults to hostname)
type NodeRegisterJob struct {
	settingService service.SettingService
	registryNodes  []string
	nodeIP         string
	nodeDesc       string
	httpClient     *http.Client
}

type registerPayload struct {
	IP          string `json:"ip"`
	Port        int    `json:"port"`
	Description string `json:"description"`
}

func NewNodeRegisterJob() *NodeRegisterJob {
	nodes := config.GetRegistryNodes()
	desc := os.Getenv("NODE_DESCRIPTION")
	if desc == "" {
		desc, _ = os.Hostname()
	}

	nodeIP := os.Getenv("NODE_IP")
	if nodeIP == "" {
		nodeIP = detectPublicIP()
	}

	return &NodeRegisterJob{
		registryNodes: nodes,
		nodeIP:        nodeIP,
		nodeDesc:      desc,
		httpClient: &http.Client{
			Timeout: 10 * time.Second,
		},
	}
}

// detectPublicIP queries external services to discover the node's public IP.
func detectPublicIP() string {
	client := &http.Client{Timeout: 5 * time.Second}
	for _, endpoint := range publicIPEndpoints {
		resp, err := client.Get(endpoint)
		if err != nil {
			continue
		}
		body, err := io.ReadAll(resp.Body)
		resp.Body.Close()
		if err != nil {
			continue
		}
		ip := strings.TrimSpace(string(body))
		if ip != "" {
			logger.Infof("NodeRegister: detected public IP: %s (via %s)", ip, endpoint)
			return ip
		}
	}
	logger.Warning("NodeRegister: failed to detect public IP from all endpoints")
	return ""
}

// IsEnabled returns true when the required configuration is present.
func (j *NodeRegisterJob) IsEnabled() bool {
	return len(j.registryNodes) > 0 && j.nodeIP != ""
}

func (j *NodeRegisterJob) Run() {
	if !j.IsEnabled() {
		return
	}

	port, err := j.settingService.GetPort()
	if err != nil {
		logger.Warning("NodeRegister: failed to get panel port:", err)
		return
	}

	portOverride := os.Getenv("NODE_PORT")
	if portOverride != "" {
		if p, e := strconv.Atoi(portOverride); e == nil && p > 0 {
			port = p
		}
	}

	payload := registerPayload{
		IP:          j.nodeIP,
		Port:        port,
		Description: j.nodeDesc,
	}

	body, err := json.Marshal(payload)
	if err != nil {
		logger.Error("NodeRegister: marshal payload failed:", err)
		return
	}

	for _, node := range j.registryNodes {
		go j.doRegister(node, body)
	}
}

func (j *NodeRegisterJob) doRegister(baseURL string, body []byte) {
	url := fmt.Sprintf("%s/nodes/api/register", baseURL)

	resp, err := j.httpClient.Post(url, "application/json", bytes.NewReader(body))
	if err != nil {
		logger.Warningf("NodeRegister: register to %s failed: %v", baseURL, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		logger.Debugf("NodeRegister: registered to %s successfully", baseURL)
	} else {
		logger.Warningf("NodeRegister: register to %s returned status %d", baseURL, resp.StatusCode)
	}
}
