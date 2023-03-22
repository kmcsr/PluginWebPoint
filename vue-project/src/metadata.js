
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

export function setMetadata({
	title=null,
	keywords=null,
	replaceKeywords=false,
	description=null,
}){
	const metas = document.head.getElementsByTagName('meta')
	var olds = {
		title: null,
		keywords: {},
		description: {},
	}
	if(title){
		olds.title = document.title
		document.title = title
	}
	if(keywords){
		if(Array.isArray(keywords)){
			keywords = {'': keywords}
		}
		for(let [lang, keyw] of Object.entries(keywords)){
			let mkeyw = getOrCreateMeta(metas, 'keywords', lang)
			olds.keywords[lang] = [mkeyw, mkeyw.content]
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
			olds.description[lang] = [mdesc, mdesc.content]
			mdesc.content = desc
		}
	}
	return {
		oldMeta: olds,
		unmount(){
			if(title){
				document.title = olds.title
			}
			for(let [node, content] of Object.values(olds.keywords)){
				node.content = content
			}
			for(let [node, content] of Object.values(olds.description)){
				node.content = content
			}
		}
	}
}
