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

// TODO: Need to use name from JSON tag
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
				panic("Field with " + wantedValue + " has a non pointer type: " + field.Name)
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

// CheckPatchReadonly checks that a JSON merge patch is not modifying a readonly struct member
// TODO: Need to use name from JSON tag
// TODO: Need to desend into objects
func CheckPatchReadonly(strucObj interface{}, patch map[string]interface{}) []string {
	t := reflect.TypeOf(strucObj)
	key := "meta"
	wantedValue := "readOnly"

	result := []string{}

	for i := 0; i < t.NumField(); i++ {
		field := t.Field(i)

		if tagValue, ok := field.Tag.Lookup(key); ok && tagValue == wantedValue {
			if _, present := patch[field.Name]; present {
				result = append(result, field.Name)
			}
		}
	}

	return result
}

// func JSONFieldNames(obj interface{}) {
// 	t := reflect.TypeOf(obj)

// 	// Struct field name -> json name
// 	result := map[string]string{}

// 	for i := 0; i < t.NumField(); i++ {
// 		field := t.Field(i)
// 		if jsonTag := field.Tag.Get("json"); jsonTag != "" && jsonTag != "-" {
// 			fieldName := field.Name
// 			if commaIdx := strings.Index(jsonTag, ","); commaIdx > 0 {
// 				fieldName = jsonTag[:commaIdx]
// 			}

// 			result[field.Name] = fieldName
// 		}
// 	}
// }

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
