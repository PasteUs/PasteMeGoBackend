# API Documentation

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
  * [Custom temporary key](#custom-temporary-key)
     * [Request](#request-2)
        * [Params](#params-4)
        * [Headers](#headers-3)
        * [Body](#body-3)
           * [Params](#params-5)
           * [Example](#example-3)
     * [Response](#response-2)
        * [Headers](#headers-4)
        * [Body](#body-4)
           * [Params](#params-6)
           * [Example](#example-4)
  * [Error](#error)
     * [Response](#response-3)
        * [Headers](#headers-5)
        * [Body](#body-5)
           * [Params](#params-7)
           * [Example](#example-5)

Table of Contents Created by [gh-md-toc](https://github.com/ekalinin/github-markdown-toc)

## Get

### Request

`GET` `api.pasteme.cn/:key`

#### Params

| Name | Type | Description |
| :---: | :---: | --- |
| key | string | Paste key |

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

#### Headers

`Content-Type: application/json`

#### Body

##### Params

| Name | Type | Description | Required |
| :---: | :---: | --- | :---: |
| lang | string | Paste content's language | Yes |
| content | text | Paste's content | Yes |
| password | string | Paste's password | No |

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
| :---: | :---: | --- |
| status | string | Request status code |
| key | string | Paste's key |

##### Example

```json
{
  "status": 201,
  "key": "100"
}
```

## Custom temporary key

### Request

`PUT` `api.pasteme.cn/:key`

#### Params

| Name | Type | Description |
| :---: | :---: | --- |
| key | string | Paste key |

#### Headers

`Content-Type: application/json`

#### Body

##### Params

| Name | Type | Description | Required |
| :---: | :---: | --- | :---: |
| lang | string | Paste content's language | Yes |
| content | text | Paste's content | Yes |
| password | string | Paste's password | No |

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
| :---: | :---: | --- |
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
| :---: | :---: | --- |
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
