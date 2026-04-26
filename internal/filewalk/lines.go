package filewalk

import (
	"bufio"
	"errors"
	"os"
)

type Line struct {
	Number int
	Text   string
}

func ScanLines(path string, visit func(Line) error) (err error) {
	file, err := os.Open(path)
	if err != nil {
		return err
	}

	scanner := bufio.NewScanner(file)
	lineNumber := 0
	for scanner.Scan() {
		lineNumber++
		if err = visit(Line{
			Number: lineNumber,
			Text:   scanner.Text(),
		}); err != nil {
			return closeFile(file, err)
		}
	}

	return closeFile(file, scanner.Err())
}

func closeFile(file *os.File, existingErr error) (err error) {
	if closeErr := file.Close(); closeErr != nil {
		return errors.Join(existingErr, closeErr)
	}

	return existingErr
}
