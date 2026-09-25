import path from "node:path";
import { fileURLToPath } from "node:url";
import { defineConfig, type UserConfig } from "vitepress";
import { SITE_ICON_DARK, SITE_ICON_LIGHT } from "./site-assets";

const __dirname = path.dirname(fileURLToPath(import.meta.url));
const ROOT = path.resolve(__dirname, "../..");
const PROD_BASE = "/lld/";

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
            items: [
              { text: "Introduction", link: "/learn/concurrency/intro" },
              { text: "Correctness", link: "/learn/concurrency/correctness" },
              { text: "Coordination", link: "/learn/concurrency/coordination" },
              { text: "Scarcity", link: "/learn/concurrency/scarcity" },
            ],
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
            items: [
              {
                text: "Functional requirements",
                link: "/problems/connect-four/functional-requirement",
              },
              { text: "Design", link: "/problems/connect-four/design" },
              { text: "Codebase", link: "/problems/connect-four/codebase" },
              { text: "Follow-up", link: "/problems/connect-four/followup" },
            ],
          },
        ],
      },
    ],
  },
];

function devRootRedirectPlugin(base: string) {
  return {
    name: "lld-dev-root-redirect",
    configureServer(server) {
      server.middlewares.use((req, res, next) => {
        const pathname = req.url?.split("?")[0] ?? "";
        if (base === "/" && (pathname === "/" || pathname === "")) {
          next();
          return;
        }
        if (pathname === "/" || pathname === "") {
          res.statusCode = 302;
          res.setHeader("Location", base);
          res.end();
          return;
        }
        next();
      });
    },
  };
}

export default defineConfig(({ command }) => {
  const base = command === "build" ? PROD_BASE : "/";

  const config: UserConfig = {
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
      plugins: [devRootRedirectPlugin(base)],
      server: {
        fs: { allow: [ROOT] },
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
          link: "/problems/connect-four/functional-requirement",
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
  };

  return config;
});
