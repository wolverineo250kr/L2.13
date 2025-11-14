package main

import (
	"bytes"
	"os/exec"
	"strings"
	"testing"
)

func TestParseFieldSelector(t *testing.T) {
	fs, err := parseFieldSelector("1,3-5,7")
	if err != nil {
		t.Fatal(err)
	}

	expected := []int{1, 3, 4, 5, 7}
	for _, v := range expected {
		if !fs.fields[v] {
			t.Errorf("ожидалось поле %d", v)
		}
	}
}

func TestCutLineSimple(t *testing.T) {
	fs, _ := parseFieldSelector("1,3")
	result := cutLine("a:b:c:d", ":", fs)

	if result != "a:c" {
		t.Errorf("ожидалось 'a:c', пришло '%s'", result)
	}
}

func TestCutLineOutOfRange(t *testing.T) {
	fs, _ := parseFieldSelector("1,10")
	result := cutLine("x:y:z", ":", fs)

	if result != "x" {
		t.Errorf("ожидалось 'x', пришло '%s'", result)
	}
}

func TestProgramExecution(t *testing.T) {
	cmd := exec.Command("go", "run", "cut.go", "-f", "2", "-d", ":")
	cmd.Stdin = strings.NewReader("a:b:c")

	var out bytes.Buffer
	cmd.Stdout = &out

	err := cmd.Run()
	if err != nil {
		t.Fatalf("потрачено: %v", err)
	}

	if strings.TrimSpace(out.String()) != "b" {
		t.Errorf("ожидалось 'b', пришло '%s'", out.String())
	}
}
