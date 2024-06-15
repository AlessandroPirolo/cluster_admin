package rpcs

type Rpc interface {
	ToString() string
	Execute()
	Encode() ([]byte, error)
	Decode(rawMex []byte) error
}
