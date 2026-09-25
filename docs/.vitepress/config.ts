import { defineConfig, type UserConfig } from "vitepress";
import { SITE_ICON_DARK, SITE_ICON_LIGHT } from "./site-assets";

const PROD_BASE = "/lld/";

/** Single sidebar everywhere: Home, full Theory tree, full Problems tree. */
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
              { text: "Overview", link: "/problems/connect-four/" },
              {
                text: "Variation 1 — In-memory OOP",
                collapsed: false,
                items: [
                  {
                    text: "Summary",
                    link: "/problems/connect-four/variation-1/",
                  },
                  {
                    text: "Requirements",
                    link: "/problems/connect-four/variation-1/requirements",
                  },
                  {
                    text: "Design",
                    link: "/problems/connect-four/variation-1/design",
                  },
                  {
                    text: "Codebase",
                    link: "/problems/connect-four/variation-1/codebase",
                  },
                ],
              },
              {
                text: "Follow-ups",
                link: "/problems/connect-four/follow-ups",
              },
            ],
          },
        ],
      },
    ],
  },
];

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
          link: "/problems/connect-four/",
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
