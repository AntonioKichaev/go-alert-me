package hasher

import (
	"bytes"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"hash"
	"io"

	"go.uber.org/zap"

	"github.com/antoniokichaev/go-alert-me/internal/logger"
)

type Hasher struct {
	key string
	h   hash.Hash
	enc io.WriteCloser
	b   *bytes.Buffer
}

func NewHasher(key string) *Hasher {
	if len(key) == 0 {
		return nil
	}
	var buf bytes.Buffer
	return &Hasher{key: key,
		h:   hmac.New(sha256.New, []byte(key)),
		b:   &buf,
		enc: base64.NewEncoder(base64.StdEncoding, &buf),
	}
}

func (h *Hasher) Sign(data []byte) string {

	h.h.Reset()
	h.h.Write(data)
	_, err := h.enc.Write(h.h.Sum(nil))
	logger.GetLogger().Error("hasher.Sign() err", zap.Error(err))
	h.enc.Close()
	sign := h.b.String()
	h.b.Reset()
	return sign
}
