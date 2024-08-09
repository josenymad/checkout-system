package main

import (
	"errors"
	"fmt"
)

type ICheckout interface {
	Scan(SKU string) error
	GetTotalPrice() (totalPrice int, err error)
}

type PricingRule struct {
	UnitPrice     int
	DiscountQty   int
	DiscountPrice int
}

type PricingService interface {
	GetPricingRule(SKU string) (PricingRule, error)
}

type FileBasedPricingService struct {
	pricingRules map[string]PricingRule
}

func NewFileBasedPricingService(filePath string) (*FileBasedPricingService, error) {
	service := &FileBasedPricingService{}
	err := service.loadPricingRules(filePath)
	if err != nil {
		return nil, err
	}
	return service, nil
}

func (s *FileBasedPricingService) loadPricingRules(filePath string) error {
	file, err := os.Open(filePath)
	if err != nil {
		return fmt.Errorf("could not open pricing file: %v", err)
	}
	defer file.Close()

	byteValue, err := io.ReadAll(file)
	if err != nil {
		return fmt.Errorf("could not read pricing file: %v", err)
	}

	var pricingRules map[string]PricingRule
	err = json.Unmarshal(byteValue, &pricingRules)
	if err != nil {
		return fmt.Errorf("could not parse pricing file: %v", err)
	}

	s.pricingRules = pricingRules
	return nil
}

func (s *FileBasedPricingService) GetPricingRule(SKU string) (PricingRule, error) {
	rule, exists := s.pricingRules[SKU]
	if !exists {
		return PricingRule{}, fmt.Errorf("no pricing rule found for SKU: %s", SKU)
	}
	return rule, nil
}

type Checkout struct {
	items        map[string]int
	pricingRules map[string]PricingRule
}

func NewCheckout(pricingRules map[string]PricingRule) *Checkout {
	return &Checkout{
		items:        make(map[string]int),
		pricingRules: pricingRules,
	}
}

func (c *Checkout) Scan(SKU string) error {
	_, exists := c.pricingRules[SKU]
	if !exists {
		return errors.New("invalid SKU")
	}
	c.items[SKU]++
	return nil
}

func (c *Checkout) GetTotalPrice() (totalPrice int, err error) {
	for SKU, count := range c.items {
		rule, exists := c.pricingRules[SKU]
		if !exists {
			return 0, fmt.Errorf("no pricing rule found for SKU: %s", SKU)
		}
		if rule.DiscountQty > 0 && count >= rule.DiscountQty {
			totalPrice += (count / rule.DiscountQty) * rule.DiscountPrice
			totalPrice += (count % rule.DiscountQty) * rule.UnitPrice
		} else {
			totalPrice += count * rule.UnitPrice
		}
	}
	return totalPrice, nil
}
