package models

type IdentifierType int32

const (
	IdentifierUnspecified IdentifierType = 0
	IdentifierEmail       IdentifierType = 1
	IdentifierCPF         IdentifierType = 2
	IdentifierCNPJ        IdentifierType = 3
	IdentifierPhone       IdentifierType = 4
)

var (
	IdentifierType_name = map[int32]string{
		0: "unspecified",
		1: "email",
		2: "CPF",
		3: "CNPJ",
		4: "phone",
	}
	IdentifierType_value = map[string]int32{
		"unspecified": 0,
		"email":       1,
		"CPF":         2,
		"CNPJ":        3,
		"phone":       4,
	}
)