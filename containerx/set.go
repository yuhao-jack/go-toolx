package containerx

// Set [K comparable]
// @Description:
type Set[K comparable] struct {
	innerMap map[K]struct{}
}

// NewSet [K comparable]
//
//	@Description:
//	@return *Set[K]
func NewSet[K comparable]() *Set[K] {
	return &Set[K]{
		innerMap: make(map[K]struct{}, 8),
	}
}

// Add
//
//	@Description:
//	@receiver r
//	@param element
//	@return *Set[K]
func (r *Set[K]) Add(element K) *Set[K] {
	r.innerMap[element] = struct{}{}
	return r
}

// Size
//
//	@Description:
//	@receiver r
//	@return int
func (r *Set[K]) Size() int {
	return len(r.innerMap)
}

// IsEmpty
//
//	@Description:
//	@receiver r
//	@return bool
func (r *Set[K]) IsEmpty() bool {
	return r.Size() == 0
}

// Remove
//
//	@Description:
//	@receiver r
//	@param element
//	@return *Set[K]
func (r *Set[K]) Remove(element K) *Set[K] {
	delete(r.innerMap, element)
	return r
}

// ForEach
//
//	@Description:
//	@receiver r
//	@param fx
func (r *Set[K]) ForEach(fx func(item K)) {
	for k, _ := range r.innerMap {
		fx(k)
	}
}

// Contains
//
//	@Description:
//	@receiver r
//	@param element
//	@return bool
func (r *Set[K]) Contains(element K) bool {
	_, ok := r.innerMap[element]
	return ok
}

// Clear
//
//	@Description:
//	@receiver r
//	@return *Set[K]
func (r *Set[K]) Clear() *Set[K] {
	for k, _ := range r.innerMap {
		delete(r.innerMap, k)
	}
	return r
}

// Elements
//
//	@Description:
//	@receiver r
//	@return elements
func (r *Set[K]) Elements() (elements []K) {
	for k, _ := range r.innerMap {
		elements = append(elements, k)
	}
	return elements
}

// AddAll
//
//	@Description:
//	@receiver r
//	@param elements
//	@return *Set[K]
func (r *Set[K]) AddAll(elements *[]K) *Set[K] {
	if elements != nil {
		for _, k := range *elements {
			r.Add(k)
		}
	}
	return r
}

// Filter
//
//	@Description:
//	@receiver r
//	@param fx
//	@return *Set[K]
func (r *Set[K]) Filter(fx func(item K) bool) *Set[K] {
	var elements []K
	for k, _ := range r.innerMap {
		if fx(k) {
			elements = append(elements, k)
		}
	}
	return NewSet[K]().AddAll(&elements)
}
