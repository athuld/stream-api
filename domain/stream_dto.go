package domain

type Data struct {
	ID           int64  `json:"id,omitempty"`
	Hash         string `json:"hash,omitempty"`
	Filename     string `json:"filename,omitempty"`
	StreamLink   string `json:"stream_link,omitempty"`
	DownloadLink string `json:"download_link,omitempty"`
	CreatedAt    string `json:"created_at,omitempty"`
}
