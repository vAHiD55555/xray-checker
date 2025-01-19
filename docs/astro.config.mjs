// @ts-check
import { defineConfig } from "astro/config";
import starlight from "@astrojs/starlight";

// https://astro.build/config
export default defineConfig({
  integrations: [
    starlight({
      title: "Xray Checker",
      social: {
        github: "https://github.com/kutovoys/xray-checker",
        telegram: "https://t.me/kutovoys",
      },
      sidebar: [
        {
          label: "Introduction",
          items: [
            { label: "Overview", slug: "index" },
            { label: "Features", slug: "intro/features" },
          ],
        },
        {
          label: "Usage",
          autogenerate: { directory: "usage" },
        },
        {
          label: "Configuration",
          items: [
            {
              label: "Environment Variables",
              slug: "configuration/envs",
            },
            {
              label: "Subscription Format",
              slug: "configuration/subscription",
            },
            { label: "Check Methods", slug: "configuration/check-methods" },
          ],
        },
        // {
        //   label: "Endpoints",
        //   items: [
        //     { label: "Health Check", link: "/endpoints/health-check" },
        //     { label: "Configs", link: "/endpoints/configs" },
        //   ],
        // },
        {
          label: "Integrations",
          items: [
            { label: "Metrics", slug: "integrations/metrics" },
            { label: "Prometheus Setup", slug: "integrations/prometheus" },
            { label: "Uptime Kuma", slug: "integrations/uptime-kuma" },
          ],
        },
        {
          label: "Contributing",
          items: [
            { label: "Development Guide", link: "/contributing/guide" },
            { label: "Code Style", link: "/contributing/code-style" },
          ],
        },
      ],
    }),
  ],
});
