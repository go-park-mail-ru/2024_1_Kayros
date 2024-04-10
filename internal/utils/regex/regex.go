package regex

import "regexp"

var (
	RegexName     = regexp.MustCompile(`^[a-zA-Zа-яА-ЯёЁ][a-zA-Zа-яА-ЯёЁ0-9]{1,19}$`)
	RegexEmail    = regexp.MustCompile(`^[a-zA-Z0-9._%+-]+@[a-zA-Z0-9.-]+\.[a-zA-Z]{2,}$`)
	RegexPassword = regexp.MustCompile(`^[A-Za-z\d!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]+$`)
	RegexPhone    = regexp.MustCompile(`^\+7 \(\d{3}\) \d{3} \d{2} \d{2}$`)
)
