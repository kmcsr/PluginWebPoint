
export function fmtSize(size){
	let unit = 'B'
	if(size > 1024){
		size /= 1024
		unit = 'KB'
		if(size > 1024){
			size /= 1024
			unit = 'MB'
			if(size > 1024){
				size /= 1024
				unit = 'GB'
				if(size > 1024){
					size /= 1024
					unit = 'TB'
				}
			}
		}
	}
	return (+size.toFixed(2)) + unit
}

export function fmtTimestamp(ts, n){
	if(n === undefined){
		n = 2;
	}
	let neg = ts < 0
	if(neg){
		ts = -ts
	}
	let unit = 'ms'
	if(ts > 1000){
		ts /= 1000
		unit = 's'
		if(ts > 60){
			ts /= 60
			unit = 'min'
			if(ts > 60){
				ts /= 60
				unit = 'h'
				if(ts > 24){
					ts /= 24
					unit = 'd'
				}
			}
		}
	}
	let res = (+ts.toFixed(n)) + unit
	if(neg){
		res = '-' + res
	}
	return res
}

export function sinceDate(date){
	return new Date() - new Date(date)
}

export function parseUintStrict(value){
	if(!value.length){
		return NaN
	}
	let number = 0
	for(let n of value){
		if(n < '0' || '9' < n){
			return NaN
		}
		number = number * 10 + Number.parseInt(n)
	}
	return number
}

export function parseIntStrict(value){
	if(!value.length){
		return NaN
	}
	let neg = value[0] == '-'
	if(neg){
		value = value.substring(1)
	}
	let number = parseUintStrict(value)
	if(Number.isNaN(number)){
		return NaN
	}
	return neg ?-number :number
}

function formatStr(fmt, value){
	let leftJustify = fmt[0] === '-'
	if(leftJustify){
		fmt = fmt.substring(1)
	}
	let fixedi = fmt.indexOf('.')
	if(fixedi < 0){
		fixedi = undefined
	}
	let minLen = fixedi === 0 ?0 :parseUintStrict(fmt.substring(0, fixedi))
	if(Number.isNaN(minLen)){
		return null
	}
	let fixedLen = fixedi === undefined ?null :parseUintStrict(fmt.substring(fixedi + 1))
	value = value.toString()
	let fill = minLen - value.length
	if(fill > 0){
		fill = ' '.repeat(fill)
		if(leftJustify){
			value += fill
		}else{
			value = fill + value
		}
	}
	if(value.length > fixedLen){
		value = value.substr(0, fixedLen)
	}
	return value
}

function formatInt(fmt, value){
	let leftJustify = fmt[0] === '-'
	if(leftJustify){
		fmt = fmt.substring(1)
	}
	let sign = fmt[0] === '+'
	if(sign){
		fmt = fmt.substring(1)
	}
	let fillZero = fmt[0] === '0'
	if(fillZero){
		fmt = fmt.substring(1)
	}
	let minLen = parseUintStrict(fmt)
	if(Number.isNaN(minLen)){
		return null
	}
	value = value.toString()
	if(sign && value[0] !== '-' && value[0] !== '+'){
		value = '+' + value
	}
	let fill = minLen - value.length
	if(fill > 0){
		fill = (fillZero ?'0' :' ').repeat(fill)
		if(leftJustify){
			value += fill
		}else{
			value = fill + value
		}
	}
	return value
}

function formatFloat(fmt, value){
	let leftJustify = fmt[0] === '-'
	if(leftJustify){
		fmt = fmt.substring(1)
	}
	let sign = fmt[0] === '+'
	if(sign){
		fmt = fmt.substring(1)
	}
	let fillZero = fmt[0] === '0'
	if(fillZero){
		fmt = fmt.substring(1)
	}
	let fixedi = fmt.indexOf('.')
	if(fixedi < 0){
		fixedi = undefined
	}
	let minLen = fixedi === 0 ?0 :parseUintStrict(fmt.substring(0, fixedi))
	if(Number.isNaN(minLen)){
		return null
	}
	if(fixedi !== undefined){
		let fixedLen = parseUintStrict(fmt.substring(fixedi + 1))
		value = value.toFixed(fixedLen)
	}
	value = value.toString()
	if(sign && value[0] !== '-' && value[0] !== '+'){
		value = '+' + value
	}
	let fill = minLen - value.length
	if(fill > 0){
		fill = (fillZero ?'0' :' ').repeat(fill)
		if(leftJustify){
			value += fill
		}else{
			value = fill + value
		}
	}
	return value
}

function format1(ch, fmt, ind, vars){
	{
		let i = fmt.indexOf('$')
		if(i >= 0 && !Number.isNaN(i = parseIntStrict(fmt.substring(0, i)))){
			ind = i
		}
	}
	let value = vars[ind];
	let out = null;
	switch(ch){
	case 's':
		out = formatStr(fmt, value)
		break
	case 'd':
		out = formatInt(fmt, value)
		break
	case 'f':
		out = formatFloat(fmt, value)
		break
	}
	return [out, ind]
}

export function format(fmt, ...vars){
	let output = ''
	let ind = 0
	for(let i = 0; i < fmt.length; i++){
		let c = fmt[i]
		if(c != '%'){
			output += c
			continue
		}
		let s = i
		while(++i < fmt.length && ((c = fmt[i]) < 'a' || 'z' < c) && (c < 'A' || 'Z' < c));
		if(i >= fmt.length){
			output += fmt.substring(s)
			break
		}
		let out
		[out, ind] = format1(c, fmt.substring(s + 1, i), ind, vars)
		if(out !== null){
			output += out
		}else{
			output += fmt.substring(s, i + 1)
		}
		ind++
	}
	return output
}

export function fmtDateTime(date){
	let d = new Date(date);
	return format('%04d-%02d-%02d %02d:%02d:%02d',
		d.getUTCFullYear(), d.getUTCMonth(), d.getUTCDate(),
		d.getUTCHours(), d.getUTCMinutes(), d.getUTCSeconds())
}
