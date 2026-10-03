import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.resolve(__dirname, "..");
const SOURCE = path.join(ROOT, "practice_with_ai");
const DESTINATION = path.join(ROOT, "docs", "practice-with-ai");
const SIDEBAR = path.join(
  ROOT,
  "docs",
  ".vitepress",
  "sidebar-practice.json",
);

function compareText(a, b) {
  return a < b ? -1 : a > b ? 1 : 0;
}

function humanize(name) {
  return name
    .replace(/\.md$/i, "")
    .replace(/^\d+[-_]/, "")
    .replace(/[-_]+/g, " ")
    .replace(/\b\w/g, (character) => character.toUpperCase());
}

function collectMarkdown(dir, prefix = "") {
  if (!fs.existsSync(dir)) {
    return [];
  }
  const files = [];
  for (const entry of fs
    .readdirSync(dir, { withFileTypes: true })
    .sort((a, b) => compareText(a.name, b.name))) {
    if (entry.name.startsWith(".") || entry.name.startsWith("_")) {
      continue;
    }
    const fullPath = path.join(dir, entry.name);
    const relativePath = prefix ? `${prefix}/${entry.name}` : entry.name;
    if (entry.isDirectory()) {
      files.push(...collectMarkdown(fullPath, relativePath));
    } else if (entry.isFile() && entry.name.endsWith(".md")) {
      files.push(relativePath);
    }
  }
  return files;
}

function publishedPath(relativePath) {
  return relativePath
    .split("/")
    .map((segment) => segment.replace(/_/g, "-"))
    .join("/");
}

function codeFence(content) {
  const longestRun = Math.max(
    0,
    ...(content.match(/`+/g) ?? []).map((run) => run.length),
  );
  return "`".repeat(Math.max(3, longestRun + 1));
}

function renderPage(page, content) {
  if (page.relativePath !== "prompt.md") {
    return content;
  }

  const fence = codeFence(content);
  const prompt = content.endsWith("\n") ? content : `${content}\n`;
  return `# Practice LLD interviews with AI

Use this prompt to turn a chatbot into a structured low-level design interviewer. It controls requirement gathering, design discussion, implementation, extensibility, and final grading without revealing answers too early.

## How to use

1. Copy the complete prompt below using the code block's copy button.
2. Change the \`LEVEL\` and \`PROBLEM\` values at the top.
3. Paste it into a new conversation with your preferred chatbot.
4. Respond as the candidate and let the chatbot conduct the interview.

## Interviewer prompt

${fence}md
${prompt}${fence}
`;
}

const files = collectMarkdown(SOURCE);
const pages = files.map((relativePath) => ({
  relativePath,
  outputRelativePath: publishedPath(relativePath),
  title:
    relativePath === "prompt.md"
      ? "LLD Interviewer Prompt"
      : humanize(path.basename(relativePath)),
  sidebarTitle:
    relativePath === "prompt.md"
      ? "Interviewer prompt"
      : humanize(path.basename(relativePath)),
}));

const publishedSources = new Map();
for (const page of pages) {
  const conflictingSource = publishedSources.get(page.outputRelativePath);
  if (conflictingSource) {
    throw new Error(
      `Practice pages "${conflictingSource}" and "${page.relativePath}" both publish to "${page.outputRelativePath}"`,
    );
  }
  publishedSources.set(page.outputRelativePath, page.relativePath);
}

const items = pages.map((page) => ({
  text: page.sidebarTitle,
  link: `/practice-with-ai/${page.outputRelativePath.replace(/\.md$/i, "")}`,
}));
const sidebarPayload = {
  items,
  firstLink: items[0]?.link ?? "/",
  pageCount: items.length,
};

const temporaryDestination = `${DESTINATION}.sync-${process.pid}-${Date.now()}`;
const temporarySidebar = `${SIDEBAR}.sync-${process.pid}-${Date.now()}`;
fs.mkdirSync(temporaryDestination, { recursive: true });

try {
  for (const page of pages) {
    const sourcePath = path.join(SOURCE, page.relativePath);
    const destinationPath = path.join(
      temporaryDestination,
      page.outputRelativePath,
    );
    const content = fs.readFileSync(sourcePath, "utf8");
    const renderedContent = renderPage(page, content);
    fs.mkdirSync(path.dirname(destinationPath), { recursive: true });
    fs.writeFileSync(
      destinationPath,
      `---\ntitle: ${JSON.stringify(page.title)}\n---\n\n${renderedContent}`,
    );
  }
  fs.writeFileSync(temporarySidebar, JSON.stringify(sidebarPayload, null, 2));

  const destinationBackup = `${DESTINATION}.previous`;
  const sidebarBackup = `${SIDEBAR}.previous`;
  fs.rmSync(destinationBackup, { recursive: true, force: true });
  fs.rmSync(sidebarBackup, { force: true });

  try {
    if (fs.existsSync(DESTINATION)) {
      fs.renameSync(DESTINATION, destinationBackup);
    }
    if (fs.existsSync(SIDEBAR)) {
      fs.renameSync(SIDEBAR, sidebarBackup);
    }
    fs.renameSync(temporaryDestination, DESTINATION);
    fs.renameSync(temporarySidebar, SIDEBAR);
  } catch (error) {
    fs.rmSync(DESTINATION, { recursive: true, force: true });
    fs.rmSync(SIDEBAR, { force: true });
    if (fs.existsSync(destinationBackup)) {
      fs.renameSync(destinationBackup, DESTINATION);
    }
    if (fs.existsSync(sidebarBackup)) {
      fs.renameSync(sidebarBackup, SIDEBAR);
    }
    throw error;
  }
  fs.rmSync(destinationBackup, { recursive: true, force: true });
  fs.rmSync(sidebarBackup, { force: true });
} finally {
  fs.rmSync(temporaryDestination, { recursive: true, force: true });
  fs.rmSync(temporarySidebar, { force: true });
}

console.log("Synced practice guides from practice_with_ai/ → docs/practice-with-ai/");
console.log(`Wrote ${SIDEBAR}`);
