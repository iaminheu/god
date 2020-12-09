package gvalid

// checkLuHn checks <value> with LUHN algorithm.
// It's usually used for bank card number validation.
func checkLuHn(value string) bool {
	var (
		sum     = 0
		nDigits = len(value)
		parity  = nDigits % 2
	)
	for i := 0; i < nDigits; i++ {
		var digit = int(value[i] - 48)
		if i%2 == parity {
			digit *= 2
			if digit > 9 {
				digit -= 9
			}
		}
		sum += digit
	}
	return sum%10 == 0
}
