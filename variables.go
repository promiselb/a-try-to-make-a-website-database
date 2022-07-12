package website

import (
	"errors"
)

const dataBasePath = "C://UsersData/"
const htmlPagesPath = "C://HtmlPages/html/"
const defaultPPPath = "C://DefaultPP/"

var Port = ":8080"

var (
	err                             error
	errInvalidPassword              = errors.New("the password is not as the verified one?")
	errNoRegisteredOrDeletedAccount = errors.New("no account registered with this gmail? Make one?")
	errLostOrDeletedData            = errors.New("lost or deleted data? Visit your actions log?")
	errAlreadyOccupiedGmail         = errors.New("an account already registered with the given gmail? User another one or check your gmail links?")
)

var (
	firstName      Name
	lastName       Name
	nation         Name
	gendre         Gendre
	gmail          Gmail
	password       Password
	verifyPassword Password
	birthday       Birthday
)
