# goadmin（裁剪版 GoAdmin）

[GoAdmin](https://github.com/GoAdminGroup/go-admin) v1.2.26 的本地 fork，给 GM 后台（同级的 `../expgm`）用。模块名 `goadmin`，用的地方在 go.mod 里 `replace goadmin => ../goadmin`（和 pdbgen 一样）。上游许可证是 Apache-2.0，见 `LICENSE`。

第一个提交是上游原样的代码（go-admin v1.2.26、themes v0.0.48 的 `sword` 和 `common`、`GoAdminGroup/html` v0.0.1），之后每一步裁剪是一个提交，`git log -p` 能看到改了什么。

## 依赖

只有 `golang.org/x/crypto`（bcrypt 存后台密码）。数据库驱动不在这里：用的程序自己注册 `sqlite3` 驱动（expgm 用纯 Go 的 modernc.org/sqlite）。

上游的直接依赖有 40 多个，编进二进制的约 35 个第三方模块（gin、xorm、goleveldb、zap、lib/pq、excelize……）。

## 去掉了什么、换成了什么

| 上游 | 现在 |
|---|---|
| gin、chi、echo、beego、iris、fiber 等 13 个框架适配器 | 只有 `adapter/nethttp`，路由挂在标准库 `http.ServeMux` 上 |
| excelize 导出 xlsx | 标准库 `encoding/csv` 导出 CSV（带 UTF-8 BOM，Excel 直接打开不乱码） |
| zap + lumberjack 日志 | 标准库 `log/slog`，写到 `slog.Default()`，落点由宿主程序配；`logger` 包的函数名不变 |
| YAML、INI 配置文件（yaml.v2、ini.v1） | 只留 Go 结构体和 JSON（`engine.AddConfig`、`AddConfigFromJSON`） |
| MySQL、MSSQL、PostgreSQL、OceanBase | 只留 SQLite（`table.DefaultConfig` 默认也改成 SQLite） |
| xorm（只用在一个没人调的 `CreateDB`） | 删掉 |
| x/text 的语言 Tag、`cases.Title` | 语言代码写成字面量；`utils.Title` 十几行 |
| google/uuid、NebulousLabs/fastrand | 标准库 `crypto/rand` |
| 安装向导、插件中心（远程插件商店）、代码生成器、`plugins/example` | 删掉，导航栏上对应的入口也删了 |
| 示例、多框架测试、要连 MySQL/MSSQL 的测试、依赖 testify 的测试 | 删掉；留下的测试只用标准库 |

另外一处行为改动：`ProhibitConfigModification: true` 时配置完全以代码为准，不读写 `goadmin_site` 表（上游总是用表里的值覆盖代码里的，改了代码不生效）。

修了一处上游的样式问题：表格的操作列是 sticky（`z-index: 1`），每格各自成层，后面行的操作格会盖住前面行弹出的操作菜单（表格有两行以上就能看到）。打开菜单的那格提到 `z-index: 2`（`td:last-child:has(.open)`，sword 和 common 的表格样式里各一处）。

sword 的页面开了 AdminLTE 的固定布局（`<body>` 加 `fixed`，`layout.tmpl` 和编进 `template.go` 的同一份各一处）：上游内容超过一屏时，滚动会把左边栏和顶栏一起滚走；现在只有右边内容滚，左边栏菜单长了自己滚（AdminLTE 自带的 slimScroll）。AdminLTE 的固定布局给 `.wrapper` 设了 `overflow:hidden`（假定 body 滚动），可这里 html、body、`.wrapper` 都是 `height:100%`，内容超过一屏就拖不下去，所以 `layout.tmpl` 里把 `.fixed .wrapper` 改回 `overflow-y: auto`，由 `.wrapper` 滚。

## 主题

只留 `themes/sword`（和它要的 `themes/common`）。adminlte 删了，`template.Default(ctx)` 的组件照常用。
