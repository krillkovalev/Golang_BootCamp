package main

import (
	"container/heap"
	"fmt"
	"log"
)

type PresentHeap []Present

type Present struct {
	Value int
	Size int
}


func main() {

	// presents := &PresentHeap{}
	// heap.Init(presents)
	// heap.Push(presents, Present{Value: 5, Size: 1})
	// heap.Push(presents, Present{Value: 4, Size: 5})
	// heap.Push(presents, Present{Value: 3, Size: 1})
	// heap.Push(presents, Present{Value: 5, Size: 2})
	
	// fmt.Println(getNCoolestPresents(*presents, 4))

	// presents := PresentHeap{}
	// presents = append(presents, Present{Value: 5, Size: 6})
	// presents = append(presents, Present{Value: 3, Size: 4})
	// presents = append(presents, Present{Value: 1, Size: 3})
	// presents = append(presents, Present{Value: 3, Size: 2})
    // presents = append(presents, Present{Value: 6, Size: 5})
	
	// fmt.Println(grabPresents(presents, 15))

	presents := PresentHeap{}
	presents = append(presents, Present{Value: 4, Size: 5})
	presents = append(presents, Present{Value: 3, Size: 1})
	presents = append(presents, Present{Value: 5, Size: 2})
	presents = append(presents, Present{Value: 5, Size: 1})
	
	fmt.Println(grabPresents(presents, 5))

}

// Len, Less, Swap для реализации интерфейса sort.Interface
func (h PresentHeap) Len() int { return len(h) }
func (h PresentHeap) Less(i, j int) bool {
	if h[i].Value > h[j].Value {
		return true
	} else if h[i].Value == h[j].Value {
		return h[i].Size < h[j].Size
	}
 
	return false
}
func (h PresentHeap) Equal(i, j int) bool { return h[i].Value == h[j].Value }
func (h PresentHeap) Swap(i, j int) { h[i], h[j] = h[j], h[i] }

func (h *PresentHeap) Push(x any) {
	// Push and Pop use pointer receivers because they modify the slice's length,
	// not just its contents.
	*h = append(*h, x.(Present))
}

func (h *PresentHeap) Pop() any {
	old := *h
	n := len(old)
	x := old[n-1]
	*h = old[0 : n-1]
	return x
}

func getNCoolestPresents(presents []Present, n int) []Present {
	if n <= 0 || n > len(presents) {
		log.Fatal("Incorrect n parameter")
	}

	result := []Present{}
	tempHeap := &PresentHeap{}
	heap.Init(tempHeap)

	for _, p := range presents {
		heap.Push(tempHeap, p)
	}
	
	for i := 0; i < n; i++ {
		result = append(result, heap.Pop(tempHeap).(Present))
	}

	return result


}

func grabPresents(presents []Present, capacity int) ([]Present, error){
	values := make([][]int, len(presents) + 1)
	for i := range values {
		values[i] = make([]int, capacity + 1)
	}

	result := []Present{}

	for i := 1; i <= len(presents); i++ {
		for w := 1; w <= capacity; w++ {
			if presents[i - 1].Size <= w {
				values[i][w] = max(values[i - 1][w], values[i - 1][w - presents[i - 1].Size] + presents[i - 1].Value)
			} else {
				values[i][w] = values[i - 1][w]
			}
		}
	}

	w := capacity
	for i := len(presents); i > 0 && w > 0; i-- {
		if values[i][w] != values[i-1][w] {
			result = append(result, presents[i-1])
			w -= presents[i-1].Size
		}
	}
	return result, nil 
}
