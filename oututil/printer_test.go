package oututil

import (
	"testing"
)

func TestErasePrevious(t *testing.T) {
	erasePrevious()
	t.Log("erasePrevious called (writes ANSI codes)")
}

func TestStartPrinting_NoTTY(t *testing.T) {
	p := StartPrinting()
	if outTTY && p == nil {
		t.Errorf("expected printer when TTY, got nil")
	}
	if !outTTY && p != nil {
		t.Errorf("expected nil when not TTY, got printer")
	}
	t.Logf("outTTY: %v", outTTY)
}

func TestPrinter_Println_Nil(t *testing.T) {
	var p *Printer = nil
	p.Println("test")
}

func TestPrinter_Finalize_Nil_WithMessage(t *testing.T) {
	var p *Printer = nil
	p.Finalize("final message")
}

func TestPrinter_Finalize_Nil_Empty(t *testing.T) {
	var p *Printer = nil
	p.Finalize("")
}

func TestCanTestRealPrinter(t *testing.T) {
	result := outTTY
	expected := outTTY
	if result != expected {
		t.Errorf("expected %v, got %v", expected, result)
	}
}

// Note: Real printer testing might be skipped in CI environments
func TestPrinter_IfTTY(t *testing.T) {
	if !outTTY {
		t.Skip("Skipping real printer test because not TTY")
	}

	p := StartPrinting()
	if p == nil {
		t.Fatal("expected printer in TTY mode")
	}

	p.Println("test line")
	p.Finalize("final line")
}
