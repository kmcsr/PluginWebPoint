import { nextTick, watch } from 'vue'
import VueCookies from 'vue-cookies'
import { createI18n } from 'vue-i18n'

const defaultOptions = { locale: 'en', fallbackLocale: 'en' }

var i18n = null;

export async function setupI18n(){
	var lang = VueCookies.get('lang')
	if(!lang){
		$cookies.set('lang', 'en_us', '30d')
		lang = 'en_us'
	}
	console.log('cached lang:', lang)
	const _i18n = createI18n({
		legacy: false,
		fallbackLocale: 'en_us',
		messages: {
			'en_us': await import('../i18n/en_us.json'),
			'zh_cn': await import('../i18n/zh_cn.json'),
		}
	})
	watch(_i18n.global.locale, (v) => {
		console.log('setting cookie:', v)
		$cookies.set('lang', v, '30d')
	})
	i18n = _i18n
	setI18nLanguage(lang)
	return _i18n
}

export function setI18nLanguage(locale){
	if(!i18n.global.availableLocales.includes(locale)){
		locale = i18n.global.fallbackLocale.value
	}
	i18n.global.locale.value = locale
	document.querySelector('html').setAttribute('lang', locale)
}

export async function loadLocaleMessages(locale){
	const messages = await import(`../i18n/${locale}.json`)
	i18n.global.setLocaleMessage(locale, messages.default)
	return nextTick()
}

export function getI18n(){
	return i18n
}
