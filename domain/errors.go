package domain

import "errors"

var ErrNotFound = errors.New("id not found")
var ErrEnoughMoney = errors.New("insufficient funds to write off")
var ErrUnavailable = errors.New("service is unavailable")
