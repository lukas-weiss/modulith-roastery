package internal

import "testing"

func TestReducesEveryProductOfAnOrder(t *testing.T) {
	stock := NewStockRepository()

	if err := stock.Reduce(map[string]int{"crema-beans": 2, "espresso-cups": 1}); err != nil {
		t.Fatal(err)
	}

	if got := stock.Get("crema-beans"); got != 10 {
		t.Errorf("crema-beans = %d, want 10", got)
	}
	if got := stock.Get("espresso-cups"); got != 5 {
		t.Errorf("espresso-cups = %d, want 5", got)
	}
}

func TestChangesNothingIfOneProductIsShort(t *testing.T) {
	stock := NewStockRepository()

	if err := stock.Reduce(map[string]int{"crema-beans": 2, "espresso-cups": 7}); err == nil {
		t.Fatal("want an error")
	}

	if got := stock.Get("crema-beans"); got != 12 {
		t.Errorf("crema-beans = %d, want 12", got)
	}
	if got := stock.Get("espresso-cups"); got != 6 {
		t.Errorf("espresso-cups = %d, want 6", got)
	}
}
