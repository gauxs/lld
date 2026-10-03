<script setup>
import { computed } from "vue";
import { useData } from "vitepress";
import practiceNavigation from "../../sidebar-practice.json";
import problemNavigation from "../../sidebar-problems.json";
import theoryNavigation from "../../sidebar-theory.json";

const { site, isDark } = useData();
const base = site.value.base;
const logoSrc = computed(
  () => `${base}icons/${isDark.value ? "site-dark" : "site"}.svg`,
);

function withBase(link) {
  return `${base}${link.replace(/^\//, "")}`;
}

const cards = [
  {
    title: "Theory",
    description:
      "Concept trails covering the foundations behind low-level design.",
    link: withBase(theoryNavigation.firstLink),
    meta: `${theoryNavigation.articleCount} articles`,
  },
  {
    title: "Problems",
    description:
      "Requirement-first extensions with optional Go reference code.",
    link: withBase(problemNavigation.firstLink),
    meta: `${problemNavigation.problemCount} problems`,
  },
  {
    title: "Practice with AI",
    description:
      "Turn a chatbot into a structured low-level design interviewer.",
    link: withBase(practiceNavigation.firstLink),
    meta: `${practiceNavigation.pageCount} prompt${practiceNavigation.pageCount === 1 ? "" : "s"}`,
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
        Low-level design notes, problem walkthroughs, and guided AI practice.
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
