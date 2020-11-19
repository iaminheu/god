package internal

import "time"

const (
	endpointsSeparator = ","
	autoSyncInterval   = time.Minute
	coolDownInterval   = time.Second
	dialTimeout        = time.Second * 5
	dialKeepAliveTime  = time.Second * 5
	requestTimeout     = time.Second * 3
	Delimiter          = '/'
)

var (
	NewClient      = DialClient
	DialTimeout    = dialTimeout
	RequestTimeout = requestTimeout
)
