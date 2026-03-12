$ErrorActionPreference = "Stop"

$projectRoot = Split-Path -Parent $PSScriptRoot
Set-Location $projectRoot

$buildDir = Join-Path $projectRoot "build"
$bindingsDir = Join-Path $projectRoot "bindings\counter"
$contractPath = Join-Path $projectRoot "contracts\Counter.sol"

if (!(Test-Path $buildDir)) {
    New-Item -Path $buildDir -ItemType Directory | Out-Null
}
if (!(Test-Path $bindingsDir)) {
    New-Item -Path $bindingsDir -ItemType Directory | Out-Null
}

if (!(Get-Command npm -ErrorAction SilentlyContinue)) {
    throw "npm is required to compile Solidity with solcjs."
}

if (!(Test-Path (Join-Path $projectRoot "package.json"))) {
    npm init -y | Out-Null
}

npm install --save-dev solc | Out-Null
npx solcjs --bin --abi $contractPath -o $buildDir

$gopath = go env GOPATH
$abigenPath = Join-Path $gopath "bin\abigen.exe"
if (!(Test-Path $abigenPath)) {
    go install github.com/ethereum/go-ethereum/cmd/abigen@v1.14.12
}

$abiPath = Join-Path $buildDir "contracts_Counter_sol_Counter.abi"
$binPath = Join-Path $buildDir "contracts_Counter_sol_Counter.bin"
$outPath = Join-Path $bindingsDir "counter.go"

& $abigenPath --abi $abiPath --bin $binPath --pkg counter --type Counter --out $outPath

Write-Host "Go binding generated at: $outPath"
