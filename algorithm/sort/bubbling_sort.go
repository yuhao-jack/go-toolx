package sort

type Compare interface {
	int | int8 | int16 | int32 | int64 | uint | uint8 | uint16 | uint32 | uint64 | float32 | float64
}

func BubblingSort[K Compare](inputArr *[]K) {
	if inputArr == nil {
		return
	}
	length := len(*inputArr)
	for i := 0; i < length; i++ {
		for j := i + 1; j < length; j++ {
			if (*inputArr)[i] > (*inputArr)[j] {
				tmp := (*inputArr)[i]
				(*inputArr)[i] = (*inputArr)[j]
				(*inputArr)[j] = tmp
			}
		}
	}
}
