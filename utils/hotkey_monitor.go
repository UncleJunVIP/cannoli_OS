package utils

import (
	"context"
	"log/slog"
	"os"
	"time"
)

type HotkeyMonitor struct {
	ctx      context.Context
	cancel   context.CancelFunc
	isActive bool
	logger   *slog.Logger
}

func NewHotkeyMonitor() *HotkeyMonitor {
	ctx, cancel := context.WithCancel(context.Background())
	return &HotkeyMonitor{
		ctx:    ctx,
		cancel: cancel,
		logger: GetLoggerInstance(),
	}
}

func (hm *HotkeyMonitor) Start(overlayClient *OverlayClient) {
	if hm.isActive {
		return
	}

	hm.isActive = true
	go hm.monitor(overlayClient)
}

func (hm *HotkeyMonitor) Stop() {
	if !hm.isActive {
		return
	}

	hm.cancel()
	hm.isActive = false
}

func (hm *HotkeyMonitor) monitor(overlayClient *OverlayClient) {
	hm.logger.Debug("Starting hotkey monitor for in-game menu...")

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-hm.ctx.Done():
			hm.logger.Debug("Hotkey monitor stopped")
			return
		case <-ticker.C:
			// Only monitor if RetroArch is still running
			if !IsRetroArchRunning() {
				hm.logger.Debug("RetroArch not running, stopping hotkey monitor")
				return
			}

			if hm.checkHotkey() {
				hm.logger.Debug("Hotkey detected - showing in-game menu")
				response, err := overlayClient.ShowMenu()
				if err != nil {
					hm.logger.Debug("Failed to show menu: %v", err)
					continue
				}

				hm.logger.Debug("Menu response: %s", response.Action)

				if response.Action == "exit_game" {
					hm.logger.Debug("User requested exit - terminating RetroArch")
					return
				}

				// Add a delay to prevent multiple triggers
				time.Sleep(500 * time.Millisecond)
			}
		}
	}
}

func (hm *HotkeyMonitor) checkHotkey() bool {
	// Simple file-based trigger for testing
	if _, err := os.Stat("/tmp/trigger_igm"); err == nil {
		os.Remove("/tmp/trigger_igm")
		return true
	}
	return false
}
