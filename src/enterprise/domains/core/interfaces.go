package core

type (
	GuidFactory interface {
		Create() (*string, error)
	}
)
