package renderer

import "testing"

func TestReadLine(t *testing.T) {
	str := "Hello there"
	pos := 0
	_, err := readLine(str, &pos)
	if err == nil {
		t.Error("Err expected as string doesnt terminate with an \\n")
	}
}
func TestReadLine2(t *testing.T) {
	str := "Hello there\n"
	pos := 0
	_, err := readLine(str, &pos)
	if err != nil {
		t.Error("String terminates!")
	}
	if pos != 12 {
		t.Error("New post must be 12, seen", pos)
	}
}
func TestReadLine3(t *testing.T) {
	str := "Hello there\nHow Are you?"
	pos := 12
	_, err := readLine(str, &pos)
	if err == nil {
		t.Error("String doesn't terminate!")
	}
	if pos != 12 {
		t.Error("Pos has to stay 12, it's", pos, "instead")
	}
}
func TestReadLine4(t *testing.T) {
	str := "Hello there\nHow Are you?\nWhat's up?"
	pos := 12
	_, err := readLine(str, &pos)
	if err != nil {
		t.Error("New line must be found!")
	}
	if pos != 25 {
		t.Error("Pos has to stay 25, it's", pos, "instead")
	}
}
func TestReadLine5(t *testing.T) {
	str := "Hello there\nHow Are you?\nWhat's up?\n"
	pos := 12
	line, err := readLine(str, &pos)
	if err != nil {
		t.Error("New line must be found!")
	}
	if line != "How Are you?\n" {
		t.Error("Expected line to be 'How Are you?' Seen", line)
	}
	if pos != 25 {
		t.Error("Pos has to stay 25, it's", pos, "instead")
	}
	line, err = readLine(str, &pos)
	if line != "What's up?\n" {
		t.Error("Expected line to be 'What's up?' Seen", line)
	}
}
