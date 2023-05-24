package env

import (
	"github.com/Scalingo/go-scalingo/v6"
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

func TestMergeVariables(t *testing.T) {
	variables1 := scalingo.Variables{
		&scalingo.Variable{
			Name:  "NODE_ENV",
			Value: "production",
		},
	}
	variables2 := scalingo.Variables{
		&scalingo.Variable{
			Name:  "NODE_ENV",
			Value: "test",
		},
		&scalingo.Variable{
			Name:  "VAR1",
			Value: "VAL1",
		},
	}
	variables := mergeVariables(variables1, variables2)
	if len(variables) != 2 {
		t.Fatal("expected 2 variable, got", len(variables))
	}
	if variables[0].Name != "NODE_ENV" {
		t.Fatal("expected NODE_ENV, got", variables[0].Name)
	}
	if variables[0].Value != "test" {
		t.Fatal("expected test, got", variables[0].Value)
	}
	if variables[1].Name != "VAR1" {
		t.Fatal("expected VAR1, got", variables[1].Name)
	}
	if variables[1].Value != "VAL1" {
		t.Fatal("expected VAL1, got", variables[1].Value)
	}
}
