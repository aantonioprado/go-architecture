package user

type Email string

func NewEmail(raw string) (Email, error) {
	if raw == "" {
		return "", ErrEmailRequired
	}

	return Email(raw), nil
}

func (e Email) String() string {
	return string(e)
}
