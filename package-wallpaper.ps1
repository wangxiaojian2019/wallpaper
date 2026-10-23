$ErrorActionPreference = "Stop"
$rootDir = "d:\WALLPAPER"
$stagingDir = Join-Path $rootDir "release-staging"
$outputZip = Join-Path $rootDir "wallpaper-docker-prod.zip"

Write-Host "========================================================"
Write-Host "  WallArt Docker Packaging for wallpaper.icangoanyway.com"
Write-Host "========================================================"

# 1. 检查 Linux 二进制
Write-Host "[1/5] Checking wallpaper-server-linux..."
$linuxBin = Join-Path $rootDir "backend\wallpaper-server-linux"
if (-not (Test-Path $linuxBin)) {
    throw "wallpaper-server-linux not found!"
}
Write-Host "      wallpaper-server-linux ready."

# 2. 准备暂存目录
Write-Host "[2/5] Preparing staging dir: $stagingDir"
if (Test-Path $stagingDir) {
    Remove-Item $stagingDir -Recurse -Force
}
New-Item -ItemType Directory -Path $stagingDir -Force | Out-Null

function Copy-FilteredDirectory {
    param (
        [string]$Source,
        [string]$Destination,
        [string[]]$ExcludeDirs
    )
    if (-not (Test-Path $Destination)) {
        New-Item -ItemType Directory -Path $Destination -Force | Out-Null
    }
    
    Get-ChildItem -Path $Source -Force | ForEach-Object {
        $itemName = $_.Name
        if ($_.PSIsContainer) {
            if ($ExcludeDirs -notcontains $itemName) {
                $subDest = Join-Path $Destination $itemName
                Copy-FilteredDirectory -Source $_.FullName -Destination $subDest -ExcludeDirs $ExcludeDirs
            }
        } else {
            Copy-Item $_.FullName -Destination (Join-Path $Destination $itemName) -Force
        }
    }
}

# 3. 组织后端服务 (backend)
Write-Host "[3/5] Staging backend (excluding uploads)..."
$backendDest = Join-Path $stagingDir "backend"
New-Item -ItemType Directory -Path $backendDest -Force | Out-Null
Copy-Item (Join-Path $rootDir "backend\Dockerfile") -Destination (Join-Path $backendDest "Dockerfile") -Force
Copy-Item $linuxBin -Destination (Join-Path $backendDest "wallpaper-server-linux") -Force

# 严禁打包数据库文件，确保解压绝不会覆盖服务器现有真实商品、订单与配置数据
Write-Host "      [PROTECTED] wallpaper.db is EXCLUDED to prevent overwriting server data."

# 排除已有图片，只创建空 uploads 目录和 .gitkeep
$uploadsDest = Join-Path $backendDest "uploads"
New-Item -ItemType Directory -Path $uploadsDest -Force | Out-Null
[System.IO.File]::WriteAllText((Join-Path $uploadsDest ".gitkeep"), "")
Write-Host "      Uploads directory cleaned (.gitkeep placeholder created)."

# 4. 组织前端 (uni-app 和 frontend，排除 node_modules 和 dist)
Write-Host "[4/5] Staging uni-app and frontend..."
$uniappDest = Join-Path $stagingDir "uni-app"
Copy-FilteredDirectory -Source (Join-Path $rootDir "uni-app") -Destination $uniappDest -ExcludeDirs @('node_modules', '.git', 'dist')

$frontendDest = Join-Path $stagingDir "frontend"
Copy-FilteredDirectory -Source (Join-Path $rootDir "frontend") -Destination $frontendDest -ExcludeDirs @('node_modules', '.git', 'dist', '.vscode')

# 5. 复制根目录编排文件与脚本
Write-Host "[5/5] Staging root compose & management scripts..."
$rootFiles = @(
    'docker-compose.yml',
    'gateway.conf',
    'start.sh',
    'stop.sh',
    'restart.sh',
    'logs.sh',
    'DEPLOY_HK.md',
    'PAYPAL_CONFIG.md',
    'docker-start.bat',
    'docker-stop.bat',
    'setup-ssl.sh',
    'fix-ports.sh'
)

foreach ($file in $rootFiles) {
    $srcPath = Join-Path $rootDir $file
    if (Test-Path $srcPath) {
        $destPath = Join-Path $stagingDir $file
        if ($file.EndsWith('.sh')) {
            $content = [System.IO.File]::ReadAllText($srcPath, [System.Text.Encoding]::UTF8)
            $content = $content.Replace("`r`n", "`n")
            [System.IO.File]::WriteAllText($destPath, $content, (New-Object System.Text.UTF8Encoding($false)))
        } else {
            Copy-Item $srcPath -Destination $destPath -Force
        }
    }
}

# 复制 SSL 证书目录与 Certbot ACME 目录
$sslSrc = Join-Path $rootDir "ssl"
if (Test-Path $sslSrc) {
    Copy-FilteredDirectory -Source $sslSrc -Destination (Join-Path $stagingDir "ssl") -ExcludeDirs @()
}
$certbotWwwSrc = Join-Path $rootDir "certbot-www"
if (Test-Path $certbotWwwSrc) {
    Copy-FilteredDirectory -Source $certbotWwwSrc -Destination (Join-Path $stagingDir "certbot-www") -ExcludeDirs @()
}

# 6. 生成 ZIP 压缩包
Write-Host "Creating zip package with tar.exe..."
if (Test-Path $outputZip) {
    Remove-Item $outputZip -Force
}

& tar.exe -a -c -f $outputZip -C $stagingDir .

# 7. 清理暂存
Start-Sleep -Milliseconds 800
if (Test-Path $stagingDir) {
    try {
        Remove-Item $stagingDir -Recurse -Force -ErrorAction SilentlyContinue
    } catch {}
}

$zipStat = Get-Item $outputZip
$sizeMb = [math]::Round($zipStat.Length / 1MB, 2)

Write-Host "========================================================"
Write-Host "Package completed successfully!"
Write-Host "Domain:  wallpaper.icangoanyway.com (and admin.icangoanyway.com)"
Write-Host "Output:  $outputZip"
Write-Host "Size:    $sizeMb MB"
Write-Host "========================================================"
