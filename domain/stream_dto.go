package domain

type Data struct {
	ID           int64  `json:"id,omitempty"`
	Hash         string `json:"hash,omitempty"`
	Filename     string `json:"filename,omitempty"`
	StreamLink   string `json:"stream_link,omitempty"`
	DownloadLink string `json:"download_link,omitempty"`
	HasThumb     int    `json:"has_thumb,omitempty"`
	ThumbUrl     string `json:"thumb_url,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
}
