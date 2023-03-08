
function fmtSize(size){
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

function fmtTimestamp(ts, n){
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

function sinceDate(date){
	return new Date() - new Date(date)
}

export {
	fmtSize,
	fmtTimestamp,
	sinceDate
}
