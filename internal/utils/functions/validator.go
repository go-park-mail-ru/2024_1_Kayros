package functions

import (
	"regexp"

	"2024_1_kayros/internal/utils/regex"
	"github.com/asaskevich/govalidator"
	"go.uber.org/zap"
)

func InitValidator(logger *zap.Logger) {
	govalidator.SetFieldsRequiredByDefault(true)

	govalidator.TagMap["user_pwd"] = func(pwd string) bool {
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
	}

	govalidator.TagMap["user_email"] = func(email string) bool {
		return regex.RegexEmail.MatchString(email)
	}

	govalidator.TagMap["user_name"] = func(name string) bool {
		return regex.RegexName.MatchString(name)
	}

	govalidator.TagMap["user_phone"] = func(phone string) bool {
		return regex.RegexPhone.MatchString(phone)
	}

	logger.Info("Кастомные теги созданы")
}
