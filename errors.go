package multiparse

const (
	ParseMonetaryStringError = "Cannot parse the string as a monetary type."
	ParseMoneyError          = ParseMonetaryStringError
	ParseMoneySeparatorError = "Cannot distinguish digit and decimal separators."
	ParseNumericError        = "Cannot parse the string as a numeric type."
)
