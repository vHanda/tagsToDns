package main

import (
	"encoding/hex"
	"io/ioutil"
	"math/rand"
	"os"
	"path/filepath"
	"testing"
	"time"
)

func TestFileCreate(t *testing.T) {
	filePath := TempFileName("hosts_file_path", "")
	os.Setenv("DISCOVERY_HOSTS_FILE_PATH", filePath)

	hostsFile := NewHostsFile()
	hostsFile.Add("1.1.1.1", []string{"example.host.com"})

	contents, _ := ioutil.ReadFile(filePath)
	if string(contents) != "1.1.1.1 example.host.com\n" {
		t.Error("Failed to create proper hosts file")
	}
}

func TempFileName(prefix, suffix string) string {
	randBytes := make([]byte, 16)
	rand.Seed(time.Now().UTC().UnixNano())
	rand.Read(randBytes)
	return filepath.Join(os.TempDir(), prefix+hex.EncodeToString(randBytes)+suffix)
}
