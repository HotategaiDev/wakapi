package utils

import (
	"errors"
	"fmt"
	"regexp"
	"time"

	"github.com/muety/wakapi/models"
)

func ParseDate(date string) (time.Time, error) {
	return time.Parse("2006-01-02 15:04:05", date)
}

func FormatDate(date time.Time) string {
	return date.Format("2006-01-02 15:04:05")
}

func FormatDateHuman(date time.Time) string {
	return date.Format("Mon, 02 Jan 2006 15:04")
}

func ParseUserAgent(ua string) (string, string, error) {
	re := regexp.MustCompile(`^wakatime\/[\d+.]+\s\((\w+).*\)\s.+\s(\w+)\/.+$`)
	groups := re.FindAllStringSubmatch(ua, -1)
	if len(groups) == 0 || len(groups[0]) != 3 {
		return "", "", errors.New("failed to parse user agent string")
	}
	return groups[0][1], groups[0][2], nil
}

func MakeConnectionString(config *models.Config) string {
	switch config.DbDialect {
	case "mysql":
		return mySqlConnectionString(config)
	case "postgres":
		return postgresConnectionString(config)
	case "sqlite3":
		return sqliteConnectionString(config)
	}
	return ""
}

func mySqlConnectionString(config *models.Config) string {
	location, _ := time.LoadLocation("Local")
	return fmt.Sprintf(
		"%s:%s@tcp(%s:%d)/%s?charset=utf8&parseTime=true&loc=%s",
		config.DbUser,
		config.DbPassword,
		config.DbHost,
		config.DbPort,
		config.DbName,
		location.String(),
	)
}

func postgresConnectionString(config *models.Config) string {
	return fmt.Sprintf("host=%s port=%d user=%s dbname=%s password=%s sslmode=disable",
		config.DbHost,
		config.DbPort,
		config.DbUser,
		config.DbName,
		config.DbPassword,
	)
}

func sqliteConnectionString(config *models.Config) string {
	return config.DbName
}
