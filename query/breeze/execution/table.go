package execution

import "fmt"

// table effectively partitions the value space first by breeze type, and then
// for each, uses an actual map. Exceptions exist for bool and null, whose
// value space is small enough to effectively hardcode with a field.
type table struct {
	// TODO: This table does not support arrays. Now, we could attempt to add yet
	// another map here for holding arrays, but that may get hairy...
	// I wonder if making this hold breeze.Concrete is a better solution?
	// In particular, I think we're going to need to be able to ensure that each
	// breeze.Concrete has a way of providing a unique string/number to identify
	// itself amongst all other breeze.Concrete besides an equivalently valued
	// one.
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
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		f64, convertOK := convertPotentialNumber(typedKey)
		if !convertOK {
			panic(fmt.Sprintf("bad type given to map: %T", typedKey))
		}
		val, ok = t.numberMap[f64]
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
	case int, int8, int16, int32, int64, uint, uint8, uint16, uint32, uint64, float32, float64:
		f64, ok := convertPotentialNumber(typedKey)
		if !ok {
			panic(fmt.Sprintf("bad type given to map: %T", typedKey))
		}
		t.numberMap[f64] = value
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
