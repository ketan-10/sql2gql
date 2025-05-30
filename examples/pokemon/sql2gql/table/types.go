// Code generated by sql2gql. DO NOT EDIT.

package table

import (
	"database/sql"

	sq "github.com/elgris/sqrl"
	"github.com/ketan-10/sql2gql/examples/pokemon/internal"
	"github.com/pkg/errors"
)

type Types struct {
	ID        int          `json:"id" db:"id"`
	TypeName  string       `json:"type_name" db:"type_name"`
	Active    bool         `json:"active" db:"active"`
	CreatedAt sql.NullTime `json:"created_at" db:"created_at"`
	UpdatedAt sql.NullTime `json:"updated_at" db:"updated_at"`
}

type TypesFilter struct {
	ID        internal.FilterOnField
	TypeName  internal.FilterOnField
	Active    internal.FilterOnField
	CreatedAt internal.FilterOnField
	UpdatedAt internal.FilterOnField
	Wheres    []sq.Sqlizer
	Joins     []sq.Sqlizer
	LeftJoins []sq.Sqlizer
	GroupBys  []string
	Havings   []sq.Sqlizer
}

func (f *TypesFilter) NewFilter() interface{} {
	if f == nil {
		return &TypesFilter{}
	}
	return f
}

func (f *TypesFilter) TableName() string {
	return "`types`"
}

func (f *TypesFilter) ModuleName() string {
	return "types"
}

func (f *TypesFilter) IsNil() bool {
	return f == nil
}
func (f *TypesFilter) AddID(filterType internal.FilterType, v interface{}) {
	f.ID = append(f.ID, map[internal.FilterType]interface{}{filterType: v})
}
func (f *TypesFilter) AddTypeName(filterType internal.FilterType, v interface{}) {
	f.TypeName = append(f.TypeName, map[internal.FilterType]interface{}{filterType: v})
}
func (f *TypesFilter) AddActive(filterType internal.FilterType, v interface{}) {
	f.Active = append(f.Active, map[internal.FilterType]interface{}{filterType: v})
}
func (f *TypesFilter) AddCreatedAt(filterType internal.FilterType, v interface{}) {
	f.CreatedAt = append(f.CreatedAt, map[internal.FilterType]interface{}{filterType: v})
}
func (f *TypesFilter) AddUpdatedAt(filterType internal.FilterType, v interface{}) {
	f.UpdatedAt = append(f.UpdatedAt, map[internal.FilterType]interface{}{filterType: v})
}

func (f *TypesFilter) Where(v sq.Sqlizer) *TypesFilter {
	f.Wheres = append(f.Wheres, v)
	return f
}

func (f *TypesFilter) Join(j sq.Sqlizer) *TypesFilter {
	f.Joins = append(f.Joins, j)
	return f
}

func (f *TypesFilter) LeftJoin(j sq.Sqlizer) *TypesFilter {
	f.LeftJoins = append(f.LeftJoins, j)
	return f
}

func (f *TypesFilter) GroupBy(gb string) *TypesFilter {
	f.GroupBys = append(f.GroupBys, gb)
	return f
}

func (f *TypesFilter) Having(h sq.Sqlizer) *TypesFilter {
	f.Havings = append(f.Havings, h)
	return f
}

type TypesCreate struct {
	TypeName string `json:"type_name" db:"type_name"`
}

// TODO: We have to exclude AutoGenerated fields
// For now I am keeping it in, as not sure how it affects
type TypesUpdate struct {
	TypeName *string // type_name
	Active   *bool   // active
}

// helper functions
func (u *TypesUpdate) ToTypesCreate() (res TypesCreate, err error) {
	if u.TypeName != nil {
		res.TypeName = *u.TypeName
	} else {
		return res, errors.New("Value Can not be NULL")
	}
	return res, nil
}

type ListTypes struct {
	TotalCount int
	Data       []Types
}

func (l *ListTypes) GetAllID() []int {
	var res []int
	for _, item := range l.Data {
		res = append(res, item.ID)
	}
	return res
}
func (l *ListTypes) GetAllTypeName() []string {
	var res []string
	for _, item := range l.Data {
		res = append(res, item.TypeName)
	}
	return res
}
func (l *ListTypes) GetAllActive() []bool {
	var res []bool
	for _, item := range l.Data {
		res = append(res, item.Active)
	}
	return res
}
func (l *ListTypes) GetAllCreatedAt() []sql.NullTime {
	var res []sql.NullTime
	for _, item := range l.Data {
		res = append(res, item.CreatedAt)
	}
	return res
}
func (l *ListTypes) GetAllUpdatedAt() []sql.NullTime {
	var res []sql.NullTime
	for _, item := range l.Data {
		res = append(res, item.UpdatedAt)
	}
	return res
}

func (l *ListTypes) Filter(f func(item Types) bool) (res ListTypes) {
	for _, item := range l.Data {
		if f(item) {
			res.Data = append(res.Data, item)
		}
	}
	res.TotalCount = len(res.Data)
	return res
}

func (l *ListTypes) Find(f func(item Types) bool) (res Types, found bool) {
	for _, item := range l.Data {
		if f(item) {
			return item, true
		}
	}
	return Types{}, false
}
func (l *ListTypes) MapByID() (m map[int]Types) {
	m = make(map[int]Types, len(l.Data))
	for _, item := range l.Data {
		m[item.ID] = item
	}
	return m
}
