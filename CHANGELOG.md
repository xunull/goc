# Changelog

## [0.2.0](https://github.com/xunull/goc/compare/v0.1.1...v0.2.0) (2026-07-02)


### Features

* **lang_ext:** add ExcludeLineCountOf value accessor ([5e74269](https://github.com/xunull/goc/commit/5e7426954f949e127f3ab9e951a7e672a5207c55))
* **lang_ext:** add ExcludeLineCountOf value accessor ([b8cf7cb](https://github.com/xunull/goc/commit/b8cf7cbfd4b7dc5e18bd7514602977ff17a935b0))
* **traverse_v2:** add WithShouldRecurse dynamic prune hook ([f884200](https://github.com/xunull/goc/commit/f884200fd05642742affdf065b643cacaba76483))
* **traverse_v2:** add WithShouldRecurse dynamic prune hook ([c89452d](https://github.com/xunull/goc/commit/c89452db1c76b5b024c620f80d9f0dc6158cf634))

## [0.1.1](https://github.com/xunull/goc/compare/v0.1.0...v0.1.1) (2026-06-03)


### Bug Fixes

* **ci:** release-please tag 去掉 goc- 前缀，符合 Go module 规范 ([78df5c5](https://github.com/xunull/goc/commit/78df5c5eccb724c942dc6b3e7d5f5ac7c23a76ca))




### Features

* **commandx:** 底层切换到 mvdan.cc/sh 内置解释器，修复 stderr 重定向等 bug ([a52de04](https://github.com/xunull/goc/commit/a52de04e335df77dfeacaab269bb9448db9a2414))
* **traverse_v2:** 补齐 v1 公开 API(GetAllPaths/GetFileCount/GetFileList + 兼容别名) ([061d2dc](https://github.com/xunull/goc/commit/061d2dcb9693b7d3f17270c896b39a8560dbca74))
* 新增 traverse_v2 包(层级签单 + race-free + ~10% 加速) ([b71291c](https://github.com/xunull/goc/commit/b71291c666ccb1b03ba5e64e1f332288133d8273))


### Bug Fixes

* **ci:** bump go-version 1.21 -&gt; 1.25 以匹配 go.mod ([3b7c4fa](https://github.com/xunull/goc/commit/3b7c4fa194e37b25ce37e3777cabace67f26fba9))
