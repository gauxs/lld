# LLDZen

Go reference implementations and interview-style write-ups. The static site is built with VitePress under `docs/`.

## Site

```bash
npm install
npm run docs:dev
```

Production build: `npm run docs:build` (output in `docs/.vitepress/dist`).

- **Learn:** concurrency trail under `docs/learn/concurrency/`
- **Problems:** multi-page trails under `docs/problems/` (Connect Four first)

Authoring template for new problems: [`_template/requirements.md`](_template/requirements.md).

## Code packages

| Package | Description |
| --- | --- |
| `connect_four/` | Two-player Connect Four (variation 1) |
| `rate_limiter/` | Fixed-window rate limiter (WIP) |
