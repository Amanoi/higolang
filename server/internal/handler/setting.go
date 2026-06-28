package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"

	"higolang/server/internal/model"
)

// SettingHandler handles site settings.
type SettingHandler struct {
	db *gorm.DB
}

// NewSettingHandler creates a new SettingHandler.
func NewSettingHandler(db *gorm.DB) *SettingHandler {
	return &SettingHandler{db: db}
}

// GetPublicSettings returns publicly visible settings.
func (h *SettingHandler) GetPublicSettings(c *gin.Context) {
	settings := h.getAllSettings()
	public := gin.H{
		"site_name":   settings["site_name"],
		"description": settings["description"],
		"logo_url":    settings["logo_url"],
	}
	OK(c, public)
}

// AdminGetSettings returns all settings (admin).
func (h *SettingHandler) AdminGetSettings(c *gin.Context) {
	OK(c, h.getAllSettings())
}

type settingsInput struct {
	Settings map[string]string `json:"settings"`
}

// AdminUpdateSettings updates site settings (admin).
func (h *SettingHandler) AdminUpdateSettings(c *gin.Context) {
	var input settingsInput
	if err := c.ShouldBindJSON(&input); err != nil {
		Fail(c, http.StatusBadRequest, err.Error())
		return
	}

	for key, value := range input.Settings {
		var setting model.SiteSetting
		result := h.db.Where("key = ?", key).First(&setting)
		if result.Error == gorm.ErrRecordNotFound {
			h.db.Create(&model.SiteSetting{Key: key, Value: value})
		} else {
			h.db.Model(&setting).Update("value", value)
		}
	}

	OK(c, nil)
}

func (h *SettingHandler) getAllSettings() map[string]string {
	var settings []model.SiteSetting
	h.db.Find(&settings)

	result := make(map[string]string)
	for _, s := range settings {
		result[s.Key] = s.Value
	}
	return result
}
