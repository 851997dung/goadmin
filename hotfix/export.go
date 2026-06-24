package hotfix

import (
	"context"
	"fmt"
	"path/filepath"
)

// ExportTables runs db2xml and xml2dbfile via runXML.app.
func ExportTables(ctx context.Context, cfg Config, progress func(string)) error {
	workDir := runXMLBinDir(cfg)
	runXML := filepath.Join(workDir, "runXML.app")

	progress("runXML.app db2xml")
	if _, err := RunCmdInDir(ctx, workDir, runXML, "db2xml"); err != nil {
		return fmt.Errorf("runXML db2xml: %w", err)
	}

	progress("runXML.app xml2dbfile")
	if _, err := RunCmdInDir(ctx, workDir, runXML, "xml2dbfile"); err != nil {
		return fmt.Errorf("runXML xml2dbfile: %w", err)
	}

	return nil
}

// ExportAndArchiveAll exports all tables and produces xml.zip.
func ExportAndArchiveAll(ctx context.Context, cfg Config, progress func(string)) (string, error) {
	if err := ExportTables(ctx, cfg, progress); err != nil {
		return "", err
	}
	workDir := runXMLBinDir(cfg)
	zipPath := filepath.Join(workDir, "xml.zip")
	progress("archiving all tables -> xml.zip")
	if err := ArchiveFile(ctx, zipPath, workDir, "xml"); err != nil {
		return "", fmt.Errorf("archive all tables: %w", err)
	}
	return zipPath, nil
}
