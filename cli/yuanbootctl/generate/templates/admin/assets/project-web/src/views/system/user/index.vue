<script setup lang="ts">
import { useUser } from "./utils/hook";
import { transformI18n } from "@/plugins/i18n";
import { PureTableBar } from "@/components/RePureTableBar";
import { useRenderIcon } from "@/components/ReIcon/src/hooks";
import { usePublicHooks } from "../hooks";
import { ref, computed, h } from "vue";
import { hasAuth } from "@/router/utils";
import {
  ElMessageBox,
  ElInput,
  ElForm,
  ElFormItem,
  ElButton,
  ElSelect,
  ElOption
} from "element-plus";
import Delete from "~icons/ep/delete";
import EditPen from "~icons/ep/edit-pen";
import Refresh from "~icons/ep/refresh";
import ArrowDown from "~icons/ep/arrow-down";
import AddFill from "~icons/ri/add-circle-line";
import Picture from "~icons/ep/picture";
import Lock from "~icons/ep/lock";
import Key from "~icons/ep/key";

defineOptions({
  name: "SystemUser"
});

const formRef = ref();
const tableRef = ref();
const dialogRef = ref();
const treeRef = ref();

const {
  form,
  loading,
  columns,
  dataList,
  treeData,
  treeLoading,
  selectedNum,
  pagination,
  buttonClass,
  onSearch,
  resetForm,
  onbatchDel,
  openDialog,
  onTreeSelect,
  handleDelete,
  handleUpload,
  handleReset,
  handleRole,
  handleSizeChange,
  onSelectionCancel,
  handleCurrentChange,
  handleSelectionChange
} = useUser(tableRef, treeRef);

const { switchStyle } = usePublicHooks();

function onFullscreen() {
  tableRef.value.setAdaptive();
}

const selectDialogRef = ref();
function handleSelect() {
  selectDialogRef.value.dialogVisible = true;
}
</script>

<template>
  <div class="main">
    <!-- 左侧部门树 -->
    <div class="flex">
      <div class="w-1/5 min-w-300px">
        <el-card shadow="never" class="h-full">
          <template #header>
            <div class="flex-bc">
              <span class="font-bold">部门列表</span>
              <el-button
                type="primary"
                link
                :icon="useRenderIcon('ri:refresh-line')"
                @click="onSearch"
              />
            </div>
          </template>
          <el-tree
            ref="treeRef"
            :data="treeData"
            :props="{ children: 'children', label: 'name' }"
            node-key="id"
            highlight-current
            default-expand-all
            :expand-on-click-node="false"
            @node-click="onTreeSelect"
          >
            <template #default="{ node }">
              <span class="flex items-center">
                <span>{{ node.label }}</span>
              </span>
            </template>
          </el-tree>
        </el-card>
      </div>

      <!-- 右侧用户表格 -->
      <div class="flex-1 ml-4">
        <!-- 搜索表单 -->
        <el-form
          ref="formRef"
          :inline="true"
          :model="form"
          class="search-form bg-bg_color w-full pl-8 pt-3 overflow-auto"
        >
          <el-form-item label="用户名称：" prop="username">
            <el-input
              v-model="form.username"
              placeholder="请输入用户名称"
              clearable
              class="w-45!"
            />
          </el-form-item>
          <el-form-item label="手机号码：" prop="phone">
            <el-input
              v-model="form.phone"
              placeholder="请输入手机号码"
              clearable
              class="w-45!"
            />
          </el-form-item>
          <el-form-item label="状态：" prop="status">
            <el-select
              v-model="form.status"
              placeholder="请选择状态"
              clearable
              class="w-45!"
            >
              <el-option label="正常" value="1" />
              <el-option label="禁用" value="0" />
            </el-select>
          </el-form-item>
          <el-form-item>
            <el-button
              type="primary"
              :icon="useRenderIcon('ri/search-line')"
              :loading="loading"
              @click="onSearch"
            >
              搜索
            </el-button>
            <el-button
              :icon="useRenderIcon(Refresh)"
              @click="resetForm(formRef)"
            >
              重置
            </el-button>
          </el-form-item>
        </el-form>

        <!-- 表格工具栏 -->
        <PureTableBar
          title="用户管理"
          :columns="columns"
          :isExpandAll="false"
          :tableRef="tableRef?.getTableRef()"
          @refresh="onSearch"
          @fullscreen="onFullscreen"
        >
          <template #buttons>
            <el-button
              type="primary"
              :icon="useRenderIcon(AddFill)"
              @click="openDialog('新增')"
            >
              新增用户
            </el-button>
          </template>

          <template v-slot="{ size, dynamicColumns }">
            <pure-table
              ref="tableRef"
              adaptive
              :adaptiveConfig="{ offsetBottom: 45 }"
              align-whole="center"
              row-key="id"
              showOverflowTooltip
              table-layout="auto"
              :loading="loading"
              :size="size"
              :data="dataList"
              :columns="dynamicColumns"
              :pagination="pagination"
              :header-cell-style="{
                background: 'var(--el-fill-color-light)',
                color: 'var(--el-text-color-primary)'
              }"
              @selection-change="handleSelectionChange"
              @page-size-change="handleSizeChange"
              @page-current-change="handleCurrentChange"
            >
              <template #operation="{ row }">
                <el-button
                  class="reset-margin"
                  link
                  type="primary"
                  :size="size"
                  :icon="useRenderIcon(EditPen)"
                  @click="openDialog('修改', row)"
                >
                  修改
                </el-button>
                <el-popconfirm
                  :title="`是否确认删除用户名称为${row.username}的这条数据？`"
                  @confirm="handleDelete(row)"
                >
                  <template #reference>
                    <el-button
                      class="reset-margin"
                      link
                      type="danger"
                      :size="size"
                      :icon="useRenderIcon(Delete)"
                    >
                      删除
                    </el-button>
                  </template>
                </el-popconfirm>
                <el-dropdown trigger="click" class="reset-margin ml-1">
                  <el-button link type="primary" :size="size">
                    更多
                    <el-icon class="ml-0.5"
                      ><component :is="ArrowDown"
                    /></el-icon>
                  </el-button>
                  <template #dropdown>
                    <el-dropdown-menu>
                      <el-dropdown-item @click="handleUpload(row)">
                        <el-icon class="mr-1"><Picture /></el-icon>修改头像
                      </el-dropdown-item>
                      <el-dropdown-item @click="handleReset(row)">
                        <el-icon class="mr-1"><Lock /></el-icon>重置密码
                      </el-dropdown-item>
                      <el-dropdown-item @click="handleRole(row)">
                        <el-icon class="mr-1"><Key /></el-icon>分配角色
                      </el-dropdown-item>
                    </el-dropdown-menu>
                  </template>
                </el-dropdown>
              </template>
            </pure-table>
          </template>
        </PureTableBar>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.main-content {
  margin: 24px 24px 0 !important;
}

.search-form {
  :deep(.el-form-item) {
    margin-bottom: 12px;
  }
}
</style>
