package internal

import "time"

const (
	endpointsSeparator = ","
	autoSyncInterval   = time.Minute
	coolDownInterval   = time.Second
	dialTimeout        = 5 * time.Second
	dialKeepAliveTime  = 5 * time.Second
	requestTimeout     = 3 * time.Second
	Delimiter          = '/'
)

var (
	NewClient      = DialClient
	DialTimeout    = dialTimeout
	RequestTimeout = requestTimeout
)
