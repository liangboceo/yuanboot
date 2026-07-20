import { defineStore } from "pinia";
import { getSystemConfig } from "@/api/system-config";
import { getConfig, store } from "../utils";
import { ref } from "vue";

export const useSystemConfigStore = defineStore("systemConfig", {
  state: () => ({
    // 系统名称
    systemName: ref(getConfig().Title || "SendEx"),
    // 系统 Logo（base64）
    systemLogo: ref(""),
    // 系统标题（登录页大标题）
    systemTitle: ref("SendEx Agent Plan+"),
    // 系统宣传语
    systemDescription: ref("Agent 全生命周期管理"),
    // 是否已加载
    loaded: false
  }),
  getters: {
    /** 侧边栏 / 导航栏使用的系统名称 */
    getName(state) {
      return state.systemName;
    },
    /** 系统 Logo URL，无自定义时返回默认 logo.svg */
    getLogo(state): string {
      return state.systemLogo || new URL("/logo.svg", import.meta.url).href;
    },
    /** 系统标题 */
    getTitle(state) {
      return state.systemTitle;
    },
    /** 宣传语 */
    getDescription(state) {
      return state.systemDescription;
    },
    /** 浏览器标签页标题使用系统名称 */
    browserTitle(state) {
      return state.systemName;
    }
  },
  actions: {
    async loadConfig() {
      if (this.loaded) return;
      try {
        const keys = [
          { key: "system_name", field: "systemName" },
          { key: "system_logo", field: "systemLogo" },
          { key: "system_title", field: "systemTitle" },
          { key: "system_description", field: "systemDescription" }
        ];

        await Promise.all(
          keys.map(async ({ key, field }) => {
            const result = (await getSystemConfig({ configKey: key })) as any;
            const data = result?.data ?? result;
            if (data?.configValue !== undefined && data.configValue !== "") {
              this[field] = data.configValue;
            }
          })
        );
        // 加载完成后设置浏览器标题
        document.title = this.systemName;
      } catch {
        // 加载失败时使用默认值
      } finally {
        this.loaded = true;
      }
    },

    updateConfig(key: string, value: string) {
      switch (key) {
        case "system_name":
          this.systemName = value;
          break;
        case "system_logo":
          this.systemLogo = value;
          break;
        case "system_title":
          this.systemTitle = value;
          break;
        case "system_description":
          this.systemDescription = value;
          break;
      }

      // 更新浏览器标题
      if (key === "system_name") {
        document.title = `${this.systemName}`;
      }
    }
  }
});

export function useSystemConfigStoreHook() {
  return useSystemConfigStore(store);
}
