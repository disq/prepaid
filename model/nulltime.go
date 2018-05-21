package model

import (
	"fmt"
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/dynamodb"
)

// NullTime represents an time.Time value that may be null, represented as an empty string in text, and null in JSON.
type NullTime struct {
	Time  time.Time
	Valid bool // Valid is true if Time is not NULL
}

const timeFormat = "2006-01-02 15:04:05"

// UnmarshalText implements the encoding.TextUnmarshaler interface.
func (n *NullTime) UnmarshalText(text []byte) error {
	n.Valid = false

	if len(text) == 0 {
		return nil
	}

	t, err := time.Parse(timeFormat, string(text))
	if err != nil {
		return err
	}

	n.Time = t
	n.Valid = true
	return nil
}

// MarshalText implements the encoding.TextMarshaler interface.
func (n NullTime) MarshalText() ([]byte, error) {
	var s string

	if n.Valid {
		s = n.Time.Format(timeFormat)
	}

	return []byte(s), nil
}

// MarshalJSON implements the encoding/json.MarshalJSON interface.
func (n NullTime) MarshalJSON() ([]byte, error) {
	s := "null"

	if n.Valid {
		s = `"` + n.Time.Format(timeFormat) + `"`
	}

	return []byte(s), nil
}

// MarshalDynamoDBAttributeValue implements the dynamodbattribute.Marshaler interface.
func (n NullTime) MarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	if !n.Valid {
		av.NULL = aws.Bool(true)
		return nil
	}

	av.N = aws.String(fmt.Sprintf("%v", n.Time.Unix()))
	return nil
}

// UnmarshalDynamoDBAttributeValue implements the dynamodbattribute.Unmarshaler interface.
func (n *NullTime) UnmarshalDynamoDBAttributeValue(av *dynamodb.AttributeValue) error {
	n.Valid = false
	if av == nil || av.N == nil || av.NULL != nil {
		return nil
	}

	i, err := strconv.ParseInt(*av.N, 10, 64)
	if err != nil {
		return err
	}

	n.Time = time.Unix(i, 0)
	n.Valid = true
	return nil
}
