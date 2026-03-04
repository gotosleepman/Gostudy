const fs = require("fs");
const path = require("path");
const { spawnSync } = require("child_process");

const ROOT = process.cwd();
const REPORTS_DIR = path.join(ROOT, "reports");
const TEST_LOG_PATH = path.join(REPORTS_DIR, "test-output.log");
const COVERAGE_FINAL_PATH = path.join(ROOT, "coverage", "coverage-final.json");
const REPORT_MD_PATH = path.join(REPORTS_DIR, "test-report.md");

function run(cmd, args) {
  const result = spawnSync(cmd, args, {
    cwd: ROOT,
    encoding: "utf8",
    shell: process.platform === "win32",
  });
  return {
    code: result.status ?? 1,
    stdout: result.stdout || "",
    stderr: result.stderr || "",
  };
}

function readJsonSafe(filePath) {
  try {
    const text = fs.readFileSync(filePath, "utf8");
    return JSON.parse(text);
  } catch {
    return null;
  }
}

function parseTestStats(output) {
  const parse = (regex) => {
    const m = output.match(regex);
    return m ? Number(m[1]) : 0;
  };

  return {
    tests: parse(/(\d+)\s+passing/),
    failures: parse(/(\d+)\s+failing/),
    pending: parse(/(\d+)\s+pending/),
    durationText: (output.match(/\(([^)]+)\)\s*$/m) || [])[1] || "N/A",
  };
}

function extractFailureItems(output) {
  const lines = output.split(/\r?\n/);
  const failures = [];
  for (let i = 0; i < lines.length; i++) {
    if (/^\s*\d+\)\s+/.test(lines[i])) {
      const title = lines[i].trim();
      const detail = lines[i + 1] ? lines[i + 1].trim() : "未知错误";
      failures.push({ title, detail });
    }
  }
  return failures;
}

function summarizeCoverage(coverageFinal) {
  if (!coverageFinal) return null;

  const total = {
    statements: { covered: 0, total: 0 },
    branches: { covered: 0, total: 0 },
    functions: { covered: 0, total: 0 },
    lines: { covered: 0, total: 0 },
  };

  for (const fileData of Object.values(coverageFinal)) {
    const s = fileData.s || {};
    const f = fileData.f || {};
    const b = fileData.b || {};
    const l = fileData.l || {};
    const statementMap = fileData.statementMap || {};

    total.statements.total += Object.keys(s).length;
    total.statements.covered += Object.values(s).filter((hit) => hit > 0).length;

    total.functions.total += Object.keys(f).length;
    total.functions.covered += Object.values(f).filter((hit) => hit > 0).length;

    for (const arr of Object.values(b)) {
      const branchHits = Array.isArray(arr) ? arr : [];
      total.branches.total += branchHits.length;
      total.branches.covered += branchHits.filter((hit) => hit > 0).length;
    }

    if (Object.keys(l).length > 0) {
      total.lines.total += Object.keys(l).length;
      total.lines.covered += Object.values(l).filter((hit) => hit > 0).length;
    } else {
      const lineHits = new Map();
      for (const [id, map] of Object.entries(statementMap)) {
        const line = map?.start?.line;
        if (!line) continue;
        const prev = lineHits.get(line) || 0;
        lineHits.set(line, Math.max(prev, s[id] || 0));
      }
      total.lines.total += lineHits.size;
      total.lines.covered += Array.from(lineHits.values()).filter((hit) => hit > 0).length;
    }
  }

  const pct = (covered, all) => (all === 0 ? 0 : Number(((covered / all) * 100).toFixed(2)));
  return {
    statements: pct(total.statements.covered, total.statements.total),
    branches: pct(total.branches.covered, total.branches.total),
    functions: pct(total.functions.covered, total.functions.total),
    lines: pct(total.lines.covered, total.lines.total),
  };
}

function formatCoverage(summary) {
  if (!summary) {
    return [
      "| 指标 | 覆盖率 |",
      "| --- | --- |",
      "| Statements | N/A |",
      "| Branches | N/A |",
      "| Functions | N/A |",
      "| Lines | N/A |",
    ].join("\n");
  }

  const pct = (name) => `${summary[name] ?? 0}%`;
  return [
    "| 指标 | 覆盖率 |",
    "| --- | --- |",
    `| Statements | ${pct("statements")} |`,
    `| Branches | ${pct("branches")} |`,
    `| Functions | ${pct("functions")} |`,
    `| Lines | ${pct("lines")} |`,
  ].join("\n");
}

function main() {
  fs.mkdirSync(REPORTS_DIR, { recursive: true });

  console.log("[report] Compiling contracts...");
  const compile = run("npx", ["hardhat", "compile"]);
  if (compile.code !== 0) {
    console.error(compile.stdout || compile.stderr);
    process.exit(1);
  }

  console.log("[report] Running tests...");
  const test = run("npx", ["hardhat", "test", "--no-compile"]);
  fs.writeFileSync(TEST_LOG_PATH, `${test.stdout}\n${test.stderr}`, "utf8");

  console.log("[report] Running coverage...");
  const coverage = run("npx", ["hardhat", "coverage", "--testfiles", "test/**/*.js"]);

  const stats = parseTestStats(`${test.stdout}\n${test.stderr}`);
  const failures = extractFailureItems(`${test.stdout}\n${test.stderr}`);
  const coverageSummary = summarizeCoverage(readJsonSafe(COVERAGE_FINAL_PATH));
  const createdAt = new Date().toISOString();

  const failureSection =
    failures.length === 0
      ? "- 无失败用例"
      : failures
          .map((f) => {
            return `- ${f.title}\n  - ${String(f.detail).replace(/\n/g, " ")}`;
          })
          .join("\n");

  const report = [
    "# 测试报告",
    "",
    `- 生成时间: ${createdAt}`,
    `- 测试命令退出码: ${test.code}`,
    `- 覆盖率命令退出码: ${coverage.code}`,
    "",
    "## 测试结果",
    "",
    "| 项目 | 数值 |",
    "| --- | --- |",
    `| 总用例 | ${stats.tests + stats.failures} |`,
    `| 通过 | ${stats.tests} |`,
    `| 失败 | ${stats.failures} |`,
    `| Pending | ${stats.pending} |`,
    `| 测试耗时 | ${stats.durationText} |`,
    "",
    "## 失败详情",
    "",
    failureSection,
    "",
    "## 覆盖率汇总",
    "",
    formatCoverage(coverageSummary),
    "",
    "## 产物文件",
    "",
    `- 测试日志: \`${path.relative(ROOT, TEST_LOG_PATH)}\``,
    `- 覆盖率明细: \`${path.relative(ROOT, COVERAGE_FINAL_PATH)}\``,
  ].join("\n");

  fs.writeFileSync(REPORT_MD_PATH, report, "utf8");
  console.log(`[report] Generated: ${path.relative(ROOT, REPORT_MD_PATH)}`);

  if (test.code !== 0 || coverage.code !== 0) {
    process.exit(1);
  }
}

main();
