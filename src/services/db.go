package services

import "sync"

var once sync.Once
var db *sync.Map

func getDB() *sync.Map {
	once.Do(func() {
		db = new(sync.Map)
	})
	return db
}

func GetValue(key string) string {
	v, ok := getDB().Load(key)
	if !ok {
		return ""
	}
	return v.(string)
}

func SetKeyValue(key, value string) {
	getDB().Store(key, value)
}
