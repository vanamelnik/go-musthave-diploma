package currency

// Add exactly adds two float numbers expecting them to be numbers with two decimal places.
func Add(a, b float32) float32 {
	aInt, bInt := int(a*100), int(b*100)
	return float32(aInt+bInt) / 100
}
