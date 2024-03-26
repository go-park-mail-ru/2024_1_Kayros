package functions

import (
	"regexp"

	"2024_1_kayros/internal/utils/regex"
	"github.com/asaskevich/govalidator"
)

func InitValidator() {
	govalidator.SetFieldsRequiredByDefault(true)

	// Добавим тег валидации для пароля
	govalidator.TagMap["user_pwd"] = govalidator.Validator(func(pwd string) bool {
		// Проверка на минимальную длину
		if len(pwd) < 8 {
			return false
		}

		// Проверка на наличие хотя бы одной буквы
		letterRegex := regexp.MustCompile(`[A-Za-z]`)
		if !letterRegex.MatchString(pwd) {
			return false
		}

		// Проверка на наличие хотя бы одной цифры
		digitRegex := regexp.MustCompile(`\d`)
		if !digitRegex.MatchString(pwd) {
			return false
		}

		// Проверка на наличие разрешенных символов
		return regex.RegexPassword.MatchString(pwd)
	})

	govalidator.TagMap["user_email"] = govalidator.Validator(func(email string) bool {
		return regex.RegexEmail.MatchString(email)
	})

	govalidator.TagMap["user_name"] = govalidator.Validator(func(name string) bool {
		return regex.RegexName.MatchString(name)
	})

	govalidator.TagMap["user_phone"] = govalidator.Validator(func(phone string) bool {
		return regex.RegexPhone.MatchString(phone)
	})
}
