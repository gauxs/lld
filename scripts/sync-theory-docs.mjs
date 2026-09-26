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

function humanizeName(name) {
  return name
    .replace(/\.md$/i, "")
    .replace(/^\d+[-_]/, "")
    .replace(/_/g, " ")
    .replace(/\b\w/g, (c) => c.toUpperCase());
}

function humanizeDir(dirName) {
  return humanizeName(dirName);
}

/** Sort key from leading `01-` / `01_` on file or folder name; unnumbered sorts last. */
function orderKey(name) {
  const num = path.basename(name).match(/^(\d+)[-_]/);
  return num ? parseInt(num[1], 10) : 5000;
}

function sortByOrderPrefix(names) {
  return [...names].sort((a, b) => {
    const oa = orderKey(a);
    const ob = orderKey(b);
    if (oa !== ob) {
      return oa - ob;
    }
    return a.localeCompare(b, undefined, { numeric: true, sensitivity: "base" });
  });
}

function sortRelPaths(relPaths) {
  return [...relPaths].sort((a, b) => {
    const oa = orderKey(path.basename(a));
    const ob = orderKey(path.basename(b));
    if (oa !== ob) {
      return oa - ob;
    }
    return a.localeCompare(b, undefined, { numeric: true, sensitivity: "base" });
  });
}

function parseSimpleFrontMatter(raw) {
  if (!raw.startsWith("---\n")) {
    return { meta: {}, body: raw };
  }
  const end = raw.indexOf("\n---\n", 4);
  if (end === -1) {
    return { meta: {}, body: raw };
  }
  const block = raw.slice(4, end);
  const body = raw.slice(end + 5);
  const meta = {};
  let key = null;

  for (const line of block.split("\n")) {
    const top = line.match(/^(\w+):\s*(.*)$/);
    if (top) {
      key = top[1];
      const rest = top[2].trim();
      if (rest.startsWith("{")) {
        try {
          meta[key] = JSON.parse(rest);
        } catch {
          meta[key] = rest;
        }
      } else if (rest) {
        meta[key] = rest;
      } else {
        meta[key] = {};
      }
      continue;
    }
    const sub = line.match(/^\s{2}(\w+):\s*(.*)$/);
    if (sub && key && typeof meta[key] === "object") {
      meta[key][sub[1]] = sub[2].trim();
    }
  }
  return { meta, body };
}

function titleFromBody(body) {
  const m = body.match(/^#\s+(.+)/m);
  return m?.[1]?.trim() ?? null;
}

function collectMarkdownFiles(dir, prefix = "") {
  const out = [];
  if (!fs.existsSync(dir)) {
    return out;
  }
  for (const ent of fs.readdirSync(dir, { withFileTypes: true })) {
    if (ent.name.startsWith("_") || ent.name.startsWith(".")) {
      continue;
    }
    const rel = prefix ? `${prefix}/${ent.name}` : ent.name;
    const full = path.join(dir, ent.name);
    if (ent.isDirectory()) {
      out.push(...collectMarkdownFiles(full, rel));
      continue;
    }
    if (ent.isFile() && ent.name.endsWith(".md")) {
      out.push(rel.replace(/\\/g, "/"));
    }
  }
  return out;
}

function listRootEntries(topicDir) {
  const files = [];
  const dirs = [];
  for (const ent of fs.readdirSync(topicDir, { withFileTypes: true })) {
    if (ent.name.startsWith("_") || ent.name.startsWith(".")) {
      continue;
    }
    if (ent.isDirectory()) {
      dirs.push(ent.name);
    } else if (ent.isFile() && ent.name.endsWith(".md")) {
      files.push(ent.name);
    }
  }

  const entries = [
    ...files.map((name) => ({ type: "file", name })),
    ...dirs.map((name) => ({ type: "dir", name })),
  ];

  entries.sort((a, b) => {
    const oa = orderKey(a.name);
    const ob = orderKey(b.name);
    if (oa !== ob) {
      return oa - ob;
    }
    return a.name.localeCompare(b.name, undefined, {
      numeric: true,
      sensitivity: "base",
    });
  });

  return entries;
}

function readingOrder(allRelPaths, rootEntries) {
  const order = [];
  for (const entry of rootEntries) {
    if (entry.type === "file") {
      if (allRelPaths.includes(entry.name)) {
        order.push(entry.name);
      }
      continue;
    }
    const inDir = allRelPaths.filter((p) => p.startsWith(`${entry.name}/`));
    order.push(...sortRelPaths(inDir));
  }
  return order;
}

function pageIdFromRel(rel) {
  return rel.replace(/\.md$/i, "");
}

function pageLink(topicId, pageId) {
  return `/learn/${topicId}/${pageId}`;
}

function resolvePageMeta(rel, rawBody) {
  const { meta: srcFm, body } = parseSimpleFrontMatter(rawBody);
  const pageId = pageIdFromRel(rel);

  const title =
    srcFm.title ?? titleFromBody(body) ?? humanizeName(path.basename(rel));

  return {
    pageId,
    rel,
    body,
    title,
    description: srcFm.description ?? null,
    sidebar: srcFm.sidebar ?? title,
    prev: srcFm.prev ?? null,
    next: srcFm.next ?? null,
    prevLabel: srcFm.prevLabel ?? null,
  };
}

function navLabel(page) {
  if (page.prevLabel) {
    return page.prevLabel;
  }
  const s = page.sidebar ?? page.title;
  return typeof s === "string" ? s.toLowerCase() : page.pageId.split("/").pop();
}

function buildSidebar(topicId, allRelPaths, pagesById, rootEntries) {
  const items = [];

  for (const entry of rootEntries) {
    if (entry.type === "file") {
      if (!allRelPaths.includes(entry.name)) {
        continue;
      }
      const p = pagesById.get(pageIdFromRel(entry.name));
      items.push({
        text: p.sidebar,
        link: pageLink(topicId, p.pageId),
      });
      continue;
    }

    const groupFiles = allRelPaths.filter((p) =>
      p.startsWith(`${entry.name}/`),
    );
    if (!groupFiles.length) {
      continue;
    }
    items.push({
      text: humanizeDir(entry.name),
      collapsed: false,
      items: sortRelPaths(groupFiles).map((rel) => {
        const p = pagesById.get(pageIdFromRel(rel));
        return {
          text: p.sidebar,
          link: pageLink(topicId, p.pageId),
        };
      }),
    });
  }

  return items;
}

function syncTopic(topicId) {
  const topicDir = path.join(THEORY, topicId);
  const allRelPaths = collectMarkdownFiles(topicDir);
  if (!allRelPaths.length) {
    return [];
  }

  const rootEntries = listRootEntries(topicDir);
  const order = readingOrder(allRelPaths, rootEntries);
  const pagesById = new Map();

  for (const rel of allRelPaths) {
    const srcPath = path.join(topicDir, rel);
    const raw = fs.readFileSync(srcPath, "utf8");
    pagesById.set(pageIdFromRel(rel), resolvePageMeta(rel, raw));
  }

  const destRoot = path.join(DOCS_LEARN, topicId);
  ensureDir(destRoot);
  const synced = new Set();

  for (let i = 0; i < order.length; i++) {
    const rel = order[i];
    const page = pagesById.get(pageIdFromRel(rel));
    const destPath = path.join(destRoot, rel);
    ensureDir(path.dirname(destPath));

    let fm = `---
title: ${String(page.title).replace(/"/g, '\\"')}
`;
    if (page.description) {
      fm += `description: ${String(page.description).replace(/"/g, '\\"')}\n`;
    }

    const prevRel = i > 0 ? order[i - 1] : null;
    const nextRel = i < order.length - 1 ? order[i + 1] : null;

    if (page.prev) {
      fm += yamlLink("prev", page.prev.text, page.prev.link);
    } else if (prevRel) {
      const prev = pagesById.get(pageIdFromRel(prevRel));
      fm += yamlLink(
        "prev",
        navLabel(prev),
        pageLink(topicId, prev.pageId),
      );
    }

    if (page.next) {
      fm += yamlLink("next", page.next.text, page.next.link);
    } else if (nextRel) {
      const next = pagesById.get(pageIdFromRel(nextRel));
      fm += yamlLink(
        "next",
        navLabel(next),
        pageLink(topicId, next.pageId),
      );
    }

    fm += `---\n\n`;
    fs.writeFileSync(destPath, fm + page.body);
    synced.add(rel.replace(/\\/g, "/"));
  }

  function walkClean(dir, prefix = "") {
    for (const name of fs.readdirSync(dir, { withFileTypes: true })) {
      const rel = prefix ? `${prefix}/${name.name}` : name.name;
      const full = path.join(dir, name.name);
      if (name.isDirectory()) {
        walkClean(full, rel);
        if (fs.readdirSync(full).length === 0) {
          fs.rmdirSync(full);
        }
        continue;
      }
      if (name.name.endsWith(".md") && !synced.has(rel.replace(/\\/g, "/"))) {
        fs.unlinkSync(full);
      }
    }
  }
  walkClean(destRoot);

  return buildSidebar(topicId, allRelPaths, pagesById, rootEntries);
}

function discoverTopics() {
  if (!fs.existsSync(THEORY)) {
    return [];
  }
  return fs
    .readdirSync(THEORY, { withFileTypes: true })
    .filter((d) => d.isDirectory() && !d.name.startsWith("_"))
    .map((d) => d.name)
    .filter((id) => collectMarkdownFiles(path.join(THEORY, id)).length > 0);
}

const sidebarPayload = {};
for (const topicId of discoverTopics()) {
  sidebarPayload[topicId] = syncTopic(topicId);
}

const sidebarOut = path.join(ROOT, "docs", ".vitepress", "sidebar-theory.json");
fs.writeFileSync(sidebarOut, JSON.stringify(sidebarPayload, null, 2));

console.log("Synced theory from theory/ → docs/learn/");
console.log(`Wrote ${sidebarOut}`);
