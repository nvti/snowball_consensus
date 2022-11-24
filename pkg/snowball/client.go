package snowball

type Client interface {
	// Returns the currently preferred choice
	Preference() []byte
}
