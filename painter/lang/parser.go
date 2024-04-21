package lang

import (
	"bufio"
	"fmt"
	"io"
	"slices"
	"strings"

	"github.com/aim4ik11/architecture-lab-3/painter"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
}

type OperationStruct struct {
	operation   painter.Operation
	commandName string
}

func complexOperation(operation painter.OperationFunc) painter.Operation {
	if operation == nil {
		return nil
	}
	return painter.OperationFunc(operation)
}

func (p *Parser) CommandParser(commandName string, opList []OperationStruct, args []string) painter.Operation {
	switch commandName {
	case "white":
		return painter.OperationFunc(painter.WhiteFill)
	case "green":
		return painter.OperationFunc(painter.GreenFill)
	case "bgrect":
		return complexOperation(painter.DrawRectangle(args))
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
	// var res []painter.Operation
	var res []OperationStruct

	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	i := 0
	bgRectLastIndex := -1
	for scanner.Scan() {
		commandLine := scanner.Text()

		sliced := strings.Split(commandLine, " ")
		args := sliced[1:]
		switch sliced[0] {
		case "reset":
			{
				bgRectLastIndex = -1
				res = []OperationStruct{}
				i = 0
			}
		case "bgrect":
			{
				fmt.Println(bgRectLastIndex)
				if bgRectLastIndex != -1 {
					res = slices.Delete[[]OperationStruct](res, bgRectLastIndex, bgRectLastIndex+1)
					i -= 1
				}
				bgRectLastIndex = i
			}
		case "update":
			{
				bgRectLastIndex = -1
			}
		}
		appendAction := p.CommandParser(sliced[0], res, args)
		curCommandOperation := OperationStruct{
			operation:   appendAction,
			commandName: sliced[0],
		}
		if appendAction != nil {
			res = append(res, curCommandOperation)
		}
		i += 1
	}

	var resOperations []painter.Operation
	for _, value := range res {
		resOperations = append(resOperations, value.operation)
	}

	return resOperations, nil
}
