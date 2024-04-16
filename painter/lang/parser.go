package lang

import (
	"bufio"
	"fmt"
	"io"
	"strings"

	"github.com/roman-mazur/architecture-lab-3/painter"
)

// Parser уміє прочитати дані з вхідного io.Reader та повернути список операцій представлені вхідним скриптом.
type Parser struct {
}

func (p *Parser) CommandParser(commandName string, args []string) painter.Operation {
	// var painterRes []painter.Operation
	switch commandName {
	case "white":
		return painter.OperationFunc(painter.WhiteFill)
	case "green":
		return painter.OperationFunc(painter.GreenFill)
	// case "bgrect":
	// 	return painter.OperationFunc(painter.)
	// case "move":
	// 	return painter.OperationFunc(painter.)
	// case "figure":
	// 	return painter.OperationFunc(painter.WhiteFill)
	case "update":
		return painter.UpdateOp
		// case "reset":
		// 	return painter.OperationFunc(painter.WhiteFill)
	}
	return nil
}

func (p *Parser) Parse(in io.Reader) ([]painter.Operation, error) {
	var res []painter.Operation
	fmt.Println("test")

	scanner := bufio.NewScanner(in)
	scanner.Split(bufio.ScanLines)

	for scanner.Scan() {
		commandLine := scanner.Text()

		sliced := strings.Split(commandLine, " ")
		args := sliced[1:]

		appendAction := p.CommandParser(sliced[0], args)
		if appendAction != nil {
			res = append(res, appendAction)
		}
	}

	return res, nil
}
