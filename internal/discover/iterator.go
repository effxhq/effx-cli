package discover

func New(val string) *Iterator {
	return &Iterator{val: val}
}

type Iterator struct {
	ptr int
	val string
}

func (i Iterator) HasNext() bool {
	return i.ptr < len(i.val)
}

func (i *Iterator) Peek() string {
	if !i.HasNext() {
		return ""
	}
	return i.val[i.ptr : i.ptr+1]
}

func (i *Iterator) Next() string {
	v := i.Peek()
	i.ptr++
	return v
}

func generateIterators(list []string) []*Iterator {
	result := []*Iterator{}

	for _, s := range list {
		result = append(result, New(s))
	}

	return result
}
