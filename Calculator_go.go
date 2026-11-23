package main

import "fmt"

func add(a, b float64) float64 {
	return a + b
}

func subtract(a, b float64) float64 {
	return a - b
}

func multiply(a, b float64) float64 {
	return a * b
}

func divide(a, b float64) (float64, error) {
	if b == 0 {
		return 0, fmt.Errorf("cannot divide by zero")
	}
	return a / b, nil
}

func main() {
	for {
		var choice int

		fmt.Println("\n1. Add")
		fmt.Println("2. Subtract")
		fmt.Println("3. Multiply")
		fmt.Println("4. Divide")
		fmt.Println("5. Exit")

		fmt.Print("Enter choice: ")
		_, err := fmt.Scan(&choice)

		if err != nil {
			fmt.Println("Invalid input")
			return
		}

		if choice == 5 {
			fmt.Println("Exiting calculator")
			break
		}

		var a, b float64

		fmt.Print("Enter two numbers: ")
		_, err = fmt.Scan(&a, &b)

		if err != nil {
			fmt.Println("Invalid numbers")
			return
		}

		switch choice {
		case 1:
			fmt.Println("Result:", add(a, b))
		case 2:
			fmt.Println("Result:", subtract(a, b))
		case 3:
			fmt.Println("Result:", multiply(a, b))
		case 4:
			result, err := divide(a, b)
			if err != nil {
				fmt.Println("Error:", err)
			} else {
				fmt.Println("Result:", result)
			}
		default:
			fmt.Println("Invalid choice")
		}
	}
}
