package chain

import (
	"errors"
	"sync"
)

type SimpleBlock struct {
	Data []byte
}

func (b *SimpleBlock) Equal(other *SimpleBlock) bool {
	return HashFunc(b.Data) == HashFunc(other.Data)
}

// Chain a Linear chain implementation
type SimpleLinearChain struct {
	Blocks []*SimpleBlock
	mu     sync.Mutex
}

func (c *SimpleLinearChain) Length() int {
	return len(c.Blocks)
}

func (c *SimpleLinearChain) Get(index int) (*SimpleBlock, error) {
	if index < 0 || index >= c.Length() {
		return nil, errors.New("index out of range")
	}

	return c.Blocks[index], nil
}

func (c *SimpleLinearChain) Add(data []byte) error {
	c.mu.Lock()
	c.Blocks = append(c.Blocks, &SimpleBlock{
		Data: data,
	})
	c.mu.Unlock()

	return nil
}

func (c *SimpleLinearChain) Set(index int, data []byte) error {
	dataBlock, err := c.Get(index)
	if err != nil {
		return err
	}

	c.mu.Lock()
	dataBlock.Data = data
	c.mu.Unlock()

	return nil
}
