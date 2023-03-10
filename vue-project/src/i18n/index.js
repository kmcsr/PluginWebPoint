import { nextTick } from 'vue'
import { createI18n } from 'vue-i18n'

const defaultOptions = { locale: 'en', fallbackLocale: 'en' }

var i18n = null;

export async function setupI18n(){
	const _i18n = createI18n({
		locale: 'en_us',
		fallbackLocale: 'en_us',
		messages: {
			'en_us': await import('../i18n/en_us.json'),
			'zh_cn': await import('../i18n/zh_cn.json'),
		}
	})
	i18n = _i18n
	console.debug('i18n:', i18n)
	setI18nLanguage(_i18n, 'en_us')
	return _i18n
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

export function getI18n(){
	return i18n
}
