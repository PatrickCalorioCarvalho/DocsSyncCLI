# DocsSyncCLI

CLI em Go para sincronizar documentação Markdown de um projeto com destinos como
[Docsaurus](https://docusaurus.io/) e [OpenWebUI](https://openwebui.com/).

## Instalação local

Baixe o binário na [página de releases](https://github.com/PatrickCalorioCarvalho/DocsSyncCLI/releases)
ou instale via Go:

```bash
go install github.com/PatrickCalorioCarvalho/DocsSyncCLI@latest
```

## Uso

```bash
# gera um docssync.yaml modelo na raiz do projeto
docssync init

# escaneia os markdowns e monta a pasta de staging (.precommit/<key>)
docssync precommit

# sincroniza a pasta de staging com os destinos habilitados e limpa o staging
docssync commit
```

Todos os comandos aceitam `-p/--path` para apontar para uma raiz de projeto diferente
do diretório atual.

### `docssync.yaml`

```yaml
project:
  key: ProjectID

scan:
  root: .
  include:
    - "**/*.md"
  exclude:
    - "**/node_modules/**"

precommit:
  baseDir: .precommit
  stripDirs:
    - docs

sync:
  docsaurus:
    enabled: true
    repoUrl: https://github.com/sua-org/seu-docsaurus.git
    repoToken: ${DOCS_REPO_TOKEN}
    repoBranch: main
    docsPath: docs

  openwebui:
    enabled: false
    apiUrl: https://api.openwebui.com/ingest
    apiKey: ${OPENWEBUI_API_KEY}
    knowledgeId: your-collection-name
```

Qualquer campo pode referenciar uma variável de ambiente com `${NOME_DA_VAR}` — útil
para não versionar tokens/chaves em texto puro. O valor é expandido a partir do
ambiente do processo no momento em que o `docssync.yaml` é lido.

## Uso como GitHub Action

```yaml
name: Sync docs

on:
  push:
    branches: [main]

jobs:
  docssync:
    runs-on: ubuntu-latest
    steps:
      - uses: actions/checkout@v4

      - uses: PatrickCalorioCarvalho/DocsSyncCLI@v0
        env:
          DOCS_REPO_TOKEN: ${{ secrets.DOCS_REPO_TOKEN }}
        with:
          path: .
```

O único input é `path` (padrão `.`), a pasta onde está o `docssync.yaml` do repositório.
Quaisquer secrets referenciados no `docssync.yaml` via `${VAR}` devem ser expostos ao
step via `env:`.

### Limitações da versão atual

- Só há binário publicado para Linux e Windows (amd64); a Action usa o binário Linux e
  portanto só roda em runners `ubuntu-*`.
- A Action espera que o `actions/checkout` já tenha rodado antes dela no job.

## Contribuindo

Releases são gerados automaticamente ao fazer merge de PR em `main`
([.github/workflows/pr-to-release.yml](.github/workflows/pr-to-release.yml)), com bump
de patch (`vX.Y.Z`) e atualização da tag flutuante `v0` usada pela Action.
