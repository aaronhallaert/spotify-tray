package tests_test

import (
	"context"
	"os"
	"os/exec"
	"testing"
	"time"
)

func TestNativeTrayLifecycle(t *testing.T) {
	if os.Getenv("SPOTIFY_TRAY_TEST_APPKIT") != "1" {
		t.Skip("set SPOTIFY_TRAY_TEST_APPKIT=1 to run the interactive AppKit test")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 30*time.Second)
	defer cancel()
	cmd := exec.CommandContext(ctx, "go", "run", "./testdata/smoke")
	if output, err := cmd.CombinedOutput(); err != nil {
		t.Fatalf("native tray test: %v\n%s", err, output)
	}
}
