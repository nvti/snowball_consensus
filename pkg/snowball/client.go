package snowball

type PreferenceType comparable

type Client[T PreferenceType] interface {
	// Returns the currently preferred choice
	// May got error: network error, etc
	Preference() (T, error)
}
