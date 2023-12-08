import { nextTick, watch } from 'vue'
import { createI18n } from 'vue-i18n'
import { LANG_COOKIE, i18nLangMap } from './consts'

export { LANG_COOKIE, i18nLangMap }

const defaultOptions = { locale: 'en', fallbackLocale: 'en' }

var i18n = null;

const onClient = typeof document !== 'undefined'

export async function setupI18n(lang){
	const cookies = onClient ?await import('vue-cookies') :null
	if(!lang && cookies){
		lang = cookies.get(LANG_COOKIE)
	}
	if(!lang){
		for(let l of window.navigator.languages){
			l = l.toLowerCase().replace('-', '_')
			if(l in i18nLangMap){
				lang = l
				break
			}
		}
		if(!lang){
			lang = 'en_us'
		}
		if(onClient){
			cookies.set(LANG_COOKIE, lang, '30d')
		}
	}
	console.debug('cached lang:', lang)
	const _i18n = createI18n({
		legacy: false,
		fallbackLocale: 'en_us',
		messages: {
			'en_us': await import('../i18n/en_us.json'),
			'zh_cn': await import('../i18n/zh_cn.json'),
		}
	})
	if(onClient){
		watch(_i18n.global.locale, (v) => {
			console.debug('setting cookie:', v)
			cookies.set(LANG_COOKIE, v, '30d')
		})
	}
	i18n = _i18n
	setI18nLanguage(lang)
	return _i18n
}

export function setI18nLanguage(locale){
	if(!i18n.global.availableLocales.includes(locale)){
		locale = i18n.global.fallbackLocale.value
	}
	i18n.global.locale.value = locale
	if(onClient){
		document.getElementsByTagName('html')[0].setAttribute('lang', locale)
	}
}

export async function loadLocaleMessages(locale){
	const messages = await import(`../i18n/${locale}.json`)
	i18n.global.setLocaleMessage(locale, messages.default)
	return nextTick()
}

export function getI18n(){
	return i18n
}
