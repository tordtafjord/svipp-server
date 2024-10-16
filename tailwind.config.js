/** @type {import('tailwindcss').Config} */
module.exports = {
  content: [
    "./assets/templates/*.{html,tmpl,gohtml}",
    // Add any other template files you're using
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
