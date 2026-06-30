// Use the dedicated PostCSS adapter for Tailwind per new packaging.
module.exports = {
  plugins: {
    // the @tailwindcss/postcss package exports a PostCSS plugin wrapper
    // that Tailwind's runtime expects when used as a PostCSS plugin.
    '@tailwindcss/postcss': {},
    autoprefixer: {},
  },
};
