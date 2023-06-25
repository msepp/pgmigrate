package root

import (
	"database/sql"
	"os"

	"github.com/spf13/cobra"

	"github.com/peterldowns/pgmigrate"
	"github.com/peterldowns/pgmigrate/cli/shared"
)

var verifyCmd = &cobra.Command{ //nolint:gochecknoglobals
	Use:   "verify",
	Short: "Verify that migrations have been applied correctly",
	RunE: func(cmd *cobra.Command, args []string) error {
		database := shared.GetDatabase()
		migrations := shared.GetMigrations()
		if err := shared.Validate(database, migrations); err != nil {
			return err
		}

		slogger, mlogger := shared.NewLogger()
		dir := os.DirFS(migrations.Value())
		db, err := sql.Open("pgx", database.Value())
		if err != nil {
			return err
		}
		defer db.Close()

		verrs, err := pgmigrate.Verify(cmd.Context(), db, dir, mlogger)
		if err != nil {
			return err
		}
		for _, verr := range verrs {
			var attrs []any
			for key, val := range verr.Fields {
				attrs = append(attrs, key, val)
			}
			slogger.With(attrs...).Warn(verr.Message)
		}
		return nil
	},
}

//nolint:gochecknoinits
func init() {
	Command.AddCommand(verifyCmd)
}
