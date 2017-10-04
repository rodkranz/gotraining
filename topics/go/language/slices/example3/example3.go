// All material is licensed under the Apache License Version 2.0, January 2004
// http://www.apache.org/licenses/LICENSE-2.0

// Sample program to show how to takes slices of slices to create different
// views of and make changes to the underlying array.
package main

import "fmt"

func main() {

	// Create a slice with a length of 5 elements and a capacity of 8.
	slice1 := make([]string, 5, 8)
	slice1[0] = "Apple"
	slice1[1] = "Orange"
	slice1[2] = "Banana"
	slice1[3] = "Grape"
	slice1[4] = "Plum"

	fmt.Print("slice 1 ")
	inspectSlice(slice1)

	// Take a slice of slice1. We want just indexes 2 and 3.
	// Parameters are [starting_index : (starting_index + length)]
	slice2 := slice1[2:4]
	fmt.Print("slice 2 ")
	inspectSlice(slice2)

	// Make a new slice big enough to hold elements of slice 1 and copy the
	// values over using the builtin copy function.
	slice3 := make([]string, len(slice1))
	copy(slice3, slice1)
	fmt.Print("slice 3 ")
	inspectSlice(slice3)

	fmt.Println("*************************")

	// Change the value of the index 0 of slice2.
	slice2[0] = "CHANGED"

	// Display the change across all existing slices.
	fmt.Print("slice 1 ")
	inspectSlice(slice1)
	fmt.Print("slice 2 ")
	inspectSlice(slice2)
	fmt.Print("slice 3 ")
	inspectSlice(slice3)
}

// inspectSlice exposes the slice header for review.
func inspectSlice(slice []string) {
	fmt.Printf("Length[%d] Capacity[%d]\n", len(slice), cap(slice))
	for i := range slice {
		fmt.Printf("[%d] %p %s\n",
			i,
			&slice[i],
			slice[i])
	}
}
