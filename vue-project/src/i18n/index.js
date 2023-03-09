import { nextTick } from 'vue'
import { createI18n } from 'vue-i18n'

export const SUPPORT_LOCALES = ['en', 'zh-cn']

const defaultOptions = { locale: 'en', fallbackLocale: 'en' }

export function setupI18n(options){
	if(!options){
		options = defaultOptions
	}else{
		for(let [v, k] in Object.entries(defaultOptions)){
			if(options[k] === undefined){
				options[k] = v
			}
		}
	}
	const i18n = createI18n(options)
	setI18nLanguage(i18n, options.locale)
	return i18n
}

export function setI18nLanguage(i18n, locale){
	if(i18n.mode === 'legacy'){
		i18n.global.locale = locale
	}else{
		i18n.global.locale.value = locale
	}
	document.querySelector('html').setAttribute('lang', locale)
}

export async function loadLocaleMessages(i18n, locale){
	const messages = await import(`../i18n/${locale}.json`)
	i18n.global.setLocaleMessage(locale, messages.default)
	return nextTick()
}
