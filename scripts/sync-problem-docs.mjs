import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.resolve(__dirname, "..");
const PROBLEMS = path.join(ROOT, "problems");

const REPO = "https://github.com/gauxs/lld";

/** problem_dir_name -> docs slug */
const SLUGS = {
  connect_four: "connect-four",
};

function ensureDir(dir) {
  fs.mkdirSync(dir, { recursive: true });
}

function yamlLink(key, label, href) {
  return `${key}:\n  text: ${label}\n  link: ${href}\n`;
}

function extensionSlug(id) {
  return id.replace(/_/g, "-");
}

function loadExtensions(problemId) {
  const jsonPath = path.join(PROBLEMS, problemId, "extensions.json");
  if (!fs.existsSync(jsonPath)) {
    return null;
  }
  const doc = JSON.parse(fs.readFileSync(jsonPath, "utf8"));
  const defaultId = doc.default ?? "baseline";
  const extMap = doc.extensions ?? {};
  let order = doc.order;
  if (!order?.length) {
    order = Object.keys(extMap);
  }
  const extensions = order.map((id) => {
    const node = extMap[id] ?? {};
    return {
      id,
      title: node.title ?? id,
      buildsOn: node.buildsOn ?? null,
    };
  });
  return { defaultId, extensions, extMap, order };
}

function mermaidExtensionGraph(extensions) {
  const lines = ["flowchart TD"];
  for (const ext of extensions) {
    const label = ext.title.replace(/"/g, "'");
    lines.push(`  ${ext.id}["${label}"]`);
  }
  for (const ext of extensions) {
    if (ext.buildsOn) {
      lines.push(`  ${ext.buildsOn} --> ${ext.id}`);
    }
  }
  return ["```mermaid", ...lines, "```"].join("\n");
}

function extensionDir(problemId, extId) {
  return path.join(PROBLEMS, problemId, "extensions", extId);
}

function hasCode(problemId, extId) {
  const dir = path.join(extensionDir(problemId, extId), "code");
  return fs.existsSync(dir) && fs.readdirSync(dir).length > 0;
}

function trailSteps(problemId, extId) {
  const steps = [
    { doc: "requirements.md", label: "requirements", pageClass: "lld-req-page" },
    { doc: "design.md", label: "design", pageClass: "" },
  ];
  if (hasCode(problemId, extId)) {
    steps.push({
      doc: "codebase.md",
      label: "codebase",
      pageClass: "lld-codebase-page",
    });
  }
  return steps;
}

function readSectionBody(srcPath) {
  const raw = fs.readFileSync(srcPath, "utf8");
  const title = raw.match(/^# (.+)/m)?.[1] ?? "Section";
  const body = raw.replace(/^# .+\n+/, "");
  return { title, body };
}

function writeExtensionPage(destDir, slug, extId, step, steps, meta) {
  const extSlug = extensionSlug(extId);
  const idx = steps.findIndex((s) => s.doc === step.doc);
  const base = `/problems/${slug}/extensions/${extSlug}`;
  const prev = idx > 0 ? steps[idx - 1] : null;
  const next = idx < steps.length - 1 ? steps[idx + 1] : null;

  let fm = `---
title: ${meta.title.replace(/"/g, '\\"')}
${step.pageClass ? `pageClass: ${step.pageClass}\n` : ""}`;

  if (meta.problem) {
    fm += `problem: ${meta.problem}\n`;
  }
  if (meta.extension) {
    fm += `extension: ${meta.extension}\n`;
  }
  if (step.doc === "requirements.md" && meta.scarcityPrev) {
    fm += yamlLink("prev", "scarcity", "/learn/concurrency/scarcity");
  }
  if (step.doc === "codebase.md") {
    fm += `aside: false
outline: false
`;
  }
  if (prev) {
    fm += yamlLink(
      "prev",
      prev.label,
      `${base}/${prev.doc.replace(/\.md$/, "")}`,
    );
  }
  if (next) {
    fm += yamlLink(
      "next",
      next.label,
      `${base}/${next.doc.replace(/\.md$/, "")}`,
    );
  }
  fm += `---\n\n`;

  fs.writeFileSync(path.join(destDir, step.doc), fm + meta.content);
}

function syncExtensionProblem(problemId, slug) {
  const { defaultId, extensions } = loadExtensions(problemId);
  const destRoot = path.join(ROOT, "docs", "problems", slug);
  ensureDir(destRoot);

  const problemTitle =
    problemId === "connect_four" ? "Connect Four" : problemId.replace(/_/g, " ");

  const hubLines = [
    `# ${problemTitle}`,
    "",
    "Extensions are **siblings** in the sidebar (baseline first). Use **Builds on** or the graph to see dependencies.",
    "",
    mermaidExtensionGraph(extensions),
    "",
    "| Extension | Builds on |",
    "| --- | --- |",
  ];

  for (const ext of extensions) {
    const extSlug = extensionSlug(ext.id);
    const buildsOn =
      ext.buildsOn == null
        ? "—"
        : `[${extensions.find((e) => e.id === ext.buildsOn)?.title ?? ext.buildsOn}](/problems/${slug}/extensions/${extensionSlug(ext.buildsOn)}/requirements)`;
    hubLines.push(
      `| [${ext.title}](/problems/${slug}/extensions/${extSlug}/requirements) | ${buildsOn} |`,
    );
  }

  fs.writeFileSync(
    path.join(destRoot, "index.md"),
    `---
title: Connect Four
---

${hubLines.join("\n")}
`,
  );

  for (const ext of extensions) {
    const srcBase = extensionDir(problemId, ext.id);
    const destDir = path.join(destRoot, "extensions", extensionSlug(ext.id));
    ensureDir(destDir);

    const steps = trailSteps(problemId, ext.id);
    const codeRel = `problems/${problemId}/extensions/${ext.id}/code`;
    const codeUrl = `${REPO}/tree/main/${codeRel}`;

    for (const step of steps) {
      if (step.doc === "codebase.md") {
        writeExtensionPage(
          destDir,
          slug,
          ext.id,
          step,
          steps,
          {
            title: `${ext.title} — codebase`,
            problem: problemId,
            extension: ext.id,
            content: `# Codebase

[${codeRel}](${codeUrl})

<ProblemCodebase problem="${problemId}" extension="${ext.id}" />
`,
          },
        );
        continue;
      }

      const srcName =
        step.doc === "requirements.md" ? "requirements.md" : "design.md";
      const srcPath = path.join(srcBase, srcName);
      if (!fs.existsSync(srcPath)) {
        console.warn(`Missing ${srcPath}`);
        continue;
      }
      const { title, body } = readSectionBody(srcPath);
      writeExtensionPage(destDir, slug, ext.id, step, steps, {
        title,
        problem: problemId,
        extension: ext.id,
        scarcityPrev:
          ext.id === defaultId && step.doc === "requirements.md",
        content: body,
      });
    }
  }

  cleanLegacyProblemPages(destRoot);
}

function cleanLegacyProblemPages(destRoot) {
  const legacy = [
    "functional-requirement.md",
    "design.md",
    "codebase.md",
    "followup.md",
  ];
  for (const name of legacy) {
    const p = path.join(destRoot, name);
    if (fs.existsSync(p)) {
      fs.unlinkSync(p);
    }
  }
  const variationDir = path.join(destRoot, "variation-1");
  if (fs.existsSync(variationDir)) {
    fs.rmSync(variationDir, { recursive: true });
  }
}

function buildSidebarExtensionItems(problemId, slug) {
  const loaded = loadExtensions(problemId);
  if (!loaded) {
    return [];
  }
  const { extensions } = loaded;

  return extensions.map((ext) => {
    const extSlug = extensionSlug(ext.id);
    const base = `/problems/${slug}/extensions/${extSlug}`;
    const items = [
      { text: "Requirements", link: `${base}/requirements` },
      { text: "Design", link: `${base}/design` },
    ];
    if (hasCode(problemId, ext.id)) {
      items.push({ text: "Codebase", link: `${base}/codebase` });
    }
    return {
      text: ext.title,
      collapsed: true,
      items,
    };
  });
}

for (const [problemId, slug] of Object.entries(SLUGS)) {
  if (loadExtensions(problemId)) {
    syncExtensionProblem(problemId, slug);
  }
}

const sidebarExtensions = buildSidebarExtensionItems(
  "connect_four",
  "connect-four",
);

const sidebarOut = path.join(ROOT, "docs", ".vitepress", "sidebar-problems.json");
fs.writeFileSync(
  sidebarOut,
  JSON.stringify({ connectFour: sidebarExtensions }, null, 2),
);

console.log("Synced problem docs from problems/ → docs/problems/");
console.log(`Wrote ${sidebarOut} (import in config.ts)`);
