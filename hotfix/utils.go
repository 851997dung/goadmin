package hotfix

import (
	"bytes"
	"context"
	"crypto/md5"
	"crypto/tls"
	"encoding/hex"
	"errors"
	"fmt"
	"io"
	"net/http"
	"os"
	"os/exec"
	"strings"
)

func HashFileMD5(filePath string) (string, error) {
	f, err := os.Open(filePath)
	if err != nil {
		return "", fmt.Errorf("hashFile open %q: %w", filePath, err)
	}
	defer f.Close()
	h := md5.New()
	if _, err := io.Copy(h, f); err != nil {
		return "", fmt.Errorf("hashFile copy %q: %w", filePath, err)
	}
	return hex.EncodeToString(h.Sum(nil)), nil
}

func RunCmd(ctx context.Context, name string, args ...string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	cmd := exec.CommandContext(ctx, name, args...)
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return buf.String(), fmt.Errorf("run %s: %w\n%s", name, err, buf.String())
	}
	return buf.String(), nil
}

func RunCmdInDir(ctx context.Context, dir, name string, args ...string) (string, error) {
	if ctx == nil {
		ctx = context.Background()
	}
	cmd := exec.CommandContext(ctx, name, args...)
	cmd.Dir = dir
	var buf bytes.Buffer
	cmd.Stdout = &buf
	cmd.Stderr = &buf
	if err := cmd.Run(); err != nil {
		return buf.String(), fmt.Errorf("run %s in %s: %w\n%s", name, dir, err, buf.String())
	}
	return buf.String(), nil
}

var httpClient = &http.Client{
	Transport: &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: true},
	},
}

func HTTPGet(url string) ([]byte, error) {
	resp, err := httpClient.Get(url)
	if err != nil {
		return nil, fmt.Errorf("http.Get %q: %w", url, err)
	}
	defer resp.Body.Close()
	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body %q: %w", url, err)
	}
	if !bytes.Equal(body, []byte("All Is OK!")) {
		return body, errors.New(string(body))
	}
	return body, nil
}

func HTTPPostForm(url string, data map[string]string) ([]byte, error) {
	vals := make([]string, 0, len(data)*2)
	for k, v := range data {
		vals = append(vals, k, v)
	}
	body := strings.NewReader(strings.ReplaceAll(
		strings.ReplaceAll(
			strings.Join(vals, "="), "=", "%3D"), "&", "%26"))

	resp, err := httpClient.Post(url, "application/x-www-form-urlencoded", body)
	if err != nil {
		return nil, fmt.Errorf("http.Post %q: %w", url, err)
	}
	defer resp.Body.Close()
	rbody, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, fmt.Errorf("read body %q: %w", url, err)
	}
	if !bytes.Equal(rbody, []byte("All Is OK!")) {
		return rbody, errors.New(string(rbody))
	}
	return rbody, nil
}

func IsContextDone(ctx context.Context) bool {
	select {
	case <-ctx.Done():
		return true
	default:
		return false
	}
}
