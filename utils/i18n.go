package utils

import (
	"encoding/json"

	"github.com/nicksnyder/go-i18n/v2/i18n"
	"golang.org/x/text/language"
)

var localizer *i18n.Localizer
var bundle *i18n.Bundle

func init() {
	bundle = i18n.NewBundle(language.English)
	bundle.RegisterUnmarshalFunc("json", json.Unmarshal)
	bundle.LoadMessageFile("resources/i18n/en.json")
	bundle.LoadMessageFile("resources/i18n/es.json")
	localizer = i18n.NewLocalizer(bundle, language.English.String(), language.Spanish.String())
}

// GetString retrieves a localized string by key
// If the key is not found, it returns the key itself as fallback
func GetString(key string) string {
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID: key,
	})
	if err != nil {
		return "I18N Error"
	}
	return msg
}

// GetStringWithData retrieves a localized string by key with template data
// templateData should be a map[string]interface{} containing values for template variables
func GetStringWithData(key string, templateData map[string]interface{}) string {
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:    key,
		TemplateData: templateData,
	})
	if err != nil {
		return "I18N Error"
	}
	return msg
}

// GetPluralString retrieves a localized string with plural support
// count determines which plural form to use
func GetPluralString(key string, count int) string {
	msg, err := localizer.Localize(&i18n.LocalizeConfig{
		MessageID:   key,
		PluralCount: count,
	})
	if err != nil {
		return "I18N Error"
	}
	return msg
}
