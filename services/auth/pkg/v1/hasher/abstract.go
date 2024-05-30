package hasher

type Hasher interface {
	Hash(s string) (string, error)
	Check(s string, hash string) bool
}
