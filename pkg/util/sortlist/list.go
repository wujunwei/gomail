package sortlist

const (
	DefaultLoadFactor = 1000
)

type SortedList[T comparable] struct {
	offset  int
	load    int
	maxes   []T
	lists   [][]T
	indexes []int //index sum tree
	size    int
	c       Compare[T]
}

func (l *SortedList[T]) Push(a T) {
	l.size++
	if len(l.maxes) == 0 {
		l.maxes = append(l.maxes, a)
		l.lists = append(l.lists, []T{a})
		return
	}
	pos := BisectLeft(l.maxes, l.c, a)
	if pos > 0 && l.maxes[pos-1] == a {
		pos--
	}
	if pos == len(l.maxes) {
		pos--
		l.maxes[pos] = a
		l.lists[pos] = append(l.lists[pos], a)
	} else {
		l.lists[pos] = InSort(l.lists[pos], l.c, a)
	}
	l.fresh(pos)
}

func (l *SortedList[T]) DeleteItem(a T) bool {
	if l.size == 0 {
		return false
	}

	pos := BisectLeft[T](l.maxes, l.c, a)
	if pos == len(l.maxes) {
		return false
	}

	var removed bool
	l.lists[pos], removed = RemoveSort(l.lists[pos], l.c, a)
	if !removed {
		return removed
	}
	l.size--
	if len(l.lists[pos]) == 0 {
		// delete maxes at pos
		copy(l.maxes[pos:], l.maxes[pos+1:])
		l.maxes = l.maxes[:len(l.maxes)-1]

		// delete lists at pos
		copy(l.lists[pos:], l.lists[pos+1:])
		l.lists = l.lists[:len(l.lists)-1]
		l.resetIndex()
	} else {
		l.maxes[pos] = l.lists[pos][len(l.lists[pos])-1]
		l.updateIndex(pos, -1)
	}
	return removed
}

func (l *SortedList[T]) Delete(index int) {
	if index >= l.size {
		return
	}
	var pos, in int
	if index == 0 {
		pos, in = 0, 0
	} else if index == l.size-1 {
		pos = len(l.lists) - 1
		in = len(l.lists[pos]) - 1
	} else {
		if len(l.indexes) == 0 {
			l.buildIndex()
		}
		pos, in = l.findPos(index)
	}
	l.size--
	l.lists[pos] = Remove(l.lists[pos], in)
	if len(l.lists[pos]) == 0 {
		// delete maxes at pos
		l.maxes = Remove(l.maxes, pos)

		// delete lists at pos
		copy(l.lists[pos:], l.lists[pos+1:])
		l.lists = l.lists[:len(l.lists)-1]
		l.resetIndex()
	} else {
		l.maxes[pos] = l.lists[pos][len(l.lists[pos])-1]
		l.updateIndex(pos, -1)
	}
}

func (l *SortedList[T]) Values() []T {
	res := make([]T, l.Size())
	i := 0
	l.Each(func(_ int, a T) {
		res[i] = a
		i++
	})
	return res
}

func (l *SortedList[T]) At(index int) (item T, found bool) {
	if index >= l.size {
		return
	}
	if index < len(l.lists[0]) {
		return l.lists[0][index], true
	}
	if index == l.size-1 {
		return l.maxes[len(l.maxes)-1], true
	}
	if len(l.indexes) == 0 {
		l.buildIndex()
	}
	pos, in := l.findPos(index)
	return l.lists[pos][in], true
}

func (l *SortedList[T]) Each(f ForEach[T]) {
	i := 0
	for _, list := range l.lists {
		for _, j := range list {
			f(i, j)
			i++
		}
	}
}

func (l *SortedList[T]) Has(a T) bool {
	if l.size == 0 {
		return false
	}
	pos := BisectLeft(l.maxes, l.c, a)
	if pos == len(l.maxes) {
		return false
	}
	index := BisectLeft(l.lists[pos], l.c, a)
	return l.lists[pos][index] == a
}

func (l *SortedList[T]) Floor(a T) (item T, ok bool) {
	if l.size == 0 {
		return
	}
	pos := BisectLeft(l.maxes, l.c, a)
	if pos == len(l.maxes) {
		return l.maxes[pos-1], true
	}
	index := BisectLeft(l.lists[pos], l.c, a)
	if index == 0 && l.lists[pos][0] != a {
		if pos == 0 {
			return
		} else {
			return l.maxes[pos-1], true
		}
	}
	if l.lists[pos][index] == a {
		return l.lists[pos][index], true
	}
	return l.lists[pos][index-1], true
}

func (l *SortedList[T]) Ceil(a T) (item T, ok bool) {
	if l.size == 0 {
		return
	}
	pos := BisectLeft(l.maxes, l.c, a)
	if pos == len(l.maxes) {
		return
	}
	index := BisectLeft(l.lists[pos], l.c, a)
	return l.lists[pos][index], true
}

// Index return the index of the position where the item to insert,and if the item exist or not.
func (l *SortedList[T]) Index(a T) (int, bool) {
	if l.size == 0 {
		return 0, false
	}
	pos := BisectLeft(l.maxes, l.c, a)
	if pos == len(l.maxes) {
		return l.size, false
	}
	if a == l.lists[0][0] {
		return 0, true
	}
	if a == l.maxes[0] {
		return len(l.lists[0]) - 1, true
	}
	if a == l.maxes[len(l.maxes)-1] {
		return l.size - 1, true
	}
	index := BisectLeft(l.lists[pos], l.c, a)
	exist := index < len(l.lists[pos]) && l.lists[pos][index] == a
	return l.locate(pos, index), exist
}

func (l *SortedList[T]) Empty() bool {
	return l.size == 0
}

func (l *SortedList[T]) Size() int {
	return l.size
}

func (l *SortedList[T]) Len() int {
	return l.size
}

func (l *SortedList[T]) Clear() {
	l.resetIndex()
	l.lists = [][]T{}
	l.maxes = []T{}
	l.size = 0
}

func (l *SortedList[T]) Top() (item T, ok bool) {
	if l.size == 0 {
		return
	}
	return l.maxes[len(l.maxes)-1], true
}

func (l *SortedList[T]) Bottom() (item T, ok bool) {
	if l.size == 0 {
		return
	}
	return l.lists[0][0], true
}

// fresh update the index and rebuild basic array if the load is greater than load factor after insert
func (l *SortedList[T]) fresh(pos int) {
	var zeroValue T
	listPosLen := len(l.lists[pos])
	if listPosLen > l.load {
		halfLen := listPosLen >> 1
		half := append([]T{}, l.lists[pos][halfLen:]...)
		l.lists[pos] = l.lists[pos][:halfLen]
		l.lists = append(l.lists, nil)
		copy(l.lists[pos+2:], l.lists[pos+1:])
		l.lists[pos+1] = half
		// update max
		l.maxes[pos] = l.lists[pos][halfLen-1]
		l.maxes = append(l.maxes, zeroValue)
		copy(l.maxes[pos+2:], l.maxes[pos+1:])
		l.maxes[pos+1] = l.lists[pos+1][len(l.lists[pos+1])-1]
		l.resetIndex()
	} else {
		l.maxes[pos] = l.lists[pos][listPosLen-1]
		l.updateIndex(pos, 1)
	}
}

func (l *SortedList[T]) buildIndex() {
	n := len(l.lists)
	rowLens := roundUpOf2((n + 1) / 2)
	l.offset = rowLens*2 - 1
	indexLens := l.offset + n

	indexes := make([]int, indexLens)
	for i, list := range l.lists { // fill row0
		indexes[len(indexes)-n+i] = len(list)
	}

	last := indexLens - n - rowLens
	for rowLens > 0 {
		for i := 0; i < rowLens; i++ {
			if (last+i)*2+1 >= indexLens {
				break
			}
			if (last+i)*2+2 >= indexLens {
				indexes[last+i] = indexes[(last+i)*2+1]
				break
			}
			indexes[last+i] = indexes[(last+i)*2+1] + indexes[(last+i)*2+2]
		}
		rowLens >>= 1
		last -= rowLens
	}
	l.indexes = indexes
}

func (l *SortedList[T]) updateIndex(pos, incr int) {
	if len(l.indexes) > 0 {
		child := l.offset + pos
		for child > 0 {
			l.indexes[child] += incr
			child = (child - 1) >> 1
		}
		l.indexes[0] += 1
	}
}

func (l *SortedList[T]) findPos(index int) (int, int) {
	if index < len(l.lists[0]) {
		return 0, index
	}
	pos := 0
	child := 1
	lenIndex := len(l.indexes)

	for child < lenIndex {
		indexChild := l.indexes[child]
		if index < indexChild {
			pos = child
		} else {
			index -= indexChild
			pos = child + 1
		}
		child = (pos << 1) + 1
	}
	return pos - l.offset, index
}

func (l *SortedList[T]) locate(pos, index int) int {
	if len(l.indexes) == 0 {
		l.buildIndex()
	}
	total := 0
	pos += l.offset
	for pos > 0 {
		if pos&1 == 0 {
			total += l.indexes[pos-1]
		}
		pos = (pos - 1) >> 1
	}
	return total + index
}
func (l *SortedList[T]) resetIndex() {
	l.indexes = []int{}
	l.offset = 0
}

func roundUpOf2(a int) int {
	i := 1
	for ; i < a; i <<= 1 {
	}
	return i
}

func NewSortedList[T comparable](c Compare[T], loadFactor int) SortedList[T] {
	if loadFactor <= 0 {
		loadFactor = DefaultLoadFactor
	}
	return SortedList[T]{load: loadFactor, c: c}
}
