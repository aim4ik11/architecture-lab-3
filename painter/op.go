package painter

import (
	"fmt"
	"image"
	"image/color"
	"strconv"

	"golang.org/x/exp/shiny/screen"
)

func getCordsByArgs(width int, height int, floatArgs []float64) []int {
	cords := make([]int, 4)

	fWidth := float64(width)
	fHeight := float64(height)

	cords[0] = int(fWidth * floatArgs[0])
	cords[1] = int(fHeight * floatArgs[1])
	cords[2] = int(fWidth * floatArgs[2])
	cords[3] = int(fHeight * floatArgs[3])

	return cords
}

func convertArgs(args []string) ([]float64, error) {
	parsedValues := make([]float64, len(args))
	for i, str := range args {
		num, err := strconv.ParseFloat(str, 64)
		if err != nil {
			return nil, err
		}
		parsedValues[i] = num
	}
	return parsedValues, nil
}

// Operation змінює вхідну текстуру.
type Operation interface {
	// Do виконує зміну операції, повертаючи true, якщо текстура вважається готовою для відображення.
	Do(t screen.Texture) (ready bool)
}

// OperationList групує список операції в одну.
type OperationList []Operation

func (ol OperationList) Do(t screen.Texture) (ready bool) {
	for _, o := range ol {
		ready = o.Do(t) || ready
	}
	return
}

// UpdateOp операція, яка не змінює текстуру, але сигналізує, що текстуру потрібно розглядати як готову.
var UpdateOp = updateOp{}

type updateOp struct{}

func (op updateOp) Do(t screen.Texture) bool { return true }

// OperationFunc використовується для перетворення функції оновлення текстури в Operation.
type OperationFunc func(t screen.Texture)

func (f OperationFunc) Do(t screen.Texture) bool {
	f(t)
	return false
}

// WhiteFill зафарбовує тестуру у білий колір. Може бути викоистана як Operation через OperationFunc(WhiteFill).
func WhiteFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.White, screen.Src)
}

// func Test(t screen.Texture) {
// 	t.Fill()
// }

// GreenFill зафарбовує тестуру у зелений колір. Може бути викоистана як Operation через OperationFunc(GreenFill).
func GreenFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.RGBA{G: 0xff, A: 0xff}, screen.Src)
}

func BlackFill(t screen.Texture) {
	t.Fill(t.Bounds(), color.Black, screen.Src)
}

func DrawRectangle(args []string) OperationFunc {
	if len(args) < 4 {
		fmt.Println("Not enough args to draw a rect")
		return nil
	}
	floatArgs, err := convertArgs(args)
	if err != nil {
		fmt.Println("Error parsing string")
		fmt.Println(err)
		return nil
	}
	return func(t screen.Texture) {

		cords := getCordsByArgs(t.Bounds().Dx(), t.Bounds().Dy(), floatArgs)

		startPoint := image.Point{int(cords[0]), int(cords[1])}
		endPoint := image.Point{int(cords[2]), int(cords[3])}
		t.Fill(image.Rectangle{startPoint, endPoint}, color.White, screen.Src)
	}
}
