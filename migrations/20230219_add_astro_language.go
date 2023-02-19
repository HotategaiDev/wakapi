package migrations

import (
	"github.com/emvi/logbuch"
	"github.com/muety/wakapi/config"
	"gorm.io/gorm"
)

func init() {
	const name = "20230219-add_astro_language"
	f := migrationFunc{
		name: name,
		f: func(db *gorm.DB, cfg *config.Config) error {
			if hasRun(name, db) {
				return nil
			}

			logbuch.Info("running migration '%s'", name)

			if err := db.Exec("UPDATE heartbeats SET language = 'Astro' where language = '' and entity like '%.astro'").Error; err != nil {
				return err
			}

			setHasRun(name, db)
			return nil
		},
	}

	registerPostMigration(f)
}
