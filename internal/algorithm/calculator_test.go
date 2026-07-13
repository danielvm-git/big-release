package algorithm

import "testing"

func TestCalculateNextVersion_FirstRelease_WithInitialVersion(t *testing.T) {
	calc := NewCalculator()
	branch := &Branch{Name: "main", Type: BranchTypeRelease}
	version, err := calc.CalculateNextVersion(nil, ReleaseTypeMinor, branch, "0.1.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != "0.1.0" {
		t.Errorf("expected 0.1.0, got %s", version)
	}
}

func TestCalculateNextVersion_FirstRelease_Default(t *testing.T) {
	calc := NewCalculator()
	branch := &Branch{Name: "main", Type: BranchTypeRelease}
	version, err := calc.CalculateNextVersion(nil, ReleaseTypeMinor, branch, "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != FirstRelease {
		t.Errorf("expected %s, got %s", FirstRelease, version)
	}
}

func TestCalculateNextVersion_IncrementMinor(t *testing.T) {
	calc := NewCalculator()
	branch := &Branch{Name: "main", Type: BranchTypeRelease}
	last := &Release{Version: "1.0.0"}
	version, err := calc.CalculateNextVersion(last, ReleaseTypeMinor, branch, "0.1.0")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if version != "1.1.0" {
		t.Errorf("expected 1.1.0, got %s", version)
	}
}
