package storage

import (
	"bytes"
	"context"
	"image"
	_ "image/jpeg"
	_ "image/png"
	"os"
	"path/filepath"
	"time"

	"github.com/google/uuid"

	apperr "github.com/ngothanhtung/go-tutorials/internal/common/errors"
	"github.com/ngothanhtung/go-tutorials/internal/config"
)

var allowedMIME = map[string]bool{
	"image/jpeg":      true,
	"image/png":       true,
	"image/webp":      true,
	"application/pdf": true,
}

type Local struct {
	dir        string
	publicBase string
	maxBytes   int64
}

func NewLocal(cfg config.UploadConfig) *Local {
	return &Local{dir: cfg.Dir, publicBase: cfg.PublicBaseURL, maxBytes: cfg.MaxBytes()}
}

func (l *Local) Save(ctx context.Context, data []byte, mime string) (*FileMeta, error) {
	if !allowedMIME[mime] {
		return nil, apperr.NewUploadUnsupported()
	}
	if int64(len(data)) > l.maxBytes {
		return nil, apperr.NewUploadTooLarge()
	}
	ext := extForMIME(mime)
	rel := time.Now().Format("2006/01") + "/" + uuid.NewString() + ext
	full := filepath.Join(l.dir, rel)
	if err := os.MkdirAll(filepath.Dir(full), 0o755); err != nil {
		return nil, apperr.NewInternal("mkdir upload", err)
	}
	if err := os.WriteFile(full, data, 0o644); err != nil {
		return nil, apperr.NewInternal("write upload", err)
	}
	meta := &FileMeta{
		URL:  l.publicBase + "/static/" + filepath.ToSlash(rel),
		Path: rel,
		Size: int64(len(data)),
		MIME: mime,
	}
	if isImage(mime) {
		if cfg, format, err := image.DecodeConfig(bytes.NewReader(data)); err == nil && format != "" {
			meta.Width = cfg.Width
			meta.Height = cfg.Height
		}
	}
	return meta, nil
}

func (l *Local) Delete(ctx context.Context, path string) error {
	full := filepath.Join(l.dir, path)
	_ = os.Remove(full)
	return nil
}

func (l *Local) URL(path string) string {
	return l.publicBase + "/static/" + filepath.ToSlash(path)
}

func extForMIME(mime string) string {
	switch mime {
	case "image/jpeg":
		return ".jpg"
	case "image/png":
		return ".png"
	case "image/webp":
		return ".webp"
	case "application/pdf":
		return ".pdf"
	default:
		return ""
	}
}

func isImage(mime string) bool {
	return mime == "image/jpeg" || mime == "image/png" || mime == "image/webp"
}
