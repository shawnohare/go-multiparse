package multiparse

// Standard parsing errors.
const (
	ParseBoolError           = "Cannot parse string as a boolean."
	ParseIntError            = "Cannot parse string as an integer."
	ParseFloatError          = "Cannot parse string as a float."
	ParseMonetaryStringError = "Cannot parse string as a monetary value."
	ParseMoneyError          = ParseMonetaryStringError
	ParseMoneySeparatorError = "Cannot distinguish digit and decimal separators."
	ParseNumericError        = "Cannot parse string as a numeric type."
	ParseTimeError           = "Cannot parse string as a time."
	ParseTypeAssertError     = "Cannot assert correct type for parsed value."
	ParseError               = "Cannot parse string as any valid type."
	MoneyFloatError          = "Cannot convert Money instance to a float."
)
