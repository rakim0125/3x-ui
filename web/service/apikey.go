package service

import (
	"crypto/sha256"
	"encoding/hex"
	"errors"
	"time"

	"github.com/mhsanaei/3x-ui/v2/database"
	"github.com/mhsanaei/3x-ui/v2/database/model"
	"github.com/mhsanaei/3x-ui/v2/util/random"
)

type ApiKeyService struct{}

func hashApiKey(raw string) string {
	h := sha256.Sum256([]byte(raw))
	return hex.EncodeToString(h[:])
}

// Create generates a new API key for the given user.
// It returns the created model (without the hash) and the raw key string
// that must be shown to the user exactly once.
func (s *ApiKeyService) Create(userId int, name string, expiresAt int64) (*model.ApiKey, string, error) {
	if name == "" {
		return nil, "", errors.New("api key name is required")
	}

	rawKey := "sk-" + random.Seq(45)
	hashed := hashApiKey(rawKey)
	prefix := rawKey[:10]

	apiKey := &model.ApiKey{
		Name:      name,
		Key:       hashed,
		Prefix:    prefix,
		UserId:    userId,
		CreatedAt: time.Now().Unix(),
		ExpiresAt: expiresAt,
	}

	db := database.GetDB()
	if err := db.Create(apiKey).Error; err != nil {
		return nil, "", err
	}

	return apiKey, rawKey, nil
}

// Validate checks an incoming raw API key.
// Returns the ApiKey record if valid, or an error if not found / expired.
func (s *ApiKeyService) Validate(rawKey string) (*model.ApiKey, error) {
	hashed := hashApiKey(rawKey)

	db := database.GetDB()
	apiKey := &model.ApiKey{}
	if err := db.Where("key = ?", hashed).First(apiKey).Error; err != nil {
		return nil, errors.New("invalid api key")
	}

	if apiKey.ExpiresAt > 0 && apiKey.ExpiresAt < time.Now().Unix() {
		return nil, errors.New("api key has expired")
	}

	return apiKey, nil
}

// List returns all API keys for a given user (raw key is never stored).
func (s *ApiKeyService) List(userId int) ([]*model.ApiKey, error) {
	db := database.GetDB()
	var keys []*model.ApiKey
	err := db.Where("user_id = ?", userId).Order("created_at DESC").Find(&keys).Error
	return keys, err
}

// Delete removes an API key by its id, scoped to the user.
func (s *ApiKeyService) Delete(userId, keyId int) error {
	db := database.GetDB()
	result := db.Where("id = ? AND user_id = ?", keyId, userId).Delete(&model.ApiKey{})
	if result.RowsAffected == 0 {
		return errors.New("api key not found")
	}
	return result.Error
}
