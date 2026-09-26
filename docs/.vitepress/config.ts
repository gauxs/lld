import path from "node:path";
import fs from "node:fs";
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
  : { connectFour: [], rateLimiter: [] };

const sidebarTheoryPath = path.join(__dirname, "sidebar-theory.json");
const sidebarTheory = fs.existsSync(sidebarTheoryPath)
  ? JSON.parse(fs.readFileSync(sidebarTheoryPath, "utf8"))
  : { concurrency: [] };

const concurrencyTheoryItems = sidebarTheory.concurrency ?? [];

const connectFourItems = [
  { text: "Overview", link: "/problems/connect-four/" },
  ...(sidebarProblems.connectFour ?? []),
];

const rateLimiterItems = [
  { text: "Overview", link: "/problems/rate-limiter/" },
  ...(sidebarProblems.rateLimiter ?? []),
];

const sidebar = [
  {
    items: [
      { text: "Home", link: "/" },
      {
        text: "Theory",
        collapsed: false,
        items: [
          {
            text: "Concurrency",
            collapsed: false,
            items: concurrencyTheoryItems,
          },
        ],
      },
      {
        text: "Problems",
        collapsed: false,
        items: [
          {
            text: "Connect Four",
            collapsed: false,
            items: connectFourItems,
          },
          {
            text: "Rate limiter",
            collapsed: false,
            items: rateLimiterItems,
          },
        ],
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
      "Low-level design: concurrency theory and problem trails with Go code",
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
      plugins: [devRootRedirectPlugin()],
      server: {
        fs: { allow: [ROOT] },
      },
      optimizeDeps: {
        needsInterop: ["fastdom"],
        include: ["fastdom", "mermaid"],
      },
      resolve: {
        alias: [
          {
            find: path.resolve(
              __dirname,
              "../../node_modules/vitepress/dist/client/theme-default/components/VPSidebarItem.vue",
            ),
            replacement: path.join(
              __dirname,
              "theme/components/VPSidebarItem.vue",
            ),
          },
          {
            find: path.resolve(
              __dirname,
              "../../node_modules/vitepress/dist/client/theme-default/components/VPNavBarTitle.vue",
            ),
            replacement: path.join(
              __dirname,
              "theme/components/VPNavBarTitle.vue",
            ),
          },
        ],
      },
    },
    themeConfig: {
      logo: {
        light: SITE_ICON_LIGHT,
        dark: SITE_ICON_DARK,
        alt: "LLDZen",
      },
      logoLink: "/",
      nav: [
        {
          text: "Theory",
          link: "/learn/concurrency/intro",
          activeMatch: "/learn/",
        },
        {
          text: "Problems",
          link: "/problems/connect-four/extensions/baseline/requirements",
          activeMatch: "/problems/",
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
