
export const clientSidePrefix = process.env.NODE_ENV === 'production' ?'/v1' :'/dev'
export const serverSidePrefix = import.meta.env.SSR ?process.env.API_V1_HOST :null

if(import.meta.env.SSR){
	if(!serverSidePrefix){
		console.error('ENV: API_V1_HOST is not defined')
		process.exit(1)
	}
	console.log('Server side API prefix:', serverSidePrefix)
}

export const prefix = import.meta.env.SSR ?serverSidePrefix :clientSidePrefix
