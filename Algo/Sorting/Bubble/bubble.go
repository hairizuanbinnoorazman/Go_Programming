package main

import "fmt"

func BubbleSort(items []int) {
	if len(items) <= 1 {
		return
	}
	for {
		sortHappened := false
		for i := 0; i < len(items)-1; i++ {
			if items[i] > items[i+1] {
				temp := items[i]
				items[i] = items[i+1]
				items[i+1] = temp
				sortHappened = true
			}
		}
		if sortHappened == false {
			break
		}
	}

}

func main() {
	values := []int{4, 3, 2, 1}
	fmt.Println(values)
	BubbleSort(values)
	fmt.Println(values)
}
