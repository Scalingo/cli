package apps

import (
	"testing"
)

func TestParseEnvVar(t *testing.T) {
	ctx := t.Context()
	runCtx := &runContext{}

	env, err := runCtx.buildEnv(ctx, []string{"TEST=abc"})
	if err != nil {
		t.Fatal(err)
	}

	if env["TEST"] != "abc" {
		t.Fatal(env["TEST"], "should be abc")
	}
}

func TestParseEnvVarWithEqualSign(t *testing.T) {
	ctx := t.Context()
	runCtx := &runContext{}

	env, err := runCtx.buildEnv(ctx, []string{"TEST=a=b"})
	if err != nil {
		t.Fatal(err)
	}

	if env["TEST"] != "a=b" {
		t.Fatal(env["TEST"], "should be a=b")
	}
}
