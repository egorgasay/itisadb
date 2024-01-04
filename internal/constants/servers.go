package constants

import "time"

const MainStorageNumber = 1

const AutoServerNumber = 0

const MaxServerTries = 3

const ServerConnectTimeout = 5 * time.Second

const (
	DeleteFromAllServers = -1
	SetToAllServers      = -1
)
