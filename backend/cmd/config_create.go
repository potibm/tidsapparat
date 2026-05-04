package cmd

import (
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"log/slog"

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
			const defaultAPIKeyLength = 32

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

func generateSecureToken(byteLength int) (string, error) {
	bytes := make([]byte, byteLength)
	if _, err := rand.Read(bytes); err != nil {
		return "", fmt.Errorf("error while generating the token: %w", err)
	}

	return hex.EncodeToString(bytes), nil
}
