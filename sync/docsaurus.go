package sync

import (
	"bytes"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"strings"
	"time"

	"github.com/PatrickCalorioCarvalho/DocsSyncCLI/config"
)

func SyncDocsaurus(cfg *config.Config, precommitDir string) error {

	d := cfg.Sync.Docsaurus
	projectKey := cfg.Project.Key

	if d.RepoUrl == "" || d.RepoToken == "" || d.RepoBranch == "" {
		return fmt.Errorf("docsaurus.repoUrl, repoToken ou repoBranch não configurados")
	}

	repoPath := filepath.Join(".docssync")

	fmt.Println("📚 Docsaurus repo:", repoPath)

	if err := ensureRepo(repoPath, d.RepoUrl, d.RepoBranch, d.RepoToken); err != nil {
		return err
	}

	if err := ensureGitIdentity(repoPath); err != nil {
		return err
	}

	docsPath := d.DocsPath
	if docsPath == "" {
		docsPath = "docs"
	}

	docsProjectPath := filepath.Join(repoPath, docsPath, projectKey)

	fmt.Println("🧹 Limpando " + docsPath + "/" + projectKey)
	_ = os.RemoveAll(docsProjectPath)

	fmt.Println("📂 Copiando precommit → docsaurus")
	if err := copyDir(precommitDir, docsProjectPath); err != nil {
		return err
	}

	if err := git(repoPath, "add", docsPath+"/"+projectKey); err != nil {
		return err
	}

	now := time.Now().Format("200601021504")
	commitMsg := fmt.Sprintf("docsSync: %s %s", now, projectKey)

	if err := gitCommitIfNeeded(repoPath, commitMsg); err != nil {
		return err
	}

	if err := git(repoPath, "push", "origin", d.RepoBranch); err != nil {
		return err
	}

	fmt.Println("✔ Docsaurus sincronizado com sucesso")
	return nil
}

func git(dir string, args ...string) error {
	cmd := exec.Command("git", args...)
	cmd.Dir = dir

	var out bytes.Buffer
	cmd.Stdout = &out
	cmd.Stderr = &out

	if err := cmd.Run(); err != nil {
		return fmt.Errorf("git %v: %s", args, out.String())
	}

	return nil
}

func ensureRepo(path, repo, branch, token string) error {

	if _, err := os.Stat(path); os.IsNotExist(err) {

		fmt.Println("📥 Clonando repositório Docsaurus")

		repoAuth := injectToken(repo, token)

		cmd := exec.Command(
			"git", "clone", "-b", branch, repoAuth, path,
		)

		out, err := cmd.CombinedOutput()
		if err != nil {
			return fmt.Errorf("git clone: %s", string(out))
		}

		return nil
	}

	fmt.Println("🔄 Atualizando repositório Docsaurus")

	if err := git(path, "fetch"); err != nil {
		return err
	}
	if err := git(path, "checkout", branch); err != nil {
		return err
	}
	if err := git(path, "pull"); err != nil {
		return err
	}

	return nil
}

func ensureGitIdentity(repoPath string) error {
	if err := git(repoPath, "config", "user.name", "DocsSync Bot"); err != nil {
		return err
	}
	return git(repoPath, "config", "user.email", "docssync@users.noreply.github.com")
}

func injectToken(repo, token string) string {

	if strings.Contains(repo, "github.com") {
		return strings.Replace(repo, "https://", "https://"+token+"@", 1)
	}

	if strings.Contains(repo, "https://") {
		return strings.Replace(repo, "https://", "https://oauth2:"+token+"@", 1)
	}

	if strings.Contains(repo, "http://") {
		return strings.Replace(repo, "http://", "http://oauth2:"+token+"@", 1)
	}

	return repo
}

func copyDir(src, dst string) error {

	return filepath.Walk(src, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		rel, _ := filepath.Rel(src, path)
		target := filepath.Join(dst, rel)

		if info.IsDir() {
			return os.MkdirAll(target, info.Mode())
		}

		data, err := os.ReadFile(path)
		if err != nil {
			return err
		}

		return os.WriteFile(target, data, info.Mode())
	})
}

func gitCommitIfNeeded(dir, message string) error {

	cmd := exec.Command("git", "status", "--porcelain")
	cmd.Dir = dir

	out, err := cmd.Output()
	if err != nil {
		return err
	}

	if len(out) == 0 {
		fmt.Println("ℹ Nenhuma alteração detectada (skip commit)")
		return nil
	}

	return git(dir, "commit", "-m", message)
}
