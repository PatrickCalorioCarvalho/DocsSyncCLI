# 🚀 DocsSyncCLI

> CLI para padronização, versionamento e sincronização automática de
> documentação entre projetos e um repositório central (Docsaurus), com
> suporte futuro a integração com bases de conhecimento (RAG /
> OpenWebUI).

------------------------------------------------------------------------

## 📌 Problema que ele resolve

Em ambientes com múltiplos projetos, a documentação costuma:

-   Ficar espalhada
-   Perder padrão
-   Ficar desatualizada
-   Não estar centralizada
-   Não estar preparada para integração com LLM / RAG

O **DocsSyncCLI** resolve isso aplicando o conceito de:

> 📖 Documentation as Code + Sync Automatizado

------------------------------------------------------------------------

## 🎯 Objetivo

-   Padronizar documentação Markdown
-   Gerar estrutura limpa via `precommit`
-   Sincronizar automaticamente com um repositório central Docsaurus
-   Permitir integração futura com base de conhecimento (OpenWebUI)
-   Funcionar com GitHub ou GitLab (via git CLI)
-   Ser executável localmente ou em CI/CD (inclusive como GitHub Action)

------------------------------------------------------------------------

## 🧠 Conceito de funcionamento

Fluxo simplificado:

Projeto │ ├── Markdown espalhado │ └── DocsSyncCLI │ ├── Precommit
(estrutura limpa) └── Sync └── Repo central Docsaurus

------------------------------------------------------------------------

## ⚙️ Instalação local

Baixe o binário na [página de releases](https://github.com/PatrickCalorioCarvalho/DocsSyncCLI/releases)
ou instale via Go:

```bash
go install github.com/PatrickCalorioCarvalho/DocsSyncCLI@latest
```

## ⚙️ Configuração

Gere um `docssync.yaml` modelo com:

```bash
docssync init
```

Ou crie manualmente na raiz do projeto:

``` yaml
project:
  key: ProjectID

scan:
  root: .
  include:
    - "**/*.md"
  exclude:
    - "**/node_modules/**"
    - "**/dist/**"
    - "**/.git/**"
    - "**/README.md"

precommit:
  baseDir: .precommit
  stripDirs:
    - Documentacao
    - docs

sync:
  docsaurus:
    enabled: true
    repoUrl: https://github.com/sua-org/seu-docsaurus.git
    repoToken: ${DOCS_REPO_TOKEN}
    repoBranch: main
    docsPath: documentation/docs

  openwebui:
    enabled: false
    apiUrl: https://api.openwebui.com/ingest
    apiKey: ${OPENWEBUI_API_KEY}
    knowledgeId: your-collection-name
```

Qualquer campo pode referenciar uma variável de ambiente com `${NOME_DA_VAR}` —
útil para não versionar tokens/chaves em texto puro. O valor é expandido a partir
do ambiente do processo no momento em que o `docssync.yaml` é lido.

------------------------------------------------------------------------

## 📂 Etapas do Processo

### 1️⃣ Scan

-   Localiza arquivos `.md`
-   Aplica filtros `include` / `exclude`

------------------------------------------------------------------------

### 2️⃣ Precommit

-   Gera estrutura limpa em:

.precommit/`<ProjectKey>`{=html}/

-   Remove diretórios definidos em `stripDirs`
-   Prepara documentação pronta para publicação

------------------------------------------------------------------------

### 3️⃣ Sync Docsaurus

Ao executar:

docssync commit --path .

O CLI:

1.  Clona ou atualiza o repositório Docsaurus
2.  Vai para a branch configurada
3.  Remove: `<docsPath>`{=html}/`<ProjectKey>`{=html}
4.  Copia conteúdo do `.precommit`
5.  Realiza commit automático
6.  Faz push

Mensagem de commit gerada automaticamente:

docsSync: `<token>`{=html} `<ProjectKey>`{=html} 202602052022

------------------------------------------------------------------------

## 🔐 Segurança

-   Autenticação via Personal Access Token
-   Compatível com:
    -   GitHub (token direto na URL: `https://<token>@github.com/...`)
    -   GitLab (formato `oauth2:<token>@...`)
-   Suporte a `${VAR}` no `docssync.yaml` para injetar tokens via variável de
    ambiente/secret de CI, em vez de texto puro no arquivo
-   Não depende de API REST específica
-   Usa git CLI (mais robusto e universal)

------------------------------------------------------------------------

## 🏗️ Estrutura do Projeto

DocsSyncCLI/ ├── config/ ├── scanner/ ├── sync/ │ ├── docsaurus.go │ └──
openwebui.go ├── cmd/ └── main.go

------------------------------------------------------------------------

## 🚀 Execução

### Rodar manualmente

```bash
go run . init --path .
go run . precommit --path .
go run . commit --path .
```

### Build binário

```bash
go build -o docssync .
./docssync commit --path .
```

------------------------------------------------------------------------

## 🤖 Uso como GitHub Action

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

O único input é `path` (padrão `.`), a pasta onde está o `docssync.yaml` do
repositório. Quaisquer secrets referenciados no `docssync.yaml` via `${VAR}`
devem ser expostos ao step via `env:`.

**Limitações da versão atual:** só há binário publicado para Linux e Windows
(amd64); a Action usa o binário Linux e portanto só roda em runners `ubuntu-*`.
Ela também espera que o `actions/checkout` já tenha rodado antes dela no job.

------------------------------------------------------------------------

## 🌍 Compatibilidade

-   Windows
-   Linux
-   macOS (CLI local; ainda sem binário de release nem suporte na Action)
-   GitHub
-   GitLab
-   Execução local ou CI/CD

------------------------------------------------------------------------

## 📈 Benefícios

✔ Centralização de documentação\
✔ Versionamento real\
✔ Padronização entre projetos\
✔ Automação total\
✔ Preparado para LLM / RAG\
✔ Independente de plataforma Git

------------------------------------------------------------------------

## 🔮 Evoluções Futuras

-   Integração com OpenWebUI (RAG)
-   Validação de documentação (modo strict)
-   Lint para imagens sem descrição
-   Publicação da Action no GitHub Marketplace
-   Build de binário macOS/arm64
-   Docker execution mode
-   Sincronização automática de base de conhecimento

------------------------------------------------------------------------

## 📄 Licença

Definir conforme necessidade do projeto.
