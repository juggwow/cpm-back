package utils

func IfThenElse(condition bool, trueValue interface{}, falseValue interface{}) interface{} {
	if condition {
		return trueValue
	}
	return falseValue
}

// func TestIfThenElse(t *testing.T) {
//     assert.Equal(t, IfThenElse(1 == 1, "Yes", false), "Yes")
//     assert.Equal(t, IfThenElse(1 != 1, nil, 1), 1)
//     assert.Equal(t, IfThenElse(1 < 2, nil, "No"), nil)
// }
