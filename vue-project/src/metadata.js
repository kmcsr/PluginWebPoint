
function getOrCreateMeta(metas, name, lang){
	for(let m of metas){
		if(m.name === name && (!lang || m.lang === lang)){
			return m
		}
	}
	var meta = document.createElement('meta')
	meta.name = name
	if(lang){
		meta.lang = lang
	}
	meta.content = ''
	document.head.appendChild(meta)
	return meta
}

function getMetaWithData(metas, obj){
	for(let m of metas){
		if(m.name === obj.name && m.content === obj.content){
			return m
		}
	}
	return null
}

export function setMetadata({
	title=null,
	keywords=null,
	replaceKeywords=false,
	description=null,
	extras=null,
}){
	const metas = document.head.getElementsByTagName('meta')
	var oldMeta = {
		title: null,
		keywords: {},
		description: {},
		cleaners: [],
		unmount(){
			if(title){
				document.title = oldMeta.title
			}
			for(let [node, content] of Object.values(oldMeta.keywords)){
				node.content = content
			}
			for(let [node, content] of Object.values(oldMeta.description)){
				node.content = content
			}
			for(let cleaner of oldMeta.cleaners){
				cleaner()
			}
		}
	}
	if(title){
		oldMeta.title = document.title
		document.title = title + ' - PluginWebPoint - MCDReforged'
	}
	if(keywords){
		if(Array.isArray(keywords)){
			keywords = {'': keywords}
		}
		for(let [lang, keyw] of Object.entries(keywords)){
			let mkeyw = getOrCreateMeta(metas, 'keywords', lang)
			oldMeta.keywords[lang] = [mkeyw, mkeyw.content]
			if(replaceKeywords){
				mkeyw.content = keyw.join(',')
			}else{
				let kwArr = mkeyw.content.split(',')
				for(let k of keyw){
					if(kwArr.indexOf(k) == -1){
						kwArr.push(k)
					}
				}
				mkeyw.content = kwArr.join(',')
			}
		}
	}
	if(description){
		if(typeof description === 'string'){
			description = {'': description}
		}
		for(let [lang, desc] of Object.entries(description)){
			let mdesc = getOrCreateMeta(metas, 'description', lang)
			oldMeta.description[lang] = [mdesc, mdesc.content]
			mdesc.content = desc
		}
	}
	if(extras){
		for(let obj of extras){
			if(!getMetaWithData(metas, obj)){
				let meta = document.createElement('meta')
				meta.name = obj.name
				meta.content = obj.content
				document.head.appendChild(meta)
				oldMeta.cleaners.push(() => { document.head.removeChild(meta) })
			}
		}
	}
	return oldMeta
}
