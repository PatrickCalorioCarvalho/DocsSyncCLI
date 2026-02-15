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
-   Ser executável localmente ou em CI/CD

------------------------------------------------------------------------

## 🧠 Conceito de funcionamento

Fluxo simplificado:

Projeto │ ├── Markdown espalhado │ └── DocsSyncCLI │ ├── Precommit
(estrutura limpa) └── Sync └── Repo central Docsaurus

------------------------------------------------------------------------

## ⚙️ Configuração

Arquivo `docssync.yaml`:

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
    repoUrl: https://gitlab.com/org/docsaurus.git
    repoToken: your-token
    repoBranch: main
    docsPath: documentation/docs
```

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
    -   GitHub
    -   GitLab
-   Não depende de API REST específica
-   Usa git CLI (mais robusto e universal)

------------------------------------------------------------------------

## 🏗️ Estrutura do Projeto

DocsSyncCLI/ ├── config/ ├── sync/ │ ├── docsaurus.go │ └── git.go ├──
cmd/ └── main.go

------------------------------------------------------------------------

## 🚀 Execução

### Rodar manualmente

go run . commit --path .

### Build binário

go build -o docssync ./docssync commit --path .

------------------------------------------------------------------------

## 🌍 Compatibilidade

-   Windows
-   Linux
-   macOS
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
-   Execução oficial como GitHub Action
-   Docker execution mode
-   Sincronização automática de base de conhecimento

------------------------------------------------------------------------

## 📄 Licença

Definir conforme necessidade do projeto.
