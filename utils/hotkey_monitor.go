package utils

import (
	"context"
	"os"
	"time"
)

type HotkeyMonitor struct {
	ctx      context.Context
	cancel   context.CancelFunc
	isActive bool
}

func NewHotkeyMonitor() *HotkeyMonitor {
	ctx, cancel := context.WithCancel(context.Background())
	return &HotkeyMonitor{
		ctx:    ctx,
		cancel: cancel,
	}
}

func (hm *HotkeyMonitor) Start(romPath string, overlayClient *OverlayClient) {
	if hm.isActive {
		return
	}

	hm.isActive = true
	go hm.monitor(romPath, overlayClient)
}

func (hm *HotkeyMonitor) Stop() {
	if !hm.isActive {
		return
	}

	hm.cancel()
	hm.isActive = false
}

func (hm *HotkeyMonitor) monitor(romPath string, overlayClient *OverlayClient) {
	Logger.Println("Starting hotkey monitor for in-game menu...")

	ticker := time.NewTicker(100 * time.Millisecond)
	defer ticker.Stop()

	for {
		select {
		case <-hm.ctx.Done():
			Logger.Println("Hotkey monitor stopped")
			return
		case <-ticker.C:
			// Only monitor if RetroArch is still running
			if !IsRetroArchRunning() {
				Logger.Println("RetroArch not running, stopping hotkey monitor")
				return
			}

			if hm.checkHotkey() {
				Logger.Println("Hotkey detected - showing in-game menu")
				response, err := overlayClient.ShowMenu(romPath)
				if err != nil {
					Logger.Printf("Failed to show menu: %v", err)
					continue
				}

				Logger.Printf("Menu response: %s", response.Action)

				if response.Action == "exit_game" {
					Logger.Println("User requested exit - terminating RetroArch")
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
