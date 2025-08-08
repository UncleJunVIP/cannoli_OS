package utils

import (
	gaba "github.com/UncleJunVIP/gabagool/pkg/gabagool"
	"time"
)

func ShowMessage(message string, delay time.Duration) {
	gaba.ProcessMessage(message, gaba.ProcessMessageOptions{}, func() (interface{}, error) {
		time.Sleep(delay)
		return nil, nil
	})
}

func ShowConfirmation(message string) bool {
	return ShowCustomConfirmation(message, "Cancel", "Confirm", gaba.ButtonA)
}

func ShowCustomConfirmation(message string, cancelText string, confirmText string, confirmButton gaba.Button) bool {
	result, err := gaba.ConfirmationMessage(message, []gaba.FooterHelpItem{
		{ButtonName: "B", HelpText: cancelText},
		{ButtonName: "A", HelpText: confirmText},
	}, gaba.MessageOptions{ConfirmButton: confirmButton})

	return err == nil && result.IsSome()
}
