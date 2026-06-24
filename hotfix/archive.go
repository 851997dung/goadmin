package hotfix

import (
	"archive/zip"
	"context"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"strings"
)

// ArchiveFile creates a zip at tgtFile containing the file/directory srcFile from srcDir.
func ArchiveFile(ctx context.Context, tgtFile, srcDir, srcFile string) error {
	return ArchiveFiles(ctx, tgtFile, srcDir, []string{srcFile})
}

// ArchiveFiles creates a zip at tgtFile containing the listed files/dirs from srcDir.
func ArchiveFiles(ctx context.Context, tgtFile, srcDir string, srcFiles []string) error {
	rawFile, err := os.Create(tgtFile)
	if err != nil {
		return fmt.Errorf("create zip %q: %w", tgtFile, err)
	}
	defer rawFile.Close()

	zipFile := zip.NewWriter(rawFile)
	defer zipFile.Close()

	for _, srcFile := range srcFiles {
		path := filepath.Join(srcDir, srcFile)
		err = filepath.Walk(path, func(path string, info os.FileInfo, err error) error {
			if err != nil {
				return err
			}
			if IsContextDone(ctx) {
				return ctx.Err()
			}
			return writeZipEntry(zipFile, info, path, srcDir)
		})
		if err != nil {
			return fmt.Errorf("archive %s: %w", srcFile, err)
		}
	}
	return nil
}

func writeZipEntry(zipFile *zip.Writer, fileInfo os.FileInfo, filePath, srcDir string) error {
	name, err := filepath.Rel(srcDir, filePath)
	if err != nil {
		return err
	}
	if name == "." {
		return nil
	}
	header, err := zip.FileInfoHeader(fileInfo)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		name += "/"
	}
	header.Method = zip.Deflate
	header.Name = strings.ReplaceAll(name, "\\", "/")
	writer, err := zipFile.CreateHeader(header)
	if err != nil {
		return err
	}
	if fileInfo.IsDir() {
		return nil
	}
	file, err := os.Open(filePath)
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(writer, file)
	return err
}

// ArchiveSeveralTables creates xml(part).zip containing only the specified table XML files.
func ArchiveSeveralTables(ctx context.Context, cfg Config, tableNames []string, progress func(string)) (string, error) {
	srcFiles := MapTableName2ConfigFiles(tableNames)
	workDir := runXMLBinDir(cfg)
	tgtFile := filepath.Join(workDir, "xml(part).zip")
	progress(fmt.Sprintf("archiving %d tables -> xml(part).zip", len(tableNames)))
	if err := ArchiveFiles(ctx, tgtFile, workDir, srcFiles); err != nil {
		return "", err
	}
	return tgtFile, nil
}
