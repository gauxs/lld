import type { Theme } from "vitepress";
import DefaultTheme from "vitepress/theme";
import Mermaid from "vitepress-plugin-mermaid/Mermaid.vue";
import Layout from "./Layout.vue";
import ProblemCodebase from "./components/ProblemCodebase.vue";
import "./custom.css";

export default {
  extends: DefaultTheme,
  Layout,
  enhanceApp({ app }) {
    app.component("Mermaid", Mermaid);
    app.component("ProblemCodebase", ProblemCodebase);
  },
} satisfies Theme;
