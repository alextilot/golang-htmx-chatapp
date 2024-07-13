/** @type {import('tailwindcss').Config} */
module.exports = {
  // mode: "jit,
  darkMode: "media",
  content: ["./web/**/*.{html,js,templ,go}"],
  theme: {
    extend: {},
  },
  plugins: [require("daisyui"), require("autoprefixer")],
  // daisyUI config (optional - here are the default values)
  daisyui: {
    // themes: false, // false: only light + dark | true: all themes | array: specific themes like this ["light", "dark", "cupcake"]
    themes: true,
    darkTheme: "dark", // name of one of the included themes for dark mode
    base: true, // applies background color and foreground color for root element by default
    styled: true, // include daisyUI colors and design decisions for all components
    utils: true, // adds responsive and modifier utility classes
    prefix: "", // prefix for daisyUI classnames (components, modifiers and responsive class names. Not colors)
    logs: true, // Shows info about daisyUI version and used config in the console when building your CSS
    themeRoot: ":root", // The element that receives theme color CSS variables
    // themes: [
    //   {
    //     dark: {
    //       primary: "#793ef9",
    //       "primary-focus": "#570df8",
    //       "primary-content": "#ffffff",
    //
    //       secondary: "#f000b8",
    //       "secondary-focus": "#bd0091",
    //       "secondary-content": "#ffffff",
    //
    //       accent: "#37cdbe",
    //       "accent-focus": "#2ba69a",
    //       "accent-content": "#ffffff",
    //
    //       neutral: "#2a2e37",
    //       "neutral-focus": "#16181d",
    //       "neutral-content": "#ffffff",
    //
    //       "base-100": "#3b424e",
    //       "base-200": "#2a2e37",
    //       "base-300": "#16181d",
    //       "base-content": "#ebecf0",
    //
    //       info: "#66c7ff",
    //       success: "#87cf3a",
    //       warning: "#e1d460",
    //       error: "#ff6b6b",
    //
    //       "--rounded-box": "1rem",
    //       "--rounded-btn": ".5rem",
    //       "--rounded-badge": "1.9rem",
    //
    //       "--animation-btn": ".25s",
    //       "--animation-input": ".2s",
    //
    //       "--btn-text-case": "uppercase",
    //       "--navbar-padding": ".5rem",
    //       "--border-btn": "1px",
    //     },
    //   },
    // ],
  },
};
