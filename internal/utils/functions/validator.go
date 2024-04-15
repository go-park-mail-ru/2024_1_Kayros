package functions

import (
	"regexp"

	"2024_1_kayros/internal/utils/regex"
	"github.com/asaskevich/govalidator"
	"go.uber.org/zap"
)

func InitValidator(logger *zap.Logger) {
	govalidator.SetFieldsRequiredByDefault(true)

	// Relation "user"
	// user_name_domain
	govalidator.TagMap["user_name_domain"] = func(name string) bool {
		return regex.Name.MatchString(name)
	}

	// phone_domain
	govalidator.TagMap["user_phone_domain"] = func(phone string) bool {
		return regex.Phone.MatchString(phone)
	}

	// email_domain
	govalidator.TagMap["user_email_domain"] = func(email string) bool {
		return len(email) >= 6 && len(email) <= 50 && regex.Email.MatchString(email)
	}

	// address_domain
	govalidator.TagMap["user_address_domain"] = func(address string) bool {
		return len(address) >= 14 || len(address) <= 100
	}

	// img_url_domain
	govalidator.TagMap["img_url_domain"] = func(imgUrl string) bool {
		return len(imgUrl) <= 60
	}

	// card_number_domain
	govalidator.TagMap["user_card_number_domain"] = func(cardNumber string) bool {
		return regex.CardNumber.MatchString(cardNumber)
	}

	// password_domain
	govalidator.TagMap["user_password_domain"] = func(pwd string) bool {
		// Check length range
		if len(pwd) < 8 || len(pwd) > 20 {
			return false
		}

		// Checking for the presence of at least one letter
		letterRegex := regexp.MustCompile(`[A-Za-z]`)
		if !letterRegex.MatchString(pwd) {
			return false
		}

		// Checking for the presence of at least one digit
		digitRegex := regexp.MustCompile(`\d`)
		if !digitRegex.MatchString(pwd) {
			return false
		}

		// Checking for the regular expression matching
		return regex.Password.MatchString(pwd)
	}

	// Relation category
	govalidator.TagMap["category_name_domain"] = func(name string) bool {
		return regex.CategoryName.MatchString(name)
	}

	// Relation restaurant
	govalidator.TagMap["rest_name_domain"] = func(name string) bool {
		return regex.RestName.MatchString(name)
	}
	logger.Info("Custom tags created")
}
