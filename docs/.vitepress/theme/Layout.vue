<script setup>
import { computed } from "vue";
import { useData } from "vitepress";
import DefaultTheme from "vitepress/theme";
import HomeHub from "./components/HomeHub.vue";

const { page } = useData();
const DefaultLayout = DefaultTheme.Layout;

const isHome = computed(() => page.value.relativePath === "index.md");
</script>

<template>
  <DefaultLayout>
    <template v-if="isHome" #doc-before>
      <HomeHub />
    </template>
  </DefaultLayout>
</template>
