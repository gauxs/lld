import fs from "node:fs";
import path from "node:path";
import { fileURLToPath } from "node:url";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.resolve(__dirname, "..");
const THEORY = path.join(ROOT, "theory");
const DOCS_LEARN = path.join(ROOT, "docs", "learn");

function ensureDir(dir) {
  fs.mkdirSync(dir, { recursive: true });
}

function yamlLink(key, label, href) {
  return `${key}:\n  text: ${label}\n  link: ${href}\n`;
}

function loadTopic(topicId) {
  const topicDir = path.join(THEORY, topicId);
  const metaPath = path.join(topicDir, "sections.json");
  if (!fs.existsSync(metaPath)) {
    return null;
  }
  const meta = JSON.parse(fs.readFileSync(metaPath, "utf8"));
  const order = meta.order ?? [];
  const pages = meta.pages ?? {};
  return { topicId, topicDir, meta, order, pages };
}

function pageBase(topicId, pageId) {
  return `/learn/${topicId}/${pageId}`;
}

function navLabel(pageId, pageMeta) {
  if (pageMeta.prevLabel) {
    return pageMeta.prevLabel;
  }
  if (pageMeta.sidebar) {
    return pageMeta.sidebar.toLowerCase();
  }
  return pageId.replace(/_/g, "-");
}

function syncTopic(topic) {
  const { topicId, topicDir, order, pages } = topic;
  const destDir = path.join(DOCS_LEARN, topicId);
  ensureDir(destDir);

  const synced = new Set();

  for (let i = 0; i < order.length; i++) {
    const pageId = order[i];
    const pageMeta = pages[pageId] ?? {};
    const srcPath = path.join(topicDir, `${pageId}.md`);
    if (!fs.existsSync(srcPath)) {
      console.warn(`Missing ${srcPath}`);
      continue;
    }

    const body = fs.readFileSync(srcPath, "utf8");
    const base = pageBase(topicId, pageId);

    let fm = `---
title: ${(pageMeta.title ?? pageId).replace(/"/g, '\\"')}
`;
    if (pageMeta.description) {
      fm += `description: ${pageMeta.description.replace(/"/g, '\\"')}\n`;
    }

    const prevId = i > 0 ? order[i - 1] : null;
    const nextId = i < order.length - 1 ? order[i + 1] : null;

    if (pageMeta.prev) {
      fm += yamlLink("prev", pageMeta.prev.text, pageMeta.prev.link);
    } else if (prevId) {
      const prevMeta = pages[prevId] ?? {};
      fm += yamlLink(
        "prev",
        navLabel(prevId, prevMeta),
        pageBase(topicId, prevId),
      );
    }

    if (pageMeta.next) {
      fm += yamlLink("next", pageMeta.next.text, pageMeta.next.link);
    } else if (nextId) {
      const nextMeta = pages[nextId] ?? {};
      fm += yamlLink(
        "next",
        navLabel(nextId, nextMeta),
        pageBase(topicId, nextId),
      );
    }

    fm += `---\n\n`;

    fs.writeFileSync(path.join(destDir, `${pageId}.md`), fm + body);
    synced.add(`${pageId}.md`);
  }

  for (const name of fs.readdirSync(destDir)) {
    if (name.endsWith(".md") && !synced.has(name)) {
      fs.unlinkSync(path.join(destDir, name));
    }
  }

  return order.map((pageId) => {
    const pageMeta = pages[pageId] ?? {};
    return {
      text: pageMeta.sidebar ?? pageMeta.title ?? pageId,
      link: pageBase(topicId, pageId),
    };
  });
}

function discoverTopics() {
  if (!fs.existsSync(THEORY)) {
    return [];
  }
  return fs
    .readdirSync(THEORY, { withFileTypes: true })
    .filter((d) => d.isDirectory() && !d.name.startsWith("_"))
    .map((d) => d.name)
    .filter((id) => loadTopic(id));
}

const sidebarPayload = {};
for (const topicId of discoverTopics()) {
  const topic = loadTopic(topicId);
  sidebarPayload[topicId] = syncTopic(topic);
}

const sidebarOut = path.join(ROOT, "docs", ".vitepress", "sidebar-theory.json");
fs.writeFileSync(sidebarOut, JSON.stringify(sidebarPayload, null, 2));

console.log("Synced theory from theory/ → docs/learn/");
console.log(`Wrote ${sidebarOut}`);
