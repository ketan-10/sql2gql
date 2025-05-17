package internal

import (
	"errors"
	"fmt"
	"strings"

	sq "github.com/elgris/sqrl"
)

func contains(list []string, value string) bool{
	for _, v := range list {
		if v == value {
			return true
		}
	}
	return false
}

func AddFilter(qb *sq.SelectBuilder, columnName string, filterOnField FilterOnField) (*sq.SelectBuilder, error) {
	sqlizer, err := FilterOnFieldToSqlizer(columnName, filterOnField)
	if err != nil {
		return nil, err
	}
	if sqlizer != nil {
		qb.Where(sqlizer)
	}
	return qb, nil
}

func FilterOnFieldToSqlizer(columnName string, filterOnField FilterOnField) (sq.Sqlizer, error) {
	var and sq.And
	for _, filterList := range filterOnField {
		for filterType, value := range filterList {
			switch filterType {
			case Eq:
				and = append(and, sq.Eq{columnName: value})
			case Neq:
				and = append(and, sq.NotEq{columnName: value})
			case Gt:
				and = append(and, sq.Gt{columnName: value})
			case Gte:
				and = append(and, sq.GtOrEq{columnName: value})
			case Lt:
				and = append(and, sq.Lt{columnName: value})
			case Lte:
				and = append(and, sq.LtOrEq{columnName: value})
			case Like:
				and = append(and, sq.Expr(columnName+" LIKE ?", value))
			case Between:
				casted, ok := value.([]interface{})
				if !ok || len(casted) != 2 {
					return nil, errors.New("invalid filter on between")
				}
			}
		}
	}
	// If nil is not returned here. in fun AddFilter the queryBuilder will add "WHERE" string to query even though there was no filter.
	if and == nil {
		return nil, nil
	}
	return and, nil
}

func AddAdditionalFilter(qb *sq.SelectBuilder, wheres, joins, leftJoins []sq.Sqlizer, groupBys []string, havings []sq.Sqlizer) (*sq.SelectBuilder, error) {

	for _, where := range wheres {
		query, args, err := where.ToSql()
		if err != nil {
			return qb, err
		}
		qb = qb.Where(query, args...)
	}

	for _, join := range joins {
		query, args, err := join.ToSql()
		if err != nil {
			return qb, err
		}
		qb = qb.Join(query, args...)
	}

	for _, leftJoin := range leftJoins {
		query, args, err := leftJoin.ToSql()
		if err != nil {
			return qb, err
		}
		qb = qb.LeftJoin(query, args...)
	}

	if groupBys != nil {
		qb = qb.GroupBy(groupBys...)
	}

	for _, item := range havings {
		query, args, err := item.ToSql()
		if err != nil {
			return qb, err
		}
		qb = qb.Having(query, args...)
	}

	return qb, nil
}

// Pagination
type Pagination struct {
	Page       *int
	PerPage    *int
	Sort       []string
	CustomSort []string
}

func AddPagination(qb *sq.SelectBuilder, pagination *Pagination, tableName string, fields []string) (*sq.SelectBuilder, error) {

	tableName = strings.Trim(tableName, "`")

	if pagination == nil {
		return qb, nil
	}
	if pagination.Page != nil && pagination.PerPage != nil {
		offset := uint64((*pagination.Page - 1) * (*pagination.PerPage))
		qb = qb.Offset(offset).Limit(uint64(*pagination.PerPage))
	}
	if pagination.CustomSort != nil {
		qb = qb.OrderBy(pagination.CustomSort...)
	}
	if pagination.Sort != nil {
		var ordersQuery []string
		for _, column := range pagination.Sort {
			if !contains(fields, column) {
				return qb, fmt.Errorf("field: %s Not found for this table", column)
			}
			var order string
			if strings.HasPrefix(column, "-") {
				order = "DESC"
			} else {
				order = "ASC"
			}
			ordersQuery = append(ordersQuery, fmt.Sprintf("`%s`.`%s` %s", tableName, column, order))
		}
		qb = qb.OrderBy(ordersQuery...)
	}
	return qb, nil
}

type ListMetadata struct {
	Count int `db:"count"`
}
