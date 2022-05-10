package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/cobra/doc"
)

func main() {
	var rootCmd = &cobra.Command{
		Use:   "petracli",
		Short: "private terraform registry cli",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
	}
	err := doc.GenMarkdownTree(rootCmd, "./")
	if err != nil {
		log.Fatal(err)
	}
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

// Declare Flags.
var (
	flagGCSBucket       string
	flagModuleDirectory string
)

func init() {
	// Flags
	rootCmd.PersistentFlags().StringVar(&flagGCSBucket, "gcs-bucket", "", "Name of the Google Cloud Storage bucket you want to use for storage (required) e.g.: my-bucket")
	rootCmd.PersistentFlags().StringVar(&flagModuleDirectory, "module-directory", "", "Directory of the module you want to upload (required) e.g.: ./modules-example/rabbitmq/")

	// Making Flags required
	err := rootCmd.MarkPersistentFlagRequired("gcs-bucket")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
	err = rootCmd.MarkPersistentFlagRequired("module-directory")
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
	}
}
