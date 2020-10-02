package xgorm

import (
	"github.com/Aoi-hosizora/ahlib/xproperty"
	"github.com/Aoi-hosizora/ahlib/xstatus"
	"gorm.io/gorm"
	"strings"
)

type Helper struct {
	db *gorm.DB
}

func WithDB(db *gorm.DB) *Helper {
	return &Helper{db: db}
}

func (h *Helper) Pagination(limit int32, page int32) *gorm.DB {
	return h.db.Limit(int(limit)).Offset(int((page - 1) * limit))
}

func IsMySQL(db *gorm.DB) bool {
	return db.Dialector.Name() == "mysql"
}

func CreateErr(rdb *gorm.DB) (xstatus.DbStatus, error) {
	if IsMySQL(rdb) && IsMySQLDuplicateEntryError(rdb.Error) {
		return xstatus.DbExisted, nil
	} else if rdb.Error != nil || rdb.RowsAffected == 0 {
		return xstatus.DbFailed, rdb.Error
	}

	return xstatus.DbSuccess, nil
}

func UpdateErr(rdb *gorm.DB) (xstatus.DbStatus, error) {
	if IsMySQL(rdb) && IsMySQLDuplicateEntryError(rdb.Error) {
		return xstatus.DbExisted, nil
	} else if rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	} else if rdb.RowsAffected == 0 {
		return xstatus.DbNotFound, nil
	}

	return xstatus.DbSuccess, nil
}

func DeleteErr(rdb *gorm.DB) (xstatus.DbStatus, error) {
	if rdb.Error != nil {
		return xstatus.DbFailed, rdb.Error
	} else if rdb.RowsAffected == 0 {
		return xstatus.DbNotFound, nil
	}

	return xstatus.DbSuccess, nil
}

func OrderByFunc(p xproperty.PropertyDict) func(source string) string {
	return func(source string) string {
		result := make([]string, 0)
		if source == "" {
			return ""
		}

		sources := strings.Split(source, ",")
		for _, src := range sources {
			if src == "" {
				continue
			}

			src = strings.TrimSpace(src)
			reverse := strings.HasSuffix(src, " desc") || strings.HasSuffix(src, " DESC")
			src = strings.Split(src, " ")[0]

			dest, ok := p[src]
			if !ok || dest == nil || len(dest.Destinations) == 0 {
				continue
			}

			if dest.Revert {
				reverse = !reverse
			}
			for _, prop := range dest.Destinations {
				if !reverse {
					prop += " ASC"
				} else {
					prop += " DESC"
				}
				result = append(result, prop)
			}
		}

		return strings.Join(result, ", ")
	}
}