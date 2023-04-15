package gtranslate

import (
	"testing"

	"gpt-client/otto"
)

func TestGetSM(t *testing.T) {
	ttk, err := sM(otto.FalseValue())
	if err != nil {
		t.Error(err)
	}
	if ttk.IsNull() {
		t.Error("ttk is null")
	}
}
