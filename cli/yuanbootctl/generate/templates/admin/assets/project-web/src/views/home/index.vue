<script setup lang="ts">
import { computed } from "vue";
import { useUserStoreHook } from "@/store/modules/user";
const username = computed(() => useUserStoreHook().username || "管理员");

const modules = [
  { title: "用户管理", description: "维护后台用户及账号状态", icon: "用" },
  { title: "角色管理", description: "配置角色与访问权限", icon: "角" },
  { title: "菜单管理", description: "维护动态菜单与按钮权限", icon: "菜" },
  { title: "系统设置", description: "管理系统基础配置", icon: "设" }
];
</script>

<template>
  <div class="home-page">
    <section class="welcome-card">
      <div>
        <p class="eyebrow">YUANBOOT ADMIN</p>
        <h1>欢迎回来，{{ username }}</h1>
        <p class="description">
          这是基于 Yuanboot 与 PureAdmin 构建的后台管理系统。
        </p>
      </div>
      <div class="brand-mark">Y</div>
    </section>

    <section class="module-grid">
      <article v-for="item in modules" :key="item.title" class="module-card">
        <div class="module-icon">
          {{ item.icon }}
        </div>
        <div>
          <h2>{{ item.title }}</h2>
          <p>{{ item.description }}</p>
        </div>
      </article>
    </section>
  </div>
</template>

<style scoped lang="scss">
.home-page {
  min-height: calc(100vh - 112px);
  padding: 24px;
  background: var(--el-bg-color-page);
}

.welcome-card {
  display: flex;
  align-items: center;
  justify-content: space-between;
  padding: 42px;
  color: #fff;
  background: linear-gradient(135deg, #2563eb, #4f46e5);
  border-radius: 18px;
  box-shadow: 0 18px 45px rgb(37 99 235 / 18%);

  .eyebrow {
    margin-bottom: 10px;
    font-size: 12px;
    font-weight: 700;
    letter-spacing: 0.18em;
    opacity: 0.72;
  }

  h1 {
    margin: 0;
    font-size: 30px;
    font-weight: 700;
  }

  .description {
    margin-top: 12px;
    font-size: 15px;
    opacity: 0.82;
  }
}

.brand-mark {
  display: grid;
  width: 88px;
  height: 88px;
  font-size: 42px;
  font-weight: 800;
  place-items: center;
  background: rgb(255 255 255 / 14%);
  border: 1px solid rgb(255 255 255 / 28%);
  border-radius: 24px;
}

.module-grid {
  display: grid;
  grid-template-columns: repeat(4, minmax(0, 1fr));
  gap: 18px;
  margin-top: 22px;
}

.module-card {
  display: flex;
  gap: 15px;
  align-items: center;
  padding: 22px;
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 14px;

  h2 {
    margin: 0;
    font-size: 16px;
    color: var(--el-text-color-primary);
  }

  p {
    margin-top: 6px;
    font-size: 13px;
    color: var(--el-text-color-secondary);
  }
}

.module-icon {
  display: grid;
  flex: 0 0 46px;
  width: 46px;
  height: 46px;
  font-size: 22px;
  color: var(--el-color-primary);
  place-items: center;
  background: var(--el-color-primary-light-9);
  border-radius: 12px;
}

@media (width <= 900px) {
  .module-grid {
    grid-template-columns: repeat(2, minmax(0, 1fr));
  }
}

@media (width <= 560px) {
  .home-page {
    padding: 14px;
  }

  .welcome-card {
    padding: 28px 22px;
  }

  .brand-mark {
    display: none;
  }

  .module-grid {
    grid-template-columns: 1fr;
  }
}
</style>
