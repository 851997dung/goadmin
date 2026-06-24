package hotfix

import (
	"context"
	"fmt"
	"net/url"
	"strings"
)

// SendGMCommands sends the full GM command sequence to hot-reload DB tables on target GS.
func SendGMCommands(ctx context.Context, cfg Config, gsIDs []uint32, tableNames []string, extParams string, progress func(string)) error {
	// Step 1: run xml2dbfile on each GS
	progress("running xml2dbfile on GS")
	cmdline := "./runXML.app xml2dbfile"
	if err := RemoteBash(cfg, gsIDs, cmdline, progress); err != nil {
		return fmt.Errorf("xml2dbfile: %w", err)
	}

	// Step 2: AdminServer HotfixDBFile
	progress("AdminServer HotfixDBFile")
	if err := sendGMCommand(cfg, gsIDs, "AdminServer", "HotfixDBFile", "", progress); err != nil {
		return fmt.Errorf("HotfixDBFile: %w", err)
	}

	// Step 3: AdminServer HotfixDBTables (with table names)
	if len(tableNames) > 0 {
		progress(fmt.Sprintf("AdminServer HotfixDBTables: %s", strings.Join(tableNames, ",")))
		if err := sendGMPostCommand(cfg, gsIDs, "AdminServer", "HotfixDBTables", tableNames, progress); err != nil {
			return fmt.Errorf("HotfixDBTables: %w", err)
		}
	}

	// Step 4-10: Conditional special-table commands
	if containsParam(extParams, "Spell") && containsParam(extParams, "SpellArgs") {
		progress("MapServer HotfixSpellRelation")
		if err := sendGMCommand(cfg, gsIDs, "MapServer", "HotfixSpellRelation", "", progress); err != nil {
			return fmt.Errorf("HotfixSpellRelation: %w", err)
		}
		progress("MapServer ClearSpellEffectArgs")
		if err := sendGMCommand(cfg, gsIDs, "MapServer", "ClearSpellEffectArgs", "", progress); err != nil {
			return fmt.Errorf("ClearSpellEffectArgs: %w", err)
		}
	}

	if containsParam(extParams, "Loot") {
		progress("MapServer HotfixLootRelation")
		if err := sendGMCommand(cfg, gsIDs, "MapServer", "HotfixLootRelation", "", progress); err != nil {
			return fmt.Errorf("HotfixLootRelation: %w", err)
		}
	}

	if containsParam(extParams, "Quest") {
		progress("MapServer HotfixQuests")
		if err := sendGMCommand(cfg, gsIDs, "MapServer", "HotfixQuests", "", progress); err != nil {
			return fmt.Errorf("HotfixQuests: %w", err)
		}
	}

	if containsParam(extParams, "Item") {
		progress("MapServer HotfixItemPrototypes")
		if err := sendGMCommand(cfg, gsIDs, "MapServer", "HotfixItemPrototypes", "", progress); err != nil {
			return fmt.Errorf("HotfixItemPrototypes: %w", err)
		}
	}

	if containsParam(extParams, "MapSpecial") {
		progress("MapServer HotfixMapSpecial")
		if err := sendGMCommand(cfg, gsIDs, "MapServer", "HotfixMapSpecial", "", progress); err != nil {
			return fmt.Errorf("HotfixMapSpecial: %w", err)
		}
	}

	if containsParam(extParams, "Activity") {
		progress("SocialServer ReloadActivity")
		if err := sendGMCommand(cfg, gsIDs, "SocialServer", "ReloadActivity", "", progress); err != nil {
			return fmt.Errorf("SocialServer ReloadActivity: %w", err)
		}
	}

	return nil
}

func sendGMCommand(cfg Config, gsIDs []uint32, serverName, cmd, args string, progress func(string)) error {
	for _, gsID := range gsIDs {
		urlStr := fmt.Sprintf("https://%s:%s/GM2S.html?name=%s&cmd=%s&gsId=%d",
			cfg.CenterHost, cfg.CenterPort, serverName, cmd, gsID)
		if args != "" {
			urlStr += "&args=" + url.QueryEscape(args)
		}
		if _, err := HTTPGet(urlStr); err != nil {
			return fmt.Errorf("gsId %d: %w", gsID, err)
		}
	}
	return nil
}

func sendGMPostCommand(cfg Config, gsIDs []uint32, serverName, cmd string, tableNames []string, progress func(string)) error {
	args := strings.Join(tableNames, ",")
	for _, gsID := range gsIDs {
		urlStr := fmt.Sprintf("https://%s:%s/GM2S.html?name=%s&cmd=%s&gsId=%d",
			cfg.CenterHost, cfg.CenterPort, serverName, cmd, gsID)
		data := map[string]string{"args": args}
		if _, err := HTTPPostForm(urlStr, data); err != nil {
			return fmt.Errorf("gsId %d: %w", gsID, err)
		}
	}
	return nil
}

func containsParam(extParams, key string) bool {
	return strings.Contains(extParams, key+"=true")
}
