package rental

import "testing"

func TestCalculateIncludesInsuranceAndLateFees(t *testing.T) {
	got, err := calculate(100, "2026-05-01", "2026-05-03", "2026-05-05", true)
	if err != nil {
		t.Fatal(err)
	}
	if got.Days != 3 {
		t.Fatalf("days = %d, want 3", got.Days)
	}
	if got.LateFee != 100 {
		t.Fatalf("late fee = %.2f, want 100.00", got.LateFee)
	}
	if got.Total != 445 {
		t.Fatalf("total = %.2f, want 445.00", got.Total)
	}
}

func TestCalculateRejectsInvalidRange(t *testing.T) {
	if _, err := calculate(100, "2026-05-03", "2026-05-01", "", false); err == nil {
		t.Fatal("expected invalid date range error")
	}
}
