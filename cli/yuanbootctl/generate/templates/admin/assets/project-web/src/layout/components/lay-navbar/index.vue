<script setup lang="ts">
import { h, reactive, ref } from "vue";
import { changePwd } from "@/api/system";
import { message } from "@/utils/message";
import { useUserStoreHook } from "@/store/modules/user";
import { ElButton, ElNotification } from "element-plus";
import type { FormInstance, FormRules } from "element-plus";
import { useNav } from "@/layout/hooks/useNav";
import LaySearch from "../lay-search/index.vue";
import LayNotice from "../lay-notice/index.vue";
import LayNavMix from "../lay-sidebar/NavMix.vue";
import { useTranslationLang } from "@/layout/hooks/useTranslationLang";
import LaySidebarFullScreen from "../lay-sidebar/components/SidebarFullScreen.vue";
import LaySidebarBreadCrumb from "../lay-sidebar/components/SidebarBreadCrumb.vue";
import LaySidebarTopCollapse from "../lay-sidebar/components/SidebarTopCollapse.vue";

import GlobalizationIcon from "@/assets/svg/globalization.svg?component";
import LogoutCircleRLine from "~icons/ri/logout-circle-r-line";
import LockPasswordLine from "~icons/ri/lock-password-line";
import Setting from "~icons/ri/settings-3-line";
import Check from "~icons/ep/check";

const {
  layout,
  device,
  logout,
  onPanel,
  pureApp,
  username,
  userAvatar,
  avatarsStyle,
  toggleSideBar,
  getDropdownItemStyle,
  getDropdownItemClass
} = useNav();

const { t, locale, translationCh, translationEn } = useTranslationLang();

const passwordFormRef = ref<FormInstance>();
const passwordDialogVisible = ref(false);
const passwordLoading = ref(false);
const passwordForm = reactive({
  oldPassword: "",
  passwordType: "auto",
  newPassword: "",
  confirmPassword: ""
});
const passwordPattern =
  /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[~!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]).{8,20}$/;
const passwordRules = reactive<FormRules>({
  oldPassword: [{ required: true, message: "请输入旧密码", trigger: "blur" }],
  newPassword: [
    { required: true, message: "请输入新密码", trigger: "blur" },
    {
      validator: (_rule, value, callback) => {
        if (passwordForm.passwordType !== "manual") {
          callback();
          return;
        }
        if (!passwordPattern.test(value || "")) {
          callback(
            new Error("密码需为 8-20 个字符，包含大小写字母、数字和特殊字符")
          );
          return;
        }
        callback();
      },
      trigger: "blur"
    }
  ],
  confirmPassword: [
    { required: true, message: "请再次确认新密码", trigger: "blur" },
    {
      validator: (_rule, value, callback) => {
        if (passwordForm.passwordType !== "manual") {
          callback();
          return;
        }
        if (value !== passwordForm.newPassword) {
          callback(new Error("两次输入的新密码不一致"));
          return;
        }
        callback();
      },
      trigger: "blur"
    }
  ]
});

function openPasswordDialog() {
  passwordDialogVisible.value = true;
}

function resetPasswordForm() {
  passwordForm.oldPassword = "";
  passwordForm.passwordType = "auto";
  passwordForm.newPassword = "";
  passwordForm.confirmPassword = "";
  passwordFormRef.value?.clearValidate();
}

function generatePassword() {
  const upper = "ABCDEFGHJKLMNPQRSTUVWXYZ";
  const lower = "abcdefghijkmnopqrstuvwxyz";
  const number = "23456789";
  const special = "~!@#$%^&*()_+-=[]{};:,.<>?";
  const all = upper + lower + number + special;
  const chars = [upper, lower, number, special].map(
    item => item[Math.floor(Math.random() * item.length)]
  );
  while (chars.length < 12) {
    chars.push(all[Math.floor(Math.random() * all.length)]);
  }
  return chars.sort(() => Math.random() - 0.5).join("");
}

function showChangePasswordResult(currentUsername: string, password: string) {
  ElNotification({
    title: "密码修改成功",
    type: "success",
    position: "top-right",
    duration: 0,
    message: () =>
      h("div", { class: "w-72" }, [
        h(
          "div",
          {
            class:
              "mb-3 rounded-lg bg-[var(--el-fill-color-light)] p-3 text-sm leading-7"
          },
          [
            h("div", { class: "flex justify-between gap-3" }, [
              h(
                "span",
                { class: "text-[var(--el-text-color-secondary)]" },
                "用户名"
              ),
              h(
                "span",
                { class: "font-medium text-[var(--el-text-color-primary)]" },
                currentUsername
              )
            ]),
            h("div", { class: "flex justify-between gap-3" }, [
              h(
                "span",
                { class: "text-[var(--el-text-color-secondary)]" },
                "新密码"
              ),
              h(
                "span",
                {
                  class:
                    "font-mono font-semibold text-[var(--el-color-primary)]"
                },
                password
              )
            ])
          ]
        ),
        h(
          ElButton,
          {
            type: "primary",
            class: "w-full",
            onClick: async () => {
              await navigator.clipboard.writeText(
                `用户名：${currentUsername}\n密码：${password}`
              );
              message("用户名和密码已复制", { type: "success" });
            }
          },
          () => "复制用户名和密码"
        )
      ])
  });
}

async function submitChangePassword() {
  if (!passwordFormRef.value) return;
  const valid = await passwordFormRef.value.validate().catch(() => false);
  if (!valid) return;
  const password =
    passwordForm.passwordType === "auto"
      ? generatePassword()
      : passwordForm.newPassword;
  passwordLoading.value = true;
  try {
    await changePwd({
      oldPassword: passwordForm.oldPassword,
      newPassword: password
    });
    const currentUsername =
      useUserStoreHook().username || username.value || "当前用户";
    message("修改密码成功", { type: "success" });
    passwordDialogVisible.value = false;
    showChangePasswordResult(currentUsername, password);
    resetPasswordForm();
  } finally {
    passwordLoading.value = false;
  }
}
</script>

<template>
  <div class="navbar bg-white shadow-xs shadow-[rgba(0,21,41,0.08)]">
    <LaySidebarTopCollapse
      v-if="device === 'mobile'"
      class="hamburger-container"
      :is-active="pureApp.sidebar.opened"
      @toggleClick="toggleSideBar"
    />

    <LaySidebarBreadCrumb
      v-if="layout !== 'mix' && device !== 'mobile'"
      class="breadcrumb-container"
    />

    <LayNavMix v-if="layout === 'mix'" />

    <div v-if="layout === 'vertical'" class="vertical-header-right">
      <!-- 菜单搜索 -->
      <LaySearch id="header-search" />
      <!-- 国际化 -->
      <el-dropdown id="header-translation" trigger="click">
        <div
          class="globalization-icon navbar-bg-hover hover:[&>svg]:animate-scale-bounce"
        >
          <IconifyIconOffline :icon="GlobalizationIcon" />
        </div>
        <template #dropdown>
          <el-dropdown-menu class="translation">
            <el-dropdown-item
              :style="getDropdownItemStyle(locale, 'zh')"
              :class="['dark:text-white!', getDropdownItemClass(locale, 'zh')]"
              @click="translationCh"
            >
              <IconifyIconOffline
                v-show="locale === 'zh'"
                class="check-zh"
                :icon="Check"
              />
              简体中文
            </el-dropdown-item>
            <el-dropdown-item
              :style="getDropdownItemStyle(locale, 'en')"
              :class="['dark:text-white!', getDropdownItemClass(locale, 'en')]"
              @click="translationEn"
            >
              <span v-show="locale === 'en'" class="check-en">
                <IconifyIconOffline :icon="Check" />
              </span>
              English
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      <!-- 全屏 -->
      <LaySidebarFullScreen id="full-screen" />
      <!-- 消息通知 -->
      <LayNotice id="header-notice" />
      <!-- 退出登录 -->
      <el-dropdown trigger="click">
        <span class="el-dropdown-link navbar-bg-hover select-none">
          <img :src="userAvatar" :style="avatarsStyle" />
          <p v-if="username" class="dark:text-white">{{ username }}</p>
        </span>
        <template #dropdown>
          <el-dropdown-menu class="logout">
            <el-dropdown-item @click="openPasswordDialog">
              <IconifyIconOffline
                :icon="LockPasswordLine"
                style="margin: 5px"
              />
              修改密码
            </el-dropdown-item>
            <el-dropdown-item @click="logout">
              <IconifyIconOffline
                :icon="LogoutCircleRLine"
                style="margin: 5px"
              />
              {{ t("buttons.pureLoginOut") }}
            </el-dropdown-item>
          </el-dropdown-menu>
        </template>
      </el-dropdown>
      <span
        class="set-icon navbar-bg-hover hover:[&>svg]:animate-scale-bounce"
        :title="t('buttons.pureOpenSystemSet')"
        @click="onPanel"
      >
        <IconifyIconOffline :icon="Setting" />
      </span>
    </div>

    <el-dialog
      v-model="passwordDialogVisible"
      title="修改密码"
      width="420px"
      @closed="resetPasswordForm"
    >
      <el-form
        ref="passwordFormRef"
        :model="passwordForm"
        :rules="passwordRules"
        label-width="96px"
      >
        <el-form-item label="旧密码" prop="oldPassword">
          <el-input
            v-model="passwordForm.oldPassword"
            type="password"
            show-password
            clearable
            placeholder="请输入旧密码"
          />
        </el-form-item>
        <el-form-item label="密码类型" prop="passwordType">
          <el-radio-group v-model="passwordForm.passwordType">
            <el-radio value="auto">自动生成密码</el-radio>
            <el-radio value="manual">手动生成密码</el-radio>
          </el-radio-group>
        </el-form-item>
        <template v-if="passwordForm.passwordType === 'manual'">
          <el-form-item label="新密码" prop="newPassword">
            <el-input
              v-model="passwordForm.newPassword"
              type="password"
              show-password
              clearable
              placeholder="请输入新密码"
            />
          </el-form-item>
          <el-form-item label="再次确认" prop="confirmPassword">
            <el-input
              v-model="passwordForm.confirmPassword"
              type="password"
              show-password
              clearable
              placeholder="请再次输入新密码"
            />
          </el-form-item>
          <div class="mb-2 ml-24 text-xs text-gray-500">
            密码需为 8-20 个字符，包含大小写字母、数字和特殊字符
          </div>
        </template>
      </el-form>
      <template #footer>
        <el-button @click="passwordDialogVisible = false">取消</el-button>
        <el-button
          type="primary"
          :loading="passwordLoading"
          @click="submitChangePassword"
        >
          修改密码
        </el-button>
      </template>
    </el-dialog>
  </div>
</template>

<style lang="scss" scoped>
.navbar {
  width: 100%;
  height: 48px;
  overflow: hidden;

  .hamburger-container {
    float: left;
    height: 100%;
    line-height: 48px;
    cursor: pointer;
  }

  .vertical-header-right {
    display: flex;
    align-items: center;
    justify-content: flex-end;
    min-width: 280px;
    height: 48px;
    color: #000000d9;

    .el-dropdown-link {
      display: flex;
      align-items: center;
      justify-content: space-around;
      height: 48px;
      padding: 10px;
      color: #000000d9;
      cursor: pointer;

      p {
        font-size: 14px;
      }

      img {
        width: 22px;
        height: 22px;
        border-radius: 50%;
      }
    }
  }

  .breadcrumb-container {
    float: left;
    margin-left: 16px;
  }
}

.translation {
  :deep(.el-dropdown-menu__item) {
    padding: 5px 40px;
  }

  .check-zh {
    position: absolute;
    left: 20px;
  }

  .check-en {
    position: absolute;
    left: 20px;
  }
}

.logout {
  width: 130px;

  :deep(.el-dropdown-menu__item) {
    display: inline-flex;
    flex-wrap: wrap;
    min-width: 100%;
  }
}
</style>
