/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./assets/templates/**/*.{html,templ}",
    "./assets/templates/**/*.go", // Include generated Go files if they contain Tailwind classes
  ],
  theme: {
    extend: {},
  },
  plugins: [
    require('daisyui'),
  ],

  daisyui: {
    themes: ["emerald"],
  },
}
