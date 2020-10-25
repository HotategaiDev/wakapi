package models

import (
	"fmt"
	"regexp"
	"time"
)

type Heartbeat struct {
	ID              uint       `gorm:"primary_key"`
	User            *User      `json:"-" gorm:"not null"`
	UserID          string     `json:"-" gorm:"not null; index:idx_time_user"`
	Entity          string     `json:"entity" gorm:"not null; index:idx_entity"`
	Type            string     `json:"type"`
	Category        string     `json:"category"`
	Project         string     `json:"project"`
	Branch          string     `json:"branch"`
	Language        string     `json:"language" gorm:"index:idx_language"`
	IsWrite         bool       `json:"is_write"`
	Editor          string     `json:"editor"`
	OperatingSystem string     `json:"operating_system"`
	Machine         string     `json:"machine"`
	Time            CustomTime `json:"time" gorm:"type:timestamp(3); default:CURRENT_TIMESTAMP(3); index:idx_time,idx_time_user"`
	languageRegex   *regexp.Regexp
}

func (h *Heartbeat) Valid() bool {
	return h.User != nil && h.UserID != "" && h.Time != CustomTime(time.Time{})
}

func (h *Heartbeat) Augment(customRules []*CustomRule) {
	for _, lang := range customRules {
		reg := fmt.Sprintf(".*%s$", lang.Extension)
		match, err := regexp.MatchString(reg, h.Entity)
		if match && err == nil {
			h.Language = lang.Language
			return
		}
	}
}
