package cmd

import (
	"bytes"
	"os"
	"syscall"
	"testing"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/kubetail-org/kubetail/modules/shared/config"
	"github.com/stretchr/testify/require"
)

func TestServeCmd_WithTestFlag(t *testing.T) {
	gin.SetMode("test")

	// Point to rootCmd
	cmd := rootCmd
	buf := &bytes.Buffer{}
	cmd.SetOut(buf)
	cmd.SetErr(buf)
	cmd.SetArgs([]string{"serve", "--test"})

	err := cmd.Execute()
	require.NoError(t, err)
	require.Contains(t, buf.String(), "ok")
}

func TestServe_ShutsDownGracefully(t *testing.T) {
	gin.SetMode("test")

	cfg := config.DefaultConfig()

	// Quit signal we control
	quit := make(chan os.Signal, 1)

	// Run serve in background
	done := make(chan struct{})
	go func() {
		serve(cfg, "127.0.0.1", true, quit)
		close(done)
	}()

	// Give server time to start
	time.Sleep(200 * time.Millisecond)

	// Trigger shutdown
	quit <- syscall.SIGINT

	select {
	case <-done:
		// success
	case <-time.After(20 * time.Second):
		t.Fatal("server did not shut down in time")
	}
}
