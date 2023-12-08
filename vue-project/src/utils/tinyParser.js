
const linkPattern = /<(https?:\/\/[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(?:\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})*?\S*?|mailto:\S+?@[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(?:\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})*?)>/
const namedLinkPattern = /\[([^\[\]]+)\]\((https?:\/\/[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(?:\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})*?\S*?|mailto:\S+?@[a-zA-Z0-9][-a-zA-Z0-9]{0,62}(?:\.[a-zA-Z0-9][-a-zA-Z0-9]{0,62})*?)\)/

const codePatten = /(`+)([^\1]+?)\1/

// export const nodeToString = (function(){
// 	const box = document.createElement('div')
// 	if('outerHTML' in box){
// 		return function _nodeToString(node){
// 			return node.outerHTML
// 		}
// 	}
// 	return function _nodeToString(node){
// 		box.replaceChildren(node)
// 		return box.innerHTML
// 	}
// })()

export function escapeHtml(content){
	return !content ?'' :content
		.replaceAll('&', '&amp;')
		.replaceAll('<', '&lt;')
		.replaceAll('>', '&gt;')
		.replaceAll('"', '&quot;')
		.replaceAll("'", '&#039;')
}

export function parseLinks(content){
	var output = ''
	while(content){
		let full, name, target
		let group = namedLinkPattern.exec(content)
		let group2 = linkPattern.exec(content)
		if(!group && !group2){
			break
		}
		if(group2 && group2.index < group.index){
			[full, target] = group = group2
			name = target
		}else{
			[full, name, target] = group
		}
		output += content.substring(0, group.index)
		content = content.substring(group.index + full.length)
		// let obj = document.createElement('a')
		// obj.innerText = name
		// obj.href = target
		// obj.rel = 'nofollow'
		// output += nodeToString(obj)
		output += `<a href="${escape(target)}" ref="nofollow">${escapeHtml(name)}</a>`
	}
	if(content){
		output += content
	}
	return output
}

export function parseCodes(content){
	var output = ''
	while(content){
		let group = codePatten.exec(content)
		if(!group){
			break
		}
		let [full, _, codes] = group
		output += content.substring(0, group.index)
		content = content.substring(group.index + full.length)
		// let obj = document.createElement('code')
		// obj.innerText = codes
		// output += nodeToString(obj)
		output += `<code>${escapeHtml(codes)}</code>`
	}
	if(content){
		output += content
	}
	return output
}

export function parse(content){
	let out = parseCodes(escapeHtml(content))
	out = parseLinks(out)
	return out
}

export default {
	parse,
}
