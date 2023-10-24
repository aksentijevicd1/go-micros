package data

import "testing"

func TestCheckValidation(t *testing.T) {

	p := &Product{
		Name:  "duke",
		Price: 1,
		SKU:   "as-asd-asds-",
	}

	err := p.Validate()

	if err != nil {

		t.Fatal(err)

	}

}
