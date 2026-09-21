package config

import "time"

var App = struct {
	Production       bool
	Domain           string
	Port             string
	MySQLHost        string
	MySQLPort        string
	MySQLUser        string
	MySQLPass        string
	MySQLDB          string
	SodiumKey        string
	TokenName        string
	TokenExpiredTime time.Duration
	TokenRefreshTime time.Duration
	Version          string
}{
	Production:       false,
	Domain:           "http://localhost:5173,http://127.0.0.1:5173",
	Port:             "8082",
	MySQLHost:        "localhost",
	MySQLPort:        "3306",
	MySQLUser:        "root",
	MySQLPass:        "",
	MySQLDB:          "",
	SodiumKey:        "6621",
	TokenName:        "app-platform",
	TokenExpiredTime: 7 * 24 * time.Hour,
	TokenRefreshTime: 3 * 24 * time.Hour,
	Version:          "1.0.0",
}
