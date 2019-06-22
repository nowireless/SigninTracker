package app

import "reflect"

// nilFields returns a list of fields that are nil
func nilFields(obj interface{}) []string {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)

	results := []string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		value := v.FieldByName(field.Name)
		if value.Kind() == reflect.Ptr {
			// If a field does not have a value then it is set to nil
			if value.IsNil() {
				results = append(results, field.Name)
			}
		}
	}

	return results
}
