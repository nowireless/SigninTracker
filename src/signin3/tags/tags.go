package tags

import (
	"reflect"

	log "github.com/sirupsen/logrus"
)

func FindFields(v interface{}, key, wantedValue string) []string {
	t := reflect.TypeOf(v)
	fields := []string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)
		if value, ok := field.Tag.Lookup(key); ok && value == wantedValue {
			log.Debug(field.Name, value)

			fields = append(fields, field.Name)
		}
	}

	return fields
}

func CheckRequiredOnCreate(obj interface{}) []string {
	t := reflect.TypeOf(obj)
	v := reflect.ValueOf(obj)
	key := "meta"
	wantedValue := "requiredOnCreate"

	missing := []string{}
	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if tagValue, ok := field.Tag.Lookup(key); ok && tagValue == wantedValue {
			value := v.FieldByName(field.Name)
			if value.Kind() != reflect.Ptr {
				panic("Field with requiredOnCreate has a non pointer type: " + field.Name)
			}

			// log.Info(field.Name, ", ", field.Type.Kind())

			// If a field does not have a value then it is set to nil
			if value.IsNil() {
				missing = append(missing, field.Name)
			}
		}
	}

	return missing
}

// Extract json field name from struct
// https://stackoverflow.com/questions/40864840/how-to-get-the-json-field-names-of-a-struct-in-golang
// func (b example) PrintFields() {
//     val := reflect.ValueOf(b)
//     for i := 0; i < val.Type().NumField(); i++ {
//         t := val.Type().Field(i)
//         fieldName := t.Name

//         if jsonTag := t.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
//             if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {
//                 fieldName = jsonTag[:commaIdx]
//             }
//         }

//         fmt.Println(fieldName)
//     }
// }
