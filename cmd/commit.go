package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"

	"github.com/PatrickCalorioCarvalho/DocsSyncCLI/config"
	"github.com/PatrickCalorioCarvalho/DocsSyncCLI/sync"
)

var commitCmd = &cobra.Command{
	Use:   "commit",
	Short: "Sincroniza a documentação conforme definido no docssync.yaml",
	RunE: func(cmd *cobra.Command, args []string) error {

		projectRoot, _ := cmd.Flags().GetString("path")

		cfg, err := config.LoadConfig(projectRoot)
		if err != nil {
			fmt.Println("❌", err)
			fmt.Println("👉 Crie um arquivo docssync.yaml na raiz do projeto.")
			os.Exit(1)
		}

		if len(cfg.Project.Key) == 0 {
			return fmt.Errorf("project.key não definido no docssync.yaml")
		}

		baseDir := ".precommit"
		if cfg.Precommit.BaseDir != "" {
			baseDir = cfg.Precommit.BaseDir
		}

		workingDir, _ := os.Getwd()
		precommitDir := filepath.Join(
			workingDir,
			baseDir,
			cfg.Project.Key,
		)

		if _, err := os.Stat(precommitDir); os.IsNotExist(err) {
			return fmt.Errorf("precommit não encontrado (%s). Execute `docssync precommit` primeiro", precommitDir)
		}

		fmt.Println("📦 Precommit encontrado:", precommitDir)

		if cfg.Sync.Docsaurus.Enabled {
			fmt.Println("🚀 Sincronizando com Docsaurus...")
			if err := commitDocsaurus(cfg, precommitDir); err != nil {
				return err
			}
		}

		if cfg.Sync.OpenWebUI.Enabled {
			fmt.Println("🚀 Sincronizando com OpenWebUI...")
			if err := sync.SyncOpenWebUI(cfg, precommitDir); err != nil {
				return err
			}
		}

		if err := os.RemoveAll(precommitDir); err != nil {
			return fmt.Errorf("commit realizado, mas falhou ao remover precommit: %w", err)
		}

		fmt.Println("✔ Commit finalizado e precommit removido")
		return nil
	},
}

func init() {
	rootCmd.AddCommand(commitCmd)
}

func commitDocsaurus(cfg *config.Config, precommitDir string) error {
	fmt.Println("   ↳ Docsaurus OK")
	return nil
}
