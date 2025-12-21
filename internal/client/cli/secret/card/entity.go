package card

import "fmt"

type inputCard struct {
	id             string
	name           string
	number         string
	expiredDate    string
	cardHolderName string
	cvv            string
	notes          string
}

func (i *inputCard) Validate() error {
	if i.id == "" {
		return fmt.Errorf("id is required")
	}
	if i.name == "" {
		return fmt.Errorf("name is required")
	}
	if i.number == "" {
		return fmt.Errorf("number is required")
	}
	if i.expiredDate == "" {
		return fmt.Errorf("expired date is required")
	}
	if i.cardHolderName == "" {
		return fmt.Errorf("card holder name is required")
	}
	if i.cvv == "" {
		return fmt.Errorf("cvv is required")
	}
	return nil
}
