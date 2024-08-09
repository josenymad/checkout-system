package main

import "testing"

func TestCheckoutImplementation(t *testing.T) {
	var _ ICheckout = (*Checkout)(nil)
}

func TestCheckout(t *testing.T) {
	mockRules := map[string]PricingRule{
		"A": {UnitPrice: 50, DiscountQty: 3, DiscountPrice: 130},
		"B": {UnitPrice: 30, DiscountQty: 2, DiscountPrice: 45},
	}

	checkout := NewCheckout(mockRules)

	// test without discount
	err := checkout.Scan("A")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	err = checkout.Scan("B")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	totalPrice, err := checkout.GetTotalPrice()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if totalPrice != 80 {
		t.Fatalf("expected total price 80, got %d", totalPrice)
	}

	// test with discount
	checkout = NewCheckout(mockRules)

	err = checkout.Scan("A")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = checkout.Scan("A")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	err = checkout.Scan("A")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	totalPrice, err = checkout.GetTotalPrice()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if totalPrice != 130 {
		t.Fatalf("expected total price 130, got %d", totalPrice)
	}
}
