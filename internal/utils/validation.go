
package utils

import (
	"fmt"

	validator "github.com/go-playground/validator/v10"
	models "restaurant/internal/models"
)

var	validate = validator.New()

func ValidateDish(dish models.Dish) error {
	// basic struct validation (tags) first
	if err := validate.Struct(dish); err != nil {
		return err
	}

	// conditional requirement: dish price must be set only if no variants provided
	if len(dish.Variants) == 0 {
		if dish.Price <= 0 {
			return fmt.Errorf("price is required when no variants are defined")
		}
	} else {
		// ensure all variants have valid prices
		for _, v := range dish.Variants {
			if v.Price <= 0 {
				return fmt.Errorf("all variants must have a positive price")
			}
		}
	}

	return nil
}
 


