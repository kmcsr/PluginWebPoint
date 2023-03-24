
const linkPattern = /<(https?:\/\/\S+(?:\.\S+)*\b\S*|mailto:\S+@\S+(?:\.\S+)*)>/
const namedLinkPattern = /\[([^\[\]]+)\]\((https?:\/\/\S+(?:\.\S+)*\b\S*|mailto:\S+@\S+(?:\.\S+)*)\)/

export const nodeToString = (function(){
	const box = document.createElement('div')
	if('outerHTML' in box){
		return function _nodeToString(node){
			return node.outerHTML
		}
	}
	return function _nodeToString(node){
		box.replaceChildren(node)
		return box.innerHTML
	}
})()

export function escapeHtml(content){
	return content
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
		output += escapeHtml(content.substring(0, group.index))
		content = content.substring(group.index + full.length)
		let obj = document.createElement('a')
		obj.innerText = name
		obj.href = target
		obj.rel = 'nofollow'
		output += nodeToString(obj)
	}
	if(content){
		output += escapeHtml(content)
	}
	return output
}

export function parse(content){
	return parseLinks(content)
}

export default {
	parse,
}
