package functions

import (
	"regexp"
	"strconv"
	"unicode/utf8"

	"2024_1_kayros/internal/utils/regex"
	"github.com/asaskevich/govalidator"
	"go.uber.org/zap"
)

func InitDtoValidator(logger *zap.Logger) {
	govalidator.SetFieldsRequiredByDefault(true)

	// Relation "user"
	// user_name_domain
	govalidator.TagMap["user_name_domain"] = func(name string) bool {
		return regex.Name.MatchString(name)
	}

	// user_phone_domain
	govalidator.TagMap["user_phone_domain"] = func(phone string) bool {
		return regex.Phone.MatchString(phone)
	}

	// user_email_domain
	govalidator.TagMap["user_email_domain"] = func(email string) bool {
		emailLen := utf8.RuneCountInString(email)
		return emailLen >= 6 && emailLen <= 50 && regex.Email.MatchString(email)
	}

	// user_address_domain
	govalidator.TagMap["user_address_domain"] = func(address string) bool {
		addressLen := utf8.RuneCountInString(address)
		return addressLen >= 14 && addressLen <= 100
	}

	// order_extra_address_domain
	govalidator.TagMap["user_extra_address_domain"] = func(address string) bool {
		addressLen := utf8.RuneCountInString(address)
		return addressLen >= 2 && addressLen <= 100
	}

	// img_url_domain
	govalidator.TagMap["img_url_domain"] = func(imgUrl string) bool {
		imgUrlLen := utf8.RuneCountInString(imgUrl)
		return imgUrlLen <= 60
	}

	// user_card_number_domain
	govalidator.TagMap["user_card_number_domain"] = func(cardNumber string) bool {
		return regex.CardNumber.MatchString(cardNumber)
	}

	// user_password_domain
	govalidator.TagMap["user_password_domain"] = func(pwd string) bool {
		// Check length range
		pwdLen := utf8.RuneCountInString(pwd)
		if pwdLen < 8 || pwdLen > 20 {
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

	govalidator.TagMap["positive"] = func(numStr string) bool {
		num, err := strconv.Atoi(numStr)
		if err != nil {
			return false
		}
		return num > 0
	}

	logger.Info("Custom tags created")
}
