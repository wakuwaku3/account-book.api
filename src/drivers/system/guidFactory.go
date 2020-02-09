package system

import (
	"github.com/google/uuid"
	"github.com/wakuwaku3/account-book.api/src/enterprise/domains/core"
)

type (
	guidFactory struct{}
)

func newGuidFactory() core.GuidFactory { return &guidFactory{} }
func (t *guidFactory) Create() (*string, error) {
	u, err := uuid.NewRandom()
	if err != nil {
		return nil, err
	}
	uu := u.String()
	return &uu, nil
}
