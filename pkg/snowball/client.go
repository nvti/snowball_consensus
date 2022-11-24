package snowball

type PreferenceType comparable

type Client[T PreferenceType] interface {
	// Returns the currently preferred choice
	Preference() T
}
