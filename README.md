# zues-lua

多语言通用组件库 - lua版

### 接入方式
#### 方式一: 通过git submodule安装
1. 支持通过git submodule的方式加载框架代码。项目根目录下创建`.gitmodules`文件 内容如下:
```
[submodule "vendor/zues"]
	path = vendor/zues
	url = ssh://git@git.intra.weibo.com:2222/wbsearch/zues-lua.git
```
2. git submodule init && git submodule update
3. 执行完git命令后,自动在项目根目录下创建/vendor/zues 并在目录下拉取框架代码

#### 方式二: 通过归档文件安装
1. 下载归档后的框架文件放置于项目根目录下的vendor目录下即可

#### 修改nginx配置
1. 通过`lua_package_path`指令添加lua include目录 `lua_package_path "app-path/vendor/?.lua;;";`
2. 通过`lua_package_cpath`指令添加so文件加载路面 `lua_package_cpath "app-path/vendor/libc/?.so;;";`

#### 组件配置
组件配置文件默认`app_path/conf/component.lua`;
配置信息：
```
{
    -- 缓存组建配置 组件ID cache
    cache = {
        class = 'zues.component.cache.NgxShareCache', //组件对应类 php: zues\component\cache\NgxShareCache
        attrs = {//创建对象时初始化的对象属性
            cacheInstance = 'my_cache',
            keyPrefix = 'was_',
        },
    },
    -- api管理组建
    apiManager = {
        class = 'zues.component.apiManager.CacheApi',
        attrs = {
            tokenFile = '/config/tauth_token_file',
            apisConfig = 'api',
            enableCache = true, -- 是否启动缓存 只有当此参数为true 并且配置了有效的cache字段时，才会真正启动缓存功能
            cache = 'cache', -- 缓存实例
            --degradeFunc = tool_common.degrade,
        },
    },
}
```

#### 加载组件
```lua
init_by_lua_block {
    ZUES_MANAGER = require('zues.component.Manager'):new('component')

    function get_manager() return ZUES_MANAGER end
}
```

#### 使用
```lua
   get_manager():get('cache'):get('dns_resolve')
```