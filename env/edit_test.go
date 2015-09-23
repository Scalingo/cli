package env

import "bytes"
import "github.com/Scalingo/go-scalingo"
import "testing"

func TestAdd(t *testing.T) {
}

func TestDelete(t *testing.T) {
}

func TestIsEnvEditValid(t *testing.T) {
	v := "VAR1=VAL1"
	if err := isEnvEditValid(v); err != nil {
		t.Fatal(v, "should be valid, got", err)
	}

	vs := []string{"VAR1=", "=VAL1", "VAR"}
	for _, v = range vs {
		if err := isEnvEditValid(v); err == nil {
			t.Fatal(v, "should not be valid")
		} else if err != setInvalidSyntaxError {
			t.Fatal("expected", setInvalidSyntaxError, "error, got", err)
		}
	}

	vs = []string{"VA R=VAL", "	VAR=VAL", "VAR=VAL", "%%%=VAL"}
	for _, v = range vs {
		if err := isEnvEditValid(v); err == nil {
			t.Fatal(v, "should not be valid")
		} else if err != invalidNameFormatError {
			t.Fatal("expected", invalidNameFormatError, "error, got", err)
		}
	}

	longName := new(bytes.Buffer)
	for i := 0; i < scalingo.EnvNameMaxLength; i++ {
		longName.WriteRune('A')
	}
	longName.WriteString("A=VAL")
	if err := isEnvEditValid(longName.String()); err == nil {
		t.Fatal(v, "should not be valid")
	} else if err != nameTooLongError {
		t.Fatal("expected", nameTooLongError, "error, got", err)
	}
}

func TestParseVariable(t *testing.T) {
	v := "VAR1=VAL1"
	name, value := parseVariable(v)
	if name != "VAR1" {
		t.Fatal("expected VAR1, got", name)
	}
	if value != "VAL1" {
		t.Fatal("expected VAL1, got", value)
	}
}
