# Changelog

## [0.1.0](https://github.com/xunull/goc/compare/v0.0.13...v0.1.0) (2026-06-03)


### Features

* **commandx:** 底层切换到 mvdan.cc/sh 内置解释器，修复 stderr 重定向等 bug ([a52de04](https://github.com/xunull/goc/commit/a52de04e335df77dfeacaab269bb9448db9a2414))
* **traverse_v2:** 补齐 v1 公开 API(GetAllPaths/GetFileCount/GetFileList + 兼容别名) ([061d2dc](https://github.com/xunull/goc/commit/061d2dcb9693b7d3f17270c896b39a8560dbca74))
* 新增 traverse_v2 包(层级签单 + race-free + ~10% 加速) ([b71291c](https://github.com/xunull/goc/commit/b71291c666ccb1b03ba5e64e1f332288133d8273))


### Bug Fixes

* **ci:** bump go-version 1.21 -&gt; 1.25 以匹配 go.mod ([3b7c4fa](https://github.com/xunull/goc/commit/3b7c4fa194e37b25ce37e3777cabace67f26fba9))
