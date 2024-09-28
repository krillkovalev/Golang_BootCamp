package main

import (
	"flag"
	"fmt"
	"math"
	"slices"
)

func main() {
	nums := make([]float64, 0)
	var d, sum float64
	for {
		_, err := fmt.Scan(&d)
		if err != nil {
			break
		} else {
			nums = append(nums, d)
			sum += d
		}
	}

	meanFlag := flag.Bool("mean", false, "show mean of array")
	medianFlag := flag.Bool("median", false, "show median of array")
	modeFlag := flag.Bool("mode", false, "show mode of array")
	sdFlag := flag.Bool("sd", false, "show standart deviation of array")

	flag.Parse()

	if *meanFlag || *medianFlag || *modeFlag || *sdFlag {
		if *meanFlag {
			fmt.Printf("Mean: %v\n", sum/float64(len(nums)))
		}
		if *medianFlag {
			fmt.Printf("Median: %v\n", medianFunc(nums))
		}
		if *modeFlag {
			fmt.Printf("Mode: %v\n", modeFunc(nums))
		}
		if *sdFlag {
			fmt.Printf("SD: %v\n", standardDeviation(nums))
		}
	} else {
		fmt.Printf("Mean: %v\n", sum/float64(len(nums)))
		fmt.Printf("Median: %v\n", medianFunc(nums))
		fmt.Printf("Mode: %v\n", modeFunc(nums))
		fmt.Printf("SD: %v\n", standardDeviation(nums))
	}

}

func medianFunc(data []float64) float64 {
	dataCopy := make([]float64, len(data))
	copy(dataCopy, data)

	slices.Sort(dataCopy)

	var median float64
	l := len(dataCopy)
	if l == 0 {
		return 0
	} else if l%2 == 0 {
		median = (dataCopy[l/2-1] + dataCopy[l/2]) / 2
	} else {
		median = dataCopy[l/2]
	}

	return median
}

func modeFunc(data []float64) float64 {
	slices.Reverse(data)
	slices.Sort(data)

	max_count := 1
	res := data[0]
	curr_count := 1
	for i := 1; i < len(data); i++ {
		if data[i] == data[i-1] {
			curr_count++
		} else {
			curr_count = 1
		}
		if curr_count > max_count {
			max_count = curr_count
			res = data[i-1]
		}
	}

	return res
}

func standardDeviation(num []float64) float64 {
	var sum, mean, sd float64
	for i := 1; i <= len(num); i++ {
		num[i-1] = float64(i) + 123
		sum += num[i-1]
	}
	mean = sum / float64(len(num))
	fmt.Println("The mean of above array is:", mean)
	for j := 0; j < len(num); j++ {
		sd += math.Pow(num[j]-mean, 2)
	}
	sd = math.Sqrt(sd / float64(len(num)))
	return sd
}
