package external

import "testing"

func TestFilterValidMeals(t *testing.T) {
	meals := []MealItem{
		{Name: "BREAKFAST", HasMenu: false},
		{Name: "LUNCH", HasMenu: true},
		{Name: "DINNER", HasMenu: true},
	}

	valid := filterValidMeals(meals)
	if len(valid) != 2 {
		t.Fatalf("expected 2 valid meals, got %d", len(valid))
	}
	if valid[0].Name != "LUNCH" || valid[1].Name != "DINNER" {
		t.Fatalf("unexpected valid meal names: %#v", valid)
	}
}

func TestSelectHoursForMealMatchesByTitle(t *testing.T) {
	hours := []EventHour{
		{EventTitle: "Breakfast"},
		{EventTitle: "Lunch"},
	}

	selected := selectHoursForMeal("LUNCH", hours, 2)
	if len(selected) != 1 {
		t.Fatalf("expected 1 selected hour, got %d", len(selected))
	}
	if selected[0].EventTitle != "Lunch" {
		t.Fatalf("unexpected selected hour: %q", selected[0].EventTitle)
	}
}

func TestSelectHoursForMealFallsBackWhenSingleValidMeal(t *testing.T) {
	hours := []EventHour{
		{EventTitle: "Open"},
	}

	selected := selectHoursForMeal("LUNCH", hours, 1)
	if len(selected) != 1 {
		t.Fatalf("expected fallback to all hours, got %d", len(selected))
	}
	if selected[0].EventTitle != "Open" {
		t.Fatalf("unexpected selected hour: %q", selected[0].EventTitle)
	}
}

func TestSelectHoursForMealNoFallbackWhenMultipleValidMeals(t *testing.T) {
	hours := []EventHour{
		{EventTitle: "Open"},
	}

	selected := selectHoursForMeal("LUNCH", hours, 2)
	if len(selected) != 0 {
		t.Fatalf("expected no selected hours, got %d", len(selected))
	}
}
