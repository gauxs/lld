import type { Theme } from "vitepress";
import DefaultTheme from "vitepress/theme";
import Layout from "./Layout.vue";
import ProblemCodebase from "./components/ProblemCodebase.vue";
import "./custom.css";

export default {
  extends: DefaultTheme,
  Layout,
  enhanceApp({ app }) {
    app.component("ProblemCodebase", ProblemCodebase);
  },
} satisfies Theme;
