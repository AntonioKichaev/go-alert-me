package mgzip

type Zipper interface {
	Compress(data []byte) ([]byte, error)
	Decompress(data []byte) ([]byte, error)
	GetEncoding() string
}
