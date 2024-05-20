package constants

import "time"

const LocalServerNumber = 1

const AutoServerNumber int32 = 0

const MaxServerTries = 3

const ServerConnectTimeout = 5 * time.Second

const (
	DeleteFromAllServers = -1
	SetToAllServers      = -1
)
