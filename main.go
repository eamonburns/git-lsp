package main

import (
	"bufio"
	"encoding/json"
	"fmt"
	"io"
	"log/slog"
	"os"
	"path/filepath"

	"github.com/eamonburns/git-lsp/lsp"
	"github.com/eamonburns/git-lsp/rpc"
)

func main() {
	homeDir, err := os.UserHomeDir()
	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	stateDir := filepath.Join(homeDir, ".local", "state", "git-lsp")
	if err := os.MkdirAll(stateDir, 0o755); err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}

	// Setup logging
	logfile, err := os.OpenFile(filepath.Join(stateDir, "git-lsp.log"), os.O_CREATE|os.O_TRUNC|os.O_WRONLY, 0o644)

	if err != nil {
		fmt.Fprintf(os.Stderr, "error: %v", err)
		os.Exit(1)
	}
	logger := slog.New(slog.NewTextHandler(logfile, nil))
	slog.SetDefault(logger)

	// Start LSP
	slog.Info("LSP Started")

	scanner := bufio.NewScanner(os.Stdin)
	scanner.Split(rpc.Split)

	state := any(nil) // TODO:
	_ = state
	writer := os.Stdout
	_ = writer

	for scanner.Scan() {
		msg := scanner.Bytes()
		method, contents, err := rpc.DecodeMessage(msg)
		if err != nil {
			slog.Info("unable to decode message", "error", err)
			continue
		}

		handleMessage(writer, state, method, contents)
	}
}

func handleMessage(writer io.Writer, state any, method string, contents []byte) {
	slog.Info("Recieved message", "method", method)

	switch method {
	case "initialize":
		var request lsp.InitializeRequest
		if err := json.Unmarshal(contents, &request); err != nil {
			slog.Error("unable to parse request", "error", err)
		}

		slog.Info(
			"connected to client",
			"name", request.Params.ClientInfo.Name,
			"version", request.Params.ClientInfo.Version,
		)

		// TODO: Respond
	}
}
