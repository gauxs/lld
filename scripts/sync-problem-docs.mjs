import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.resolve(__dirname, "..");
const PROBLEMS = path.join(ROOT, "problems");

/** problem_dir_name -> docs slug */
const REPO = "https://github.com/gauxs/lld";

const SLUGS = {
  connect_four: "connect-four",
};

const CODE_PATHS = {
  connect_four: `${REPO}/tree/main/problems/connect_four/code`,
};

const TRAILS = {
  "connect-four": [
    { doc: "functional-requirement.md", label: "functional-requirement" },
    { doc: "design.md", label: "design" },
    { doc: "codebase.md", label: "codebase", extra: { problem: "connect_four" } },
    { doc: "followup.md", label: "followup" },
  ],
};

const SECTION_FILES = {
  "functional-requirement.md": "functional_requirement.md",
  "design.md": "design.md",
  "followup.md": "followup.md",
};

function ensureDir(dir) {
  fs.mkdirSync(dir, { recursive: true });
}

function yamlLink(key, label, href) {
  return `${key}:\n  text: ${label}\n  link: ${href}\n`;
}

function writeTrailPage(destDir, slug, step, content, title) {
  const trail = TRAILS[slug];
  const idx = trail.findIndex((s) => s.doc === step.doc);
  const base = `/problems/${slug}`;
  const prev = idx > 0 ? trail[idx - 1] : null;
  const next = idx < trail.length - 1 ? trail[idx + 1] : null;

  const pageClass =
    step.doc === "codebase.md"
      ? "lld-codebase-page"
      : step.doc === "design.md"
        ? ""
        : "lld-req-page";
  let fm = `---
title: ${title.replace(/"/g, '\\"')}
${pageClass ? `pageClass: ${pageClass}\n` : ""}`;
  if (step.doc === "functional-requirement.md" && slug === "connect-four") {
    fm += yamlLink("prev", "scarcity", "/learn/concurrency/scarcity");
  }
  if (step.extra?.problem) {
    fm += `problem: ${step.extra.problem}\n`;
  }
  if (step.doc === "codebase.md") {
    fm += `aside: false
outline: false
`;
  }
  if (prev) {
    fm += yamlLink("prev", prev.label, `${base}/${prev.doc.replace(/\.md$/, "")}`);
  }
  if (next) {
    fm += yamlLink("next", next.label, `${base}/${next.doc.replace(/\.md$/, "")}`);
  }
  fm += `---\n\n`;

  fs.writeFileSync(path.join(destDir, step.doc), fm + content);
}

function copySection(problemId, slug) {
  const srcDir = path.join(PROBLEMS, problemId);
  const destDir = path.join(ROOT, "docs", "problems", slug);
  ensureDir(destDir);

  const indexPath = path.join(destDir, "index.md");
  if (fs.existsSync(indexPath)) {
    fs.rmSync(indexPath);
  }

  for (const step of TRAILS[slug]) {
    if (step.doc === "codebase.md") {
      writeTrailPage(
        destDir,
        slug,
        step,
        `# Codebase

[problems/${problemId}/code](${CODE_PATHS[problemId]})

<ProblemCodebase problem="connect_four" />
`,
        "Codebase",
      );
      continue;
    }

    const srcName = SECTION_FILES[step.doc];
    const raw = fs.readFileSync(path.join(srcDir, srcName), "utf8");
    const title = raw.match(/^# (.+)/m)?.[1] ?? step.label;
    const body = raw.replace(/^# .+\n+/, "");
    writeTrailPage(destDir, slug, step, body, title);
  }
}

for (const [problemId, slug] of Object.entries(SLUGS)) {
  copySection(problemId, slug);
}

console.log("Synced problem docs from problems/ → docs/problems/");
