package env

import (
	"testing"
)

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
		} else if err != errSetInvalidSyntax {
			t.Fatal("expected", errSetInvalidSyntax, "error, got", err)
		}
	}

	vs = []string{"VA R=VAL", "	VAR=VAL", "VAR=VAL", "%%%=VAL"}
	for _, v = range vs {
		if err := isEnvEditValid(v); err == nil {
			t.Fatal(v, "should not be valid")
		} else if err != errInvalidNameFormat {
			t.Fatal("expected", errInvalidNameFormat, "error, got", err)
		}
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
