package datastructure

type KeySpaceStat struct {
	Key     int64
	Expires int64
}

var HashKeySpace = KeySpaceStat{Key: 0, Expires: 0}
