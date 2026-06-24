package hotfix

import (
	"context"
	"fmt"
	"net/url"
)

// SkipFlags controls which pipeline stages to skip.
type SkipFlags struct {
	IsSkipBuildTool bool
	IsSkipConfigs   bool
}

// RunFullPipeline executes all 5 stages of the hotfix data table pipeline.
func RunFullPipeline(session *PipelineSession, cfg Config, gsIDs []uint32, tableNames []string, skip SkipFlags) error {
	progress := session.AddMsg

	// Stage 1: Build & deploy tool
	if !skip.IsSkipBuildTool {
		progress("=== Stage 1: Build Tool ===")
		zipPath, err := BuildTool(context.Background(), cfg, progress)
		if err != nil {
			return fmt.Errorf("build tool: %w", err)
		}

		progress("Uploading runXML.zip")
		if err := UploadFile(cfg, zipPath); err != nil {
			return fmt.Errorf("upload runXML.zip: %w", err)
		}

		progress("Deploying runXML.zip to GS")
		if err := DeployZipToGS(cfg, gsIDs, "runXML.zip", progress); err != nil {
			return fmt.Errorf("deploy runXML.zip: %w", err)
		}

		progress("chmod +x runXML.app on GS")
		cmdline := "chmod +x runXML.app"
		if err := RemoteBash(cfg, gsIDs, cmdline, progress); err != nil {
			return fmt.Errorf("chmod runXML.app: %w", err)
		}
	}

	// Stage 2: Canonicalize table names
	progress("=== Canonicalize Table Names ===")
	if err := CanonicalizeTableNames(cfg, tableNames); err != nil {
		return fmt.Errorf("canonicalize: %w", err)
	}
	progress(fmt.Sprintf("Resolved table names: %v", tableNames))

	// Stage 3: Export, archive, upload, deploy configs
	if !skip.IsSkipConfigs {
		progress("=== Stage 2-4: Export -> Archive -> Upload -> Deploy ===")

		progress("Exporting DB tables")
		if err := ExportTables(context.Background(), cfg, progress); err != nil {
			return fmt.Errorf("export: %w", err)
		}

		// Clean spell/loot directories before archiving
		for _, name := range tableNames {
			if name == "$Spell" {
				cmdline := "rm -rf xml/Spell"
				urlVals := fmt.Sprintf("gsIds=%s", joinUint32s(gsIDs))
				urlStr := fmt.Sprintf("http://%s/S/RunBashSilent.html?cmdline=%s&%s",
					cfg.ServerAddr, url.QueryEscape(cmdline), urlVals)
				if _, err := HTTPGet(urlStr); err != nil {
					return fmt.Errorf("rm xml/Spell: %w", err)
				}
			}
			if name == "$Loot" {
				cmdline := "rm -rf xml/Loot"
				urlVals := fmt.Sprintf("gsIds=%s", joinUint32s(gsIDs))
				urlStr := fmt.Sprintf("http://%s/S/RunBashSilent.html?cmdline=%s&%s",
					cfg.ServerAddr, url.QueryEscape(cmdline), urlVals)
				if _, err := HTTPGet(urlStr); err != nil {
					return fmt.Errorf("rm xml/Loot: %w", err)
				}
			}
		}

		progress("Archiving selected tables")
		zipPath, err := ArchiveSeveralTables(context.Background(), cfg, tableNames, progress)
		if err != nil {
			return fmt.Errorf("archive: %w", err)
		}

		progress("Uploading xml(part).zip")
		if err := UploadFile(cfg, zipPath); err != nil {
			return fmt.Errorf("upload xml(part).zip: %w", err)
		}

		progress("Deploying xml(part).zip to GS")
		if err := DeployZipToGS(cfg, gsIDs, "xml(part).zip", progress); err != nil {
			return fmt.Errorf("deploy xml(part).zip: %w", err)
		}
	}

	// Stage 5: GM commands
	progress("=== Stage 5: GM Commands ===")
	remaining, extParams := ParseHotfixDBTableArgs(tableNames)
	if extParams != "" {
		progress(fmt.Sprintf("Special table params: %s", extParams[1:]))
	}
	if err := SendGMCommands(context.Background(), cfg, gsIDs, remaining, extParams, progress); err != nil {
		return fmt.Errorf("GM commands: %w", err)
	}

	return nil
}

// RunGMOnlyPipeline runs only the GM command stage (for re-triggering hotfix on already deployed files).
func RunGMOnlyPipeline(session *PipelineSession, cfg Config, gsIDs []uint32, tableNames []string) error {
	progress := session.AddMsg

	if err := CanonicalizeTableNames(cfg, tableNames); err != nil {
		return fmt.Errorf("canonicalize: %w", err)
	}

	remaining, extParams := ParseHotfixDBTableArgs(tableNames)
	if extParams != "" {
		progress(fmt.Sprintf("Special table params: %s", extParams[1:]))
	}
	return SendGMCommands(context.Background(), cfg, gsIDs, remaining, extParams, progress)
}
