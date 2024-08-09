package main

import "testing"

func TestCheckoutImplementation(t *testing.T) {
	var _ ICheckout = (*Checkout)(nil)
}
