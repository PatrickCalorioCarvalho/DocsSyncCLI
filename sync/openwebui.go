package sync

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"mime/multipart"
	"net/http"
	"os"
	"path/filepath"
	"time"

	"github.com/PatrickCalorioCarvalho/DocsSyncCLI/config"
)

type uploadResponse struct {
	ID string `json:"id"`
}

type processStatusResponse struct {
	Status string `json:"status"`
	Error  string `json:"error,omitempty"`
}

type knowledgeFilesResponse struct {
	Items []struct {
		ID       string `json:"id"`
		Filename string `json:"filename"`
		Hash     string `json:"hash"`
	} `json:"items"`
}

func SyncOpenWebUI(cfg *config.Config, precommitDir string) error {

	ow := cfg.Sync.OpenWebUI

	if ow.ApiUrl == "" || ow.ApiKey == "" || ow.KnowledgeId == "" {
		return fmt.Errorf("openwebui.apiUrl, apiKey ou knowledgeId não configurados")
	}

	files := []string{}

	err := filepath.Walk(precommitDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}
		if !info.IsDir() && filepath.Ext(path) == ".md" {
			files = append(files, path)
		}
		return nil
	})
	if err != nil {
		return err
	}

	if len(files) == 0 {
		fmt.Println("⚠ Nenhum arquivo Markdown para enviar ao OpenWebUI")
		return nil
	}
	if err := clearKnowledge(ow); err != nil {
		return err
	}
	for _, file := range files {
		fmt.Println("   📤 upload:", filepath.Base(file))

		fileID, err := uploadFile(ow, file)
		if err != nil {
			return err
		}

		if err := waitForProcessing(ow, fileID, 5*time.Minute); err != nil {
			return err
		}

		if err := addToKnowledge(ow, fileID); err != nil {
			return err
		}

		fmt.Println("   📚 adicionado à knowledge base:", filepath.Base(file))
	}

	fmt.Println("   ↳ OpenWebUI ingest finalizado")
	return nil
}

func uploadFile(cfg config.OpenWebUIConfig, filePath string) (string, error) {

	file, err := os.Open(filePath)
	if err != nil {
		return "", err
	}
	defer file.Close()

	body := &bytes.Buffer{}
	writer := multipart.NewWriter(body)

	part, err := writer.CreateFormFile("file", filepath.Base(filePath))
	if err != nil {
		return "", err
	}

	if _, err := io.Copy(part, file); err != nil {
		return "", err
	}

	writer.Close()

	req, err := http.NewRequest(
		"POST",
		cfg.ApiUrl+"/api/v1/files/",
		body,
	)
	if err != nil {
		return "", err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.ApiKey)
	req.Header.Set("Content-Type", writer.FormDataContentType())
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return "", fmt.Errorf("erro upload (%d): %s", resp.StatusCode, string(respBody))
	}

	var result uploadResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	if result.ID == "" {
		return "", fmt.Errorf("upload retornou ID vazio")
	}

	return result.ID, nil
}

func waitForProcessing(cfg config.OpenWebUIConfig, fileID string, timeout time.Duration) error {

	fmt.Println("   ⏳ processando:", fileID)

	start := time.Now()

	for time.Since(start) < timeout {

		req, err := http.NewRequest(
			"GET",
			cfg.ApiUrl+"/api/v1/files/"+fileID+"/process/status",
			nil,
		)
		if err != nil {
			return err
		}

		req.Header.Set("Authorization", "Bearer "+cfg.ApiKey)
		req.Header.Set("Accept", "application/json")

		resp, err := http.DefaultClient.Do(req)
		if err != nil {
			return err
		}

		var status processStatusResponse
		err = json.NewDecoder(resp.Body).Decode(&status)
		resp.Body.Close()

		if err != nil {
			return err
		}

		switch status.Status {
		case "completed":
			fmt.Println("   ✅ processamento concluído")
			return nil
		case "failed":
			return fmt.Errorf("processamento falhou: %s", status.Error)
		default:
			time.Sleep(2 * time.Second)
		}
	}

	return fmt.Errorf("timeout aguardando processamento do arquivo %s", fileID)
}

func addToKnowledge(cfg config.OpenWebUIConfig, fileID string) error {

	payload := map[string]string{
		"file_id": fileID,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(
		"POST",
		cfg.ApiUrl+"/api/v1/knowledge/"+cfg.KnowledgeId+"/file/add",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao adicionar à knowledge (%d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}

func clearKnowledge(cfg config.OpenWebUIConfig) error {

	fmt.Println("🧹 Limpando knowledge base...")

	fileIDs, err := listKnowledgeFiles(cfg)
	if err != nil {
		return err
	}

	if len(fileIDs) == 0 {
		fmt.Println("   ↳ knowledge já está vazia")
		return nil
	}

	for _, id := range fileIDs {
		if err := deleteKnowledgeFile(cfg, id); err != nil {
			return err
		}
		fmt.Println("   🗑 removido:", id)
	}

	return nil
}


func listKnowledgeFiles(cfg config.OpenWebUIConfig) ([]string, error) {

	req, err := http.NewRequest(
		"GET",
		cfg.ApiUrl+"/api/v1/knowledge/"+cfg.KnowledgeId+"/files",
		nil,
	)
	if err != nil {
		return nil, err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")


	resp, err := http.DefaultClient.Do(req)

	if err != nil {
		return nil, err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		body, _ := io.ReadAll(resp.Body)
		return nil, fmt.Errorf(
			"erro ao listar knowledge (%d): %s",
			resp.StatusCode,
			string(body),
		)
	}

	var result knowledgeFilesResponse
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return nil, err
	}

	ids := []string{}
	for _, f := range result.Items {
		ids = append(ids, f.ID)
	}
	return ids, nil
}
func deleteKnowledgeFile(cfg config.OpenWebUIConfig, fileID string) error {
	payload := map[string]string{
		"file_id": fileID,
	}

	body, _ := json.Marshal(payload)

	req, err := http.NewRequest(
		"POST",
		cfg.ApiUrl+"/api/v1/knowledge/"+cfg.KnowledgeId+"/file/remove",
		bytes.NewBuffer(body),
	)
	if err != nil {
		return err
	}

	req.Header.Set("Authorization", "Bearer "+cfg.ApiKey)
	req.Header.Set("Content-Type", "application/json")
	req.Header.Set("Accept", "application/json")

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return err
	}

	defer resp.Body.Close()

	if resp.StatusCode >= 300 {
		respBody, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("erro ao adicionar à knowledge (%d): %s", resp.StatusCode, string(respBody))
	}

	return nil
}
			