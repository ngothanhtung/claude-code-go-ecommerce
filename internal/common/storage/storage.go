package storage

import "context"

// FileMeta describes an uploaded file.
type FileMeta struct {
	URL    string `json:"url"`
	Path   string `json:"path"`
	Size   int64  `json:"size"`
	MIME   string `json:"mime"`
	Width  int    `json:"width,omitempty"`
	Height int    `json:"height,omitempty"`
}

// Storage abstracts a file backend.
type Storage interface {
	Save(ctx context.Context, data []byte, mime string) (*FileMeta, error)
	Delete(ctx context.Context, path string) error
	URL(path string) string
}
