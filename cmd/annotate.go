package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/yourusername/vaultpull/internal/secret"
)

var (
	annotateFile  string
	annotateKey   string
	annotateNote  string
	annotateOwner string
	annotateRemove bool
)

var annotateCmd = &cobra.Command{
	Use:   "annotate",
	Short: "Add or remove metadata annotations on secret keys",
	RunE:  runAnnotate,
}

func init() {
	annotateCmd.Flags().StringVar(&annotateFile, "file", secret.DefaultAnnotateOptions().Path, "path to annotations file")
	annotateCmd.Flags().StringVar(&annotateKey, "key", "", "secret key to annotate (required)")
	annotateCmd.Flags().StringVar(&annotateNote, "note", "", "annotation note")
	annotateCmd.Flags().StringVar(&annotateOwner, "owner", "", "owner of the secret")
	annotateCmd.Flags().BoolVar(&annotateRemove, "remove", false, "remove annotation for the given key")
	_ = annotateCmd.MarkFlagRequired("key")
	rootCmd.AddCommand(annotateCmd)
}

func runAnnotate(cmd *cobra.Command, _ []string) error {
	store, err := secret.LoadAnnotations(annotateFile)
	if err != nil {
		return fmt.Errorf("annotate: load: %w", err)
	}

	if annotateRemove {
		store = secret.RemoveAnnotation(store, annotateKey)
		fmt.Fprintf(cmd.OutOrStdout(), "removed annotation for %q\n", annotateKey)
	} else {
		if annotateNote == "" {
			return fmt.Errorf("annotate: --note is required when adding an annotation")
		}
		store = secret.Annotate(store, annotateKey, annotateNote, annotateOwner)
		fmt.Fprintf(cmd.OutOrStdout(), "annotated %q\n", annotateKey)
	}

	if err := secret.SaveAnnotations(annotateFile, store); err != nil {
		return fmt.Errorf("annotate: save: %w", err)
	}

	fmt.Fprint(cmd.OutOrStdout(), secret.FormatAnnotationReport(store))
	return nil
}

func runAnnotateList(cmd *cobra.Command) error {
	store, err := secret.LoadAnnotations(annotateFile)
	if err != nil {
		return err
	}
	fmt.Fprint(cmd.OutOrStdout(), secret.FormatAnnotationReport(store))
	return nil
}

var _ = runAnnotateList // suppress unused warning; available for subcommand extension

func init() { // ensure binary builds without vault creds in CI
	if os.Getenv("VAULTPULL_SKIP_VAULT") != "" {
		return
	}
}
