package execution

import "fmt"

// table effectively partitions the value space first by breeze type, and then
// for each, uses an actual map. Exceptions exist for bool and null, whose
// value space is small enough to effectively hardcode with a field.
type table struct {
	stringMap    map[string]interface{}
	numberMap    map[float64]interface{}
	boolMap      map[bool]interface{}
	nullKeyValue interface{}
	// Unfortunately, we need this here because otherwise we cannot distinguish
	// between a nil value that was explicitly set for the null key vs. a nil
	// value that is set because the null key does not exist.
	nullKeyExists bool
}

func newTable() *table {
	return &table{
		stringMap:     make(map[string]interface{}),
		numberMap:     make(map[float64]interface{}),
		boolMap:       make(map[bool]interface{}),
		nullKeyValue:  nil,
		nullKeyExists: false,
	}
}

func (t *table) Get(key interface{}) interface{} {
	val, _ := t.GetOK(key)
	return val
}

func (t *table) GetOK(key interface{}) (interface{}, bool) {
	if key == nil {
		// This can be written quicker as return t.nullKeyValue, t.nullKeyExists,
		// but I think this approach is much easier to read to myself in the future
		// when I'm not as clever.
		if t.nullKeyExists {
			return t.nullKeyValue, true
		}

		return nil, false
	}

	var val interface{}
	var ok bool
	switch typedKey := key.(type) {
	case string:
		val, ok = t.stringMap[typedKey]
	case float64:
		val, ok = t.numberMap[typedKey]
	case bool:
		val, ok = t.boolMap[typedKey]
	default:
		panic(fmt.Sprintf("bad type given to map: %T", typedKey))
	}

	return val, ok
}

func (t *table) Set(key interface{}, value interface{}) {
	if key == nil {
		t.nullKeyValue = value
		t.nullKeyExists = true
		return
	}

	switch typedKey := key.(type) {
	case string:
		t.stringMap[typedKey] = value
	case float64:
		t.numberMap[typedKey] = value
	case bool:
		t.boolMap[typedKey] = value
	default:
		panic(fmt.Sprintf("bad type given to map: %T", typedKey))
	}
}

func (t *table) Has(key interface{}) bool {
	_, ok := t.GetOK(key)
	return ok
}

func (t *table) Keys() []interface{} {
	numKeys := len(t.stringMap) + len(t.numberMap) + len(t.boolMap)
	if t.nullKeyExists {
		numKeys++
	}
	keys := make([]interface{}, 0, numKeys)

	for k := range t.stringMap {
		keys = append(keys, k)
	}

	for k := range t.numberMap {
		keys = append(keys, k)
	}

	for k := range t.boolMap {
		keys = append(keys, k)
	}

	if t.nullKeyExists {
		keys = append(keys, nil)
	}

	return keys
}
