// Code generated by sql2gql. DO NOT EDIT.

package repo

import (
	"context"

	sq "github.com/elgris/sqrl"
	"github.com/google/wire"
	"github.com/ketan-10/sql2gql/examples/pokemon/internal"
	"github.com/ketan-10/sql2gql/examples/pokemon/sql2gql/table"
)

type IPokemonEvolutionMatchupRepository interface {
	IPokemonEvolutionMatchupRepositoryQueryBuilder

	InsertPokemonEvolutionMatchup(ctx context.Context, pem table.PokemonEvolutionMatchupCreate) (*table.PokemonEvolutionMatchup, error)
	InsertPokemonEvolutionMatchupWithSuffix(ctx context.Context, pem table.PokemonEvolutionMatchupCreate, suffix sq.Sqlizer) (*table.PokemonEvolutionMatchup, error)
	InsertPokemonEvolutionMatchupIDResult(ctx context.Context, pem table.PokemonEvolutionMatchupCreate, suffix sq.Sqlizer) (int64, error)

	UpdatePokemonEvolutionMatchupByFields(ctx context.Context, id int, pem table.PokemonEvolutionMatchupUpdate) (*table.PokemonEvolutionMatchup, error)
	UpdatePokemonEvolutionMatchup(ctx context.Context, pem table.PokemonEvolutionMatchup) (*table.PokemonEvolutionMatchup, error)

	DeletePokemonEvolutionMatchup(ctx context.Context, pem table.PokemonEvolutionMatchup) error
	DeletePokemonEvolutionMatchupByID(ctx context.Context, id int) (bool, error)

	FindAllPokemonEvolutionMatchup(ctx context.Context, pem *table.PokemonEvolutionMatchupFilter, pagination *internal.Pagination) (*table.ListPokemonEvolutionMatchup, error)
	FindAllPokemonEvolutionMatchupWithSuffix(ctx context.Context, pem *table.PokemonEvolutionMatchupFilter, pagination *internal.Pagination, suffixes ...sq.Sqlizer) (*table.ListPokemonEvolutionMatchup, error)
	PokemonEvolutionMatchupByID(ctx context.Context, iD int, filter *table.PokemonEvolutionMatchupFilter) (*table.PokemonEvolutionMatchup, error)

	PokemonEvolutionMatchupByIDWithSuffix(ctx context.Context, iD int, filter *table.PokemonEvolutionMatchupFilter, suffixes ...sq.Sqlizer) (*table.PokemonEvolutionMatchup, error)

	PokemonEvolutionMatchupByPokemonID(ctx context.Context, pokemonID int, filter *table.PokemonEvolutionMatchupFilter, pagination *internal.Pagination) (*table.ListPokemonEvolutionMatchup, error)

	PokemonEvolutionMatchupByPokemonIDWithSuffix(ctx context.Context, pokemonID int, filter *table.PokemonEvolutionMatchupFilter, pagination *internal.Pagination, suffixes ...sq.Sqlizer) (*table.ListPokemonEvolutionMatchup, error)
}

type IPokemonEvolutionMatchupRepositoryQueryBuilder interface {
	FindAllPokemonEvolutionMatchupBaseQuery(ctx context.Context, filter *table.PokemonEvolutionMatchupFilter, fields string, suffix ...sq.Sqlizer) (*sq.SelectBuilder, error)
	AddPagination(ctx context.Context, qb *sq.SelectBuilder, pagination *internal.Pagination) (*sq.SelectBuilder, error)
}

type PokemonEvolutionMatchupRepository struct {
	DB           internal.IDb
	QueryBuilder IPokemonEvolutionMatchupRepositoryQueryBuilder
}

type PokemonEvolutionMatchupRepositoryQueryBuilder struct {
}

var NewPokemonEvolutionMatchupRepository = wire.NewSet(
	wire.Struct(new(PokemonEvolutionMatchupRepository), "*"),
	wire.Struct(new(PokemonEvolutionMatchupRepositoryQueryBuilder), "*"),
	wire.Bind(new(IPokemonEvolutionMatchupRepository), new(*PokemonEvolutionMatchupRepository)),
	wire.Bind(new(IPokemonEvolutionMatchupRepositoryQueryBuilder), new(*PokemonEvolutionMatchupRepositoryQueryBuilder)),
)

func (pemr *PokemonEvolutionMatchupRepository) InsertPokemonEvolutionMatchup(ctx context.Context, pem table.PokemonEvolutionMatchupCreate) (*table.PokemonEvolutionMatchup, error) {
	return pemr.InsertPokemonEvolutionMatchupWithSuffix(ctx, pem, nil)
}

func (pemr *PokemonEvolutionMatchupRepository) InsertPokemonEvolutionMatchupWithSuffix(ctx context.Context, pem table.PokemonEvolutionMatchupCreate, suffix sq.Sqlizer) (*table.PokemonEvolutionMatchup, error) {
	var err error

	id, err := pemr.InsertPokemonEvolutionMatchupIDResult(ctx, pem, suffix)
	if err != nil {
		return nil, err
	}
	newpem := table.PokemonEvolutionMatchup{}
	qb := sq.Select("*").From(`pokemon_evolution_matchup`)

	qb.Where(sq.Eq{"`id`": id})
	err = pemr.DB.Get(ctx, &newpem, qb)

	if err != nil {
		return nil, err
	}
	return &newpem, nil
}

func (pemr *PokemonEvolutionMatchupRepository) InsertPokemonEvolutionMatchupIDResult(ctx context.Context, pem table.PokemonEvolutionMatchupCreate, suffix sq.Sqlizer) (int64, error) {
	var err error

	qb := sq.Insert("`pokemon_evolution_matchup`").Columns(
		"`pokemon_id`",
		"`evolves_from_species_id`",
		"`habitat`",
		"`gender_rate`",
		"`capture_rate`",
		"`base_happiness`",
	).Values(
		pem.PokemonID,
		pem.EvolvesFromSpeciesID,
		pem.Habitat,
		pem.GenderRate,
		pem.CaptureRate,
		pem.BaseHappiness,
	)
	if suffix != nil {
		suffixQuery, suffixArgs, suffixErr := suffix.ToSql()
		if suffixErr != nil {
			return 0, suffixErr
		}
		qb.Suffix(suffixQuery, suffixArgs...)
	}

	// run query
	res, err := pemr.DB.Exec(ctx, qb)
	if err != nil {
		return 0, err
	}

	// retrieve id
	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (pemr *PokemonEvolutionMatchupRepository) UpdatePokemonEvolutionMatchupByFields(ctx context.Context, id int, pem table.PokemonEvolutionMatchupUpdate) (*table.PokemonEvolutionMatchup, error) {
	var err error

	updateMap := map[string]interface{}{}
	if pem.PokemonID != nil {
		updateMap["`pokemon_id`"] = *pem.PokemonID
	}
	if pem.EvolvesFromSpeciesID != nil {
		updateMap["`evolves_from_species_id`"] = *pem.EvolvesFromSpeciesID
	}
	if pem.Habitat != nil {
		updateMap["`habitat`"] = *pem.Habitat
	}
	if pem.GenderRate != nil {
		updateMap["`gender_rate`"] = *pem.GenderRate
	}
	if pem.CaptureRate != nil {
		updateMap["`capture_rate`"] = *pem.CaptureRate
	}
	if pem.BaseHappiness != nil {
		updateMap["`base_happiness`"] = *pem.BaseHappiness
	}
	if pem.Active != nil {
		updateMap["`active`"] = *pem.Active
	}

	qb := sq.Update(`pokemon_evolution_matchup`).SetMap(updateMap).Where(sq.Eq{"`id`": id})

	_, err = pemr.DB.Exec(ctx, qb)
	if err != nil {
		return nil, err
	}

	selectQb := sq.Select("*").From("`pokemon_evolution_matchup`")

	selectQb = selectQb.Where(sq.Eq{"`id`": id})

	result := table.PokemonEvolutionMatchup{}
	err = pemr.DB.Get(ctx, &result, selectQb)
	if err != nil {
		return nil, err
	}

	return &result, nil

}

func (pemr *PokemonEvolutionMatchupRepository) UpdatePokemonEvolutionMatchup(ctx context.Context, pem table.PokemonEvolutionMatchup) (*table.PokemonEvolutionMatchup, error) {
	var err error

	// sql query
	qb := sq.Update("`pokemon_evolution_matchup`").SetMap(map[string]interface{}{
		"`pokemon_id`":              pem.PokemonID,
		"`evolves_from_species_id`": pem.EvolvesFromSpeciesID,
		"`habitat`":                 pem.Habitat,
		"`gender_rate`":             pem.GenderRate,
		"`capture_rate`":            pem.CaptureRate,
		"`base_happiness`":          pem.BaseHappiness,
		"`active`":                  pem.Active,
	}).Where(sq.Eq{"`id`": pem.ID})

	// run query
	_, err = pemr.DB.Exec(ctx, qb)
	if err != nil {
		return nil, err
	}

	selectQb := sq.Select("*").From("`pokemon_evolution_matchup`")
	selectQb = selectQb.Where(sq.Eq{"`id`": pem.ID})

	result := table.PokemonEvolutionMatchup{}
	err = pemr.DB.Get(ctx, &result, selectQb)
	if err != nil {
		return nil, err
	}

	return &result, nil
}

func (pemr *PokemonEvolutionMatchupRepository) DeletePokemonEvolutionMatchup(ctx context.Context, pem table.PokemonEvolutionMatchup) error {
	_, err := pemr.DeletePokemonEvolutionMatchupByID(ctx, pem.ID)
	return err
}

func (pemr *PokemonEvolutionMatchupRepository) DeletePokemonEvolutionMatchupByID(ctx context.Context, id int) (bool, error) {
	var err error

	qb := sq.Update("`pokemon_evolution_matchup`").Set("active", false)

	qb = qb.Where(sq.Eq{"`id`": id})

	_, err = pemr.DB.Exec(ctx, qb)
	if err != nil {
		return false, err
	}
	return true, nil
}

func (pemr *PokemonEvolutionMatchupRepository) FindAllPokemonEvolutionMatchupBaseQuery(ctx context.Context, filter *table.PokemonEvolutionMatchupFilter, fields string, suffixes ...sq.Sqlizer) (*sq.SelectBuilder, error) {
	return pemr.QueryBuilder.FindAllPokemonEvolutionMatchupBaseQuery(ctx, filter, fields, suffixes...)
}

func (pemr *PokemonEvolutionMatchupRepositoryQueryBuilder) FindAllPokemonEvolutionMatchupBaseQuery(ctx context.Context, filter *table.PokemonEvolutionMatchupFilter, fields string, suffixes ...sq.Sqlizer) (*sq.SelectBuilder, error) {
	var err error
	qb := sq.Select(fields).From("`pokemon_evolution_matchup`")
	if filter != nil {
		if qb, err = internal.AddFilter(qb, "`pokemon_evolution_matchup`.`id`", filter.ID); err != nil {
			return qb, err
		}
		if qb, err = internal.AddFilter(qb, "`pokemon_evolution_matchup`.`pokemon_id`", filter.PokemonID); err != nil {
			return qb, err
		}
		if qb, err = internal.AddFilter(qb, "`pokemon_evolution_matchup`.`evolves_from_species_id`", filter.EvolvesFromSpeciesID); err != nil {
			return qb, err
		}
		if qb, err = internal.AddFilter(qb, "`pokemon_evolution_matchup`.`habitat`", filter.Habitat); err != nil {
			return qb, err
		}
		if qb, err = internal.AddFilter(qb, "`pokemon_evolution_matchup`.`gender_rate`", filter.GenderRate); err != nil {
			return qb, err
		}
		if qb, err = internal.AddFilter(qb, "`pokemon_evolution_matchup`.`capture_rate`", filter.CaptureRate); err != nil {
			return qb, err
		}
		if qb, err = internal.AddFilter(qb, "`pokemon_evolution_matchup`.`base_happiness`", filter.BaseHappiness); err != nil {
			return qb, err
		}
		if filter.Active == nil {
			if qb, err = internal.AddFilter(qb, "`pokemon_evolution_matchup`.`active`", internal.FilterOnField{{internal.Eq: true}}); err != nil {
				return qb, err
			}
		} else {
			if qb, err = internal.AddFilter(qb, "`pokemon_evolution_matchup`.`active`", filter.Active); err != nil {
				return qb, err
			}
		}
		if qb, err = internal.AddFilter(qb, "`pokemon_evolution_matchup`.`created_at`", filter.CreatedAt); err != nil {
			return qb, err
		}
		if qb, err = internal.AddFilter(qb, "`pokemon_evolution_matchup`.`updated_at`", filter.UpdatedAt); err != nil {
			return qb, err
		}
		qb, err = internal.AddAdditionalFilter(qb, filter.Wheres, filter.Joins, filter.LeftJoins, filter.GroupBys, filter.Havings)
		if err != nil {
			return qb, err
		}
	} else {
		if qb, err = internal.AddFilter(qb, "`pokemon_evolution_matchup`.`active`", internal.FilterOnField{{internal.Eq: true}}); err != nil {
			return qb, err
		}
	}

	for _, suffix := range suffixes {
		query, args, err := suffix.ToSql()
		if err != nil {
			return qb, err
		}
		qb.Suffix(query, args...)
	}
	return qb, nil
}

func (pemr *PokemonEvolutionMatchupRepository) AddPagination(ctx context.Context, qb *sq.SelectBuilder, pagination *internal.Pagination) (*sq.SelectBuilder, error) {
	return pemr.QueryBuilder.AddPagination(ctx, qb, pagination)
}

func (pem *PokemonEvolutionMatchupRepositoryQueryBuilder) AddPagination(ctx context.Context, qb *sq.SelectBuilder, pagination *internal.Pagination) (*sq.SelectBuilder, error) {
	fields := []string{
		"id",
		"pokemon_id",
		"evolves_from_species_id",
		"habitat",
		"gender_rate",
		"capture_rate",
		"base_happiness",
		"active",
		"created_at",
		"updated_at",
	}
	return internal.AddPagination(qb, pagination, "pokemon_evolution_matchup", fields)
}

func (pemr *PokemonEvolutionMatchupRepository) FindAllPokemonEvolutionMatchup(ctx context.Context, filter *table.PokemonEvolutionMatchupFilter, pagination *internal.Pagination) (*table.ListPokemonEvolutionMatchup, error) {
	return pemr.FindAllPokemonEvolutionMatchupWithSuffix(ctx, filter, pagination)
}

func (pemr *PokemonEvolutionMatchupRepository) FindAllPokemonEvolutionMatchupWithSuffix(ctx context.Context, filter *table.PokemonEvolutionMatchupFilter, pagination *internal.Pagination, suffixes ...sq.Sqlizer) (*table.ListPokemonEvolutionMatchup, error) {
	var list table.ListPokemonEvolutionMatchup
	qb, err := pemr.FindAllPokemonEvolutionMatchupBaseQuery(ctx, filter, "`pokemon_evolution_matchup`.*", suffixes...)
	if err != nil {
		return &list, err
	}
	qb, err = pemr.AddPagination(ctx, qb, pagination)
	if err != nil {
		return &list, err
	}

	err = pemr.DB.Select(ctx, &list.Data, qb)

	if err != nil {
		return &list, err
	}

	if pagination == nil || pagination.PerPage == nil || pagination.Page == nil {
		list.TotalCount = len(list.Data)
		return &list, nil
	}

	var listMeta internal.ListMetadata
	if qb, err = pemr.FindAllPokemonEvolutionMatchupBaseQuery(ctx, filter, "COUNT(1) AS count"); err != nil {
		return &table.ListPokemonEvolutionMatchup{}, err
	}
	if filter != nil && len(filter.GroupBys) > 0 {
		qb = sq.Select("COUNT(1) AS count").FromSelect(qb, "a")
	}
	err = pemr.DB.Get(ctx, &listMeta, qb)

	list.TotalCount = listMeta.Count

	return &list, err
}
func (pemr *PokemonEvolutionMatchupRepository) PokemonEvolutionMatchupByID(ctx context.Context, iD int, filter *table.PokemonEvolutionMatchupFilter) (*table.PokemonEvolutionMatchup, error) {
	return pemr.PokemonEvolutionMatchupByIDWithSuffix(ctx, iD, filter)
}

func (pemr *PokemonEvolutionMatchupRepository) PokemonEvolutionMatchupByIDWithSuffix(ctx context.Context, iD int, filter *table.PokemonEvolutionMatchupFilter, suffixes ...sq.Sqlizer) (*table.PokemonEvolutionMatchup, error) {
	var err error

	// sql query
	qb, err := pemr.FindAllPokemonEvolutionMatchupBaseQuery(ctx, filter, "`pokemon_evolution_matchup`.*", suffixes...)
	if err != nil {
		return &table.PokemonEvolutionMatchup{}, err
	}
	qb = qb.Where(sq.Eq{"`pokemon_evolution_matchup`.`id`": iD})

	// run query
	pem := table.PokemonEvolutionMatchup{}
	err = pemr.DB.Get(ctx, &pem, qb)
	if err != nil {
		return &table.PokemonEvolutionMatchup{}, err
	}
	return &pem, nil
}

func (pemr *PokemonEvolutionMatchupRepository) PokemonEvolutionMatchupByPokemonID(ctx context.Context, pokemonID int, filter *table.PokemonEvolutionMatchupFilter, pagination *internal.Pagination) (*table.ListPokemonEvolutionMatchup, error) {
	return pemr.PokemonEvolutionMatchupByPokemonIDWithSuffix(ctx, pokemonID, filter, pagination)
}

func (pemr *PokemonEvolutionMatchupRepository) PokemonEvolutionMatchupByPokemonIDWithSuffix(ctx context.Context, pokemonID int, filter *table.PokemonEvolutionMatchupFilter, pagination *internal.Pagination, suffixes ...sq.Sqlizer) (*table.ListPokemonEvolutionMatchup, error) {

	var list table.ListPokemonEvolutionMatchup
	// sql query
	qb, err := pemr.FindAllPokemonEvolutionMatchupBaseQuery(ctx, filter, "`pokemon_evolution_matchup`.*", suffixes...)
	if err != nil {
		return &list, err
	}
	qb = qb.Where(sq.Eq{"`pokemon_evolution_matchup`.`pokemon_id`": pokemonID})

	if qb, err = pemr.AddPagination(ctx, qb, pagination); err != nil {
		return &list, err
	}

	// run query
	if err = pemr.DB.Select(ctx, &list.Data, qb); err != nil {
		return &list, err
	}

	if pagination == nil || pagination.PerPage == nil || pagination.Page == nil {
		list.TotalCount = len(list.Data)
		return &list, nil
	}

	var listMeta internal.ListMetadata
	if qb, err = pemr.FindAllPokemonEvolutionMatchupBaseQuery(ctx, filter, "COUNT(1) AS count"); err != nil {
		return &list, err
	}
	if filter != nil && len(filter.GroupBys) > 0 {
		qb = sq.Select("COUNT(1) AS count").FromSelect(qb, "a")
	}
	qb = qb.Where(sq.Eq{"`pokemon_evolution_matchup`.`pokemon_id`": pokemonID})
	if err = pemr.DB.Get(ctx, &listMeta, qb); err != nil {
		return &list, err
	}

	list.TotalCount = listMeta.Count

	return &list, nil

}
