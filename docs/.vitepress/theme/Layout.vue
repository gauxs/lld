<script setup>
import { computed, nextTick, watch } from "vue";
import { useData, useRoute } from "vitepress";
import DefaultLayout from "vitepress/dist/client/theme-default/Layout.vue";
import HomeHub from "./components/HomeHub.vue";
import { useSidebarWidthSync } from "./composables/useSidebarWidth.js";

const { page } = useData();
const route = useRoute();
const { sync: syncSidebarWidth } = useSidebarWidthSync();

const isHome = computed(() => page.value.relativePath === "index.md");

watch(
  () => route.path,
  () => nextTick(() => syncSidebarWidth()),
);
</script>

<template>
  <DefaultLayout>
    <template v-if="isHome" #doc-before>
      <HomeHub />
    </template>
  </DefaultLayout>
</template>
