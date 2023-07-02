package root

import (
	"github.com/spf13/cobra"

	"github.com/peterldowns/pgmigrate/cli/shared"
)

var debugCmd = &cobra.Command{
	Use:     "debug",
	Short:   "parse and log the current configuration",
	GroupID: "dev",
	RunE: func(_ *cobra.Command, _ []string) error {
		logger, _ := shared.State.Logger()
		configfile := shared.State.Configfile()

		logger.Info(configfile.Name(), "is_set", configfile.IsSet(), "value", configfile.Value())

		shared.State.Parse()

		database := shared.State.Database()
		logformat := shared.State.LogFormat()
		migrations := shared.State.Migrations()

		logger.Info(migrations.Name(), "is_set", migrations.IsSet(), "value", migrations.Value())
		logger.Info(database.Name(), "is_set", database.IsSet(), "value", database.Value())
		logger.Info(logformat.Name(), "is_set", logformat.IsSet(), "value", logformat.Value())

		return nil
	},
}
