package models

import (
	"database/sql/driver"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"math"
	"strconv"
	"strings"
	"time"
)

const (
	UserKey       = "user"
	ImprintKey    = "imprint"
	AuthCookieKey = "wakapi_auth"
)

type MigrationFunc func(db *gorm.DB) error

type KeyStringValue struct {
	Key   string `gorm:"primary_key"`
	Value string `gorm:"type:text"`
}

type Filters struct {
	Project  string
	OS       string
	Language string
	Editor   string
	Machine  string
}

func NewFiltersWith(entity uint8, key string) *Filters {
	switch entity {
	case SummaryProject:
		return &Filters{Project: key}
	case SummaryOS:
		return &Filters{Project: key}
	case SummaryLanguage:
		return &Filters{Project: key}
	case SummaryEditor:
		return &Filters{Project: key}
	case SummaryMachine:
		return &Filters{Project: key}
	}
	return &Filters{}
}

func (f *Filters) First() (bool, uint8, string) {
	if f.Project != "" {
		return true, SummaryProject, f.Project
	} else if f.OS != "" {
		return true, SummaryOS, f.OS
	} else if f.Language != "" {
		return true, SummaryLanguage, f.Language
	} else if f.Editor != "" {
		return true, SummaryEditor, f.Editor
	} else if f.Machine != "" {
		return true, SummaryMachine, f.Machine
	}
	return false, 0, ""
}

type CustomTime time.Time

func (j *CustomTime) UnmarshalJSON(b []byte) error {
	s := strings.Replace(strings.Trim(string(b), "\""), ".", "", 1)
	i, err := strconv.ParseInt(s, 10, 64)
	if err != nil {
		return err
	}
	t := time.Unix(0, i*int64(math.Pow10(19-len(s))))
	*j = CustomTime(t)
	return nil
}

func (j *CustomTime) Scan(value interface{}) error {
	switch value.(type) {
	case string:
		t, err := time.Parse("2006-01-02 15:04:05-07:00", value.(string))
		if err != nil {
			return errors.New(fmt.Sprintf("unsupported date time format: %s", value))
		}
		*j = CustomTime(t)
	case int64:
		*j = CustomTime(time.Unix(value.(int64), 0))
		break
	case time.Time:
		*j = CustomTime(value.(time.Time))
		break
	default:
		return errors.New(fmt.Sprintf("unsupported type: %T", value))
	}
	return nil
}

func (j CustomTime) Value() (driver.Value, error) {
	return time.Time(j), nil
}

func (j CustomTime) String() string {
	t := time.Time(j)
	return t.Format("2006-01-02 15:04:05.000")
}

func (j CustomTime) Time() time.Time {
	return time.Time(j)
}
