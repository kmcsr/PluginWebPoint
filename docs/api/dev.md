
# Dev API _(Unstable)_

Development API is the latest API, but it's unstable. It's **NOT backward compatible**, and we won't ensure this doc is match the API.

The prefix of dev API is `/dev/`, you should add this prefix to every routes below if there are no specific comment.

The charset is `utf8`

## Error response

- If there is an error occur when requesting, the response will be:
	- StatusCode: `4xx` (should never be `200`)
	- Content-Type: `application/json`
	- Payload:
		```js
		{
			"status": "error", // should always be error
			"error": String, // The name of the error
			"message": String, // Message from the error
		}
		```

## `/`

- Description:
	This entry point just for check if the API is available
- Request:
	- Method: `GET`
	- Payload: *None*
- Response:
	- StatusCode: `200` OK
	- Content-Type: `application/json`
	- Payload:
		```js
		{
			"status": "ok", // should always be ok
			"time": String, // Access time in UTC time zone. Formatted
			"version": Number, // The developing version, should be a non-negative integer
		}
		```

## `/plugins/`

- Description:
	Get the plugin list that saved in our database, then filter and sorted by the params you sent.
- Request:
	- Method: `GET`
	- URLParams _(optional, priority is higher than the json payload)_:
		- `filterBy`: The text filter in the search box, support query specific fields by adding case insensitive prefixs:
			`id:` : Match by plugin id
			`n:`, `name:` : Match by plugin name
			`@`, `a:`, `author:`, `authors:` : Match by authors, split authors by comma `,`
			`d:`, `desc:`, `description:` : Match by the short description
			No prefix: Query all of above
		- `tags`: The filter tags, split by comma(`,`).
			Elements should be `information`, `tool`, `management`, or `api` _(case-insensitive)_
		- `sortBy`: Sort by which field.
			Could be None or empty string, `id`, `name`, `authors`, `createAt`, `lastUpdate`, `downloads`
		- `reversed`: Reversed the output
		- `offset`: Return plugins from the offset, use when split page
		- `limit`: The plugin list limit, use when split page
	- Content-Type: `application/json` or *None*
	- Payload _(optional)_:
		```js
		{
			"filterBy": String, // A string, same as URLParams above with name `filterBy`
			"tags": [String], // List of string, same as above `tags` but use string list instead string split with comma
			"sortBy": String, // A string, same as above `sortBy`
			"reversed": Boolean, // A boolean, same as above `reversed`
			"offset": Number, // A positive integer or zero, same as above `offset`
			"limit": Number, // A positive integer, same as above `limit`
		}
		```
- Response:
	- StatusCode: `200` OK
	- Content-Type: `application/json`
	- Payload:
		```js
		{
			"status": "ok", // should be ok
			"data": [ // List of plugin, see below `/plugin/{id:string}/info`
			]
		}
		```

## `/plugins/count`

- Description:
	Get the plugin count that matched the filters
- Request:
	- Method: `GET`
	- URLParams: As same as `/plugins` above
	- Payload: As same as `/plugins` above
- Response:
	- StatusCode: `200` OK
	- Content-Type: `application/json`
	- Payload:
		```js
		{
			"status": "ok", // should be ok
			"data": {
				"total": Number, // total plugin count
				"information": Number, // Count of plugin with tag information
				"tool": Number, // Count of plugin with tag tool
				"management": Number, // Count of plugin with tag management
				"api": Number, // Count of plugin with tag api
			}
		}
		```

## `/plugin/{id:string}/info`

- Description:
	Get the plugin info by `id`
- Request:
	- Method: `GET`
	- Payload: *None*
- Response:
	- StatusCode: `200` OK, `404` if plugin not found
	- Content-Type: `application/json`
	- Payload:
		```js
		{
			"status": "ok", // should be ok
			"data": {
				"id": String, // Plugin's ID
				"name": String, // Plugin's display name
				"version": String, // Plugin's current version
				"authors": [String], // The authors of the plugin, in a string list
				"desc": String | undefined, // The plugin's description, mostly in English, could be none
				"desc_zhCN": String | undefined, // The plugin's description in Chinese, could be none
				"createAt": String, // When is the plugin be added into the database. Formatted
				"lastUpdate": String, // The plugin last update time, not used now. Formatted
				"repo": String, // Github repo link for the plugin
				"link": String, // Github main page link
				"labels": { // is this plugin labeled by the key, maybe undefined
					"information": Boolean | undefined,
					"tool": Boolean | undefined,
					"management": Boolean | undefined,
					"api": Boolean | undefined,
				},
				"downloads": Number, // The total download count of the plugin releases, synced from github, maybe delayed
				"dependencies": { // The dependent map
					"<plugin id>": "<version condition>", // for version condition, please see <https://mcdreforged.readthedocs.io/en/latest/plugin_dev/metadata.html#dependencies>
				},
				"github_sync": Boolean, // Is the plugin synced from github or not
				"last_sync": String | undefined, // Last time it synced and updated from github. Maybe undefined if it's not synced from github. Formatted
			}
		}
		```

## `/plugin/{id:string}/releases`

- Description:
	Get the plugin releases
- Request:
	- Method: `GET`
	- Payload: *None*
- Response:
	- StatusCode: `200` OK, `404` if plugin not found
	- Content-Type: `application/json`
	- Payload:
		```js
		{
			"status": "ok",
			"data": [ // List of releases, see below `/plugin/{id:string}/release/{tag:string}/`
			]
		}
		```


## `/plugin/{id:string}/release/{tag:string}/`

- Description:
	Get the specific release info for the plugin
- Request:
	- Method: `GET`
	- Payload: *None*
- Response:
	- StatusCode: `200` OK, `404` if release not found
	- Content-Type: `application/json`
	- Payload:
		```js
		{
			"status": "ok",
			"data": {
				"id": String, // Plugin's ID
				"tag": String, // The release's version
				"enabled": Boolean, // Is this release enabled, not used
				"stable": Boolean, // Is this release stable. If it's `false` means this release a prerelease
				"size": Number, // The release file size
				"uploaded": String, // Time the release uploaded
				"filename": String, // The filename for the release asset
				"downloads": Number, // The download count for the release asset
				"github_url": String, // The Github download URL for the release
			}
		}
		```

## `/plugin/{id:string}/release/{tag:string}/asset`

- Description:
	Download the plugin release file
- Request:
	- Method: `GET`
	- Payload: *None*
- Response:
	- StatusCode: `200` OK
	- Content-Type: `*/*`
	- Payload: The asset file

