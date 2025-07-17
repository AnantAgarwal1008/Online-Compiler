package handler

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"os/exec"
	"path/filepath"
)

//Input structure expected from frontend

type CodeRequest struct {
	Language string `json:"language"`
	Code     string `json:"code"`
	Input    string `json:"input"`
}

type CodeResponse struct {
	Success bool   `json:"success"`
	Output  string `json:"output"`
	Error   string `json:"error"`
}

func RunCodeHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var req CodeRequest
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, `{"success":false,"error":"Invalid request"}`, http.StatusBadRequest)
		return
	}

	//create temporary directory

	tempDir, err := ioutil.TempDir("", "code-runner-*")
	if err != nil {
		respondError(w, "Failed to create temp directory", err)
		return
	}
	defer os.RemoveAll(tempDir)

	// save code+input to files
	filename := getFileName(req.Language)
	codePath := filepath.Join(tempDir, filename)
	inputPath := filepath.Join(tempDir, "input.txt")

	if err := ioutil.WriteFile(codePath, []byte(req.Code), 0644); err != nil {
		respondError(w, "Failed to save code file", err)
		return
	}
	if err := ioutil.WriteFile(inputPath, []byte(req.Input), 0644); err != nil {
		respondError(w, "Failed to save input file", err)
		return
	}
	dockerCmd := exec.Command("docker", "run",
		"--rm",
		"-v", tempDir+":/app",
		getImage(req.Language),
		"/app/"+filename,
		"/app/input.txt",
	)
	output, err := dockerCmd.CombinedOutput()
	if err != nil {
		json.NewEncoder(w).Encode(CodeResponse{
			Success: false,
			Output:  "",
			Error:   string(output),
		})
		return
	}

	json.NewEncoder(w).Encode(CodeResponse{
		Success: true,
		Output:  string(output),
		Error:   "",
	})
}

func getFileName(lang string) string {
	switch lang {
	case "python":
		return "main.py"
	case "cpp":
		return "main.cpp"
	case "java":
		return "Main.java"
	case "go":
		return "main.go"
	default:
		return "main.txt"
	}
}

func getImage(lang string) string {
	switch lang {
	case "python":
		return "online-python"
	case "cpp":
		return "online-cpp"
	case "java":
		return "online-java"
	case "go":
		return "online-go"
	default:
		return "online-base"
	}
}

func respondError(w http.ResponseWriter, msg string, err error) {
	log.Println("ERROR:", msg, err)
	json.NewEncoder(w).Encode(CodeResponse{
		Success: false,
		Output:  "",
		Error:   msg + ": " + err.Error(),
	})
}
