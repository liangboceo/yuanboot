<script setup lang="ts">
import { ref, onMounted } from "vue";
import { ElMessage } from "element-plus";
import { message } from "@/utils/message";
import { useRenderIcon } from "@/components/ReIcon/src/hooks";
import { getSystemConfig, saveSystemConfig } from "@/api/system-config";
import Save from "~icons/ep/check";

defineOptions({ name: "SystemSystemConfig" });

const loading = ref(false);
const saving = ref(false);
const systemName = ref("SendEx");
const systemLogo = ref("");
const systemTitle = ref("SendEx Admin");
const systemDescription = ref("Yuanboot 管理后台");
const logoUploadRef = ref<HTMLInputElement>();

async function loadConfig(key: string, target: { value: string }) {
  const result = (await getSystemConfig({ configKey: key })) as any;
  const data = result?.data ?? result;
  if (data?.configValue) target.value = data.configValue;
}

async function loadData() {
  loading.value = true;
  try {
    await Promise.all([
      loadConfig("system_name", systemName),
      loadConfig("system_logo", systemLogo),
      loadConfig("system_title", systemTitle),
      loadConfig("system_description", systemDescription)
    ]);
  } finally {
    loading.value = false;
  }
}

function handleLogoUpload() {
  logoUploadRef.value?.click();
}

function handleLogoFileChange(event: Event) {
  const file = (event.target as HTMLInputElement).files?.[0];
  if (!file) return;
  if (!file.type.startsWith("image/")) return ElMessage.warning("请选择图片文件");
  if (file.size > 2 * 1024 * 1024) return ElMessage.warning("图片大小不能超过2MB");
  const reader = new FileReader();
  reader.onload = () => (systemLogo.value = reader.result as string);
  reader.readAsDataURL(file);
}

async function handleSave() {
  saving.value = true;
  try {
    await Promise.all([
      saveSystemConfig({ configKey: "system_name", configValue: systemName.value, description: "系统名称" }),
      saveSystemConfig({ configKey: "system_logo", configValue: systemLogo.value, description: "系统Logo" }),
      saveSystemConfig({ configKey: "system_title", configValue: systemTitle.value, description: "系统标题" }),
      saveSystemConfig({ configKey: "system_description", configValue: systemDescription.value, description: "系统宣传语" })
    ]);
    message("系统基础设置保存成功", { type: "success" });
  } catch {
    ElMessage.error("保存失败");
  } finally {
    saving.value = false;
  }
}

onMounted(loadData);
</script>

<template>
  <div class="main">
    <div class="config-header"><h2 class="config-title">系统配置</h2></div>
    <div v-loading="loading" class="config-body">
      <div class="tab-content">
        <el-form label-width="120px" class="tab-form">
          <el-form-item label="系统名称">
            <el-input v-model="systemName" placeholder="请输入系统名称" />
          </el-form-item>
          <el-form-item label="系统Logo">
            <div class="logo-upload-area">
              <div v-if="systemLogo" class="logo-preview"><img :src="systemLogo" alt="Logo" /></div>
              <div v-else class="logo-placeholder" @click="handleLogoUpload"><span>+</span><span>自定义上传</span></div>
              <input ref="logoUploadRef" type="file" accept="image/*" hidden @change="handleLogoFileChange" />
              <el-button v-if="systemLogo" plain @click="handleLogoUpload">更换Logo</el-button>
              <el-button v-if="systemLogo" plain type="danger" @click="systemLogo = ''">移除</el-button>
            </div>
          </el-form-item>
          <el-form-item label="系统标题"><el-input v-model="systemTitle" placeholder="请输入系统标题" /></el-form-item>
          <el-form-item label="宣传语"><el-input v-model="systemDescription" type="textarea" :rows="2" placeholder="请输入系统宣传语" /></el-form-item>
          <el-form-item>
            <el-button type="primary" :icon="useRenderIcon(Save)" :loading="saving" @click="handleSave">保存设置</el-button>
          </el-form-item>
        </el-form>
      </div>
    </div>
  </div>
</template>

<style scoped lang="scss">
.main { min-height: 100%; }
.config-header { padding: 16px 24px; border-bottom: 1px solid var(--el-border-color-light); }
.config-title { font-size: 18px; font-weight: 600; }
.config-body { min-height: calc(100vh - 60px); padding: 24px; background: #f5f7fa; }
.tab-content { max-width: 800px; padding: 28px 32px; background: #fff; border-radius: 12px; box-shadow: 0 1px 3px rgb(0 0 0 / 6%); }
.tab-form { max-width: 720px; }
.logo-upload-area { display: flex; gap: 12px; align-items: center; }
.logo-preview { display: flex; align-items: center; justify-content: center; width: 120px; height: 64px; overflow: hidden; border: 1px solid var(--el-border-color); border-radius: 8px; }
.logo-preview img { max-width: 100%; max-height: 100%; object-fit: contain; }
.logo-placeholder { display: flex; flex-direction: column; align-items: center; justify-content: center; width: 120px; height: 120px; cursor: pointer; border: 1.5px dashed #dcdfe6; border-radius: 8px; color: #909399; }
</style>
