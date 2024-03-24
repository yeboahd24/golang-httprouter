package data

import (
  "fmt"
  "strconv"
  "errors"
  "strings"
)

type Runtime int32

var ErrInvalidRuntimeFormat = errors.New("invalid runtime format")

func (r Runtime) MarshalJSON() ([]byte, error) {

  jsonValue := fmt.Sprintf("%d mins", r)

  quotedValue := strconv.Quote(jsonValue)

  return []byte(quotedValue), nil

}


func (r *Runtime) UnmarshalJSON(jsonValue []byte) error {

  unquotedValue, err := strconv.Unquote(string(jsonValue))
  if err != nil {
    return ErrInvalidRuntimeFormat
  }

  parts := strings.Split(unquotedValue, " ")
  
  i, err := strconv.ParseInt(parts[0], 10, 32)
  if err != nil {
    return ErrInvalidRuntimeFormat
  }
  

  *r = Runtime(i)

  return nil
}
