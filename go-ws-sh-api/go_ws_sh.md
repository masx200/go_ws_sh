---
title: go_ws_sh
language_tabs:
  - shell: Shell
  - http: HTTP
  - javascript: JavaScript
  - ruby: Ruby
  - python: Python
  - php: PHP
  - java: Java
  - go: Go
toc_footers: []
includes: []
search: true
code_clipboard: true
highlight_theme: darkula
headingLevel: 2
generator: "@tarslib/widdershins v4.0.29"

---

# go_ws_sh

Base URLs:

# Authentication

# Default

## POST 创建令牌

POST /tokens

> Body 请求参数

```json
{
  "authorization": {
    "identifer": "59",
    "username": "通鹏",
    "token": "dolore Lorem ut",
    "password": "u1H84FL1dygv6H7",
    "type": "do veniam consectetur aliqua amet"
  },
  "token": {
    "description": "联太为养复近造只手适。展志斯。子世年存种行成车交。决么林快效不原回团。老度史半在命质。术者照会接。是导响飞产提论我育定。",
    "username": "束艳",
    "identifer": "59",
    "token": "Excepteur commodo enim ullamco ex"
  }
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|x-HTTP-method-override|header|string| 否 |POST|
|body|body|object| 否 |none|
|» authorization|body|[authorization](#schemaauthorization)| 是 |none|
|»» identifer|body|string| 是 |none|
|»» username|body|string| 是 |none|
|»» token|body|string| 是 |none|
|»» password|body|string| 是 |none|
|»» type|body|string| 是 |none|
|» token|body|[token](#schematoken)| 是 |none|
|»» description|body|string| 是 |none|
|»» username|body|string| 是 |none|
|»» identifer|body|string| 是 |none|
|»» token|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "message": "laborum id",
  "token": {
    "description": "常应儿到统型整就这。本务记子权方。计选增酸单带断。况样南起据报题象门。然法深六认委两量。山事体。交满据至被了。",
    "username": "蒲梓豪",
    "identifer": "68",
    "token": "laborum nulla nostrud"
  },
  "username": "空紫林"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» token|[token](#schematoken)|true|none||none|
|»» description|string|true|none||none|
|»» username|string|true|none||none|
|»» identifer|string|true|none||none|
|»» token|string|true|none||none|
|» username|string|true|none||none|

## PUT 修改令牌

PUT /tokens

> Body 请求参数

```json
{
  "token": {
    "identifer": "37",
    "description": "细矿义。合认便情用南别法重同。解思周。把话老权道选带全百。院条实际展。提确备最广具术去包制。果流口布正大。温你任人然算。"
  },
  "authorization": {
    "identifer": "2",
    "username": "表志明",
    "token": "id sit",
    "password": "rizLm0jfLMpJoCE",
    "type": "velit Lorem voluptate in"
  }
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» token|body|object| 是 |none|
|»» identifer|body|string| 是 |none|
|»» description|body|string| 是 |none|
|» authorization|body|[authorization](#schemaauthorization)| 是 |none|
|»» identifer|body|string| 是 |none|
|»» username|body|string| 是 |none|
|»» token|body|string| 是 |none|
|»» password|body|string| 是 |none|
|»» type|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "message": "sed deserunt ut",
  "username": "及志国",
  "token": {
    "description": "物除支斗保了研况。合和世米。教素华老代石究。",
    "username": "冀梓豪",
    "identifer": "44"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» username|string|true|none||none|
|» token|[token](#schematoken)|true|none||none|
|»» description|string|true|none||none|
|»» username|string|true|none||none|
|»» identifer|string|true|none||none|
|»» token|string|true|none||none|

## DELETE 删除令牌

DELETE /tokens

> Body 请求参数

```json
{
  "token": {
    "description": "也难各况老须山办何线。酸线给名。好东声间候下制备。受你质白东本处次较命。示难这名组合。专市该资。",
    "username": "始国强",
    "identifer": "59",
    "token": "in ipsum"
  },
  "authorization": {
    "identifer": "98",
    "username": "斋国珍",
    "token": "reprehenderit sit id",
    "password": "hkjNbMjFHCBrkbx",
    "type": "sunt officia nulla ut ut"
  }
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» token|body|[token](#schematoken)| 是 |none|
|»» description|body|string| 是 |none|
|»» username|body|string| 是 |none|
|»» identifer|body|string| 是 |none|
|»» token|body|string| 是 |none|
|» authorization|body|[authorization](#schemaauthorization)| 是 |none|
|»» identifer|body|string| 是 |none|
|»» username|body|string| 是 |none|
|»» token|body|string| 是 |none|
|»» password|body|string| 是 |none|
|»» type|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "message": "in labore aliqua",
  "username": "所丽芳",
  "token": {
    "description": "形争多历必去包队。正管铁记样。打必手京华从他者年这。",
    "username": "卞国良",
    "identifer": "7",
    "token": "et amet nostrud"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» username|string|true|none||none|
|» token|[token](#schematoken)|true|none||none|
|»» description|string|true|none||none|
|»» username|string|true|none||none|
|»» identifer|string|true|none||none|
|»» token|string|true|none||none|

## GET 显示令牌

GET /tokens

> Body 请求参数

```json
{
  "authorization": {
    "identifer": "26",
    "username": "希颖",
    "token": "cillum ad",
    "password": "dzIi7wxDnv2KLGS",
    "type": "elit proident non in"
  }
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|x-HTTP-method-override|header|string| 是 |GET|
|body|body|object| 否 |none|
|» authorization|body|[authorization](#schemaauthorization)| 是 |none|
|»» identifer|body|string| 是 |none|
|»» username|body|string| 是 |none|
|»» token|body|string| 是 |none|
|»» password|body|string| 是 |none|
|»» type|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "message": "qui",
  "username": "马超栋",
  "tokens": [
    {
      "description": "支起天器步管同日处。上单两率利元对学。队加议广及我效光标。类看进。门展科不写义县次该她。次斗变。人火车立己进。手机验类重百然。",
      "username": "龙家豪",
      "identifer": "68",
      "token": "officia deserunt in est enim"
    },
    {
      "description": "前全过过式行保。和感查。及年性成把。积五无特你及。华第也省。",
      "username": "步建军",
      "identifer": "58",
      "token": "proident"
    }
  ]
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» username|string|true|none||none|
|» tokens|[[token](#schematoken)]|true|none||none|
|»» description|string|true|none||none|
|»» username|string|true|none||none|
|»» identifer|string|true|none||none|
|»» token|string|true|none||none|

## PUT 修改密码

PUT /credentials

> Body 请求参数

```json
{
  "authorization": {
    "identifer": "11",
    "username": "京万佳",
    "token": "nisi proident",
    "password": "HQHsYZOyvSkcizB",
    "type": "velit"
  },
  "credential": {
    "username": "明苡沫",
    "password": "ZpWWRl4fkBCmcZG"
  }
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» authorization|body|[authorization](#schemaauthorization)| 是 |none|
|»» identifer|body|string| 是 |none|
|»» username|body|string| 是 |none|
|»» token|body|string| 是 |none|
|»» password|body|string| 是 |none|
|»» type|body|string| 是 |none|
|» credential|body|[credential](#schemacredential)| 是 |none|
|»» username|body|string| 是 |none|
|»» password|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "message": "id ut proident Ut",
  "username": "隗秀珍",
  "credential": {
    "username": "羊国华",
    "password": "GRTGxa6qmWPT8y5"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» username|string|true|none||none|
|» credential|[credential](#schemacredential)|true|none||none|
|»» username|string|true|none||none|
|»» password|string|true|none||none|

## DELETE 删除用户

DELETE /credentials

> Body 请求参数

```json
{
  "authorization": {
    "identifer": "29",
    "username": "阎晨",
    "token": "eu",
    "password": "Ax1RzKvcJA7FQZj",
    "type": "minim fugiat"
  },
  "credential": {
    "username": "衡刚",
    "password": "xvJw7L2YnmNOjag"
  }
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» authorization|body|[authorization](#schemaauthorization)| 是 |none|
|»» identifer|body|string| 是 |none|
|»» username|body|string| 是 |none|
|»» token|body|string| 是 |none|
|»» password|body|string| 是 |none|
|»» type|body|string| 是 |none|
|» credential|body|[credential](#schemacredential)| 是 |none|
|»» username|body|string| 是 |none|
|»» password|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "message": "dolore",
  "username": "百丽",
  "credential": {
    "username": "恭玲",
    "password": "viPAfChquGKeQ6K"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» username|string|true|none||none|
|» credential|[credential](#schemacredential)|true|none||none|
|»» username|string|true|none||none|
|»» password|string|true|none||none|

## GET 显示用户

GET /credentials

> Body 请求参数

```json
{
  "authorization": {
    "identifer": "55",
    "username": "姒静怡",
    "token": "enim pariatur culpa",
    "password": "q2fFVw6c6uY_gwQ",
    "type": "non dolor"
  }
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|x-HTTP-method-override|header|string| 是 |GET|
|body|body|object| 否 |none|
|» authorization|body|[authorization](#schemaauthorization)| 是 |none|
|»» identifer|body|string| 是 |none|
|»» username|body|string| 是 |none|
|»» token|body|string| 是 |none|
|»» password|body|string| 是 |none|
|»» type|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "credentials": [
    {
      "username": "仵晨",
      "password": "sEyDDGTf5uIBI3G"
    },
    {
      "username": "寻子豪",
      "password": "HCOvPlHMWA1vBU_"
    },
    {
      "username": "丁国华",
      "password": "XJ08qHeYUGFOLV5"
    }
  ],
  "message": "sed adipisicing dolore eiusmod",
  "username": "敛鹏"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» credentials|[[credential](#schemacredential)]|true|none||none|
|»» username|string|true|none||none|
|»» password|string|true|none||none|
|» message|string|true|none||none|
|» username|string|true|none||none|

## POST 创建用户

POST /credentials

> Body 请求参数

```json
{
  "authorization": {
    "identifer": "8",
    "username": "苌雨桐",
    "token": "dolor incididunt pariatur",
    "password": "SmuJxITZXdNYS__",
    "type": "non"
  },
  "credential": {
    "username": "僪语桐",
    "password": "zXWuNkm6lQJABvj"
  }
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|x-HTTP-method-override|header|string| 否 |POST|
|body|body|object| 否 |none|
|» authorization|body|[authorization](#schemaauthorization)| 是 |none|
|»» identifer|body|string| 是 |none|
|»» username|body|string| 是 |none|
|»» token|body|string| 是 |none|
|»» password|body|string| 是 |none|
|»» type|body|string| 是 |none|
|» credential|body|[credential](#schemacredential)| 是 |none|
|»» username|body|string| 是 |none|
|»» password|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "credential": {
    "username": "窦宇航",
    "password": "l4toCzQ3SzXMK2Q"
  },
  "message": "proident irure commodo ipsum laborum",
  "username": "本政君"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» credential|[credential](#schemacredential)|true|none||none|
|»» username|string|true|none||none|
|»» password|string|true|none||none|
|» message|string|true|none||none|
|» username|string|true|none||none|

## PUT 修改会话

PUT /sessions

> Body 请求参数

```json
{
  "authorization": {
    "identifer": "9",
    "username": "郦桂兰",
    "token": "ad nisi fugiat eu nulla",
    "password": "vhkN6hb1ZubYJw5",
    "type": "consectetur enim"
  },
  "session": {
    "name": "盈雯静",
    "args": [
      "dolor Lorem eu aliquip",
      "fugiat dolor",
      "irure aute commodo magna"
    ],
    "dir": "laboris elit",
    "cmd": "laboris sunt sint"
  }
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» authorization|body|[authorization](#schemaauthorization)| 是 |none|
|»» identifer|body|string| 是 |none|
|»» username|body|string| 是 |none|
|»» token|body|string| 是 |none|
|»» password|body|string| 是 |none|
|»» type|body|string| 是 |none|
|» session|body|[session](#schemasession)| 是 |none|
|»» name|body|string| 否 |none|
|»» args|body|[string]| 否 |none|
|»» dir|body|string| 否 |none|
|»» cmd|body|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "message": "elit eu nostrud",
  "username": "镇馨羽",
  "session": {
    "name": "斯娟",
    "cmd": "culpa occaecat est",
    "dir": "amet adipisicing dolor non",
    "args": [
      "adipisicing",
      "sed do Lorem Duis"
    ]
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» username|string|true|none||none|
|» session|[session](#schemasession)|true|none||none|
|»» name|string|false|none||none|
|»» args|[string]|false|none||none|
|»» dir|string|false|none||none|
|»» cmd|string|false|none||none|

## DELETE 删除会话

DELETE /sessions

> Body 请求参数

```json
{
  "session": {
    "name": "全丽芳",
    "args": [
      "ullamco culpa",
      "amet aliqua est",
      "quis Ut aliqua ex"
    ],
    "dir": "laboris",
    "cmd": "amet incididunt do culpa labore"
  },
  "authorization": {
    "identifer": "31",
    "username": "龙秀英",
    "token": "nulla dolore officia",
    "password": "EhQLnzz1kTcnyz6",
    "type": "pariatur aute labore enim"
  }
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|body|body|object| 否 |none|
|» session|body|[session](#schemasession)| 是 |none|
|»» name|body|string| 否 |none|
|»» args|body|[string]| 否 |none|
|»» dir|body|string| 否 |none|
|»» cmd|body|string| 否 |none|
|» authorization|body|[authorization](#schemaauthorization)| 是 |none|
|»» identifer|body|string| 是 |none|
|»» username|body|string| 是 |none|
|»» token|body|string| 是 |none|
|»» password|body|string| 是 |none|
|»» type|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "message": "in",
  "username": "容诗雨",
  "session": {
    "name": "繁国英",
    "args": [
      "eiusmod anim sit",
      "nulla exercitation aliqua culpa",
      "in ad velit cupidatat"
    ],
    "dir": "dolor",
    "cmd": "exercitation"
  }
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» message|string|true|none||none|
|» username|string|true|none||none|
|» session|[session](#schemasession)|true|none||none|
|»» name|string|false|none||none|
|»» args|[string]|false|none||none|
|»» dir|string|false|none||none|
|»» cmd|string|false|none||none|

## GET 显示会话

GET /sessions

> Body 请求参数

```json
{
  "authorization": {
    "identifer": "44",
    "username": "晋哲新",
    "token": "ut esse",
    "password": "Yj4QEiNgfJyF3cJ",
    "type": "in sint aliqua qui"
  }
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|x-HTTP-method-override|header|string| 是 |GET|
|body|body|object| 否 |none|
|» authorization|body|[authorization](#schemaauthorization)| 是 |none|
|»» identifer|body|string| 是 |none|
|»» username|body|string| 是 |none|
|»» token|body|string| 是 |none|
|»» password|body|string| 是 |none|
|»» type|body|string| 是 |none|

> 返回示例

> 200 Response

```json
{
  "sessions": [
    {
      "name": "蓬国英",
      "args": [
        "sint reprehenderit exercitation nostrud"
      ],
      "dir": "dolore in",
      "cmd": "anim et"
    },
    {
      "name": "双馥君",
      "args": [
        "irure Duis",
        "labore irure",
        "amet velit fugiat ea"
      ],
      "dir": "ipsum",
      "cmd": "veniam deserunt aliqua"
    }
  ],
  "message": "mollit id Duis magna proident",
  "username": "鱼浩晨"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» sessions|[[session](#schemasession)]|true|none||none|
|»» name|string|false|none||none|
|»» args|[string]|false|none||none|
|»» dir|string|false|none||none|
|»» cmd|string|false|none||none|
|» message|string|true|none||none|
|» username|string|true|none||none|

## POST 创建会话

POST /sessions

> Body 请求参数

```json
{
  "authorization": {
    "identifer": "90",
    "username": "郁勇",
    "token": "do eiusmod",
    "password": "mXpnuFnILrjhuGJ",
    "type": "dolore nulla Ut ut proident"
  },
  "session": {
    "name": "祈雨欣",
    "args": [
      "proident",
      "cillum esse dolore ut amet",
      "laborum sint culpa"
    ],
    "dir": "fugiat dolore ut reprehenderit sed",
    "cmd": "anim"
  }
}
```

### 请求参数

|名称|位置|类型|必选|说明|
|---|---|---|---|---|
|x-HTTP-method-override|header|string| 否 |POST|
|body|body|object| 否 |none|
|» authorization|body|[authorization](#schemaauthorization)| 是 |none|
|»» identifer|body|string| 是 |none|
|»» username|body|string| 是 |none|
|»» token|body|string| 是 |none|
|»» password|body|string| 是 |none|
|»» type|body|string| 是 |none|
|» session|body|[session](#schemasession)| 是 |none|
|»» name|body|string| 否 |none|
|»» args|body|[string]| 否 |none|
|»» dir|body|string| 否 |none|
|»» cmd|body|string| 否 |none|

> 返回示例

> 200 Response

```json
{
  "session": {
    "description": "白住利适广的确压断院。一三就具出管。育程成及她不治都带只。科变利前。被接业把五。形备委都门立理场色资。志立候头消改流需报。条万业前风。",
    "username": "星国平",
    "identifer": "42",
    "token": "Excepteur dolore sit amet"
  },
  "message": "do",
  "username": "夔蒙"
}
```

### 返回结果

|状态码|状态码含义|说明|数据模型|
|---|---|---|---|
|200|[OK](https://tools.ietf.org/html/rfc7231#section-6.3.1)|none|Inline|

### 返回数据结构

状态码 **200**

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|» session|[token](#schematoken)|true|none||none|
|»» description|string|true|none||none|
|»» username|string|true|none||none|
|»» identifer|string|true|none||none|
|»» token|string|true|none||none|
|» message|string|true|none||none|
|» username|string|true|none||none|

# 数据模型

<h2 id="tocS_authorization">authorization</h2>

<a id="schemaauthorization"></a>
<a id="schema_authorization"></a>
<a id="tocSauthorization"></a>
<a id="tocsauthorization"></a>

```json
{
  "identifer": "string",
  "username": "string",
  "token": "string",
  "password": "string",
  "type": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|identifer|string|true|none||none|
|username|string|true|none||none|
|token|string|true|none||none|
|password|string|true|none||none|
|type|string|true|none||none|

<h2 id="tocS_token">token</h2>

<a id="schematoken"></a>
<a id="schema_token"></a>
<a id="tocStoken"></a>
<a id="tocstoken"></a>

```json
{
  "description": "string",
  "username": "string",
  "identifer": "string",
  "token": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|description|string|true|none||none|
|username|string|true|none||none|
|identifer|string|true|none||none|
|token|string|true|none||none|

<h2 id="tocS_session">session</h2>

<a id="schemasession"></a>
<a id="schema_session"></a>
<a id="tocSsession"></a>
<a id="tocssession"></a>

```json
{
  "name": "string",
  "args": [
    "string"
  ],
  "dir": "string",
  "cmd": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|name|string|false|none||none|
|args|[string]|false|none||none|
|dir|string|false|none||none|
|cmd|string|false|none||none|

<h2 id="tocS_credential">credential</h2>

<a id="schemacredential"></a>
<a id="schema_credential"></a>
<a id="tocScredential"></a>
<a id="tocscredential"></a>

```json
{
  "username": "string",
  "password": "string"
}

```

### 属性

|名称|类型|必选|约束|中文名|说明|
|---|---|---|---|---|---|
|username|string|true|none||none|
|password|string|true|none||none|

