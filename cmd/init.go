package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

const initTemplate = `# Configuração do DocsSync
# Documentação: https://github.com/PatrickCalorioCarvalho/DocsSyncCLI

project:
  # Identificador único do projeto (usado como subpasta nos destinos de sync)
  key: ProjectID

scan:
  # Pasta raiz a partir de onde os arquivos serão escaneados
  root: .
  # Padrões (glob) de arquivos a incluir
  include:
    - "**/*.md"
  # Padrões (glob) de arquivos/pastas a excluir
  exclude:
    - "**/node_modules/**"
    - "**/dist/**"
    - "**/.git/**"
    - "**/README.md"

precommit:
  # Pasta local onde o staging é montado antes do sync
  baseDir: .precommit
  # Nomes de pastas a remover do caminho relativo ao copiar para o staging
  stripDirs:
    - docs

sync:
  docsaurus:
    enabled: false
    repoUrl: https://github.com/sua-org/seu-docsaurus.git
    # Use uma variável de ambiente (ex: secret do CI) em vez de um token em texto puro
    repoToken: ${DOCS_REPO_TOKEN}
    repoBranch: main
    docsPath: docs

  openwebui:
    enabled: false
    apiUrl: https://api.openwebui.com/ingest
    # Use uma variável de ambiente (ex: secret do CI) em vez de uma chave em texto puro
    apiKey: ${OPENWEBUI_API_KEY}
    knowledgeId: your-collection-name
`

var initForce bool

var initCmd = &cobra.Command{
	Use:   "init",
	Short: "Gera um docssync.yaml modelo na raiz do projeto",
	RunE: func(cmd *cobra.Command, args []string) error {
		projectRoot, _ := cmd.Flags().GetString("path")

		configPath := filepath.Join(projectRoot, "docssync.yaml")

		if _, err := os.Stat(configPath); err == nil && !initForce {
			return fmt.Errorf("docssync.yaml já existe em %s (use --force para sobrescrever)", configPath)
		}

		if err := os.WriteFile(configPath, []byte(initTemplate), 0644); err != nil {
			return err
		}

		fmt.Println("✔ docssync.yaml gerado em:", configPath)
		fmt.Println("👉 Edite os campos e habilite os destinos de sync que quiser usar.")
		return nil
	},
}

func init() {
	initCmd.Flags().BoolVar(&initForce, "force", false, "Sobrescreve o docssync.yaml existente")
	rootCmd.AddCommand(initCmd)
}
