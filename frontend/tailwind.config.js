/** @type {import('tailwindcss').Config} */
module.exports = {
  prefix: 'tw-', // <-- avoids class conflicts with Quasar
  content: [
    './src/**/*.{vue,js,ts,jsx,tsx}',   // your frontend source files
  ],
  theme: {
    extend: {},
  },
  plugins: [],
}
