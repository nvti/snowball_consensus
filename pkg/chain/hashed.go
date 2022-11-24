package chain

import (
	"errors"
	"sync"
)

type HashedBlock struct {
	Data []byte
	Hash [HashSize]byte
}

func NewHashedBlock(data []byte, lastHash [HashSize]byte) *HashedBlock {
	dataBlock := &HashedBlock{
		Data: data,
	}
	dataBlock.CalculateHash(lastHash)

	return dataBlock
}

func (d *HashedBlock) CalculateHash(lastHash [HashSize]byte) {
	d.Hash = HashFunc(append(lastHash[:], d.Data...))
}

func (d *HashedBlock) Equal(other *HashedBlock) bool {
	return d.Hash == other.Hash
}

// Chain a Linear chain implementation
type HashedLinearChain struct {
	Blocks []*HashedBlock
	mu     sync.Mutex
}

func (c *HashedLinearChain) Length() int {
	return len(c.Blocks)
}

func (c *HashedLinearChain) Get(index int) (*HashedBlock, error) {
	if index < 0 || index >= c.Length() {
		return nil, errors.New("index out of range")
	}

	return c.Blocks[index], nil
}

func (c *HashedLinearChain) getLastHash(index int) (lastHash [HashSize]byte, err error) {
	if index == 0 {
		return [HashSize]byte{0}, nil
	}

	lastBlock, err := c.Get(index - 1)
	if err != nil {
		return
	}

	return lastBlock.Hash, nil
}

func (c *HashedLinearChain) Add(data []byte) error {
	lastHash, err := c.getLastHash(c.Length())
	if err != nil {
		return err
	}

	c.mu.Lock()
	c.Blocks = append(c.Blocks, NewHashedBlock(data, lastHash))
	c.mu.Unlock()

	return nil
}

func (c *HashedLinearChain) Set(index int, data []byte) error {
	dataBlock, err := c.Get(index)
	if err != nil {
		return err
	}

	c.mu.Lock()
	defer c.mu.Unlock()
	dataBlock.Data = data
	lastHash, err := c.getLastHash(index)
	if err != nil {
		return err
	}

	// Re-calculating hashes
	for i := index; i < c.Length(); i++ {
		dataBlock, err = c.Get(i)
		if err != nil {
			return err
		}

		dataBlock.CalculateHash(lastHash)
		lastHash = dataBlock.Hash
	}

	return nil
}
