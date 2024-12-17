package apps

import (
	"testing"
)

func TestParseEnvVar(t *testing.T) {
	ctx := &runContext{}
	if env, err := ctx.buildEnv([]string{"TEST=abc"}); err != nil {
		t.Fatal(err)
	} else if env["TEST"] != "abc" {
		t.Fatal(env["TEST"], "should be abc")
	}
}

func TestParseEnvVarWithEqualSign(t *testing.T) {
	ctx := &runContext{}
	if env, err := ctx.buildEnv([]string{"TEST=a=b"}); err != nil {
		t.Fatal(err)
	} else if env["TEST"] != "a=b" {
		t.Fatal(env["TEST"], "should be a=b")
	}
}
