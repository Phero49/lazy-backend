Tailwind + Quasar setup

This project ships with Tailwind CSS configured to use the prefix `tw-` to avoid class name collisions with Quasar's CSS.

What I changed

- Added `postcss.config.cjs` to run Tailwind via PostCSS (Vite will pick it up).
- Removed the nonstandard Tailwind Vite plugin from `vite.config.ts` so PostCSS handles Tailwind processing.
- `tailwind.config.js` already sets `prefix: 'tw-'`.

How to use

- In your components or templates, use Tailwind classes with the `tw-` prefix. Example:

```html
<div class="tw-bg-blue-500 tw-text-white tw-p-4">Tailwind box</div>
```

- Quasar classes (e.g., `q-btn`, `q-pa-md`) remain unchanged and will not be affected by Tailwind.

Dev commands

Run these from `frontend/`:

```bash
npm install
npm run dev
```

Notes

- If you prefer no prefix, remove `prefix` from `tailwind.config.js` but watch for style collisions.
- If Tailwind styles don't seem to apply, ensure `./src/**/*.{vue,js,ts,jsx,tsx}` matches your file locations and restart the dev server.
