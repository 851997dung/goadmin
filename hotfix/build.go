package hotfix

import (
	"context"
	"fmt"
	"path/filepath"
)

func runXMLBinDir(cfg Config) string {
	return filepath.Join(cfg.ProjectDir, "RunXML", "bin")
}

// BuildTool builds runXML.app from source and returns the path to runXML.zip.
// Caller should deploy runXML.zip to target game servers.
func BuildTool(ctx context.Context, cfg Config, progress func(string)) (string, error) {
	extraFiles := []string{
		"XClient/ClientTools/db2file/AllDBTable.list",
		"XServer/server/GameData/DB/AllDBTable.list",
		"DesignResource/导表工具/游戏表格导表/excel/~$文本ID.xlsx",
		"XServer/global",
	}
	for _, f := range extraFiles {
		p := filepath.Join(cfg.ProjectDir, f)
		progress(fmt.Sprintf("svn update %s", p))
		if _, err := RunCmdInDir(ctx, cfg.ProjectDir, "svn", "update", p); err != nil {
			return "", fmt.Errorf("svn update %s: %w", f, err)
		}
	}

	runXMLDir := filepath.Join(cfg.ProjectDir, "RunXML")
	progress(fmt.Sprintf("svn update %s", runXMLDir))
	if _, err := RunCmdInDir(ctx, cfg.ProjectDir, "svn", "update", runXMLDir); err != nil {
		return "", fmt.Errorf("svn update RunXML: %w", err)
	}

	progress("python3 scripts/allTable.py")
	if _, err := RunCmdInDir(ctx, runXMLDir, "python3", "scripts/allTable.py"); err != nil {
		return "", fmt.Errorf("run allTable.py: %w", err)
	}

	progress("python3 scripts/tableText.py")
	if _, err := RunCmdInDir(ctx, runXMLDir, "python3", "scripts/tableText.py"); err != nil {
		return "", fmt.Errorf("run tableText.py: %w", err)
	}

	progress("svn commit")
	if _, err := RunCmdInDir(ctx, runXMLDir, "svn", "commit", "-m", "auto code."); err != nil {
		return "", fmt.Errorf("svn commit: %w", err)
	}

	workDir := runXMLBinDir(cfg)
	progress(fmt.Sprintf("go build -o runXML.app %s", cfg.BuildGoModule))
	if _, err := RunCmdInDir(ctx, workDir, "go", "build", "-o", "runXML.app", cfg.BuildGoModule); err != nil {
		return "", fmt.Errorf("build runXML.app: %w", err)
	}

	progress("archiving runXML.app -> runXML.zip")
	zipPath := filepath.Join(workDir, "runXML.zip")
	if err := ArchiveFile(ctx, zipPath, workDir, "runXML.app"); err != nil {
		return "", fmt.Errorf("archive runXML.app: %w", err)
	}

	return zipPath, nil
}
