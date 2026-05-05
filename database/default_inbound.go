package database

import (
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"log"
	"slices"
	"time"

	"github.com/google/uuid"
	"github.com/mhsanaei/3x-ui/v2/database/model"
	"github.com/mhsanaei/3x-ui/v2/util/crypto"
	"github.com/mhsanaei/3x-ui/v2/util/random"
	"github.com/mhsanaei/3x-ui/v2/xray"
)

const defaultVLESSRealitySeeder = "DefaultVLESSRealityInbound"

// initDefaultVLESSRealityInbound seeds one VLESS + REALITY inbound on first install
// (empty inbounds table), matching the panel's recommended 443 / vision / chrome preset.
func initDefaultVLESSRealityInbound() error {
	var ran []string
	if err := db.Model(&model.HistoryOfSeeders{}).Pluck("seeder_name", &ran).Error; err != nil {
		return err
	}
	if slices.Contains(ran, defaultVLESSRealitySeeder) {
		return nil
	}

	var inboundCount int64
	if err := db.Model(&model.Inbound{}).Count(&inboundCount).Error; err != nil {
		return err
	}
	if inboundCount > 0 {
		return db.Create(&model.HistoryOfSeeders{SeederName: defaultVLESSRealitySeeder}).Error
	}

	var user model.User
	if err := db.Order("id asc").First(&user).Error; err != nil {
		log.Printf("default inbound: no user yet, skip: %v", err)
		return db.Create(&model.HistoryOfSeeders{SeederName: defaultVLESSRealitySeeder}).Error
	}

	priv, pub, err := crypto.GenerateRealityX25519KeyPair()
	if err != nil {
		log.Printf("default inbound: x25519 keygen failed: %v", err)
		return err
	}

	shortIDBuf := make([]byte, 4)
	if _, err = rand.Read(shortIDBuf); err != nil {
		log.Printf("default inbound: short id rand failed: %v", err)
		return err
	}
	shortID := hex.EncodeToString(shortIDBuf)

	now := time.Now().UnixMilli()
	clientID := uuid.New().String()
	email := random.Seq(8)
	subID := random.Seq(16)

	settings := map[string]any{
		"clients": []model.Client{
			{
				ID:         clientID,
				Flow:       "xtls-rprx-vision",
				Email:      email,
				LimitIP:    0,
				TotalGB:    0,
				ExpiryTime: 0,
				Enable:     true,
				TgID:       0,
				SubID:      subID,
				Comment:    "",
				Reset:      0,
				CreatedAt:  now,
				UpdatedAt:  now,
			},
		},
		"decryption": "none",
		"encryption": "none",
		"testseed":   []int{900, 500, 900, 256},
	}
	settingsJSON, err := json.MarshalIndent(settings, "", "  ")
	if err != nil {
		return err
	}

	stream := map[string]any{
		"network":         "tcp",
		"security":        "reality",
		"externalProxy":   []any{},
		"realitySettings": buildRealityStream(priv, pub, shortID),
		"tcpSettings": map[string]any{
			"acceptProxyProtocol": false,
			"header": map[string]any{
				"type": "none",
			},
		},
	}
	streamJSON, err := json.MarshalIndent(stream, "", "  ")
	if err != nil {
		return err
	}

	sniffing := map[string]any{
		"enabled":      true,
		"destOverride": []string{"http", "tls", "quic", "fakedns"},
		"metadataOnly": false,
		"routeOnly":    false,
	}
	sniffingJSON, err := json.MarshalIndent(sniffing, "", "  ")
	if err != nil {
		return err
	}

	inbound := &model.Inbound{
		UserId:         user.Id,
		Remark:         "Reality_HTTPS",
		Enable:         true,
		Listen:         "",
		Port:           443,
		Protocol:       model.VLESS,
		Settings:       string(settingsJSON),
		StreamSettings: string(streamJSON),
		Tag:            "inbound-443",
		Sniffing:       string(sniffingJSON),
	}

	tx := db.Begin()
	if err := tx.Create(inbound).Error; err != nil {
		tx.Rollback()
		log.Printf("default inbound: create failed: %v", err)
		return err
	}

	ct := xray.ClientTraffic{
		InboundId:  inbound.Id,
		Email:      email,
		Total:      0,
		ExpiryTime: 0,
		Enable:     true,
		Reset:      0,
		Up:         0,
		Down:       0,
	}
	if err := tx.Create(&ct).Error; err != nil {
		tx.Rollback()
		log.Printf("default inbound: client traffic create failed: %v", err)
		return err
	}

	if err := tx.Create(&model.HistoryOfSeeders{SeederName: defaultVLESSRealitySeeder}).Error; err != nil {
		tx.Rollback()
		return err
	}

	return tx.Commit().Error
}

func buildRealityStream(privateKey, publicKey, shortID string) map[string]any {
	return map[string]any{
		"show":         false,
		"xver":         0,
		"target":       "www.microsoft.com:443",
		"serverNames":  []string{"www.microsoft.com"},
		"privateKey":   privateKey,
		"minClientVer": "",
		"maxClientVer": "",
		"maxTimediff":  0,
		"shortIds":     []string{shortID},
		"mldsa65Seed":  "",
		"settings": map[string]any{
			"publicKey":       publicKey,
			"fingerprint":     "chrome",
			"serverName":      "",
			"spiderX":         "/",
			"mldsa65Verify":   "",
		},
	}
}
