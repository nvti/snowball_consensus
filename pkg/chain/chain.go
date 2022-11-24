package snowball

import (
	"crypto/sha256"
	"errors"
	"sync"
)

var (
	HashFunc = sha256.Sum256
)

const (
	HashSize = sha256.Size
)

type DataBlock struct {
	Data []byte
}

// Chain a Linear chain implementation
type LinearChain struct {
	Blocks []*DataBlock
	mu     sync.Mutex
}

func (c *LinearChain) Length() int {
	return len(c.Blocks)
}

func (c *LinearChain) Get(index int) (*DataBlock, error) {
	if index < 0 || index >= c.Length() {
		return nil, errors.New("index out of range")
	}

	return c.Blocks[index], nil
}

func (c *LinearChain) Add(data []byte) error {
	c.mu.Lock()
	c.Blocks = append(c.Blocks, &DataBlock{
		Data: data,
	})
	c.mu.Unlock()

	return nil
}

func (c *LinearChain) Set(index int, data []byte) error {
	dataBlock, err := c.Get(index)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	dataBlock.Data = data

	return nil
}
