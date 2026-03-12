package main

import (
	"bytes"
	"strings"
	"testing"
)

func TestRootHelpIncludesFlags(t *testing.T) {
	buf := new(bytes.Buffer)

	oldOut := rootCmd.OutOrStdout()
	oldErr := rootCmd.ErrOrStderr()
	defer rootCmd.SetOut(oldOut)
	defer rootCmd.SetErr(oldErr)

	rootCmd.SetOut(buf)
	rootCmd.SetErr(buf)

	defer rootCmd.SetArgs(nil)
	rootCmd.SetArgs([]string{"--help"})

	if err := rootCmd.Execute(); err != nil {
		t.Fatalf("Execute returned error: %v", err)
	}

	help := buf.String()
	for _, flag := range []string{"--pretty", "--verbose"} {
		if !strings.Contains(help, flag) {
			t.Fatalf("expected %s in help output\nhelp:\n%s", flag, help)
		}
	}
}
