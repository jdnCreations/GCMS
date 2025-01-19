import type { Config } from 'tailwindcss';

export default {
  content: [
    './src/pages/**/*.{js,ts,jsx,tsx,mdx}',
    './src/components/**/*.{js,ts,jsx,tsx,mdx}',
    './src/app/**/*.{js,ts,jsx,tsx,mdx}',
  ],
  theme: {
    extend: {
      colors: {
        nook: {
          beige: '#D8CBAF',
          'light-olive': '#C8D9C3',
          olive: '#A8BBA1',
          charcoal: '#4A4A4A',
          'light-charcoal': '#7d7d7d',
          rose: '#E4B7B2',
          'dark-rose': '#C77B7A',
        },
        background: 'var(--background)',
        foreground: 'var(--foreground)',
      },
    },
  },
  plugins: [],
} satisfies Config;
