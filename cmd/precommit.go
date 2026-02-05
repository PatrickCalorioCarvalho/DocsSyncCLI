package cmd

import (
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/spf13/cobra"

	"github.com/PatrickCalorioCarvalho/DocsSyncCLI/config"
	"github.com/PatrickCalorioCarvalho/DocsSyncCLI/scanner"
)

var precommitCmd = &cobra.Command{
	Use:   "precommit",
	Short: "Gera a estrutura final de documentação centralizada",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectRoot, _ := cmd.Flags().GetString("path")

		cfg, err := config.LoadConfig(projectRoot)
		if err != nil {
			fmt.Println("❌", err)
			fmt.Println("👉 Crie um arquivo docssync.yaml na raiz do projeto.")
			os.Exit(1)
		}

		if len(cfg.Project.Key) == 0 {
			return fmt.Errorf("chave do projeto não definida em 'project.key' no docssync.yaml")
		}

		files, err := scanner.Scan(cfg, projectRoot)
		if err != nil {
			return err
		}

		if len(files) == 0 {
			fmt.Println("⚠ Nenhum arquivo Markdown encontrado")
			return nil
		}

		workingDir, _ := os.Getwd()

		baseDir := ".precommit"
		if cfg.Precommit.BaseDir != "" {
			baseDir = cfg.Precommit.BaseDir
		}

		precommitDir := filepath.Join(
			workingDir,
			baseDir,
			cfg.Project.Key,
		)

		_ = os.RemoveAll(precommitDir)

		if err := os.MkdirAll(precommitDir, 0755); err != nil {
			return err
		}

		for _, file := range files {
			relPath, err := filepath.Rel(projectRoot, file)
			if err != nil {
				continue
			}

			relPath = filepath.ToSlash(relPath)

			cleanPath := stripDirs(relPath, cfg.Precommit.StripDirs)

			destPath := filepath.Join(precommitDir, cleanPath)

			if err := os.MkdirAll(filepath.Dir(destPath), 0755); err != nil {
				return err
			}

			if err := copyFile(file, destPath); err != nil {
				return err
			}
		}

		fmt.Printf("✔ precommit gerado em: %s\n", precommitDir)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(precommitCmd)
}

func stripDirs(path string, dirs []string) string {
	parts := strings.Split(path, "/")
	var cleaned []string

	for _, part := range parts {
		shouldStrip := false
		for _, dir := range dirs {
			if part == dir {
				shouldStrip = true
				break
			}
		}
		if !shouldStrip {
			cleaned = append(cleaned, part)
		}
	}

	return strings.Join(cleaned, "/")
}

func copyFile(src, dst string) error {
	in, err := os.Open(src)
	if err != nil {
		return err
	}
	defer in.Close()

	out, err := os.Create(dst)
	if err != nil {
		return err
	}
	defer out.Close()

	if _, err := io.Copy(out, in); err != nil {
		return err
	}

	return out.Sync()
}
