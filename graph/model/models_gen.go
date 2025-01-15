// Code generated by github.com/99designs/gqlgen, DO NOT EDIT.

package model

type AuthResult struct {
	SessionToken string `json:"sessionToken"`
	RefreshToken string `json:"refreshToken"`
}

type LocationInput struct {
	Country    string  `json:"country"`
	PostalCode *string `json:"postalCode,omitempty"`
	City       *string `json:"city,omitempty"`
}

type LoginInput struct {
	Address   string `json:"address"`
	Message   string `json:"message"`
	Signature string `json:"signature"`
}

type Mutation struct {
}

type NonceInput struct {
	Address string `json:"address"`
}

type NonceResult struct {
	Nonce string `json:"nonce"`
}

type Query struct {
}

type RefreshInput struct {
	RefreshToken string `json:"refreshToken"`
}

type SignupInput struct {
	Login   *LoginInput       `json:"login"`
	Profile *UserProfileInput `json:"profile"`
}

type UserProfileInput struct {
	Firstname    string         `json:"firstname"`
	Lastname     string         `json:"lastname"`
	PrimaryEmail string         `json:"primaryEmail"`
	Location     *LocationInput `json:"location,omitempty"`
	Website      *string        `json:"website,omitempty"`
	PersonalGoal *string        `json:"personalGoal,omitempty"`
	About        *string        `json:"about,omitempty"`
}
