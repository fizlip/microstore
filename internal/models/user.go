// package models contains definitions of common types and models used accross the API.
package models

type User interface {
	Register()
	Login()
	Logout()
}
