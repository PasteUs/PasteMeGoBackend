Table of Contents
=================

   * [API Documentation](#api-documentation)
      * [Get](#get)
         * [Request](#request)
            * [Params](#params)
         * [Response](#response)
            * [Headers](#headers)
            * [Body](#body)
               * [Params](#params-1)
               * [Example](#example)
      * [Create](#create)
         * [Request](#request-1)
            * [Headers](#headers-1)
            * [Body](#body-1)
               * [Params](#params-2)
               * [Example](#example-1)
         * [Response](#response-1)
            * [Headers](#headers-2)
            * [Body](#body-2)
               * [Params](#params-3)
               * [Example](#example-2)
      * [Error](#error)
         * [Response](#response-2)
            * [Headers](#headers-3)
            * [Body](#body-3)
               * [Params](#params-4)
               * [Example](#example-3)
   * [Table of Contents](#table-of-contents)

Table of Contents Created by [gh-md-toc](https://github.com/ekalinin/github-markdown-toc)

# API Documentation

## Get

### Request

`GET` `api.pasteme.cn/:key`

#### Params

| Name | Type | Description |
| --- | --- | --- |
| key | string | Paste ID |

### Response 

#### Headers

`Content-Type: application/json`

#### Body

##### Params

| Name | Type | Description |
| :---: | :---: | --- |
| status | int | request status code |
| lang | string | Paste content's language |
| content | text | Paste content |

##### Example

```json
{
  "status": 200,
  "lang": "bash",
  "content": "echo Hello"
}
```

## Create

### Request

| Method | URL | Description |
| :---: | --- | --- |
| `POST` | `api.pasteme.cn` | Create a permanent paste |
| `POST` | `api.pasteme.cn/once` | Create a temporary paste with random key |
| `PUT` | `api.pasteme.cn` | Create a temporary paste with specific key |

#### Headers

`Content-Type: application/json`

#### Body

##### Params

| Name | Type | Description |
| --- | --- | --- |
| lang | string | Paste content's language |
| content | text | Paste's content |
| password | string | Paste's password |

##### Example

```json
{
  "lang": "bash",
  "content": "echo Hello"
}
```

or

```json
{
  "lang": "bash",
  "content": "echo Hello",
  "password": "password"
}
```

### Response 

#### Headers

`Content-Type: application/json`

#### Body

##### Params

| Name | Type | Description |
| --- | --- | --- |
| status | string | Request status code |
| key | string | Paste's key |

##### Example

```json
{
  "status": 201,
  "key": "100"
}
```

## Error

### Response 

#### Headers

`Content-Type: application/json`

#### Body

##### Params

| Name | Type | Description |
| --- | --- | --- |
| status | int | Request status code |
| error | string | Error content |
| message | string | Error's detail |

##### Example

```json
{
  "status": 401,
  "error":   "unauthorized",
  "message": "wrong password"
}
```
