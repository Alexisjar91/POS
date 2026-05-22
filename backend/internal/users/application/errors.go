package application

import "errors"

// ErrEmailInvalido se produce cuando el email no tiene un formato válido.
var ErrEmailInvalido = errors.New("el email no tiene un formato válido")
