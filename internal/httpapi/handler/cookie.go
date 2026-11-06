package handler

// CookieConfig относится к HTTP-границе; User UseCase не зависит от него.
type CookieConfig struct {
	Name   string
	Secure bool
}

// TODO HW15-06: реализация установки и удаления cookie находится в auth handlers.
