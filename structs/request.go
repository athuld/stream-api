package structs

type CacheHeaders struct {
	CacheHit   bool   `json:"cache_hit,omitempty"`
	CacheKey   string `json:"cache_key,omitempty"`
}

type QueryParams struct {
	Hash      string `form:"hash" json:"hash,omitempty"`
	IPAddress string `form:"ip_address" json:"ip_address,omitempty"`
	Action    string `form:"action" json:"action,omitempty"`
}