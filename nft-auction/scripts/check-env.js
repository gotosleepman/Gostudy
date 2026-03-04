require("dotenv").config();

const requiredVarsByNetwork = {
  sepolia: ["SEPOLIA_RPC_URL", "SEPOLIA_PRIVATE_KEY"],
};

function isMissing(name) {
  const value = process.env[name];
  return value === undefined || value.trim() === "";
}

function main() {
  const network = process.argv[2] || "sepolia";
  const required = requiredVarsByNetwork[network] || [];
  const missing = required.filter(isMissing);

  if (missing.length > 0) {
    console.error(`\n[env-check] 检测到 ${network} 网络缺少必填环境变量:`);
    for (const name of missing) {
      console.error(`- ${name}`);
    }
    console.error("\n请在项目根目录 .env 中补全后重试。示例:");
    console.error("SEPOLIA_RPC_URL=https://sepolia.infura.io/v3/<YOUR_KEY>");
    console.error("SEPOLIA_PRIVATE_KEY=0x<YOUR_PRIVATE_KEY>\n");
    process.exit(1);
  }

  if (network === "sepolia") {
    const rpc = process.env.SEPOLIA_RPC_URL || "";
    const pk = process.env.SEPOLIA_PRIVATE_KEY || "";
    const privateKeyPattern = /^0x[0-9a-fA-F]{64}$/;

    if (!/^https?:\/\//.test(rpc)) {
      console.error("\n[env-check] SEPOLIA_RPC_URL 格式不正确，必须是 http/https 开头的完整 URL。\n");
      process.exit(1);
    }

    if (!privateKeyPattern.test(pk)) {
      console.error(
        "\n[env-check] SEPOLIA_PRIVATE_KEY 格式不正确，必须是 0x + 64 位十六进制字符串（32字节私钥）。\n",
      );
      process.exit(1);
    }
  }

  console.log(`[env-check] ${network} 环境变量检查通过`);
}

main();
