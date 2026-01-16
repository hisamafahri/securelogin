package models

func AllModels() []interface{} {
	return []interface{}{
		&Application{},
		&AuthenticationProvider{},
		&AuthenticationRequest{},
		&AuthorizationCode{},
		&Session{},
		&User{},
	}
}
