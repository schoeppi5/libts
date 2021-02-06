package communication

import "github.com/mitchellh/mapstructure"

// MarshalRequest takes an interface and decodes it to a map
// can be used to send structs of arguments to the query
func MarshalRequest(i interface{}) (map[string]interface{}, error) {
	m := map[string]interface{}{}
	err := mapstructure.Decode(i, &m)
	if err != nil {
		return nil, err
	}
	return m, nil
}
