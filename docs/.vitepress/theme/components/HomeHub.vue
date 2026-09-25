<script setup>
import { computed } from "vue";
import { useData } from "vitepress";

const { site, isDark } = useData();
const base = site.value.base;
const logoSrc = computed(
  () => `${base}icons/${isDark.value ? "site-dark" : "site"}.svg`,
);

const cards = [
  {
    title: "Theory",
    description:
      "Concurrency fundamentals—shared memory, correctness, coordination, and scarcity.",
    link: `${base}learn/concurrency/intro`,
    meta: "4 articles",
  },
  {
    title: "Problems",
    description:
      "Requirement-first trails with variations and Go reference code.",
    link: `${base}problems/connect-four/`,
    meta: "Connect Four",
  },
];
</script>

<template>
  <div class="zen-home">
    <header class="zen-home-header">
      <img
        class="zen-home-logo"
        :src="logoSrc"
        width="56"
        height="56"
        alt=""
      />
      <h1 class="zen-home-title">LLDZen</h1>
      <p class="zen-home-tagline">
        Low-level design notes—multi-page theory trails and problem walkthroughs.
      </p>
    </header>

    <div class="zen-home-cards">
      <a
        v-for="card in cards"
        :key="card.title"
        class="zen-home-card"
        :href="card.link"
      >
        <span class="zen-home-card-meta">{{ card.meta }}</span>
        <h2 class="zen-home-card-title">{{ card.title }}</h2>
        <p class="zen-home-card-desc">{{ card.description }}</p>
        <span class="zen-home-card-cta">Open trail →</span>
      </a>
    </div>
  </div>
</template>
