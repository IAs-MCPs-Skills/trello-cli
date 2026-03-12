package contract_test

import (
	"testing"

	"github.com/brettmcdowell/trello-cli/internal/contract"
)

func TestRequireFlagPresent(t *testing.T) {
	err := contract.RequireFlag("board", "abc123")
	if err != nil {
		t.Errorf("RequireFlag() with value should return nil, got %v", err)
	}
}

func TestRequireFlagMissing(t *testing.T) {
	err := contract.RequireFlag("board", "")
	if err == nil {
		t.Fatal("RequireFlag() with empty value should return error")
	}

	ce, ok := err.(*contract.ContractError)
	if !ok {
		t.Fatal("RequireFlag() should return *ContractError")
	}
	if ce.Code != contract.ValidationError {
		t.Errorf("Code = %q, want %q", ce.Code, contract.ValidationError)
	}
}
