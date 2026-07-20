<script setup lang="ts">
import { computed, ref } from "vue";
import ReCol from "@/components/ReCol";
import { formRules } from "../utils/rule";
import { FormProps } from "../utils/types";
import { usePublicHooks } from "../../hooks";
import { uploadUserAvatar } from "@/api/system";
import { message } from "@/utils/message";
import userAvatar from "@/assets/user.jpg";

const props = withDefaults(defineProps<FormProps>(), {
  formInline: () => ({
    title: "新增",
    higherDeptOptions: [],
    parentId: 0,
    deptId: 0,
    nickname: "",
    username: "",
    password: "",
    phone: "",
    email: "",
    sex: "",
    status: 1,
    roleIds: [],
    avatar: "",
    remark: ""
  })
});

const sexOptions = [
  {
    value: 0,
    label: "男"
  },
  {
    value: 1,
    label: "女"
  }
];
const ruleFormRef = ref();
const { switchStyle } = usePublicHooks();
const newFormInline = ref(props.formInline);
const avatarUploading = ref(false);

function handleAvatarFileChange(event: Event) {
  const input = event.target as HTMLInputElement;
  const file = input?.files?.[0];
  if (!file) return;

  const allowed = ["image/jpeg", "image/png", "image/gif", "image/webp"];
  if (!allowed.includes(file.type)) {
    message("头像仅支持 jpg/png/gif/webp 格式", { type: "error" });
    return;
  }
  if (file.size > 5 * 1024 * 1024) {
    message("头像文件不超过 5MB", { type: "error" });
    return;
  }

  avatarUploading.value = true;
  uploadUserAvatar(file)
    .then((res: any) => {
      const data = res?.data ?? res;
      newFormInline.value.avatar = data?.avatar ?? "";
      message("头像上传成功", { type: "success" });
    })
    .catch(() => {
      message("头像上传失败", { type: "error" });
    })
    .finally(() => {
      avatarUploading.value = false;
      input.value = ""; // 允许重复上传同一文件
    });
}

function getRef() {
  return ruleFormRef.value;
}

function getFormData() {
  return newFormInline.value;
}

defineExpose({ getRef, getFormData });
</script>

<template>
  <el-form
    ref="ruleFormRef"
    :model="newFormInline"
    :rules="formRules"
    label-width="82px"
  >
    <el-row :gutter="30">
      <re-col :value="24">
        <el-form-item label="用户头像">
          <div class="flex items-center gap-3">
            <img
              :src="newFormInline.avatar || userAvatar"
              class="avatar-preview"
              alt="用户头像"
            />
            <label class="avatar-upload-box">
              <input
                type="file"
                accept="image/jpeg,image/png,image/gif,image/webp"
                class="hidden"
                :disabled="avatarUploading"
                @change="handleAvatarFileChange"
              />
              <template v-if="avatarUploading">
                <span class="avatar-upload-icon is-loading" />
                <span class="avatar-upload-text">上传中</span>
              </template>
              <template v-else>
                <span class="avatar-upload-icon">+</span>
                <span class="avatar-upload-text">自定义上传</span>
              </template>
            </label>
          </div>
        </el-form-item>
      </re-col>
      <re-col :value="12" :xs="24" :sm="24">
        <el-form-item label="用户昵称" prop="nickname">
          <el-input
            v-model="newFormInline.nickname"
            clearable
            placeholder="请输入用户昵称"
          />
        </el-form-item>
      </re-col>
      <re-col :value="12" :xs="24" :sm="24">
        <el-form-item label="用户名称" prop="username">
          <el-input
            v-model="newFormInline.username"
            clearable
            placeholder="请输入用户名称"
          />
        </el-form-item>
      </re-col>

      <re-col :value="12" :xs="24" :sm="24">
        <el-form-item label="手机号" prop="phone">
          <el-input
            v-model="newFormInline.phone"
            clearable
            placeholder="请输入手机号"
          />
        </el-form-item>
      </re-col>

      <re-col :value="12" :xs="24" :sm="24">
        <el-form-item label="邮箱" prop="email">
          <el-input
            v-model="newFormInline.email"
            clearable
            placeholder="请输入邮箱"
          />
        </el-form-item>
      </re-col>
      <re-col :value="12" :xs="24" :sm="24">
        <el-form-item label="用户性别">
          <el-select
            v-model="newFormInline.sex"
            placeholder="请选择用户性别"
            class="w-full"
            clearable
          >
            <el-option
              v-for="(item, index) in sexOptions"
              :key="index"
              :label="item.label"
              :value="item.value"
            />
          </el-select>
        </el-form-item>
      </re-col>

      <re-col :value="12" :xs="24" :sm="24">
        <el-form-item label="归属部门">
          <el-cascader
            v-model="newFormInline.deptId"
            class="w-full"
            :options="newFormInline.higherDeptOptions"
            :props="{
              value: 'id',
              label: 'name',
              emitPath: false,
              checkStrictly: true
            }"
            clearable
            filterable
            placeholder="请选择归属部门"
          >
            <template #default="{ node, data }">
              <span>{{ data.name }}</span>
              <span v-if="!node.isLeaf"> ({{ data.children.length }}) </span>
            </template>
          </el-cascader>
        </el-form-item>
      </re-col>
      <re-col
        v-if="newFormInline.title === '新增'"
        :value="12"
        :xs="24"
        :sm="24"
      >
        <el-form-item label="用户状态">
          <el-switch
            v-model="newFormInline.status"
            inline-prompt
            :active-value="1"
            :inactive-value="0"
            active-text="启用"
            inactive-text="停用"
            :style="switchStyle"
          />
        </el-form-item>
      </re-col>

      <re-col>
        <el-form-item label="备注">
          <el-input
            v-model="newFormInline.remark"
            placeholder="请输入备注信息"
            type="textarea"
          />
        </el-form-item>
      </re-col>
    </el-row>
  </el-form>

</template>

<style lang="scss" scoped>
.avatar-preview {
  width: 72px;
  height: 72px;
  object-fit: cover;
  border-radius: 50%;
}

.avatar-upload-box {
  display: flex;
  flex-direction: column;
  align-items: center;
  justify-content: center;
  width: 96px;
  height: 96px;
  cursor: pointer;
  background: var(--el-fill-color-lighter);
  border: 2px dashed var(--el-border-color);
  border-radius: 10px;
  transition: all 0.2s ease;

  &:hover {
    background: var(--el-fill-color-light);
    border-color: var(--el-color-primary);
  }
}

.avatar-upload-icon {
  margin-bottom: 4px;
  font-size: 28px;
  line-height: 1;
  color: var(--el-text-color-secondary);
  transition: color 0.2s ease;

  .avatar-upload-box:hover & {
    color: var(--el-color-primary);
  }

  &.is-loading {
    width: 24px;
    height: 24px;
    margin-bottom: 6px;
    font-size: 0;
    border: 2px solid var(--el-border-color);
    border-top-color: var(--el-color-primary);
    border-radius: 50%;
    animation: spin 0.8s linear infinite;
  }
}

@keyframes spin {
  to {
    transform: rotate(360deg);
  }
}

.avatar-upload-text {
  font-size: 12px;
  color: var(--el-text-color-regular);
}

</style>
