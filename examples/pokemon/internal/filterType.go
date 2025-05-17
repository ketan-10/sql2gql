package internal

import (
	"encoding/json"
	"io"
)

type FilterType string

const (
	Eq      FilterType = "eq"
	Neq     FilterType = "Neq"
	Gt      FilterType = "gt"
	Gte     FilterType = "gte"
	Lt      FilterType = "lt"
	Lte     FilterType = "lte"
	Like    FilterType = "like"
	Between FilterType = "between"
)

type FilterOnField []map[FilterType]interface{}

// UnmarshalGQL implements the graphql.Marshaler interface
func (f *FilterOnField) UnmarshalGQL(v interface{}) error {
	var err error
	vjson, _ := json.Marshal(v)
	if json.Unmarshal(vjson, f) == nil {
		return nil
	}
	singleMap := map[FilterType]interface{}{}
	err = json.Unmarshal(vjson, &singleMap)
	if err == nil {
		*f = []map[FilterType]interface{}{
			singleMap,
		}
		return nil
	}
	*f = []map[FilterType]interface{}{
		{Eq: v},
	}
	return nil
}


// MarshalGQL implements the graphql.Marshaler interface
func (f FilterOnField) MarshalGQL(w io.Writer) {
	data, err := json.Marshal(f)
	if err != nil {
		w.Write([]byte(`null`))
	} else {
		w.Write([]byte(`"` + string(data) + `"`))
	}
}

// type FilterType int

// const (
// 	Eq      FilterType = iota
// 	Neq
// 	Gt
// 	Gte
// 	Lt
// 	Lte
// 	Like
// 	Between
// )

// func (ft *FilterType) String() string {
// 	switch *ft {
// 	case Eq:
// 		return "eq"
// 	}
// 	...
// 	return ""
// }

// func (ft *FilterType) UnMarshal(value string) err {
// 	switch value {
// 	case "eq":
// 		*ft = Eq
// 		return nil
// 	}
// 	...
// 	return errors.New("invalid value")
// }
