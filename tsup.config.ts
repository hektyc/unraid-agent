import { defineConfig } from "tsup";

export default defineConfig({
  entry: ["src/index.ts"],
  format: ["cjs", "esm"],
  dts: true,
  sourcemap: true,
  clean: true,
  minify: false,
  splitting: false,
  bundle: true,
  target: "node20",
  platform: "node",
  outDir: "dist",
});
