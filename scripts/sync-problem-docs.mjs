import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.resolve(__dirname, "..");
const PROBLEMS = path.join(ROOT, "problems");
const DOCS_PROBLEMS = path.join(ROOT, "docs", "problems");

const REPO = "https://github.com/gauxs/lld";

function humanizeName(name) {
  return name
    .replace(/[-_]+/g, " ")
    .replace(/\b\w/g, (c) => c.toUpperCase());
}

function problemTitle(problemId, configuredTitle) {
  return configuredTitle ?? humanizeName(problemId.split("/").at(-1));
}

function ensureDir(dir) {
  fs.mkdirSync(dir, { recursive: true });
}

function yamlLink(key, label, href) {
  return `${key}:\n  text: ${label}\n  link: ${href}\n`;
}

function extensionSlug(id) {
  return id.replace(/_/g, "-");
}

function problemSlug(problemId) {
  return problemId.split("/").map(extensionSlug).join("/");
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
  return {
    title: doc.title ?? null,
    defaultId,
    extensions,
    extMap,
    order,
  };
}

function discoverProblemIds(dir = PROBLEMS, prefix = "") {
  const ids =
    prefix && fs.existsSync(path.join(dir, "extensions.json")) ? [prefix] : [];
  for (const entry of fs
    .readdirSync(dir, { withFileTypes: true })
    .sort((a, b) => a.name.localeCompare(b.name))) {
    if (
      !entry.isDirectory() ||
      entry.name.startsWith(".") ||
      entry.name.startsWith("_")
    ) {
      continue;
    }
    const childPrefix = prefix ? `${prefix}/${entry.name}` : entry.name;
    ids.push(...discoverProblemIds(path.join(dir, entry.name), childPrefix));
  }
  return ids;
}

function mermaidExtensionGraph(extensions) {
  const lines = ["flowchart TD"];
  for (const ext of extensions) {
    const label = ext.title.replace(/"/g, "'");
    lines.push(`  ${ext.id}["${label}"]`);
  }
  for (const ext of extensions) {
    if (ext.buildsOn) {
      lines.push(`  ${ext.id} -->|BuildsOn| ${ext.buildsOn}`);
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

function syncExtensionProblem(problemId, slug, outputRoot = DOCS_PROBLEMS) {
  const { extensions, title: configuredTitle } = loadExtensions(problemId);
  const destRoot = path.join(outputRoot, slug);
  ensureDir(destRoot);

  const title = problemTitle(problemId, configuredTitle);

  const hubLines = [
    `# ${title}`,
    "",
    "Extensions are **siblings** in the sidebar (baseline first). Dependencies:",
    "",
    mermaidExtensionGraph(extensions),
  ];

  fs.writeFileSync(
    path.join(destRoot, "index.md"),
    `---
title: ${title}
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

function buildProblemSidebar(problemId) {
  const loaded = loadExtensions(problemId);
  const slug = problemSlug(problemId);
  return {
    text: problemTitle(problemId, loaded.title),
    collapsed: false,
    items: [
      { text: "Overview", link: `/problems/${slug}/` },
      ...buildSidebarExtensionItems(problemId, slug),
    ],
  };
}

function buildSidebarTree(problemIds) {
  const root = { sections: new Map(), problems: [] };

  for (const problemId of problemIds) {
    const segments = problemId.split("/");
    const problemName = segments.pop();
    let node = root;
    for (const section of segments) {
      if (!node.sections.has(section)) {
        node.sections.set(section, { sections: new Map(), problems: [] });
      }
      node = node.sections.get(section);
    }
    node.problems.push(
      segments.length ? `${segments.join("/")}/${problemName}` : problemName,
    );
  }

  function render(node) {
    const sections = [...node.sections.entries()]
      .sort(([a], [b]) => a.localeCompare(b))
      .map(([name, child]) => ({
        text: humanizeName(name),
        collapsed: false,
        items: render(child),
      }));
    const problems = node.problems
      .sort((a, b) => a.localeCompare(b))
      .map(buildProblemSidebar);
    return [...sections, ...problems];
  }

  return render(root);
}

function validateProblem(problemId) {
  const loaded = loadExtensions(problemId);
  if (!loaded) {
    throw new Error(`Missing extensions.json for ${problemId}`);
  }

  const extensionIds = new Set(loaded.extensions.map((ext) => ext.id));
  if (!extensionIds.has(loaded.defaultId)) {
    throw new Error(
      `${problemId}: default extension "${loaded.defaultId}" is not in order`,
    );
  }

  for (const extension of loaded.extensions) {
    if (extension.buildsOn && !extensionIds.has(extension.buildsOn)) {
      throw new Error(
        `${problemId}/${extension.id}: unknown buildsOn "${extension.buildsOn}"`,
      );
    }
    for (const document of ["requirements.md", "design.md"]) {
      const source = path.join(extensionDir(problemId, extension.id), document);
      if (!fs.existsSync(source)) {
        throw new Error(`Missing ${source}`);
      }
    }
  }
}

function replaceGeneratedProblems(tempRoot, tempSidebar, sidebarOut) {
  const backupRoot = `${DOCS_PROBLEMS}.previous`;
  const backupSidebar = `${sidebarOut}.previous`;
  fs.rmSync(backupRoot, { recursive: true, force: true });
  fs.rmSync(backupSidebar, { force: true });
  try {
    if (fs.existsSync(DOCS_PROBLEMS)) {
      fs.renameSync(DOCS_PROBLEMS, backupRoot);
    }
    if (fs.existsSync(sidebarOut)) {
      fs.renameSync(sidebarOut, backupSidebar);
    }
    fs.renameSync(tempRoot, DOCS_PROBLEMS);
    fs.renameSync(tempSidebar, sidebarOut);
    fs.rmSync(backupRoot, { recursive: true, force: true });
    fs.rmSync(backupSidebar, { force: true });
  } catch (error) {
    fs.rmSync(DOCS_PROBLEMS, { recursive: true, force: true });
    fs.rmSync(sidebarOut, { force: true });
    if (fs.existsSync(backupRoot)) {
      fs.renameSync(backupRoot, DOCS_PROBLEMS);
    }
    if (fs.existsSync(backupSidebar)) {
      fs.renameSync(backupSidebar, sidebarOut);
    }
    throw error;
  }
}

const problemIds = discoverProblemIds();
for (const problemId of problemIds) {
  validateProblem(problemId);
}

const firstProblemId = problemIds[0] ?? null;
const firstLoaded = firstProblemId ? loadExtensions(firstProblemId) : null;
const sidebarPayload = {
  items: buildSidebarTree(problemIds),
  firstLink:
    firstProblemId && firstLoaded
      ? `/problems/${problemSlug(firstProblemId)}/extensions/${extensionSlug(firstLoaded.defaultId)}/requirements`
      : "/",
  problemCount: problemIds.length,
};
const sidebarOut = path.join(ROOT, "docs", ".vitepress", "sidebar-problems.json");
const tempSidebar = `${sidebarOut}.${process.pid}-${Date.now()}.tmp`;
fs.writeFileSync(tempSidebar, JSON.stringify(sidebarPayload, null, 2));

const tempProblems = path.join(
  ROOT,
  "docs",
  `.problems-sync-${process.pid}-${Date.now()}`,
);
ensureDir(tempProblems);
try {
  for (const problemId of problemIds) {
    syncExtensionProblem(problemId, problemSlug(problemId), tempProblems);
  }
  replaceGeneratedProblems(tempProblems, tempSidebar, sidebarOut);
} finally {
  fs.rmSync(tempProblems, { recursive: true, force: true });
  fs.rmSync(tempSidebar, { force: true });
}

console.log("Synced problem docs from problems/ → docs/problems/");
console.log(`Wrote ${sidebarOut}`);
