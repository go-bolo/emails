package migrations_email

import (
	"fmt"

	"github.com/go-bolo/bolo"
	"gorm.io/gorm"
)

func GetInitMigration() *bolo.Migration {
	queries := []struct {
		table string
		up    string
		down  string
	}{
		{
			table: "emails",
			up: `CREATE TABLE IF NOT EXISTS emails (
				id int NOT NULL AUTO_INCREMENT,
				emailId text,
				` + "`from`" + ` text NOT NULL,
				` + "`to`" + ` text,
				cc text,
				bcc text,
				replyTo text,
				inReplyTo text,
				variables text,
				subject text,
				` + "`text`" + `  text,
				html text,
				` + "`type`" + ` varchar(255) DEFAULT NULL,
				status varchar(255) NOT NULL DEFAULT 'added',
				createdAt datetime NOT NULL,
				updatedAt datetime NOT NULL,
				PRIMARY KEY (id)
			)`,
		},
		{
			table: "email_templates",
			up: `CREATE TABLE IF NOT EXISTS email_templates (
				id int NOT NULL AUTO_INCREMENT,
				subject text,
				` + "`text`" + ` text,
				css text,
				html text,
				type varchar(255) DEFAULT NULL,
				createdAt datetime NOT NULL,
				updatedAt datetime NOT NULL,
				PRIMARY KEY (id)
			)`,
		},
	}

	return &bolo.Migration{
		Name: "init",
		Up: func(app bolo.App) error {
			db := app.GetDB()
			return db.Transaction(func(tx *gorm.DB) error {
				for _, q := range queries {
					err := tx.Exec(q.up).Error
					if err != nil {
						return fmt.Errorf("failed to create "+q.table+" table: %w", err)
					}
				}

				return nil
			})
		},
		Down: func(app bolo.App) error {
			return nil
		},
	}
}
