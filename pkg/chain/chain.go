package chain

import (
	"crypto/sha256"
)

var (
	HashFunc = sha256.Sum256
)

const (
	HashSize = sha256.Size
)

type Block interface {
	Equal(Block) bool
}

type Chain interface {
	Length() int
	Get(int) (Block, error)
	Add(Block) error
	Set(int, Block) error
}

type Client[T any] interface {
	Preference(i int) (T, error)
}
