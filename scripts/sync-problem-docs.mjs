import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.resolve(__dirname, "..");
const PROBLEMS = path.join(ROOT, "problems");

/** problem_dir_name -> docs slug */
const SLUGS = {
  connect_four: "connect-four",
};

const SECTIONS = [
  { file: "functional_requirement.md", doc: "functional-requirement.md" },
  { file: "design.md", doc: "design.md" },
  { file: "followup.md", doc: "followup.md" },
];

function ensureDir(dir) {
  fs.mkdirSync(dir, { recursive: true });
}

function copySection(problemId, slug) {
  const srcDir = path.join(PROBLEMS, problemId);
  const destDir = path.join(ROOT, "docs", "problems", slug);
  ensureDir(destDir);

  const readme = fs.readFileSync(path.join(srcDir, "README.md"), "utf8");
  const body = readme.replace(/^# .+\n+/, "").trim();
  fs.writeFileSync(
    path.join(destDir, "index.md"),
    `---
title: Connect Four
description: Two-player Connect Four — requirements, design, and Go code.
next: /problems/connect-four/functional-requirement
---

# Connect Four

${body.split("## Layout")[0].trim()}
`,
  );

  const chain = [
    "/problems/connect-four/",
    "/problems/connect-four/functional-requirement",
    "/problems/connect-four/design",
    "/problems/connect-four/codebase",
    "/problems/connect-four/followup",
  ];

  for (let i = 0; i < SECTIONS.length; i++) {
    const { file, doc } = SECTIONS[i];
    const idx = i + 1;
    const content = fs.readFileSync(path.join(srcDir, file), "utf8");
    const title = content.match(/^# (.+)/m)?.[1] ?? doc;
    const next =
      doc === "followup.md" ? undefined : chain[idx + 1] ?? chain[chain.length - 1];
    const nextLine = next ? `next: ${next}\n` : "";
    fs.writeFileSync(
      path.join(destDir, doc),
      `---
title: ${title.replace(/"/g, '\\"')}
prev: ${chain[idx - 1]}
${nextLine}---

${content.replace(/^# .+\n+/, "")}`,
    );
  }

  fs.writeFileSync(
    path.join(destDir, "codebase.md"),
    `---
title: Codebase
description: Source files under problems/connect_four/code/
prev: /problems/connect-four/design
next: /problems/connect-four/followup
problem: connect_four
---

# Codebase

Reference implementation (\`package connectfour\`). Files are shown exactly as in the repo.

<ProblemCodebase problem="connect_four" />
`,
  );
}

for (const [problemId, slug] of Object.entries(SLUGS)) {
  copySection(problemId, slug);
}

console.log("Synced problem docs from problems/ → docs/problems/");
