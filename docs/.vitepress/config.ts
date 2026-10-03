import path from "node:path";
import fs from "node:fs";
import { execFile } from "node:child_process";
import { fileURLToPath } from "node:url";
import { defineConfig, type UserConfig } from "vitepress";
import { withMermaid } from "vitepress-mermaid-viewer";
import { SITE_ICON_DARK, SITE_ICON_LIGHT } from "./site-assets";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.resolve(__dirname, "../..");
const BASE = "/lld";

const sidebarProblemsPath = path.join(__dirname, "sidebar-problems.json");
const sidebarProblems = fs.existsSync(sidebarProblemsPath)
  ? JSON.parse(fs.readFileSync(sidebarProblemsPath, "utf8"))
  : { items: [], firstLink: "/", problemCount: 0 };

const sidebarTheoryPath = path.join(__dirname, "sidebar-theory.json");
const sidebarTheory = fs.existsSync(sidebarTheoryPath)
  ? JSON.parse(fs.readFileSync(sidebarTheoryPath, "utf8"))
  : { items: [], firstLink: "/", articleCount: 0, topicCount: 0 };

const sidebarPracticePath = path.join(__dirname, "sidebar-practice.json");
const sidebarPractice = fs.existsSync(sidebarPracticePath)
  ? JSON.parse(fs.readFileSync(sidebarPracticePath, "utf8"))
  : { items: [], firstLink: "/", pageCount: 0 };

const sidebar = [
  {
    items: [
      { text: "Home", link: "/" },
      {
        text: "Theory",
        collapsed: true,
        items: sidebarTheory.items ?? [],
      },
      {
        text: "Problems",
        collapsed: true,
        items: sidebarProblems.items ?? [],
      },
      {
        text: "Practice with AI",
        collapsed: true,
        items: sidebarPractice.items ?? [],
      },
    ],
  },
];

function devRootRedirectPlugin() {
  return {
    name: "lld-dev-root-redirect",
    configureServer(server) {
      server.middlewares.use((req, res, next) => {
        const pathname = req.url?.split("?")[0] ?? "";
        if (pathname === "/" || pathname === "") {
          res.statusCode = 302;
          res.setHeader("Location", `${BASE}/`);
          res.end();
          return;
        }
        next();
      });
    },
  };
}

function authoringSyncPlugin() {
  const theoryRoot = path.join(ROOT, "theory");
  const problemsRoot = path.join(ROOT, "problems");
  const practiceRoot = path.join(ROOT, "practice_with_ai");
  const theoryScript = path.join(ROOT, "scripts", "sync-theory-docs.mjs");
  const problemsScript = path.join(ROOT, "scripts", "sync-problem-docs.mjs");
  const practiceScript = path.join(ROOT, "scripts", "sync-practice-docs.mjs");
  let timer: ReturnType<typeof setTimeout> | undefined;
  const pendingScripts = new Set<string>();
  let syncQueue = Promise.resolve();

  return {
    name: "lld-authoring-sync",
    configureServer(server) {
      function runScript(script: string) {
        return new Promise<void>((resolve, reject) => {
          execFile(process.execPath, [script], (error, stdout, stderr) => {
            if (stdout) {
              server.config.logger.info(stdout.trim());
            }
            if (stderr) {
              server.config.logger.warn(stderr.trim());
            }
            if (error) {
              reject(error);
              return;
            }
            resolve();
          });
        });
      }

      function flushPendingScripts() {
        const scripts = [...pendingScripts];
        pendingScripts.clear();
        syncQueue = syncQueue
          .then(async () => {
            let completed = false;
            for (const script of scripts) {
              try {
                await runScript(script);
                completed = true;
              } catch (error) {
                server.config.logger.error(error.message);
              }
            }
            if (completed) {
              await server.restart();
            }
          })
          .catch((error) => {
            server.config.logger.error(error.message);
          });
      }

      server.watcher.add([theoryRoot, problemsRoot, practiceRoot]);
      server.watcher.on("all", (_event, file) => {
        const script = file.startsWith(theoryRoot)
          ? theoryScript
          : file.startsWith(problemsRoot)
            ? problemsScript
            : file.startsWith(practiceRoot)
              ? practiceScript
            : null;
        if (!script) {
          return;
        }

        pendingScripts.add(script);
        clearTimeout(timer);
        timer = setTimeout(flushPendingScripts, 100);
      });
    },
  };
}

export default defineConfig(() => {
  const base = `${BASE}/`;

  return withMermaid({
    mermaid: {
      theme: "base",
      // Only size-related vars here — light text/edge colors fight Mermaid's `dark` theme
      // (edge labels like "Instinct" / "Better Approach" stay readable in dark mode).
      themeVariables: {
        fontSize: "16px",
      },
      flowchart: {
        useMaxWidth: false,
        htmlLabels: true,
        padding: 16,
        nodeSpacing: 48,
        rankSpacing: 56,
      },
    },
    title: "LLDZen",
    description:
      "Low-level design theory, problem trails, and AI interview practice",
    base,
    cleanUrls: true,
    head: [
      [
        "link",
        {
          rel: "icon",
          type: "image/svg+xml",
          href: `${base}${SITE_ICON_LIGHT.replace(/^\//, "")}`,
          media: "(prefers-color-scheme: light)",
        },
      ],
      [
        "link",
        {
          rel: "icon",
          type: "image/svg+xml",
          href: `${base}${SITE_ICON_DARK.replace(/^\//, "")}`,
          media: "(prefers-color-scheme: dark)",
        },
      ],
    ],
    vite: {
      plugins: [devRootRedirectPlugin(), authoringSyncPlugin()],
      server: {
        fs: { allow: [ROOT] },
      },
      optimizeDeps: {
        needsInterop: ["fastdom"],
        include: ["fastdom", "mermaid"],
      },
    },
    themeConfig: {
      logo: {
        light: SITE_ICON_LIGHT,
        dark: SITE_ICON_DARK,
        alt: "LLDZen",
      },
      logoLink: `${BASE}/`,
      nav: [
        {
          text: "Theory",
          link: sidebarTheory.firstLink,
          activeMatch: "/learn/",
        },
        {
          text: "Problems",
          link: sidebarProblems.firstLink,
          activeMatch: "/problems/",
        },
        {
          text: "Practice with AI",
          link: sidebarPractice.firstLink,
          activeMatch: "/practice-with-ai/",
        },
      ],
      sidebar,
      outline: { level: [2, 3] },
      docFooter: { prev: "Previous", next: "Next" },
      socialLinks: [
        {
          icon: "github",
          link: "https://github.com/gauxs/lld",
        },
      ],
    },
  } satisfies UserConfig);
});
