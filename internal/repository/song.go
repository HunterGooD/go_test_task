package repository

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"strings"

	"github.com/HunterGooD/go_test_task/internal/domain/entity"
	"github.com/HunterGooD/go_test_task/internal/domain/interfaces"
	"github.com/HunterGooD/go_test_task/pkg/utils"
	"github.com/jmoiron/sqlx"
)

type songRepository struct {
	db *sqlx.DB
	tx *sqlx.Tx

	log interfaces.Logger
}

func NewSongRepository(conn *sqlx.DB, logger interfaces.Logger) *songRepository {
	logger.Info("Creat song repo")
	return &songRepository{conn, nil, logger}
}

func (s *songRepository) WithTransaction(tx *sqlx.Tx) interfaces.SongRepository {
	s.log.Info("Creat song repo transaction")
	return &songRepository{s.db, tx, s.log}
}

func (s *songRepository) CreateSong(ctx context.Context, group_id int64, songInput *entity.SongRequest) (*entity.Song, error) {
	var err error
	songReturn := &entity.Song{}
	query := `INSERT INTO public.songs(
		m_name, m_link, m_text, m_release_date, group_id)
		VALUES ($1, $2, $3, $4, $5) RETURNING id, m_name, m_link, m_text, m_release_date, group_id`

	s.log.Info("Create song started...")
	s.log.Debug("Create song started...", map[string]any{
		"group_id":   group_id,
		"song_input": songInput,
	})

	// if transaction activ exec in transaction else db exec
	if s.tx != nil {
		s.log.Info("Create Song with transaction", map[string]any{
			"query":    query,
			"group_id": group_id,
		})
		err = s.tx.GetContext(ctx, songReturn, query, songInput.Song, songInput.Link, songInput.Text, songInput.ReleaseDate, group_id)
	} else {
		s.log.Info("Create song", map[string]any{
			"query":    query,
			"group_id": group_id,
		})
		err = s.db.GetContext(ctx, songReturn, query, songInput.Song, songInput.Link, songInput.Text, songInput.ReleaseDate, group_id)
	}

	if err != nil {
		s.log.Error("Song creating error", map[string]any{
			"err": err,
		})
		return nil, err
	}

	s.log.Info("Song creating success")
	s.log.Debug("Song creating success", map[string]any{
		"creating_song": songReturn,
	})
	return songReturn, nil
}

// GetListSong filters any fields with values for where sql
// map is expected value Name|Link|Text|ReleaseDate|GroupID
// filters ? Why not use map ?? XD
func (s *songRepository) GetListSong(ctx context.Context, offset, limit int, with_deleted bool, filters *entity.SongFilters) ([]entity.Song, error) {
	res := make([]entity.Song, 0)
	whereStatemant := ""
	query := ""
	queryTemplate := `SELECT s.*, g.id AS "group.id", g.g_name AS "group.g_name",
		 	g.created_at AS "group.created_at", g.update_at AS "group.update_at",
			g.deleted_at AS "group.deleted_at" 
		FROM groups g Inner join songs s ON g.id=s.group_id 
		%s 
		%s 
		LIMIT $1 OFFSET $2;`
	delWhere := "WHERE s.deleted_at IS NULL AND g.deleted_at IS NULL"

	params := make([]any, 0)
	params = append(params, limit, offset)

	s.log.Info("Getting list song", map[string]any{
		"diaposon":           fmt.Sprintf("%d-%d", offset, limit),
		"with_deleted_field": with_deleted,
	})
	s.log.Debug("Getting list song", map[string]any{
		"diaposon":           fmt.Sprintf("%d-%d", offset, limit),
		"with_deleted_field": with_deleted,
		"filters":            filters,
	})

	if filters != nil {
		var arg []any
		s.log.Info("Creating where statement from filter ...")

		whereStatemant, arg = utils.GetFilterString(len(params)+1, filters) //TODO: to usecase package without heavy depends

		s.log.Info("Create where statement from filter ...", map[string]any{
			"where":    whereStatemant,
			"len_args": len(arg),
		})
		s.log.Debug("Create where statement from filter ...", map[string]any{
			"where":    whereStatemant,
			"len_args": len(arg),
			"args":     arg,
		})

		if len(arg) > 0 {
			params = append(params, arg...)
		}
	}

	if whereStatemant != "" {
		if with_deleted {
			whereStatemant = " WHERE " + whereStatemant
		} else {
			whereStatemant = " AND " + whereStatemant
		}
	}

	s.log.Info("select list where statemet", map[string]any{
		"where": whereStatemant,
	})

	if with_deleted {
		query = fmt.Sprintf(queryTemplate, "", whereStatemant)
	} else {
		query = fmt.Sprintf(queryTemplate, delWhere, whereStatemant)
	}

	s.log.Info("selection query", map[string]any{
		"query":      query,
		"len_params": len(params),
	})
	s.log.Debug("selection query", map[string]any{
		"query":      query,
		"len_params": len(params),
		"params":     params,
	})

	err := s.db.SelectContext(ctx, &res, query, params...)
	if err != nil {
		s.log.Error("not success get list song", map[string]any{
			"err":    err,
			"query":  query,
			"params": params,
		})
		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}
	s.log.Info("Get list song success", map[string]any{
		"getting_songs":      len(res),
		"diaposon":           fmt.Sprintf("%d-%d", offset, limit),
		"with_deleted_field": with_deleted,
	})
	return res, err
}

// TODO: filter with reverse param and add to query ORDER BY s.id DESC
func (s *songRepository) GetReverseListSongs(ctx context.Context, offset, limit int, with_deleted bool, filters *entity.SongFilters) ([]entity.Song, error) {
	res := make([]entity.Song, 0)
	params := make([]any, 0)
	params = append(params, limit, offset)
	query := ""
	delete := "WHERE s.deleted_at IS NULL AND g.deleted_at IS NULL"
	queryTemplate := `SELECT s.*, g.id AS "group.id", g.g_name AS "group.g_name",
			g.created_at AS "group.created_at", g.update_at AS "group.update_at",
		g.deleted_at AS "group.deleted_at" 
		FROM groups g Inner join songs s ON g.id=s.group_id 
		%s 
		%s 
		ORDER BY s.id DESC
		LIMIT $1 OFFSET $2;`

	s.log.Info("Geting reverse songs", map[string]any{
		"diaposon":           fmt.Sprintf("%d-%d", offset, limit),
		"with_deleted_field": with_deleted,
	})

	whereStatemant, arg := utils.GetFilterString(len(params)+1, filters) //TODO: to usecase package without heavy depends
	if len(arg) > 0 {
		params = append(params, arg...)
	}

	if whereStatemant != "" {
		if with_deleted {
			whereStatemant = " WHERE " + whereStatemant
		} else {
			whereStatemant = " AND " + whereStatemant
		}
	}

	s.log.Info("Create where statement for reverse songs", map[string]any{
		"where":         whereStatemant,
		"length_params": len(params),
	})

	if with_deleted {
		query = fmt.Sprintf(queryTemplate, "", whereStatemant)
	} else {
		query = fmt.Sprintf(queryTemplate, delete, whereStatemant)
	}

	s.log.Info("selection query for reverse songs", map[string]any{
		"query":         query,
		"length_params": len(params),
	})
	s.log.Debug("selection query for reverse songs", map[string]any{
		"query":         query,
		"params":        params,
		"length_params": len(params),
	})
	err := s.db.SelectContext(ctx, &res, query, params...)
	if err != nil {
		s.log.Error("not success get reverse list songs", map[string]any{
			"err":    err,
			"query":  query,
			"params": params,
		})

		if errors.Is(err, sql.ErrNoRows) {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	s.log.Info("Get list songs success", map[string]any{
		"getting_songs":      len(res),
		"diaposon":           fmt.Sprintf("%d-%d", offset, limit),
		"with_deleted_field": with_deleted,
	})
	s.log.Debug("Get list songs success", map[string]any{
		"getting_songs":      len(res),
		"diaposon":           fmt.Sprintf("%d-%d", offset, limit),
		"with_deleted_field": with_deleted,
		"songs":              res,
	})

	return res, nil
}

func (s *songRepository) Total(ctx context.Context, with_deleted bool, filters *entity.SongFilters) (int, error) {
	var total int
	delete := "WHERE deleted_at IS NULL"
	query := ""
	queryTemplate := "SELECT COUNT(*) FROM songs s %s"

	s.log.Info("Geting total songs", map[string]any{
		"with_deleted_field": with_deleted,
	})
	s.log.Debug("Geting total songs", map[string]any{
		"with_deleted_field": with_deleted,
		"filters":            filters,
	})

	whereStatemant, arg := utils.GetFilterString(1, filters)
	if whereStatemant != "" {
		if with_deleted {
			whereStatemant = " WHERE " + whereStatemant
		} else {
			whereStatemant = " AND " + whereStatemant
		}
	}

	if with_deleted {
		query = fmt.Sprintf(queryTemplate, whereStatemant)
	} else {
		query = fmt.Sprintf(queryTemplate, delete+whereStatemant)
	}

	s.log.Info("query for getting total songs", map[string]any{
		"where":         whereStatemant,
		"query":         query,
		"length_params": len(arg),
	})
	s.log.Debug("query for getting total songs", map[string]any{
		"where":         whereStatemant,
		"query":         query,
		"length_params": len(arg),
		"params":        arg,
	})

	err := s.db.GetContext(ctx, &total, query, arg...)
	if err != nil {
		s.log.Error("query for getting total songs", map[string]any{
			"err":    err,
			"query":  query,
			"params": arg,
		})
		if err == sql.ErrNoRows {
			return total, entity.ErrNotFound
		}
		return total, err
	}

	s.log.Info("Get total songs", map[string]any{
		"with_deleted_field": with_deleted,
		"total":              total,
	})

	return total, nil
}

func (s *songRepository) GetSongTextByID(ctx context.Context, songID int64) (string, error) {
	var text string
	query := `SELECT m_text FROM songs WHERE id = $1 ` //AND deleted_at IS NULL` // TODO: with delete param

	s.log.Info("query for getting total songs", map[string]any{
		"query":   query,
		"song_id": songID,
	})

	err := s.db.GetContext(ctx, &text, query, songID)
	if err != nil {
		s.log.Error("query for getting total songs", map[string]any{
			"err":     err,
			"query":   query,
			"song_id": songID,
		})

		if err == sql.ErrNoRows {
			return "", entity.ErrNotFound
		}
		return "", err

	}

	s.log.Info("get text songs", map[string]any{
		"song_id":     songID,
		"text_length": len(text),
	})
	s.log.Debug("get text songs", map[string]any{
		"song_id":     songID,
		"text_length": len(text),
		"text":        text,
	})

	return text, nil
}

func (s *songRepository) GetByNames(ctx context.Context, song_name, group_name string) (*entity.Song, error) {
	return nil, nil
}

func (s *songRepository) GetByName(ctx context.Context, song_name string, group_id int64) (*entity.Song, error) {
	var err error
	song := &entity.Song{}
	query := `SELECT *
				FROM songs
				WHERE m_name Like $1 AND group_id = $2;`

	s.log.Info("get song by name and group_id", map[string]any{
		"group_id":  group_id,
		"song_name": song_name,
	})

	if s.tx != nil {
		s.log.Info("query with transaction get song by name and group_id", map[string]any{
			"query":     query,
			"group_id":  group_id,
			"song_name": song_name,
		})

		// exec sql
		err = s.tx.GetContext(ctx, song, query, song_name+"%", group_id)
	} else {
		s.log.Info("query without get song by name and group_id", map[string]any{
			"query":     query,
			"group_id":  group_id,
			"song_name": song_name,
		})

		// exec sql
		err = s.db.GetContext(ctx, song, query, song_name+"%", group_id)
	}

	if err != nil {
		s.log.Error("query without get song by name and group_id", map[string]any{
			"err":       err,
			"query":     query,
			"group_id":  group_id,
			"song_name": song_name,
		})

		if err == sql.ErrNoRows {
			return nil, entity.ErrNotFound
		}
		return nil, err
	}

	s.log.Info("get song by name and group_id", map[string]any{
		"group_id":  group_id,
		"song_name": song_name,
		"song_id":   song.ID,
	})
	s.log.Debug("get song by name and group_id", map[string]any{
		"group_id":  group_id,
		"song_name": song_name,
		"song":      song,
	})

	return song, nil
}

func (s *songRepository) GetByID(ctx context.Context, id int64) (*entity.Song, error) {
	var err error
	song := &entity.Song{}
	query := `SELECT *
				FROM songs
				WHERE m_name id = $1;`

	s.log.Info("get song by id", map[string]any{
		"song_id": id,
	})

	if s.tx != nil {
		s.log.Info("transaction get song by id", map[string]any{
			"song_id": id,
			"query":   query,
		})

		err = s.tx.GetContext(ctx, song, query, id)
	} else {
		s.log.Info("without transaction get song by id", map[string]any{
			"song_id": id,
			"query":   query,
		})

		err = s.db.GetContext(ctx, song, query, id)
	}

	if err != nil {
		s.log.Error("query without get song by name and group_id", map[string]any{
			"err":     err,
			"query":   query,
			"song_id": id,
		})

		if err == sql.ErrNoRows {
			return nil, entity.ErrNotFound
		}
		return nil, err

	}
	s.log.Info("get song by id", map[string]any{
		"song_id":   id,
		"song_name": song.Name,
	})
	s.log.Debug("get song by id", map[string]any{
		"query":   query,
		"song_id": id,
		"song":    song,
	})

	return song, nil
}

func (s *songRepository) UpdateFromMapByID(ctx context.Context, id int64, song *entity.Song, fields map[string]any) error {
	queryTemplate := `UPDATE songs
						SET %s
						WHERE id = $%d RETURNING id, m_name, m_link, m_text, m_release_date, group_id`

	s.log.Info("update song from map", map[string]any{
		"song_id":    id,
		"fields_num": len(fields),
	})
	s.log.Info("update song from map", map[string]any{
		"song_id":     id,
		"fields_num":  len(fields),
		"song_before": song,
	})

	updateString, params := s.getUpdateString(fields)
	params = append(params, id)

	s.log.Info("update song from map get update string", map[string]any{
		"song_id":       id,
		"update_string": updateString,
		"parms_length":  len(params) - 1,
	})
	s.log.Debug("update song from map get update string", map[string]any{
		"song_id":       id,
		"update_string": updateString,
		"parms_length":  len(params) - 1,
		"parms":         params,
	})

	query := fmt.Sprintf(queryTemplate, updateString, len(params))

	s.log.Info("update song from map get update string", map[string]any{
		"song_id":      id,
		"query":        query,
		"parms_length": len(params),
	})
	s.log.Debug("update song from map get update string", map[string]any{
		"song_id":      id,
		"query":        query,
		"parms_length": len(params),
		"parms":        params,
	})

	err := s.db.GetContext(ctx, song, query, params...)
	if err != nil {
		s.log.Error("update song from map", map[string]any{
			"err":     err,
			"query":   query,
			"song_id": id,
			"params":  params,
		})

		if err == sql.ErrNoRows {
			return entity.ErrNotFound
		}
		return err
	}

	s.log.Info("update song from map success", map[string]any{
		"song_id": id,
	})
	s.log.Debug("update song from map success", map[string]any{
		"song_id": id,
		"song":    song,
	})

	return nil
}

func (s *songRepository) getUpdateString(fields map[string]any) (string, []any) {
	fieldsUpdate := make([]string, 0)
	params := make([]any, 0)
	idPlaceholder := 1
	for k, v := range fields {
		fieldsUpdate = append(fieldsUpdate, fmt.Sprintf("%s = $%d", k, idPlaceholder))
		params = append(params, v)
		idPlaceholder++
	}
	return strings.Join(fieldsUpdate, ", "), params
}

func (s *songRepository) UpdateFromMapByNames(ctx context.Context, song_name, group_name string, fields map[string]string) (*entity.Song, error) {
	// TODO: Do implement

	return nil, nil
}

func (s *songRepository) DeleteSoftByID(ctx context.Context, id int64) error {
	query := `UPDATE songs
		SET deleted_at=CURRENT_TIMESTAMP
		WHERE id= $1;`

	s.log.Info("delete song soft", map[string]any{
		"song_id": id,
		"query":   query,
	})

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		s.log.Error("delete song soft", map[string]any{
			"err":     err,
			"query":   query,
			"song_id": id,
		})

		if err == sql.ErrNoRows {
			return entity.ErrNotFound
		}
		return err
	}

	s.log.Info("delete song soft success", map[string]any{
		"song_id": id,
		"query":   query,
	})

	return nil
}

func (s *songRepository) DeleteSoftByGroupID(ctx context.Context, id int64) error {
	// TODO: Do implement

	return nil
}

func (s *songRepository) DeleteSoftByNames(ctx context.Context, song_name, group_name string) error {
	// TODO: Do implement

	return nil
}

func (s *songRepository) DeleteForceByID(ctx context.Context, id int64) error {
	query := `DELETE FROM songs WHERE id = $1`

	s.log.Info("delete song force", map[string]any{
		"song_id": id,
		"query":   query,
	})

	_, err := s.db.ExecContext(ctx, query, id)
	if err != nil {
		s.log.Error("delete song force", map[string]any{
			"err":     err,
			"query":   query,
			"song_id": id,
		})

		if err == sql.ErrNoRows {
			return entity.ErrNotFound
		}
		return err
	}

	s.log.Info("delete song force success", map[string]any{
		"song_id": id,
		"query":   query,
	})

	return nil
}

func (s *songRepository) DeleteForceByGroupID(ctx context.Context, id int64) error {
	// TODO: Do implement

	return nil
}
