package hotfix

import (
	"bufio"
	"bytes"
	"context"
	"fmt"
	"path/filepath"
	"strings"
)

var allSpellTableNames = []string{
	"spell_info", "spell_level_info", "spell_level_effect_info",
	"SpellInfo", "SpellLevelInfo", "SpellLevelEffectInfo",
}

var allLootTableNames = []string{
	"loot_set", "loot_set_group", "loot_set_group_item",
	"LootSet", "LootSetGroup", "LootSetGroupItem",
}

func IsSpellTable(name string) bool {
	for _, n := range allSpellTableNames {
		if n == name {
			return true
		}
	}
	return false
}

func IsLootTable(name string) bool {
	for _, n := range allLootTableNames {
		if n == name {
			return true
		}
	}
	return false
}

// MapTableName2ConfigFiles maps table names to XML file paths relative to RunXML/bin.
func MapTableName2ConfigFiles(tableNames []string) []string {
	var files []string
	for _, name := range tableNames {
		switch name {
		case "$Spell":
			files = append(files, "xml/Spell")
		case "$Loot":
			files = append(files, "xml/Loot")
		default:
			files = append(files, fmt.Sprintf("xml/%s.xml", name))
		}
	}
	return files
}

// CanonicalizeTableNames runs runXML.app listTable and resolves/normalizes user-provided names.
func CanonicalizeTableNames(cfg Config, tableNames []string) error {
	workDir := runXMLBinDir(cfg)
	runXML := filepath.Join(workDir, "runXML.app")

	output, err := RunCmdInDir(context.Background(), workDir, runXML, "listTable")
	if err != nil {
		return fmt.Errorf("runXML listTable: %w", err)
	}

	allTableNames, err := parseAllTableNames([]byte(output))
	if err != nil {
		return err
	}
	allTableNames = append(allTableNames, []string{"scriptable", "Scriptable"})

	for i, name := range tableNames {
		if IsSpellTable(name) {
			tableNames[i] = "$Spell"
		} else if IsLootTable(name) {
			tableNames[i] = "$Loot"
		} else {
			target := queryTableName(allTableNames, name)
			if target != "" {
				tableNames[i] = target
			} else {
				return fmt.Errorf("invalid table name %q", name)
			}
		}
	}
	return nil
}

func parseAllTableNames(data []byte) ([][]string, error) {
	var result [][]string
	scanner := bufio.NewScanner(bytes.NewReader(data))
	for scanner.Scan() {
		names := strings.Fields(scanner.Text())
		result = append(result, names)
	}
	if err := scanner.Err(); err != nil {
		return nil, fmt.Errorf("parse table names: %w", err)
	}
	return result, nil
}

func queryTableName(allTableNames [][]string, target string) string {
	match := func(name, target string) bool {
		if strings.EqualFold(name, target) {
			return true
		}
		if strings.HasPrefix(name, "auto_") && strings.EqualFold(name[5:], target) {
			return true
		}
		return false
	}
	for _, row := range allTableNames {
		for _, name := range row {
			if match(name, target) {
				return row[len(row)-1]
			}
		}
	}
	return ""
}

func stringSliceContains(slice []string, target string) bool {
	for _, s := range slice {
		if s == target {
			return true
		}
	}
	return false
}

func stringSliceIndex(slice []string, target string) int {
	for i, s := range slice {
		if s == target {
			return i
		}
	}
	return -1
}

func stringSliceDelete(slice []string, i, j int) []string {
	return append(slice[:i], slice[j:]...)
}

// ParseHotfixDBTableArgs detects special tables and returns remaining table names + URL extension params.
func ParseHotfixDBTableArgs(tableNames []string) ([]string, string) {
	var extParams string
	names := make([]string, len(tableNames))
	copy(names, tableNames)

	if i := stringSliceIndex(names, "$Spell"); i != -1 {
		names = stringSliceDelete(names, i, i+1)
		extParams += "&Spell=true&SpellArgs=true"
	}
	if i := stringSliceIndex(names, "$Loot"); i != -1 {
		names = stringSliceDelete(names, i, i+1)
		extParams += "&Loot=true"
	}
	if i := stringSliceIndex(names, "QuestPrototype"); i != -1 {
		names = stringSliceDelete(names, i, i+1)
		extParams += "&Quest=true"
	}
	if i := stringSliceIndex(names, "ItemPrototype"); i != -1 {
		names = stringSliceDelete(names, i, i+1)
		extParams += "&Item=true"
	}
	if i := stringSliceIndex(names, "MapSpecial"); i != -1 {
		names = stringSliceDelete(names, i, i+1)
		extParams += "&MapSpecial=true"
	}
	if i := stringSliceIndex(names, "operating_activities"); i != -1 {
		names = stringSliceDelete(names, i, i+1)
		extParams += "&Activity=true"
	}
	return names, extParams
}
