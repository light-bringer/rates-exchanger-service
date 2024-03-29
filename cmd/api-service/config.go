package main

import "time"

const (
	SyncURL        = "https://www.ecb.europa.eu/stats/eurofxref/eurofxref-hist-90d.xml"
	SyncInterval   = 15 * time.Second
	DeleteInterval = 1 * time.Minute
	ServerTimeout  = 15 * time.Second
	DeletionDays   = 30
)
