package main

import (
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
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
	items          map[string]int
	pricingService PricingService
}

func NewCheckout(pricingService PricingService) *Checkout {
	return &Checkout{
		items:          make(map[string]int),
		pricingService: pricingService,
	}
}

func (c *Checkout) Scan(SKU string) error {
	_, err := c.pricingService.GetPricingRule(SKU)
	if err != nil {
		return err
	}
	c.items[SKU]++
	return nil
}

func (c *Checkout) GetTotalPrice() (totalPrice int, err error) {
	if len(c.items) == 0 {
		return 0, errors.New("no items have been scanned")
	}

	for SKU, count := range c.items {
		rule, err := c.pricingService.GetPricingRule(SKU)
		if err != nil {
			return 0, err
		}
		if rule.UnitPrice <= 0 || (rule.DiscountQty > 0 && rule.DiscountPrice <= 0) {
			return 0, fmt.Errorf("invalid pricing rule for SKU: %s", SKU)
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

func main() {
	pricingService, err := NewFileBasedPricingService("pricing_rules.json")
	if err != nil {
		fmt.Println("Error loading pricing service:", err)
		return
	}

	checkout := NewCheckout(pricingService)

	reader := bufio.NewReader(os.Stdin)

	fmt.Println("Welcome to the checkout system!")
	fmt.Println("Type SKU letters to scan items. Type 'checkout' to finish and see the total price.")

	for {
		fmt.Print("Enter SKU: ")
		input, _ := reader.ReadString('\n')
		input = strings.TrimSpace(input)
		input = strings.ToUpper(input)

		if input == "CHECKOUT" {
			break
		}

		err = checkout.Scan(input)
		if err != nil {
			fmt.Println("Error scanning item:", err)
		} else {
			fmt.Printf("Scanned %s\n", input)
		}
	}

	totalPrice, err := checkout.GetTotalPrice()
	if err != nil {
		fmt.Println("Error calculating total price:", err)
	} else {
		fmt.Printf("Total Price: %d\n", totalPrice)
	}
}
