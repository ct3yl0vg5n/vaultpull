package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"vaultpull/internal/secret"
)

var generateCmd = &cobra.Command{
	Use:   "generate",
	Short: "Generate a random secret value",
	RunE:  runGenerate,
}

func init() {
	generateCmd.Flags().IntP("length", "l", 32, "length of the generated secret")
	generateCmd.Flags().StringP("charset", "c", secret.Alphanumeric, "charset to use (alphanumeric or special)")
	generateCmd.Flags().Bool("base64", false, "generate base64-encoded random bytes")
	generateCmd.Flags().StringP("key", "k", "", "if set, output as KEY=VALUE env line")
	RootCmd.AddCommand(generateCmd)
}

func runGenerate(cmd *cobra.Command, _ []string) error {
	length, _ := cmd.Flags().GetInt("length")
	charset, _ := cmd.Flags().GetString("charset")
	b64, _ := cmd.Flags().GetBool("base64")
	key, _ := cmd.Flags().GetString("key")

	opts := secret.GenerateOptions{
		Length:  length,
		Charset: charset,
		Base64:  b64,
	}

	v, err := secret.Generate(opts)
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v\n", err)
		return err
	}

	if key != "" {
		fmt.Printf("%s=%s\n", key, v)
	} else {
		fmt.Println(v)
	}
	return nil
}
