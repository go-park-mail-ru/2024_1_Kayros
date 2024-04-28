package regex

import "regexp"

var (
<<<<<<< HEAD
	RegexName     = regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ][a-zA-Zа-яА-ЯёЁ0-9]{1,19}$`)
	RegexEmail    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	RegexPassword = regexp.MustCompile(`^[A-Za-z\d!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]+$`)
	RegexPhone    = regexp.MustCompile(`^\+7 \(\d{3}\) \d{3} \d{2} \d{2}$`)
	RegexURI      = regexp.MustCompile(`^/api/v1/order/food/delete/\d+$`)
=======
	Name         = regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ][a-zA-Zа-яА-ЯёЁ0-9]{2,20}$`)
	Email        = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	Password     = regexp.MustCompile(`^[A-Za-z\d!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]+$`)
	Phone        = regexp.MustCompile(`^\+7 \(\d{3}\) \d{3} \d{2} \d{2}$`)
	CardNumber   = regexp.MustCompile(`\d{4} \d{4} \d{4} \d{4}`)
	CategoryName = regexp.MustCompile(`[a-zA-Zа-яА-ЯёЁ]{2,30}`)
	RestName     = regexp.MustCompile(`[a-zA-Zа-яА-ЯёЁ]{2,30}`)
>>>>>>> fix_csrf_test
)
