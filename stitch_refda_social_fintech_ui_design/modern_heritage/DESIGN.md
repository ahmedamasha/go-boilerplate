---
name: Modern Heritage
colors:
  surface: '#f9f9f7'
  surface-dim: '#dadad8'
  surface-bright: '#f9f9f7'
  surface-container-lowest: '#ffffff'
  surface-container-low: '#f4f4f1'
  surface-container: '#eeeeec'
  surface-container-high: '#e8e8e6'
  surface-container-highest: '#e2e3e0'
  on-surface: '#1a1c1b'
  on-surface-variant: '#404945'
  inverse-surface: '#2f3130'
  inverse-on-surface: '#f1f1ef'
  outline: '#717975'
  outline-variant: '#c0c8c3'
  surface-tint: '#3a6758'
  primary: '#134235'
  on-primary: '#ffffff'
  primary-container: '#2d5a4c'
  on-primary-container: '#a0cfbe'
  inverse-primary: '#a1d1bf'
  secondary: '#675d4e'
  on-secondary: '#ffffff'
  secondary-container: '#efe0cd'
  on-secondary-container: '#6d6354'
  tertiary: '#735c00'
  on-tertiary: '#ffffff'
  tertiary-container: '#cca72f'
  on-tertiary-container: '#4e3d00'
  error: '#ba1a1a'
  on-error: '#ffffff'
  error-container: '#ffdad6'
  on-error-container: '#93000a'
  primary-fixed: '#bcedda'
  primary-fixed-dim: '#a1d1bf'
  on-primary-fixed: '#002118'
  on-primary-fixed-variant: '#214f41'
  secondary-fixed: '#efe0cd'
  secondary-fixed-dim: '#d2c4b2'
  on-secondary-fixed: '#221a0f'
  on-secondary-fixed-variant: '#4f4538'
  tertiary-fixed: '#ffe088'
  tertiary-fixed-dim: '#e9c349'
  on-tertiary-fixed: '#241a00'
  on-tertiary-fixed-variant: '#574500'
  background: '#f9f9f7'
  on-background: '#1a1c1b'
  surface-variant: '#e2e3e0'
typography:
  headline-xl:
    fontFamily: Manrope
    fontSize: 32px
    fontWeight: '700'
    lineHeight: 40px
    letterSpacing: -0.02em
  headline-lg:
    fontFamily: Manrope
    fontSize: 24px
    fontWeight: '600'
    lineHeight: 32px
    letterSpacing: -0.01em
  body-md:
    fontFamily: Manrope
    fontSize: 16px
    fontWeight: '400'
    lineHeight: 24px
    letterSpacing: 0em
  body-sm:
    fontFamily: Manrope
    fontSize: 14px
    fontWeight: '400'
    lineHeight: 20px
    letterSpacing: 0em
  label-md:
    fontFamily: Manrope
    fontSize: 12px
    fontWeight: '600'
    lineHeight: 16px
    letterSpacing: 0.05em
rounded:
  sm: 0.5rem
  DEFAULT: 1rem
  md: 1.5rem
  lg: 2rem
  xl: 3rem
  full: 9999px
spacing:
  unit: 4px
  container-margin: 24px
  stack-gap: 16px
  section-gap: 32px
  inline-padding: 12px
---

## Brand & Style

This design system is anchored in the concept of "Digital Majlis"—a space where Saudi hospitality meets cutting-edge financial technology. The brand personality is generous, sophisticated, and deeply rooted in cultural trust. It avoids the coldness of traditional fintech by infusing "Cultural Warmth" through its palette and spatial relationships.

The visual style is a refined **Minimalism** blended with **Tactile** elements. It prioritizes clarity for financial transactions while using organic, rounded forms to mirror the softness of social interaction and gifting. Every interaction should feel like a premium concierge service: effortless, quiet, and reliable.

## Colors

The palette is a dialogue between the lush 'Calm Green' of growth and the 'Desert Sand' of the landscape. 

- **Primary (Calm Green):** Used for primary actions, success states, and brand-heavy moments. It signals stability and financial security.
- **Secondary (Desert Sand):** Acts as the primary surface color for cards and containers, providing a softer, more inviting alternative to pure white.
- **Tertiary (Gold Accent):** A subtle metallic used sparingly for "Premium" features or "Special Gift" highlights.
- **Neutral:** A deep charcoal-green is used for text instead of pure black to maintain the sophisticated, organic feel.
- **Background:** A very light off-white with a hint of warmth to reduce eye strain and enhance the feeling of premium paper.

## Typography

The typography system is designed for an **Arabic-first** experience. While the tokens utilize **Manrope** for its modern, geometric-yet-refined Latin support, it must be paired with **IBM Plex Sans Arabic** in implementation. 

The hierarchy prioritizes legibility in financial data and warmth in social messaging. Large headlines use a tighter tracking to feel more "editorial," while body text maintains generous line heights (1.5x) to accommodate the distinctive ascenders and descenders of the Arabic script. Labels utilize a slightly heavier weight to ensure they remain clear against the soft secondary background colors.

## Layout & Spacing

This design system employs a **Fluid Grid** model optimized for mobile-first delivery. It utilizes an 8pt rhythm with a 4pt baseline for micro-adjustments. 

- **Margins:** A generous 24px side margin is mandatory to create a sense of exclusivity and "breathing room."
- **Stacking:** Elements are grouped in vertical stacks with 16px gaps, while major functional sections are separated by 32px to provide clear visual anchors.
- **Negative Space:** Whitespace is treated as a design element, not "empty space." It should be used aggressively to separate the social "Feed" from the "Wallet" or "Registry" functions to prevent cognitive overload.

## Elevation & Depth

To achieve the "Modern Tech meets Cultural Warmth" aesthetic, depth is created through **Ambient Shadows** and **Tonal Layers** rather than heavy borders.

- **The Base Surface:** The main background is the 'Warm Desert Sand' at its lightest tint.
- **Elevated Cards:** Gift cards and registry items use a pure white surface with a very soft, diffused shadow (Blur: 20px, Y: 8px, Color: Primary Green at 4% opacity). This tinting of the shadow makes the elevation feel integrated with the brand color.
- **Interactive Depth:** Buttons utilize a slight "pressed" state (lowering elevation) to provide tactile feedback, mimicking the physical act of handing over a gift.

## Shapes

The shape language is defined by **High Roundedness**. A standard radius of 20px or higher is applied to all primary containers and cards to evoke a friendly, safe, and modern atmosphere. 

- **Primary Buttons:** Fully pill-shaped (rounded-full) to encourage interaction.
- **Feature Cards:** Use the 24px radius to frame high-quality product imagery.
- **Input Fields:** Utilize a 16px radius to maintain a distinct but complementary language to the more rounded action buttons.
- **Imagery:** All product placeholders and user avatars must feature a soft corner radius or circular masks to prevent sharp edges from breaking the "Calm" aesthetic.

## Components

The component library focuses on high-touch social-fintech interactions:

- **The Registry Card:** A premium container featuring a 70/30 split between a high-quality product image and gift progress. It uses a "soft progress bar" in Calm Green against a Desert Sand track.
- **Gifting Buttons:** The primary CTA is a pill-shaped button in 'Calm Green' with white text. Secondary actions (e.g., "Add to Registry") use a 'Desert Sand' fill with 'Calm Green' text.
- **Input Fields:** Minimalist design with a 1px border in a lightened Primary Green, shifting to a 2px stroke on focus. Labels are always floating or external to ensure the input remains the hero.
- **Status Chips:** Small, highly rounded badges (e.g., "Gifted," "Pending") use a translucent version of the state color (Success/Warning) to maintain the soft UI aesthetic.
- **Contribution Slider:** A custom component for partial gifting, featuring a tactile thumb and a subtle haptic-feedback track.
- **Lists:** Flat lists are avoided; instead, use "Inset Groups" with rounded corners that match the 24px margin of the screen.