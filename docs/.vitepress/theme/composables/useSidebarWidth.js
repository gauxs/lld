import { onMounted, onUnmounted, nextTick } from "vue";

const DESKTOP_MQ = "(min-width: 960px)";

/**
 * Set --vp-sidebar-width from the widest sidebar label (any depth, including
 * collapsed branches). Width is fixed for the current page until viewport
 * resize or route change — expanding sections does not remeasure.
 */
export function useSidebarWidthSync() {
  let mediaQuery = null;
  let raf = 0;
  let measuring = false;
  let lockedWidthPx = null;

  function measureContentWidth(sidebar) {
    measuring = true;
    sidebar.style.setProperty("width", "max-content", "important");
    sidebar.style.setProperty("max-width", "none", "important");

    const hiddenItems = sidebar.querySelectorAll(
      ".VPSidebarItem.collapsed > .items",
    );
    for (const el of hiddenItems) {
      el.style.setProperty("display", "block", "important");
    }

    const nav = sidebar.querySelector("#VPSidebarNav");
    let content = 0;
    if (nav) {
      for (const el of nav.querySelectorAll(".item")) {
        content = Math.max(content, el.scrollWidth);
      }
      for (const el of nav.querySelectorAll(".text")) {
        content = Math.max(content, el.scrollWidth);
      }
      content = Math.max(content, nav.scrollWidth);
    }

    for (const el of hiddenItems) {
      el.style.removeProperty("display");
    }

    const styles = getComputedStyle(sidebar);
    const pad =
      (parseFloat(styles.paddingLeft) || 0) +
      (parseFloat(styles.paddingRight) || 0);

    const width = Math.ceil(content + pad);

    sidebar.style.removeProperty("width");
    sidebar.style.removeProperty("max-width");
    measuring = false;
    return width;
  }

  function sync(force = false) {
    if (typeof document === "undefined" || measuring) return;

    const root = document.documentElement;
    if (!window.matchMedia(DESKTOP_MQ).matches) {
      lockedWidthPx = null;
      root.style.removeProperty("--vp-sidebar-width");
      return;
    }

    const sidebar = document.querySelector(".VPSidebar");
    if (!sidebar) return;

    if (!force && lockedWidthPx != null) {
      root.style.setProperty("--vp-sidebar-width", `${lockedWidthPx}px`);
      return;
    }

    lockedWidthPx = measureContentWidth(sidebar);
    root.style.setProperty("--vp-sidebar-width", `${lockedWidthPx}px`);
  }

  function scheduleSync(force = false) {
    if (measuring) return;
    cancelAnimationFrame(raf);
    raf = requestAnimationFrame(() => sync(force));
  }

  function onMediaChange() {
    lockedWidthPx = null;
    nextTick(() => scheduleSync(true));
  }

  function onWindowResize() {
    lockedWidthPx = null;
    scheduleSync(true);
  }

  onMounted(() => {
    mediaQuery = window.matchMedia(DESKTOP_MQ);
    mediaQuery.addEventListener("change", onMediaChange);
    window.addEventListener("resize", onWindowResize);
    nextTick(() => scheduleSync(true));
  });

  onUnmounted(() => {
    cancelAnimationFrame(raf);
    mediaQuery?.removeEventListener("change", onMediaChange);
    window.removeEventListener("resize", onWindowResize);
  });

  return {
    sync: () => scheduleSync(true),
  };
}
