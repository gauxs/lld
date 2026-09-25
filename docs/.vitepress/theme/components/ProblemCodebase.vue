<script setup>
import { computed } from "vue";

const props = defineProps({
  problem: { type: String, required: true },
});

const modules = import.meta.glob(
  "../../../../problems/*/code/**/*.go",
  { query: "?raw", import: "default", eager: true },
);

const files = computed(() => {
  const prefix = `../../../../problems/${props.problem}/code/`;
  return Object.entries(modules)
    .filter(([path]) => path.startsWith(prefix))
    .map(([path, content]) => ({
      path: path.slice(prefix.length),
      content,
    }))
    .sort((a, b) => a.path.localeCompare(b.path));
});
</script>

<template>
  <div class="problem-codebase">
    <p v-if="files.length === 0" class="problem-codebase-empty">
      No source files found for <code>{{ problem }}</code>.
    </p>
    <section
      v-for="file in files"
      :key="file.path"
      class="problem-codebase-file"
    >
      <h2 class="problem-codebase-path">{{ file.path }}</h2>
      <pre class="problem-codebase-pre"><code>{{ file.content }}</code></pre>
    </section>
  </div>
</template>
