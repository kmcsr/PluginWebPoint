
function mustGetEnv(name){
	let value = process.env[name]
	if(typeof value === 'undefined'){
		console.error(`ENV "${name}" is not defined`)
		process.exit(1)
	}
	return value
}

const PRODUCTION = process.env.NODE_ENV === 'production'

export const clientSidePrefix = PRODUCTION ?'/v1' :'/dev'
export const serverSidePrefix = import.meta.env.SSR
	?(PRODUCTION ?mustGetEnv('API_V1_HOST') :mustGetEnv('API_DEV_HOST'))
	:null

if(import.meta.env.SSR){
	console.log('Server side API prefix:', serverSidePrefix)
}

export const prefix = import.meta.env.SSR ?serverSidePrefix :clientSidePrefix
