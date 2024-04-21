package lang

import (
	"bufio"
	"fmt"
	"image"
	"image/color"
	"io"
	"strconv"
	"strings"

	"github.com/aim4ik11/architecture-lab-3/painter"
	"golang.org/x/exp/shiny/screen"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
}

func convertArgs(args []string) []float64 {
	parsedValues := make([]float64, len(args))
	for i, str := range args {
		num, err := strconv.ParseFloat(str, 64)
		if err != nil {
			fmt.Println("Error parsing string '%s': %v", str, err)
			return nil
		}
		parsedValues[i] = num
	}
	return parsedValues
}

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

func drawRectangle(args []float64) func(t screen.Texture) {
	return func(t screen.Texture) {

		cords := getCordsByArgs(t.Bounds().Dx(), t.Bounds().Dy(), args)

		startPoint := image.Point{int(cords[0]), int(cords[1])}
		endPoint := image.Point{int(cords[2]), int(cords[3])}
		t.Fill(image.Rectangle{startPoint, endPoint}, color.White, screen.Src)
	}
}

func (p *Parser) CommandParser(commandName string, args []string) painter.Operation {
	switch commandName {
	case "white":
		return painter.OperationFunc(painter.WhiteFill)
	case "green":
		return painter.OperationFunc(painter.GreenFill)
	case "bgrect":
		{
			if len(args) < 4 {
				fmt.Println("Not enough args to draw a rect")
				return nil
			}
			floatArgs := convertArgs(args)
			operation := drawRectangle(floatArgs)
			if operation == nil {
				return nil
			}
			return painter.OperationFunc(operation)
		}
	// case "move":
	// 	return painter.OperationFunc(painter.)
	// case "figure":
	// 	return painter.OperationFunc(painter.WhiteFill)
	case "update":
		return painter.UpdateOp
	case "reset":
		return painter.OperationFunc(painter.BlackFill)
	}
	return nil
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	var res []painter.Operation

	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		commandLine := scanner.Text()

		sliced := strings.Split(commandLine, " ")
		args := sliced[1:]
		if sliced[0] == "reset" {
			res = []painter.Operation{}
		}
		appendAction := p.CommandParser(sliced[0], args)
		if appendAction != nil {
			res = append(res, appendAction)
		}
	}

	return res, nil
}
