package lang

import (
	"bufio"
	"io"
	"strings"

	"github.com/aim4ik11/architecture-lab-3/painter"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
}

func ahuetOperation(operation painter.OperationFunc) painter.Operation {
	if operation == nil {
		return nil
	}
	return painter.OperationFunc(operation)
}

func (p *Parser) CommandParser(commandName string, args []string) painter.Operation {
	switch commandName {
	case "white":
		return painter.OperationFunc(painter.WhiteFill)
	case "green":
		return painter.OperationFunc(painter.GreenFill)
	case "bgrect":
		return ahuetOperation(painter.DrawRectangle(args))
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
