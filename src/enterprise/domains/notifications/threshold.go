package notifications

type (
	threshold struct {
		value int
	}
	Threshold interface {
		Get() int
		Set(int)
	}
)

func NewThreshold(value int) Threshold      { return &threshold{value} }
func (t *threshold) Get() int               { return t.value }
func (t *threshold) Set(value int)          { t.value = value }
func (t *threshold) Equal(o Threshold) bool { return t.value == o.Get() }
