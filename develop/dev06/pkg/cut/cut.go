package cut

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"
)

type Fields []string

func (f *Fields) String() string {
	output := strings.Builder{}
	for i, v := range *f {
		if i != 0 && i != len(*f) {
			output.WriteString(",")
		}
		output.WriteString(v)
	}
	return output.String()
}

func (f *Fields) Set(str string) error {
	slice := strings.Split(str, ",")
	if len(slice) == 0 {
		return errors.New("no elements provided as fields")
	}
	for _, v := range slice {
		v = strings.TrimSpace(v)
		*f = append(*f, v)
	}
	return nil
}

type Cut struct {
	Source *os.File
	F      Fields
	D      string
	S      bool
}

func (c *Cut) getFields(header string) ([]int, error) {
	output := make([]int, 0, len(c.F))
	col := strings.Split(header, c.D)
	for i, v := range col {
		for _, field := range c.F {
			if v == field {
				output = append(output, i)
			}
		}
	}
	if len(col) == 0 || len(output) == 0 {
		return nil, errors.New("no columns found with such fields flag:" + c.F.String())
	}
	return output, nil

}

func (c *Cut) Print() error {
	toPrint, err := c.Write()
	if err != nil {
		return errors.New(err.Error())
	}
	for _, v := range toPrint {
		fmt.Println(v)
	}
	return nil
}

func (c *Cut) Write() ([]string, error) {
	if c.Source == nil {
		return nil, errors.New("no file provided")
	}
	fileStats, err := c.Source.Stat()
	sliceSize := 2 * fileStats.Size() / 80 // Примерно 2 байта в одном символе, примерно 80 символов в строке
	if err != nil {
		sliceSize = 10
	}
	if sliceSize < 1 {
		sliceSize = 1
	}
	output := make([]string, 0, sliceSize)
	scanner := bufio.NewScanner(c.Source)
	var columns []int
	for scanner.Scan() {
		text := scanner.Text()
		if columns == nil {
			columns, err = c.getFields(text)
			if err != nil {
				return nil, errors.New("could not set fields for cut. Error: " + err.Error())
			}
		}
		textSplit := strings.Split(text, c.D)
		if c.S && len(textSplit) < len(c.F) {
			continue
		}
		outputStr := strings.Builder{}
		for i, v := range columns {
			val := "NONE"
			if v < len(textSplit) {
				val = textSplit[v]
			}
			outputStr.WriteString(val)
			if i != len(columns)-1 {
				outputStr.WriteString(c.D)
			}
		}
		output = append(output, outputStr.String())
	}
	return output, nil
}
