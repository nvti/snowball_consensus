package chain

import (
	"errors"
	"sync"
)

type SimpleBlock[T any] struct {
	Data T
}

// Chain a Linear chain implementation
type SimpleLinearChain[T any] struct {
	Blocks []*SimpleBlock[T]
	mu     sync.Mutex
}

func (c *SimpleLinearChain[T]) Length() int {
	return len(c.Blocks)
}

func (c *SimpleLinearChain[T]) Get(index int) (*SimpleBlock[T], error) {
	if index < 0 || index >= c.Length() {
		return nil, errors.New("index out of range")
	}

	return c.Blocks[index], nil
}

func (c *SimpleLinearChain[T]) Add(data T) error {
	c.mu.Lock()
	c.Blocks = append(c.Blocks, &SimpleBlock[T]{
		Data: data,
	})
	c.mu.Unlock()

	return nil
}

func (c *SimpleLinearChain[T]) Set(index int, data T) error {
	dataBlock, err := c.Get(index)
	if err != nil {
		return err
	}

	c.mu.Lock()
	dataBlock.Data = data
	c.mu.Unlock()

	return nil
}
