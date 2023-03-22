
[English](./dev.md) | 中文

# Dev API _(Unstable)_

开发API是最新的API, 但是并不稳定. 开发API**不是向后兼容的**, 除了特性预览外, 不建议您的程序使用它. 我们也不保证此文档与实际API始终匹配.

开发API的URL前缀是 `/dev/`, 当您查看下面的路由时, 应始终添加该前缀, 除非有特殊说明.

字符集为 `utf8`

## 错误响应

- 若处理请求时发生了错误, 会返回以下响应体
	- StatusCode: `4xx` (不应该为 `200`)
	- Content-Type: `application/json`
	- 负载:
		```js
		{
			"status": "error", // 应该永远为 "error" 字符串常量
			"error": String, // 错误名称
			"message": String, // 错误信息
		}
		```

## `/`

- 描述:
	此入口点仅用于检查该API是否可用
- 请求:
	- Method: `GET`
	- 负载: *None*
- 响应:
	- StatusCode: `200` OK
	- Content-Type: `application/json`
	- 负载:
		```js
		{
			"status": "ok",
			"time": String,
			"version": Number, // 版本号. 非负整数
		}
		```

## `/plugins/`

- 描述:
	使用指定的过滤器从数据库中获取插件列表, 并返回排序后的数据
- 请求:
	- Method: `GET`
	- URLParams _(可选, 优先级高于json负载)_:
		- `filterBy`: 文字过滤器, 支持以下几种前缀:
			`id:` : 匹配插件ID
			`n:`, `name:` : 匹配插件名称
			`@`, `a:`, `author:`, `authors:` : 匹配作者, 使用逗号(`,`)分割不同作者.
			`d:`, `desc:`, `description:` : 匹配描述
			无前缀: 查询以上所有类型
		- `tags`: 过滤标签, 使用逗号(`,`)分割.
			元素应为 `information`, `tool`, `management`, 或 `api` _(case-insensitive)_
		- `sortBy`: 排序方式.
			可能不存在, 为空字符串, 或为: `id`, `name`, `authors`, `createAt`, `lastUpdate`, `downloads`
		- `reversed`: 反向排序
		- `offset`: 从该偏移开始返回插件列表, 用于分页
		- `limit`: 插件数量限制, 用于分页
	- Content-Type: `application/json` 或 *None*
	- Payload _(可选)_:
		```js
		{
			"filterBy": String, // 同上 `filterBy`
			"tags": [String], // 同上 `tags` 但使用列表, 而不是逗号分隔作者
			"sortBy": String, // 同上 `sortBy`
			"reversed": Boolean, // 同上 `reversed`
			"offset": Number, // 一个非负整数, 同上 `offset`
			"limit": Number, // 一个正整数, 同上 `limit`
		}
		```
- 响应:
	- StatusCode: `200` OK
	- Content-Type: `application/json`
	- 负载:
		```js
		{
			"status": "ok",
			"data": [ // 插件列表
				{
					"id": String, // 插件ID
					"name": String, // 插件名称
					"version": String, // 插件目前的版本
					"authors": [String], // 插件作者, 为一个字符串列表
					"desc": String | undefined, // 插件描述, 大部分为英文版本, 可能未定义
					"desc_zhCN": String | undefined, // 插件中文描述, 可能未定义
					"createAt": String, // 插件加入数据库的时间
					"lastUpdate": String, // 插件最后一次发布新版本的时间. 若未定义, 则同上`createAt`
					"labels": { // 插件标签列表
						"information": Boolean | undefined,
						"tool": Boolean | undefined,
						"management": Boolean | undefined,
						"api": Boolean | undefined,
					},
					"downloads": Number, // 插件总下载数量, 由于从github同步, 所以可能会有延迟
					"github_sync": Boolean, // 插件数据是否是从Github仓库同步而来
					"last_sync": String | undefined, // 插件最后一次从Github同步的时间
				}
			]
		}
		```

## `/plugins/count`

- 描述:
	使用过滤器获取插件数量
- 请求:
	- Method: `GET`
	- URLParams: 同上 `/plugins`
	- 负载: 同上 `/plugins`
- 响应:
	- StatusCode: `200` OK
	- Content-Type: `application/json`
	- 负载:
		```js
		{
			"status": "ok",
			"data": {
				"total": Number, // 符合条件的插件总数
				"information": Number, // 符合条件且具有 `information` 标签的插件数量
				"tool": Number, // 符合条件且具有 `tool` 标签的插件数量
				"management": Number, // 符合条件且具有 `management` 标签的插件数量
				"api": Number, // 符合条件且具有 `api` 标签的插件数量
			}
		}
		```

## `/plugin/{id:string}/info`

- 描述:
	使用 `id` 获取指定插件的数据
- 请求:
	- Method: `GET`
	- 负载: *None*
- 响应:
	- StatusCode: `200` OK, `404` if plugin not found
	- Content-Type: `application/json`
	- 负载:
		```js
		{
			"status": "ok",
			"data": {
				"id": String, // 插件ID
				"name": String, // 插件名称
				"version": String, // 插件目前的版本
				"authors": [String], // 插件作者, 为一个字符串列表
				"desc": String | undefined, // 插件描述, 大部分为英文版本, 可能未定义
				"desc_zhCN": String | undefined, // 插件中文描述, 可能未定义
				"createAt": String, // 插件加入数据库的时间
				"lastUpdate": String, // 插件最后一次发布新版本的时间. 若未定义, 则同上`createAt`
				"repo": String, // Github仓库地址
				"link": String, // Github插件主页
				"labels": { // 插件标签列表
					"information": Boolean | undefined,
					"tool": Boolean | undefined,
					"management": Boolean | undefined,
					"api": Boolean | undefined,
				},
				"downloads": Number, // 插件总下载数量, 由于从github同步, 所以可能会有延迟
				"dependencies": { // 插件依赖列表
					"<plugin id>": "<version condition>", // 见 <https://mcdreforged.readthedocs.io/en/latest/plugin_dev/metadata.html#dependencies>
				},
				"requirements": { // Python包依赖列表
					"<package name>": "<version condition>",
				},
				"github_sync": Boolean, // 插件数据是否是从Github仓库同步而来
				"ghRepoOwner": String | undefined, // Github仓库所有者
				"ghRepoName": String | undefined, // Github仓库名称
				"last_sync": String | undefined, // 插件最后一次从Github同步的时间
			}
		}
		```

## `/plugin/{id:string}/readme`

- Description:
	获取插件README文件
- Request:
	- Method: `GET`
	- URLParams:
		`render`: Boolean. 将README从markdown格式转换为html. (默认: false)
	- Payload: *None*
- Response:
	- StatusCode: `200` OK
	- Content-Type: `text/plain`, `text/html`, `text/markdown`
	- Payload: *README文件内容*

## `/plugin/{id:string}/releases`

- 描述:
	获取插件发布信息
- 请求:
	- Method: `GET`
	- 负载: *None*
- 响应:
	- StatusCode: `200` OK, `404` if plugin not found
	- Content-Type: `application/json`
	- 负载:
		```js
		{
			"status": "ok",
			"data": [ // 发布文件列表, 见下 `/plugin/{id:string}/release/{tag:string}/`
			]
		}
		```


## `/plugin/{id:string}/release/{tag:string}/`

- 描述:
	使用`tag`获取插件指定发布信息
- 请求:
	- Method: `GET`
	- 负载: *None*
- 响应:
	- StatusCode: `200` OK, `404` if release not found
	- Content-Type: `application/json`
	- 负载:
		```js
		{
			"status": "ok",
			"data": {
				"id": String, // 插件ID
				"tag": String, // 发布的插件版本
				"enabled": Boolean, // 此发布是否启用, 目前没有明确定义
				"stable": Boolean, // 此发布是否为稳定版本. 如果值为`false`代表是一个预览发布
				"size": Number, // 发布的文件大小
				"uploaded": String, // 发布时间
				"filename": String, // 发布的文件名称
				"downloads": Number, // 发布文件下载次数
				"github_url": String, // Github下载链接
			}
		}
		```

## `/plugin/{id:string}/release/{tag:string}/asset`

- 描述:
	下载指定的插件发布文件
- 请求:
	- Method: `GET`
	- 负载: *None*
- 响应:
	- StatusCode: `200` OK
	- Content-Type: `*/*`
	- 负载: 该发布文件

