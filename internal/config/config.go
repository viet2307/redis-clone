package config

var Protocol = "tcp"
var Port = ":3000"

var MaxKeyNum int = 1000000
var EvictionRatio = 0.1
var EvictionPolicy string = "allkeys-lru"
