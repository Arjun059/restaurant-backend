
package utils

import  validator	"github.com/go-playground/validator/v10"
import	models	"restaurant/internal/models"

var	validate = validator.New()

func ValidateDish(dish models.Dish) error {
	return validate.Struct(dish)
}
 


