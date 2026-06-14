package task

import (
	"bufio"
	"encoding/json"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sort"
)

type TaskSource string

const (
	SourceNPM  TaskSource = "npm"
	SourcePNPM TaskSource = "pnpm"
	SourceYarn TaskSource = "yarn"
	SourceBun  TaskSource = "bun"
	SourceMake TaskSource = "make"
	SourceGo   TaskSource = "go"
)

type Task struct {
	Name    string
	Command string
	Source  TaskSource
}

func ScanTasks(dir string) []Task {
	var tasks []Task
	tasks = append(tasks, parsePackageJSON(dir)...)
	tasks = append(tasks, parseMakefile(dir)...)
	tasks = append(tasks, parseGoMod(dir)...)
	sort.Slice(tasks, func(i, j int) bool {
		if tasks[i].Source != tasks[j].Source {
			return tasks[i].Source < tasks[j].Source
		}
		return tasks[i].Name < tasks[j].Name
	})
	return tasks
}

func parsePackageJSON(dir string) []Task {
	data, err := os.ReadFile(filepath.Join(dir, "package.json"))
	if err != nil {
		return nil
	}
	var pkg struct {
		Scripts map[string]string `json:"scripts"`
	}
	if json.Unmarshal(data, &pkg) != nil || len(pkg.Scripts) == 0 {
		return nil
	}

	pm := detectPackageManager(dir)
	tasks := make([]Task, 0, len(pkg.Scripts))
	for name := range pkg.Scripts {
		tasks = append(tasks, Task{
			Name:    name,
			Command: string(pm) + " run " + name,
			Source:  pm,
		})
	}
	return tasks
}

func detectPackageManager(dir string) TaskSource {
	checks := []struct {
		file   string
		source TaskSource
	}{
		{"bun.lock", SourceBun},
		{"bun.lockb", SourceBun},
		{"pnpm-lock.yaml", SourcePNPM},
		{"yarn.lock", SourceYarn},
	}
	for _, c := range checks {
		if _, err := os.Stat(filepath.Join(dir, c.file)); err == nil {
			return c.source
		}
	}
	return SourceNPM
}

var makeTargetRe = regexp.MustCompile(`^([a-zA-Z_][a-zA-Z0-9_-]*)\s*:`)

func parseMakefile(dir string) []Task {
	f, err := os.Open(filepath.Join(dir, "Makefile"))
	if err != nil {
		return nil
	}
	defer f.Close()

	var tasks []Task
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := scanner.Text()
		m := makeTargetRe.FindStringSubmatch(line)
		if m == nil {
			continue
		}
		target := m[1]
		tasks = append(tasks, Task{
			Name:    target,
			Command: "make " + target,
			Source:  SourceMake,
		})
	}
	return tasks
}

func parseGoMod(dir string) []Task {
	if _, err := os.Stat(filepath.Join(dir, "go.mod")); err != nil {
		return nil
	}
	if _, err := exec.LookPath("go"); err != nil {
		return nil
	}
	return []Task{
		{Name: "build", Command: "go build ./...", Source: SourceGo},
		{Name: "test", Command: "go test ./...", Source: SourceGo},
		{Name: "vet", Command: "go vet ./...", Source: SourceGo},
		{Name: "fmt", Command: "go fmt ./...", Source: SourceGo},
	}
}
