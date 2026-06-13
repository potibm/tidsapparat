package cmd

import (
	"fmt"
	"log/slog"

	"github.com/potibm/tidsapparat/internal/app/config"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	configCreateForce    bool
	configCreateFilename string
)

func NewConfigCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create a new configuration file with default values",
		Annotations: map[string]string{
			skipConfigValidationAnnotation: "true",
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			const defaultFrontendURL = "https://localhost:3200"

			viper.SetDefault("app.frontend_url", defaultFrontendURL)
			viper.SetDefault("app.cors_allow_origins", []string{defaultFrontendURL})

			viper.SetDefault("s3_client.access_key_id", "accesskey")
			viper.SetDefault("s3_client.secret_access_key", "secretkey")
			viper.SetDefault("s3_client.region", "us-east-1")
			viper.SetDefault("s3_client.endpoint", "http://localhost:9000")
			viper.SetDefault("s3_client.use_path_style", true)

			viper.SetDefault("party.default_address", "Prins Jørgens Gård 5, 1218 København K, Denmark")

			viper.SetDefault("exporter", []config.ExporterConfig{
				{
					Name:        "ical_to_file_exporter",
					Type:        "ical",
					Destination: "file",
					Filename:    "calendar",
					Options: map[string]string{
						"dir": "./exports",
					},
					Enabled: false,
				}, {
					Name:        "ical_to_s3_exporter",
					Type:        "ical",
					Destination: "s3",
					Filename:    "calendar",
					Options: map[string]string{
						"bucket": "my-bucket",
					},
					Enabled: false,
				},
			})

			viper.SetDefault("event_durations", []int{0, 15, 30, 60, 90, 120})

			viper.SetDefault("auth", config.AuthConfig{
				Type:      "oidc",
				Name:      "Dex",
				Authority: "https://dex.tidsapparat.test/dex",
				ClientID:  "react-admin-client",
			})

			filename := configCreateFilename

			var writeErr error
			if configCreateForce {
				writeErr = viper.WriteConfigAs(filename)
			} else {
				writeErr = viper.SafeWriteConfigAs(filename)
			}

			if writeErr != nil {
				if _, ok := writeErr.(viper.ConfigFileNotFoundError); !ok {
					return fmt.Errorf(
						"file %s already exists or was not able to be created: %w",
						filename,
						writeErr,
					)
				}

				return fmt.Errorf("error writing the config: %w", writeErr)
			}

			slog.Info("Configuration file created successfully", "filename", filename)

			return nil
		},
	}

	cmd.Flags().BoolVarP(&configCreateForce, "force", "f", false, "Overwrite existing config file if it already exists")
	cmd.Flags().
		StringVarP(&configCreateFilename, "output", "o", "config/config.yaml", "Filename for the generated config file")

	return cmd
}
