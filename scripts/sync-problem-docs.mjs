import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.resolve(__dirname, "..");
const PROBLEMS = path.join(ROOT, "problems");
const DOCS_PROBLEMS = path.join(ROOT, "docs", "problems");

const REPO = "https://github.com/gauxs/lld";

function compareText(a, b) {
  return a < b ? -1 : a > b ? 1 : 0;
}

function encodeUrlSegment(value) {
  return encodeURIComponent(value).replace(/[!'()*]/g, (character) =>
    `%${character.charCodeAt(0).toString(16).toUpperCase()}`,
  );
}

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
    .sort((a, b) => compareText(a.name, b.name))) {
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

const LANGUAGE_BY_EXTENSION = {
  ".c": "c",
  ".cc": "cpp",
  ".cjs": "javascript",
  ".cpp": "cpp",
  ".cs": "csharp",
  ".cts": "typescript",
  ".css": "css",
  ".dart": "dart",
  ".go": "go",
  ".gradle": "groovy",
  ".graphql": "graphql",
  ".h": "c",
  ".hpp": "cpp",
  ".html": "html",
  ".java": "java",
  ".js": "javascript",
  ".json": "json",
  ".jsx": "jsx",
  ".kt": "kotlin",
  ".kts": "kotlin",
  ".md": "markdown",
  ".mjs": "javascript",
  ".mts": "typescript",
  ".php": "php",
  ".proto": "protobuf",
  ".py": "python",
  ".rb": "ruby",
  ".rs": "rust",
  ".scss": "scss",
  ".sh": "bash",
  ".sql": "sql",
  ".svelte": "svelte",
  ".swift": "swift",
  ".toml": "toml",
  ".ts": "typescript",
  ".tsx": "tsx",
  ".vue": "vue",
  ".xml": "xml",
  ".yaml": "yaml",
  ".yml": "yaml",
};

const LANGUAGE_BY_FILENAME = {
  Dockerfile: "dockerfile",
  Makefile: "makefile",
  "go.mod": "go",
  "go.sum": "text",
};

const INCLUDED_DOTFILES = new Set([".env.example", ".gitignore"]);
const UTF8_DECODER = new TextDecoder("utf-8", { fatal: true });

function isTextFile(buffer) {
  if (buffer.includes(0)) {
    return false;
  }
  let controlBytes = 0;
  for (const byte of buffer) {
    if (byte < 7 || (byte > 13 && byte < 32)) {
      controlBytes++;
    }
  }
  if (buffer.length > 0 && controlBytes / buffer.length >= 0.1) {
    return false;
  }
  try {
    UTF8_DECODER.decode(buffer);
    return true;
  } catch {
    return false;
  }
}

function collectCodeFiles(dir, prefix = "") {
  if (!fs.existsSync(dir)) {
    return [];
  }
  const files = [];
  for (const entry of fs
    .readdirSync(dir, { withFileTypes: true })
    .sort((a, b) => compareText(a.name, b.name))) {
    if (entry.name.startsWith(".") && !INCLUDED_DOTFILES.has(entry.name)) {
      continue;
    }
    const fullPath = path.join(dir, entry.name);
    const relativePath = prefix ? `${prefix}/${entry.name}` : entry.name;
    if (entry.isDirectory()) {
      files.push(...collectCodeFiles(fullPath, relativePath));
      continue;
    }
    if (!entry.isFile()) {
      continue;
    }
    const buffer = fs.readFileSync(fullPath);
    if (isTextFile(buffer)) {
      files.push({ path: relativePath, content: buffer.toString("utf8") });
    }
  }
  return files;
}

const codeFileCache = new Map();

function codeFiles(problemId, extId) {
  const key = `${problemId}:${extId}`;
  if (!codeFileCache.has(key)) {
    codeFileCache.set(
      key,
      collectCodeFiles(path.join(extensionDir(problemId, extId), "code")),
    );
  }
  return codeFileCache.get(key);
}

function hasCode(problemId, extId) {
  return codeFiles(problemId, extId).length > 0;
}

function codeLanguage(filePath) {
  const filename = path.basename(filePath);
  if (filename === "Dockerfile" || filename.startsWith("Dockerfile.")) {
    return "dockerfile";
  }
  return (
    LANGUAGE_BY_FILENAME[filename] ??
    LANGUAGE_BY_EXTENSION[path.extname(filename).toLowerCase()] ??
    "text"
  );
}

function codeFence(content) {
  const longestRun = Math.max(
    0,
    ...(content.match(/`+/g) ?? []).map((run) => run.length),
  );
  return "`".repeat(Math.max(3, longestRun + 1));
}

function inlineCode(value) {
  const normalized = value.replace(/[\r\n]+/g, " ");
  const longestRun = Math.max(
    0,
    ...(normalized.match(/`+/g) ?? []).map((run) => run.length),
  );
  const delimiter = "`".repeat(longestRun + 1);
  return `${delimiter} ${normalized} ${delimiter}`;
}

function directoryTree(files) {
  const root = new Map();
  for (const file of files) {
    let children = root;
    for (const [index, segment] of file.path.split("/").entries()) {
      if (!children.has(segment)) {
        children.set(segment, {
          file: index === file.path.split("/").length - 1,
          children: new Map(),
        });
      }
      children = children.get(segment).children;
    }
  }

  const lines = ["code/"];
  function render(children, prefix) {
    const entries = [...children.entries()].sort(
      ([nameA, nodeA], [nameB, nodeB]) =>
        Number(nodeA.file) - Number(nodeB.file) || compareText(nameA, nameB),
    );
    entries.forEach(([name, node], index) => {
      const last = index === entries.length - 1;
      const displayName = name.replace(/[\r\n]+/g, " ");
      lines.push(
        `${prefix}${last ? "└──" : "├──"} ${displayName}${node.file ? "" : "/"}`,
      );
      if (!node.file) {
        render(node.children, `${prefix}${last ? "    " : "│   "}`);
      }
    });
  }
  render(root, "");
  return lines.join("\n");
}

function codebaseMarkdown(files, codeRel, codeUrl) {
  const tree = directoryTree(files);
  const treeFence = codeFence(tree);
  const sections = files.map((file) => {
    const fence = codeFence(file.content);
    const content = file.content.endsWith("\n")
      ? file.content
      : `${file.content}\n`;
    return `## ${inlineCode(file.path)}

${fence}${codeLanguage(file.path)}
${content}${fence}`;
  });

  return `# Codebase

Source: [GitHub](${codeUrl})

Path: ${inlineCode(codeRel)}

## Directory structure

${treeFence}text
${tree}
${treeFence}

${sections.join("\n\n")}
`;
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
      pageClass: "",
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
    const encodedCodePath = codeRel
      .split("/")
      .map(encodeUrlSegment)
      .join("/");
    const codeUrl = `${REPO}/tree/main/${encodedCodePath}`;

    for (const step of steps) {
      if (step.doc === "codebase.md") {
        const files = codeFiles(problemId, ext.id);
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
            content: codebaseMarkdown(files, codeRel, codeUrl),
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
      items.push({ text: "Code", link: `${base}/codebase` });
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
    collapsed: true,
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
      .sort(([a], [b]) => compareText(a, b))
      .map(([name, child]) => ({
        text: humanizeName(name),
        collapsed: true,
        items: render(child),
      }));
    const problems = node.problems
      .sort(compareText)
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
