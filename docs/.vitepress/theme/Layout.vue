<script setup>
import { computed } from "vue";
import { useData } from "vitepress";
import DefaultLayout from "vitepress/dist/client/theme-default/Layout.vue";
import HomeHub from "./components/HomeHub.vue";

const { page } = useData();

const isHome = computed(() => page.value.relativePath === "index.md");
</script>

<template>
  <DefaultLayout>
    <template v-if="isHome" #doc-before>
      <HomeHub />
    </template>
  </DefaultLayout>
</template>
