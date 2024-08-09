package main

type ICheckout interface {
	Scan(SKU string) error
	GetTotalPrice() (totalPrice int, err error)
}

type PricingRule struct {
	UnitPrice     int
	DiscountQty   int
	DiscountPrice int
}

type Checkout struct {
	items map[string]int
}

func NewCheckout() *Checkout {
	return &Checkout{
		items: make(map[string]int),
	}
}

func (c *Checkout) Scan(SKU string) error {
	c.items[SKU]++
	return nil
}

func (c *Checkout) GetTotalPrice() (totalPrice int, err error) {
	for _, count := range c.items {
		totalPrice += count * 50
	}
	return totalPrice, nil
}
