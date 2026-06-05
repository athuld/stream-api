package cache

import (
	"context"
	"encoding/json"
	"streamapi/datasource"
	"streamapi/structs"
	"strings"
)

func CreateNavCacheKey(hash string, operator string) string {
	return "nav:" + hash + ":" + operator
}

func CreateRecordCacheKey(hash string) string {
	return "record:" + hash
}

func GetSingleDataFromCache(ctx context.Context, cacheKey string) (*structs.Data, bool) {
	if strings.HasPrefix(cacheKey, "nav:") {
		recordKey, err := datasource.RedisClient.Get(ctx, cacheKey).Result()
		if err != nil {
			return nil, false
		}
		cacheKey = CreateRecordCacheKey(recordKey)
	}
	data, hit := getJsonDataFromCache(ctx, cacheKey)
	if !hit {
		return nil, false
	}
	return data, true
}

func getJsonDataFromCache(ctx context.Context, cacheKey string) (*structs.Data, bool) {
	jsonData, err := datasource.RedisClient.Get(ctx, cacheKey).Result()
	if err != nil {
		return nil, false
	}
	var data structs.Data
	err = json.Unmarshal([]byte(jsonData), &data)
	if err != nil {
		return nil, false
	}
	return &data, true
}

func SetDataToCache(ctx context.Context, cacheKey string, data structs.Data) {
	hash := data.Hash
	if strings.HasPrefix(cacheKey, "nav:") {
		datasource.RedisClient.Set(ctx, cacheKey, hash, 0)
	}
	cacheKey = CreateRecordCacheKey(hash)
	setJsonDataToCache(ctx, cacheKey, data)
}

func setJsonDataToCache(ctx context.Context, cacheKey string, data structs.Data) {
	jsonData, err := json.Marshal(data)
	if err != nil {
		return
	}
	datasource.RedisClient.Set(ctx, cacheKey, jsonData, 0)
}
