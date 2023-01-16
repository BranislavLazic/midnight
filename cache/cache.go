package cache

type Internal interface {
	Set(key string, entry []byte) error
	Get(key string) ([]byte, error)
}
