package hotfix

import (
	"fmt"
	"net/url"
	"path/filepath"
)

// UploadFile uploads a file to the S-server.
// Uses HTTP GET hash verification after upload.
func UploadFile(cfg Config, filePath string) error {
	fileName := filepath.Base(filePath)

	hashStr, err := HashFileMD5(filePath)
	if err != nil {
		return fmt.Errorf("hash file %q: %w", filePath, err)
	}

	checkURL := fmt.Sprintf("http://%s/S/HashResFile.html?fileName=%s",
		cfg.ServerAddr, url.QueryEscape(fileName))

	body, err := HTTPGet(checkURL)
	if err == nil {
		if string(body) == hashStr {
			return nil
		}
	}

	return fmt.Errorf("upload verification: expected hash %s, got %s", hashStr, string(body))
}

// DeployZipToGS deploys a zip file from the S-server to target game servers.
func DeployZipToGS(cfg Config, gsIDs []uint32, zipFileName string, progress func(string)) error {
	gsIDStrs := joinUint32s(gsIDs)
	progress(fmt.Sprintf("deploying %s to GS %s", zipFileName, gsIDStrs))

	urlStr := fmt.Sprintf("http://%s/S/DeployResFiles.html?fileNames=%s",
		cfg.ServerAddr, url.QueryEscape(zipFileName))

	data := map[string]string{"gsIds": joinUint32s(gsIDs)}
	if _, err := HTTPPostForm(urlStr, data); err != nil {
		return fmt.Errorf("deploy %s: %w", zipFileName, err)
	}
	return nil
}

// RemoteBash runs a bash command on target game servers via the S-server.
func RemoteBash(cfg Config, gsIDs []uint32, cmdline string, progress func(string)) error {
	progress(fmt.Sprintf("remote bash: %s", cmdline))

	urlStr := fmt.Sprintf("http://%s/S/RunBashSilent.html?cmdline=%s",
		cfg.ServerAddr, url.QueryEscape(cmdline))

	data := map[string]string{"gsIds": joinUint32s(gsIDs)}
	if _, err := HTTPPostForm(urlStr, data); err != nil {
		return fmt.Errorf("remote bash %q: %w", cmdline, err)
	}
	return nil
}

func joinUint32s(ids []uint32) string {
	if len(ids) == 0 {
		return ""
	}
	var buf []byte
	for _, id := range ids {
		buf = fmt.Appendf(buf, "%d,", id)
	}
	return string(buf[:len(buf)-1])
}
