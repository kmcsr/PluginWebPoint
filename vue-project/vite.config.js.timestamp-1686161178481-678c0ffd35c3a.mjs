// vite.config.js
import { fileURLToPath, URL } from "node:url";
import { defineConfig } from "file:///Users/ckpn/Mine/projects/golang/src/github.com/kmcsr/PluginWebPoint/vue-project/node_modules/vite/dist/node/index.js";
import vue from "file:///Users/ckpn/Mine/projects/golang/src/github.com/kmcsr/PluginWebPoint/vue-project/node_modules/@vitejs/plugin-vue/dist/index.mjs";
import vueJsx from "file:///Users/ckpn/Mine/projects/golang/src/github.com/kmcsr/PluginWebPoint/vue-project/node_modules/@vitejs/plugin-vue-jsx/dist/index.mjs";
import ssr from "file:///Users/ckpn/Mine/projects/golang/src/github.com/kmcsr/PluginWebPoint/vue-project/node_modules/vite-plugin-ssr/dist/cjs/node/plugin/index.js";
var __vite_injected_original_import_meta_url = "file:///Users/ckpn/Mine/projects/golang/src/github.com/kmcsr/PluginWebPoint/vue-project/vite.config.js";
var vite_config_default = defineConfig(async ({ command, mode }) => {
  console.log(command, mode);
  const isdev = mode === "development";
  const minify = isdev ? "" : "esbuild";
  return {
    plugins: [vue(), vueJsx(), ssr()],
    base: "/",
    resolve: {
      alias: {
        "@": fileURLToPath(new URL("./src", __vite_injected_original_import_meta_url))
      }
    },
    mode,
    build: {
      minify
    },
    esbuild: {
      pure: mode === "production" ? ["console.debug"] : []
    }
  };
});
export {
  vite_config_default as default
};
//# sourceMappingURL=data:application/json;base64,ewogICJ2ZXJzaW9uIjogMywKICAic291cmNlcyI6IFsidml0ZS5jb25maWcuanMiXSwKICAic291cmNlc0NvbnRlbnQiOiBbImNvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9kaXJuYW1lID0gXCIvVXNlcnMvY2twbi9NaW5lL3Byb2plY3RzL2dvbGFuZy9zcmMvZ2l0aHViLmNvbS9rbWNzci9QbHVnaW5XZWJQb2ludC92dWUtcHJvamVjdFwiO2NvbnN0IF9fdml0ZV9pbmplY3RlZF9vcmlnaW5hbF9maWxlbmFtZSA9IFwiL1VzZXJzL2NrcG4vTWluZS9wcm9qZWN0cy9nb2xhbmcvc3JjL2dpdGh1Yi5jb20va21jc3IvUGx1Z2luV2ViUG9pbnQvdnVlLXByb2plY3Qvdml0ZS5jb25maWcuanNcIjtjb25zdCBfX3ZpdGVfaW5qZWN0ZWRfb3JpZ2luYWxfaW1wb3J0X21ldGFfdXJsID0gXCJmaWxlOi8vL1VzZXJzL2NrcG4vTWluZS9wcm9qZWN0cy9nb2xhbmcvc3JjL2dpdGh1Yi5jb20va21jc3IvUGx1Z2luV2ViUG9pbnQvdnVlLXByb2plY3Qvdml0ZS5jb25maWcuanNcIjtpbXBvcnQgeyBmaWxlVVJMVG9QYXRoLCBVUkwgfSBmcm9tICdub2RlOnVybCdcblxuaW1wb3J0IHsgZGVmaW5lQ29uZmlnIH0gZnJvbSAndml0ZSdcbmltcG9ydCB2dWUgZnJvbSAnQHZpdGVqcy9wbHVnaW4tdnVlJ1xuaW1wb3J0IHZ1ZUpzeCBmcm9tICdAdml0ZWpzL3BsdWdpbi12dWUtanN4J1xuaW1wb3J0IHNzciBmcm9tICd2aXRlLXBsdWdpbi1zc3IvcGx1Z2luJ1xuXG4vLyBodHRwczovL3ZpdGVqcy5kZXYvY29uZmlnL1xuZXhwb3J0IGRlZmF1bHQgZGVmaW5lQ29uZmlnKGFzeW5jICh7IGNvbW1hbmQsIG1vZGUgfSkgPT4ge1xuXHRjb25zb2xlLmxvZyhjb21tYW5kLCBtb2RlKTtcblx0Y29uc3QgaXNkZXYgPSBtb2RlID09PSAnZGV2ZWxvcG1lbnQnO1xuXHRjb25zdCBtaW5pZnkgPSBpc2RldiA/JycgOidlc2J1aWxkJztcblxuXHRyZXR1cm4ge1xuXHRcdHBsdWdpbnM6IFt2dWUoKSwgdnVlSnN4KCksIHNzcigpXSxcblx0XHRiYXNlOiAnLycsXG5cdFx0cmVzb2x2ZToge1xuXHRcdFx0YWxpYXM6IHtcblx0XHRcdFx0J0AnOiBmaWxlVVJMVG9QYXRoKG5ldyBVUkwoJy4vc3JjJywgaW1wb3J0Lm1ldGEudXJsKSlcblx0XHRcdH1cblx0XHR9LFxuXHRcdG1vZGU6IG1vZGUsXG5cdFx0YnVpbGQ6IHtcblx0XHRcdG1pbmlmeTogbWluaWZ5LFxuXHRcdH0sXG5cdFx0ZXNidWlsZDoge1xuXHRcdFx0cHVyZTogbW9kZSA9PT0gJ3Byb2R1Y3Rpb24nID8gWydjb25zb2xlLmRlYnVnJ10gOiBbXSxcblx0XHR9XG5cdH1cbn0pXG4iXSwKICAibWFwcGluZ3MiOiAiO0FBQWthLFNBQVMsZUFBZSxXQUFXO0FBRXJjLFNBQVMsb0JBQW9CO0FBQzdCLE9BQU8sU0FBUztBQUNoQixPQUFPLFlBQVk7QUFDbkIsT0FBTyxTQUFTO0FBTHdQLElBQU0sMkNBQTJDO0FBUXpULElBQU8sc0JBQVEsYUFBYSxPQUFPLEVBQUUsU0FBUyxLQUFLLE1BQU07QUFDeEQsVUFBUSxJQUFJLFNBQVMsSUFBSTtBQUN6QixRQUFNLFFBQVEsU0FBUztBQUN2QixRQUFNLFNBQVMsUUFBTyxLQUFJO0FBRTFCLFNBQU87QUFBQSxJQUNOLFNBQVMsQ0FBQyxJQUFJLEdBQUcsT0FBTyxHQUFHLElBQUksQ0FBQztBQUFBLElBQ2hDLE1BQU07QUFBQSxJQUNOLFNBQVM7QUFBQSxNQUNSLE9BQU87QUFBQSxRQUNOLEtBQUssY0FBYyxJQUFJLElBQUksU0FBUyx3Q0FBZSxDQUFDO0FBQUEsTUFDckQ7QUFBQSxJQUNEO0FBQUEsSUFDQTtBQUFBLElBQ0EsT0FBTztBQUFBLE1BQ047QUFBQSxJQUNEO0FBQUEsSUFDQSxTQUFTO0FBQUEsTUFDUixNQUFNLFNBQVMsZUFBZSxDQUFDLGVBQWUsSUFBSSxDQUFDO0FBQUEsSUFDcEQ7QUFBQSxFQUNEO0FBQ0QsQ0FBQzsiLAogICJuYW1lcyI6IFtdCn0K
