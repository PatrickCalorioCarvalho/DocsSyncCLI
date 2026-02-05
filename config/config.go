package config

type Config struct {
	Project   ProjectConfig   `yaml:"project"`
	Scan      ScanConfig      `yaml:"scan"`
	Precommit PrecommitConfig `yaml:"precommit"`
	Sync      SyncConfig      `yaml:"sync"`
}

type ProjectConfig struct {
	Key string `yaml:"key"`
}

type ScanConfig struct {
	Root    string   `yaml:"root"`
	Include []string `yaml:"include"`
	Exclude []string `yaml:"exclude"`
}

type PrecommitConfig struct {
	BaseDir   string   `yaml:"baseDir"`
	StripDirs []string `yaml:"stripDirs"`
}

type SyncConfig struct {
	Docsaurus DocsaurusConfig `yaml:"docsaurus"`
	OpenWebUI OpenWebUIConfig `yaml:"openwebui"`
}

type DocsaurusConfig struct {
	Enabled bool	 `yaml:"enabled"`
	GitlabRepo string `yaml:"gitlabRepo"`
	GitlabToken string `yaml:"gitlabToken"`
	GitlabBranch string `yaml:"gitlabBranch"`
}

type OpenWebUIConfig struct {
	Enabled   bool   `yaml:"enabled"`
	ApiUrl    string `yaml:"apiUrl"`
	ApiKey	string `yaml:"apiKey"`
	KnowledgeId string `yaml:"knowledgeId"`
}