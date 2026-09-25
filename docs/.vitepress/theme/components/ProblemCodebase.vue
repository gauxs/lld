<script setup>
import { computed, onMounted, onUnmounted, ref } from "vue";

const props = defineProps({
  problem: { type: String, required: true },
  extension: { type: String, default: "baseline" },
});

const modules = import.meta.glob(
  "../../../../problems/**/code/**/*.go",
  { query: "?raw", import: "default", eager: true },
);

const activeId = ref("");

const files = computed(() => {
  const extPrefix = `../../../../problems/${props.problem}/extensions/${props.extension}/code/`;
  const legacyPrefix = `../../../../problems/${props.problem}/code/`;
  const prefix = Object.keys(modules).some((p) => p.startsWith(extPrefix))
    ? extPrefix
    : legacyPrefix;
  return Object.entries(modules)
    .filter(([p]) => p.startsWith(prefix))
    .map(([p, content]) => {
      const rel = p.slice(prefix.length);
      const id = rel.replace(/[^\w]+/g, "-").replace(/^-|-$/g, "");
      const shortName = rel.split("/").pop();
      return { path: rel, shortName, id, content };
    })
    .sort((a, b) => a.path.localeCompare(b.path));
});

function onScroll() {
  const offset = 120;
  let current = files.value[0]?.id ?? "";
  for (const file of files.value) {
    const el = document.getElementById(`file-${file.id}`);
    if (el && el.getBoundingClientRect().top <= offset) {
      current = file.id;
    }
  }
  activeId.value = current;
}

onMounted(() => {
  activeId.value = files.value[0]?.id ?? "";
  window.addEventListener("scroll", onScroll, { passive: true });
  onScroll();
});

onUnmounted(() => {
  window.removeEventListener("scroll", onScroll);
});
</script>

<template>
  <div class="lld-codebase">
    <div class="lld-codebase-main">
      <p v-if="files.length === 0" class="lld-codebase-empty">
        No source files found for <code>{{ problem }}</code>.
      </p>
      <section
        v-for="file in files"
        :id="`file-${file.id}`"
        :key="file.path"
        class="lld-codebase-file"
      >
        <h2 class="lld-codebase-path">{{ file.path }}</h2>
        <pre class="lld-codebase-pre"><code>{{ file.content }}</code></pre>
      </section>
    </div>

    <nav v-if="files.length" class="lld-codebase-toc" aria-label="Source files">
      <p class="lld-codebase-toc-title">Files</p>
      <a
        v-for="file in files"
        :key="file.id"
        class="lld-codebase-toc-link"
        :class="{ active: activeId === file.id }"
        :href="`#file-${file.id}`"
      >
        {{ file.shortName }}
      </a>
    </nav>
  </div>
</template>
