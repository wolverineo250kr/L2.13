package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"strconv"
	"strings"
)

type FieldSelector struct {
	fields map[int]bool
}

func parseFieldSelector(input string) (*FieldSelector, error) {
	fs := &FieldSelector{fields: make(map[int]bool)}

	parts := strings.Split(input, ",")
	for _, part := range parts {
		part = strings.TrimSpace(part)

		if part == "" {
			continue
		}

		// Диапазон X-Y
		if strings.Contains(part, "-") {
			r := strings.Split(part, "-")
			if len(r) != 2 {
				return nil, fmt.Errorf("неправльный диапозон: %s", part)
			}

			start, err1 := strconv.Atoi(r[0])
			end, err2 := strconv.Atoi(r[1])
			if err1 != nil || err2 != nil || start < 1 || end < start {
				return nil, fmt.Errorf("неправльный диапозон: %s", part)
			}

			for i := start; i <= end; i++ {
				fs.fields[i] = true
			}
		} else {
			val, err := strconv.Atoi(part)
			if err != nil || val < 1 {
				return nil, fmt.Errorf("неверное поле: %s", part)
			}
			fs.fields[val] = true
		}
	}
	return fs, nil
}

func cutLine(line, delimiter string, fs *FieldSelector) string {
	parts := strings.Split(line, delimiter)
	var out []string

	for i := 1; i <= len(parts); i++ {
		if fs.fields[i] {
			out = append(out, parts[i-1])
		}
	}

	return strings.Join(out, delimiter)
}

func main() {
	fieldsFlag := flag.String("f", "", "поля (пр. 1,3-5)")
	delimiterFlag := flag.String("d", "\t", "delimiter")
	separatedFlag := flag.Bool("s", false, "only print lines with delimiter")

	// Дополнительный флаг, похожий на GNU cut
	outputDelimiter := flag.String("output-delimiter", "", "set output delimiter")

	flag.Parse()

	if *fieldsFlag == "" {
		fmt.Fprintln(os.Stderr, "Ошибка: -f flag is required")
		os.Exit(1)
	}

	fs, err := parseFieldSelector(*fieldsFlag)
	if err != nil {
		fmt.Fprintln(os.Stderr, "Ошибка парсинга -f:", err)
		os.Exit(1)
	}

	delimiter := *delimiterFlag
	outDelim := delimiter
	if *outputDelimiter != "" {
		outDelim = *outputDelimiter
	}

	scanner := bufio.NewScanner(os.Stdin)

	for scanner.Scan() {
		line := scanner.Text()

		if *separatedFlag && !strings.Contains(line, delimiter) {
			continue
		}

		parts := strings.Split(line, delimiter)
		var result []string

		for i := 1; i <= len(parts); i++ {
			if fs.fields[i] {
				result = append(result, parts[i-1])
			}
		}

		fmt.Println(strings.Join(result, outDelim))
	}

	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "ошибка ввода:", err)
	}
}
