package locale

import (
	"github.com/nicksnyder/go-i18n/i18n"
)

var (
	defaultLang = "en-US"
	defaultFunc i18n.TranslateFunc
)

// Init i18n initialize
func Init() {
	loadLanguage("en-US.all.yaml", defaultLang)
	loadLanguage("zh-CN.all.yaml", "zh-CN")
	defaultFunc, _ = i18n.Tfunc(defaultLang)
}

func loadLanguage(filename, lang string) {
	fileBytes, err := ReadFile(filename)
	if err != nil {
		panic(err)
	}
	err = i18n.ParseTranslationFileBytes(filename, fileBytes)
	if err != nil {
		panic(err)
	}
}

// Locate 获取对应语言的翻译方法
func Locate(lang string) i18n.TranslateFunc {
	tfunc, err := i18n.Tfunc(lang)
	if err != nil {
		return defaultFunc
	}
	return tfunc
}
