# WallArt 4K 壁纸独立站 - 核心项目规约与安全铁律

## 1. 生产数据零破坏原则（CRITICAL）
- **数据库保护**：`backend/wallpaper.db` 是生产环境唯一存储真实商品、订单、支付流水、用户数据的 SQLite 数据库。
  - **严禁** 在任何打包脚本（如 `package-wallpaper.ps1`）、Docker 构建、Git 提交、启动脚本中覆盖或重置此文件。
  - 发布包中必须严格排除 `backend/wallpaper.db`。
- **图片资源保护**：`backend/uploads/` 存储生产上传的 4K 壁纸、动效预览图及实机演练资源，打包与部署时严禁清空生产上的已有文件。

## 2. 真实业务交付与严禁 Mock 兜底原则
- **零 Mock 铁律**：严禁在前端（Uni-app 或 Vue3 前后台）的 API 报错、空数据或 catch 块中静默注入 Demo/Mock 假数据（如 `mock-1`、`wallpaper_sample`、`8888`、`vip8`）。
- **异常展示规范**：接口调用失败时，必须真实提示用户错误信息或提供“网络异常，点击重试”按钮，绝不允许展示虚假壁纸让买家产生困惑。
- **网盘交付动态性**：下载页（`/pages/download/index`）必须基于后台录入的实际网盘通道（夸克、百度、Google Drive、自定义网盘）进行条件渲染；未配置的通道绝不渲染，绝不允许写死 Demo 链接。

## 3. 端口与反向代理架构铁律
- **端口映射定义**：
  - **Uni-app 买家端**：Docker 容器内部端口 `80` -> 宿主机端口 `5173`（绑绑定域名 `wallpaper.icangoanyway.com`）。
  - **Vue3 管理后台**：Docker 容器内部端口 `80` -> 宿主机端口 `5174`（绑定域名 `admin.icangoanyway.com`）。
  - **Go 后端 API**：Docker 容器内部端口 `8080` -> 宿主机端口 `8080`（内部/网关代理路径 `/api`）。
- **严禁端口倒置**：严禁颠倒 5173 与 5174 的域名与服务映射。部署前必须运行 `fix-ports.sh` 或使用 `restart.sh` 自动纠偏。

## 4. 货币与价格呈现一致性
- **币种规则**：
  - 国内通道（微信、支付宝、自定义转账收款）产生的订单：货币单位必须统一记录并展示为 `CNY`（货币符号 `¥`）。
  - 国际通道（PayPal）产生的订单：货币单位必须统一记录并展示为 `USD`（货币符号 `$`）。
  - 提交订单价格即为后台展示价格，严禁未经汇率换算随意篡改展示数值。

## 5. 编译与打包标准
- **Linux 编译**：
  必须指定环境变量 `$env:GOOS="linux"; $env:GOARCH="amd64"; $env:GOCACHE="d:\WALLPAPER\backend\.cache"`，编译输出为 `backend/wallpaper-server-linux`。
- **标准打包**：
  统一使用 `powershell -ExecutionPolicy Bypass -File .\package-wallpaper.ps1` 进行打包，输出为根目录下的 `wallpaper-docker-prod.zip`。
