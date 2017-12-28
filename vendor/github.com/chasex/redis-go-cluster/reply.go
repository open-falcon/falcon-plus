// Copyright 2015 Joel Wu
// Copyright 2012 Gary Burd
//
// Licensed under the Apache License, Version 2.0 (the "License"): you may
// not use this file except in compliance with the License. You may obtain
// a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package redis

import (
	"fmt"
	"errors"
	"strconv"
	"reflect"
)

// ErrNil indicates that a reply value is nil.
var ErrNil = errors.New("nil reply")

// Int is a helper that converts a command reply to an integer. If err is not
// equal to nil, then Int returns 0, err. Otherwise, Int converts the
// reply to an int as follows:
//
//  Reply type    Result
//  integer       int(reply), nil
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
func Int(reply interface{}, err error) (int, error) {
    if err != nil {
	return 0, err
    }
    switch reply := reply.(type) {
    case int64:
	x := int(reply)
	if int64(x) != reply {
	    return 0, strconv.ErrRange
	}
	return x, nil
    case []byte:
	n, err := strconv.ParseInt(string(reply), 10, 0)
	return int(n), err
    case nil:
	return 0, ErrNil
    case redisError:
	return 0, reply
    }
    return 0, fmt.Errorf("unexpected type %T for Int", reply)
}

// Int64 is a helper that converts a command reply to 64 bit integer. If err is
// not equal to nil, then Int returns 0, err. Otherwise, Int64 converts the
// reply to an int64 as follows:
//
//  Reply type    Result
//  integer       reply, nil
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
func Int64(reply interface{}, err error) (int64, error) {
    if err != nil {
	return 0, err
    }
    switch reply := reply.(type) {
    case int64:
	return reply, nil
    case []byte:
	n, err := strconv.ParseInt(string(reply), 10, 64)
	return n, err
    case nil:
	return 0, ErrNil
    case redisError:
	return 0, reply
    }
    return 0, fmt.Errorf("unexpected type %T for Int64", reply)
}

// Float64 is a helper that converts a command reply to 64 bit float. If err is
// not equal to nil, then Float64 returns 0, err. Otherwise, Float64 converts
// the reply to an int as follows:
//
//  Reply type    Result
//  bulk string   parsed reply, nil
//  nil           0, ErrNil
//  other         0, error
func Float64(reply interface{}, err error) (float64, error) {
    if err != nil {
	return 0, err
    }
    switch reply := reply.(type) {
    case []byte:
	n, err := strconv.ParseFloat(string(reply), 64)
	return n, err
    case nil:
	return 0, ErrNil
    case redisError:
	return 0, reply
    }
    return 0, fmt.Errorf("unexpected type %T for Float64", reply)
}

// String is a helper that converts a command reply to a string. If err is not
// equal to nil, then String returns "", err. Otherwise String converts the
// reply to a string as follows:
//
//  Reply type      Result
//  bulk string     string(reply), nil
//  simple string   reply, nil
//  nil             "",  ErrNil
//  other           "",  error
func String(reply interface{}, err error) (string, error) {
    if err != nil {
	return "", err
    }
    switch reply := reply.(type) {
    case []byte:
	return string(reply), nil
    case string:
	return reply, nil
    case nil:
	return "", ErrNil
    case redisError:
	return "", reply
    }
    return "", fmt.Errorf("unexpected type %T for String", reply)
}

// Bytes is a helper that converts a command reply to a slice of bytes. If err
// is not equal to nil, then Bytes returns nil, err. Otherwise Bytes converts
// the reply to a slice of bytes as follows:
//
//  Reply type      Result
//  bulk string     reply, nil
//  simple string   []byte(reply), nil
//  nil             nil, ErrNil
//  other           nil, error
func Bytes(reply interface{}, err error) ([]byte, error) {
    if err != nil {
	return nil, err
    }
    switch reply := reply.(type) {
    case []byte:
	return reply, nil
    case string:
	return []byte(reply), nil
    case nil:
	return nil, ErrNil
    case redisError:
	return nil, reply
    }
    return nil, fmt.Errorf("unexpected type %T for Bytes", reply)
}

// Bool is a helper that converts a command reply to a boolean. If err is not
// equal to nil, then Bool returns false, err. Otherwise Bool converts the
// reply to boolean as follows:
//
//  Reply type      Result
//  integer         value != 0, nil
//  bulk string     strconv.ParseBool(reply)
//  nil             false, ErrNil
//  other           false, error
func Bool(reply interface{}, err error) (bool, error) {
    if err != nil {
	return false, err
    }
    switch reply := reply.(type) {
    case int64:
	return reply != 0, nil
    case []byte:
	return strconv.ParseBool(string(reply))
    case nil:
	return false, ErrNil
    case redisError:
	return false, reply
    }
    return false, fmt.Errorf("unexpected type %T for Bool", reply)
}

// Values is a helper that converts an array command reply to a []interface{}.
// If err is not equal to nil, then Values returns nil, err. Otherwise, Values
// converts the reply as follows:
//
//  Reply type      Result
//  array           reply, nil
//  nil             nil, ErrNil
//  other           nil, error
func Values(reply interface{}, err error) ([]interface{}, error) {
    if err != nil {
	return nil, err
    }
    switch reply := reply.(type) {
    case []interface{}:
	return reply, nil
    case nil:
	return nil, ErrNil
    case redisError:
	return nil, reply
    }
    return nil, fmt.Errorf("unexpected type %T for Values", reply)
}

// Ints is a helper that converts an array command reply to a []int. 
// If err is not equal to nil, then Ints returns nil, err.
func Ints(reply interface{}, err error) ([]int, error) {
    values, err := Values(reply, err)
    if err != nil {
	return nil, err
    }

    ints := make([]int, len(values))
    slice := make([]interface{}, len(values))
    for i, _ := range ints {
	slice[i] = &ints[i]
    }

    if _, err = Scan(values, slice...); err != nil {
	return nil, err
    }

    return ints, nil
}

// Strings is a helper that converts an array command reply to a []string. If
// err is not equal to nil, then Strings returns nil, err. Nil array items are
// converted to "" in the output slice. Strings returns an error if an array
// item is not a bulk string or nil.
func Strings(reply interface{}, err error) ([]string, error) {
    values, err := Values(reply, err)
    if err != nil {
	return nil, err
    }

    strings := make([]string, len(values))
    slice := make([]interface{}, len(values))
    for i, _ := range strings {
	slice[i] = &strings[i]
    }

    if _, err = Scan(values, slice...); err != nil {
	return nil, err
    }

    return strings, nil
}

// StringMap is a helper that converts an array of strings (alternating key, value)
// into a map[string]string. The HGETALL and CONFIG GET commands return replies in this format.
// Requires an even number of values in result.
func StringMap(result interface{}, err error) (map[string]string, error) {
    values, err := Values(result, err)
    if err != nil {
	return nil, err
    }
    if len(values) % 2 != 0 {
	return nil, errors.New("expect even number elements for StringMap")
    }

    m := make(map[string]string, len(values) / 2)
    for i := 0; i < len(values); i += 2 {
	key, okKey := values[i].([]byte)
	value, okValue := values[i + 1].([]byte)
	if !okKey || !okValue {
	    return nil, errors.New("expect bulk string for StringMap")
	}
	m[string(key)] = string(value)
    }

    return m, nil
}

// Scan copies from src to the values pointed at by dest.
//
// The values pointed at by dest must be an integer, float, boolean, string,
// []byte, interface{} or slices of these types. Scan uses the standard strconv
// package to convert bulk strings to numeric and boolean types.
//
// If a dest value is nil, then the corresponding src value is skipped.
//
// If a src element is nil, then the corresponding dest value is not modified.
//
// To enable easy use of Scan in a loop, Scan returns the slice of src
// following the copied values.
func Scan(src []interface{}, dst ...interface{}) ([]interface{}, error) {
    if len(src) < len(dst) {
	return nil, errors.New("mismatch length of source and dest")
    }
    var err error
    for i, d := range dst {
	err = convertAssign(d, src[i])
	if err != nil {
	    break
	}
    }
    return src[len(dst):], err
}

func ensureLen(d reflect.Value, n int) {
    if n > d.Cap() {
	d.Set(reflect.MakeSlice(d.Type(), n, n))
    } else {
	d.SetLen(n)
    }
}

func cannotConvert(d reflect.Value, s interface{}) error {
    return fmt.Errorf("redigo: Scan cannot convert from %s to %s",
	reflect.TypeOf(s), d.Type())
}

func convertAssignBytes(d reflect.Value, s []byte) (err error) {
    switch d.Type().Kind() {
    case reflect.Float32, reflect.Float64:
	var x float64
	x, err = strconv.ParseFloat(string(s), d.Type().Bits())
	d.SetFloat(x)
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	var x int64
	x, err = strconv.ParseInt(string(s), 10, d.Type().Bits())
	d.SetInt(x)
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	var x uint64
	x, err = strconv.ParseUint(string(s), 10, d.Type().Bits())
	d.SetUint(x)
    case reflect.Bool:
	var x bool
	x, err = strconv.ParseBool(string(s))
	d.SetBool(x)
    case reflect.String:
	d.SetString(string(s))
    case reflect.Slice:
	if d.Type().Elem().Kind() != reflect.Uint8 {
	    err = cannotConvert(d, s)
	} else {
	    d.SetBytes(s)
	}
    default:
	err = cannotConvert(d, s)
    }
    return
}

func convertAssignInt(d reflect.Value, s int64) (err error) {
    switch d.Type().Kind() {
    case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
	d.SetInt(s)
	if d.Int() != s {
	    err = strconv.ErrRange
	    d.SetInt(0)
	}
    case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
	if s < 0 {
	    err = strconv.ErrRange
	} else {
	    x := uint64(s)
	    d.SetUint(x)
	    if d.Uint() != x {
		err = strconv.ErrRange
		d.SetUint(0)
	    }
	}
    case reflect.Bool:
	d.SetBool(s != 0)
    default:
	err = cannotConvert(d, s)
    }
    return
}

func convertAssignValue(d reflect.Value, s interface{}) (err error) {
    switch s := s.(type) {
    case []byte:
	err = convertAssignBytes(d, s)
    case int64:
	err = convertAssignInt(d, s)
    default:
	err = cannotConvert(d, s)
    }
    return err
}

func convertAssignValues(d reflect.Value, s []interface{}) error {
    if d.Type().Kind() != reflect.Slice {
	return cannotConvert(d, s)
    }
    ensureLen(d, len(s))
    for i := 0; i < len(s); i++ {
	if err := convertAssignValue(d.Index(i), s[i]); err != nil {
	    return err
	}
    }
    return nil
}

func convertAssign(d interface{}, s interface{}) (err error) {
    // Handle the most common destination types using type switches and
    // fall back to reflection for all other types.
    switch s := s.(type) {
    case nil:
	// ingore
    case []byte:
	switch d := d.(type) {
	case *string:
	    *d = string(s)
	case *int:
	    *d, err = strconv.Atoi(string(s))
	case *int64:
	    *d, err = strconv.ParseInt(string(s), 10, 64)
	case *bool:
	    *d, err = strconv.ParseBool(string(s))
	case *[]byte:
	    *d = s
	case *interface{}:
	    *d = s
	case nil:
	    // skip value
	default:
	    if d := reflect.ValueOf(d); d.Type().Kind() != reflect.Ptr {
		err = cannotConvert(d, s)
	    } else {
		err = convertAssignBytes(d.Elem(), s)
	    }
	}
    case int64:
	switch d := d.(type) {
	case *int:
	    x := int(s)
	    if int64(x) != s {
		err = strconv.ErrRange
		x = 0
	    }
	    *d = x
	case *int64:
	    *d = s
	case *bool:
	    *d = s != 0
	case *interface{}:
	    *d = s
	case nil:
	    // skip value
	default:
	    if d := reflect.ValueOf(d); d.Type().Kind() != reflect.Ptr {
		err = cannotConvert(d, s)
	    } else {
		err = convertAssignInt(d.Elem(), s)
	    }
	}
    case []interface{}:
	switch d := d.(type) {
	case *[]interface{}:
	    *d = s
	case *interface{}:
	    *d = s
	case nil:
	    // skip value
	default:
	    if d := reflect.ValueOf(d); d.Type().Kind() != reflect.Ptr {
		err = cannotConvert(d, s)
	    } else {
		err = convertAssignValues(d.Elem(), s)
	    }
	}
    case redisError:
	err = s
    default:
	err = cannotConvert(reflect.ValueOf(d), s)
    }
    return
}
