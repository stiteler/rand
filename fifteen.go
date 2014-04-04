//boiler plate code ported from cs50x
package main 

import(
	"fmt"
	"time"
	"strconv"
	"os"
	"log"
)

const (
	MIN = 3
	MAX = 9
)

//lesson learned from this problem... CAPTIALIZE YOUR
//GO GLOBAL VARIABLES... hehe.
var D int

var Board [MAX][MAX]int

func main() {
	greet()

	if len(os.Args) != 2 {
		fmt.Println("Usage: fifteen d\n")
		log.Fatal(1)
	}


	//TODO so what's the better way than declaring a temp var? 
	//pointers? 
	tmp, err := strconv.Atoi(os.Args[1])
	D = tmp
	if err != nil {	panic(err) }
	if D < MIN || D > MAX {
		fmt.Printf("Board must be between %i x %i and %i x %i, inclusive.\n",
			MIN, MIN, MAX, MAX)
		log.Fatal(2)
	}

	initialize()

	for {
		draw()

		if won() {
			fmt.Println("ftw\n")
			break
		}

		tile := getInput()

		if !move(tile) {
			fmt.Println("\nIllegal Move")
			time.Sleep(time.Millisecond * 100)
		}

		time.Sleep(time.Millisecond * 100)
	}
	
}

func getInput() int {

	var value int
	fmt.Println("Tile to move: ")
	_, err := fmt.Scanf("%d", &value)
	if err != nil || value < 1 || value > D * D - 1 {
		fmt.Println("Illegal move")
		return getInput()
	} else  { return value }
}

func clear() {
	fmt.Printf("\033[2J")
	fmt.Printf("\033[%d;%dH", 0, 0)
}

func greet() {
	clear()
	fmt.Println("GAME OF FIFTEEN\n")
	time.Sleep(time.Millisecond * 500)
}

func initialize() {
	arr := make([]int, D*D, D*D)

	index := 0

	for i := D * D - 1; i > 0; i-- {
		arr[index] = i
		index++ 
	}

	//set 0 value to -1 so as not to display

	arr[D * D - 1] = -1

	//if array is even swap 2, 1
	if D % 2 == 0 {
		tmp := arr[D * D - 2];
		arr[D*D-2] = arr[D*D - 3]
		arr[D * D - 3] = tmp
	}

	//now populate the model
	arrayIndex := 0
	for row := 0; row < D; row ++ {
		for col := 0; col < D; col++ {
			Board[row][col] = arr[arrayIndex]
			arrayIndex++
		}
	}
}

func draw() {
	c := '\t'
	fmt.Printf("\n\n\n ")
	for row := 0; row < D; row++ {
		for col := 0; col < D; col ++ {
			if Board[row][col] == -1 {
				if col % 2 == 0 {
					fmt.Printf("%c<_<", c)
				} else {
					fmt.Printf("%c>_>", c)
				}
			} else {
				if Board[row][col] < 10 {
					fmt.Printf("%c %d", c, Board[row][col])
				}
				if Board[row][col] >= 10 {
					fmt.Printf("%c%d", c, Board[row][col])
				}
			}
		}
		fmt.Println("\n\n\n")
	}
}

func move(tile int) bool {
	var thisRow, thisCol int

	for row := 0; row < D; row++ {
		for col := 0; col < D; col++ {
			if Board[row][col] == tile {
				thisRow = row
				thisCol = col
				break;
			}
		}
	}

	

	//not sure if this is too "brute forcey", but I'm sure I
	//Could at least expand this by calling functions: 
	//ex: findRowCol(tile) above... etc
	if thisCol + 1 < D && Board[thisRow][thisCol + 1] == -1 {
        Board[thisRow][thisCol] = -1
        Board[thisRow][thisCol + 1] = tile
        return true
    } else if thisCol - 1 >= 0 && Board[thisRow][thisCol - 1] == -1 {
        Board[thisRow][thisCol] = -1
        Board[thisRow][thisCol - 1] = tile
        return true
    } else if thisRow + 1 < D && Board[thisRow + 1][thisCol] == -1 {
        Board[thisRow][thisCol] = -1
        Board[thisRow + 1][thisCol] = tile
        return true
    } else if thisRow - 1 >= 0 && Board[thisRow - 1][thisCol] == -1 {
        Board[thisRow][thisCol] = -1
        Board[thisRow - 1][thisCol] = tile
        return true
    } else {
        return false
    }
}

func won() bool {

	check := 0

	for i := 0; i < D; i++ {
		for j := 0; j < D; j++ {
			if i == (D-1) && j == (D-1) {
				return true
			}

			if Board[i][j] == check + 1 {
				check++
			} else {
				return false
			}
		}
	}

	return false
}