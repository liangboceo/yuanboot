<script setup lang="ts">
import { computed, ref } from "vue";
import type {
  ToolTraceItem,
  StatusEvent,
  FileOutputEvent,
  GraphNodeEvent
} from "@/api/session";

import Receiving from "~icons/ri/edit-circle-line";
import Processing from "~icons/ri/brain-line";
import ToolExec from "~icons/ri/tools-line";
import FileDone from "~icons/ri/file-list-3-line";
import Complete from "~icons/ri/check-double-line";
import ErrorIcon from "~icons/ri/close-circle-line";
import ChevronRight from "~icons/ri/arrow-right-s-line";

const props = defineProps<{
  /** 当前处理状态阶段 */
  statusPhase: string | null;
  /** 工具调用追踪列表 */
  toolTraces: ToolTraceItem[];
  /** 文件输出列表 */
  fileOutputs: FileOutputEvent[];
  /** 是否有错误 */
  hasError: boolean;
  /** 流程节点执行链路 */
  graphNodes?: GraphNodeEvent[];
}>();

const emit = defineEmits<{
  download: [docId: string, filename: string];
}>();

/** 节点定义 */
interface FlowNode {
  id: string;
  label: string;
  icon: any;
  status: "pending" | "active" | "completed" | "error";
  phase: string;
  detail?: string;
}

/** 阶段排序（决定节点显示顺序） */
const PHASE_ORDER = ["received", "processing", "tool_executing", "completed"];

const graphFlowNodes = computed<FlowNode[]>(() => {
  const graphNodes = props.graphNodes || [];
  if (!graphNodes.length) return [];
  return graphNodes.map(node => ({
    id: node.nodeId,
    label: node.nodeName || node.nodeId,
    icon:
      node.status === "error"
        ? ErrorIcon
        : node.nodeId === "execute_tools"
          ? ToolExec
          : node.nodeId === "finalize"
            ? Complete
            : Processing,
    status: node.status || "pending",
    phase: node.phase,
    detail: node.detail
  }));
});

/** 计算每个阶段的状态 */
const flowNodes = computed<FlowNode[]>(() => {
  if (graphFlowNodes.value.length) return graphFlowNodes.value;

  const currentPhase = props.statusPhase || "";
  const currentIndex = PHASE_ORDER.indexOf(currentPhase);
  const hasToolCalls = props.toolTraces.length > 0;
  const hasFiles = props.fileOutputs.length > 0;

  const nodes: FlowNode[] = [];

  // 节点1: 接收请求
  nodes.push({
    id: "received",
    label: "接收请求",
    icon: Receiving,
    status: getNodeStatus("received", currentPhase, currentIndex),
    phase: "received",
    detail: "请求已接收，准备处理"
  });

  // 节点2: 分析处理
  nodes.push({
    id: "processing",
    label: "分析处理",
    icon: Processing,
    status: getNodeStatus("processing", currentPhase, currentIndex),
    phase: "processing",
    detail:
      hasToolCalls && currentIndex >= PHASE_ORDER.indexOf("processing")
        ? "模型正在分析并调用工具..."
        : "模型正在生成回复..."
  });

  // 节点3: 工具调用 (有工具调用时显示，或状态阶段正处在 tool_executing)
  const showToolNode = hasToolCalls || currentPhase === "tool_executing";
  if (showToolNode) {
    const toolDetail =
      props.toolTraces.length > 0
        ? props.toolTraces
            .map(t => `${t.hasError ? "✗" : "✓"} ${t.toolName}`)
            .join(", ")
        : "正在准备工具调用...";
    nodes.push({
      id: "tool_executing",
      label: `工具调用${hasToolCalls ? ` (${props.toolTraces.length})` : ""}`,
      icon: ToolExec,
      status: getNodeStatus("tool_executing", currentPhase, currentIndex),
      phase: "tool_executing",
      detail: toolDetail
    });
  }

  // 节点4: 文件输出 (仅当有文件输出时显示)
  if (hasFiles) {
    nodes.push({
      id: "file_output",
      label: `文件输出 (${props.fileOutputs.length})`,
      icon: FileDone,
      status:
        currentIndex >= PHASE_ORDER.indexOf("completed")
          ? "completed"
          : "pending",
      phase: "file_output",
      detail: props.fileOutputs.map(f => f.filename).join(", ")
    });
  }

  // 节点5: 完成
  nodes.push({
    id: "completed",
    label: props.hasError ? "处理出错" : "处理完成",
    icon: props.hasError ? ErrorIcon : Complete,
    status: props.hasError
      ? "error"
      : getNodeStatus("completed", currentPhase, currentIndex),
    phase: "completed",
    detail: props.hasError ? "请求处理失败" : "所有步骤已完成"
  });

  return nodes;
});

function getNodeStatus(
  nodePhase: string,
  currentPhase: string,
  currentIndex: number
): "pending" | "active" | "completed" | "error" {
  if (props.hasError && nodePhase !== "completed") return "pending";
  const nodeIndex = PHASE_ORDER.indexOf(nodePhase);
  if (nodeIndex < currentIndex) return "completed";
  if (nodeIndex === currentIndex) {
    // 完成节点应直接标记为 completed，而非 active（"进行中"）
    return currentPhase === "completed" ? "completed" : "active";
  }
  return "pending";
}

/** 工具链路展开状态 */
const expandedTools = ref(false);

/** 文件输出区域的展开状态 */
const expandedFiles = ref(true);

function toggleTools() {
  expandedTools.value = !expandedTools.value;
}

function handleDownload(docId: string, filename: string) {
  emit("download", docId, filename);
}

function formatToolValue(value: any): string {
  if (value === null || value === undefined || value === "") return "-";
  if (typeof value === "string") return value;
  if (typeof value === "number" || typeof value === "boolean")
    return String(value);
  if (Array.isArray(value)) {
    if (!value.length) return "空列表";
    return value
      .map(item => formatToolValue(item))
      .filter(Boolean)
      .join("、");
  }
  if (typeof value === "object") {
    return Object.entries(value)
      .map(([key, val]) => `${key}: ${formatToolValue(val)}`)
      .join("；");
  }
  return String(value);
}

function formatToolData(data: any) {
  if (data === null || data === undefined || data === "") {
    return [{ label: "内容", value: "-" }];
  }
  if (typeof data !== "object" || Array.isArray(data)) {
    return [{ label: "内容", value: formatToolValue(data) }];
  }
  const entries = Object.entries(data).filter(
    ([, value]) => value !== undefined
  );
  if (!entries.length) return [{ label: "内容", value: "-" }];
  return entries.map(([label, value]) => ({
    label,
    value: formatToolValue(value)
  }));
}

function formatToolResult(trace: ToolTraceItem) {
  return formatToolData(
    trace.hasError ? trace.result?.error || "未知错误" : trace.result
  );
}
</script>

<template>
  <div class="flow-node-diagram">
    <div
      v-for="(node, index) in flowNodes"
      :key="node.id"
      :class="[
        'flow-node',
        `flow-node--${node.status}`,
        { 'is-last': index === flowNodes.length - 1 }
      ]"
    >
      <!-- 竖线连接 -->
      <div v-if="index < flowNodes.length - 1" class="flow-line" />

      <!-- 节点图标 -->
      <div :class="['flow-dot', `flow-dot--${node.status}`]">
        <IconifyIconOffline :icon="node.icon" class="flow-dot-icon" />
      </div>

      <!-- 节点内容 -->
      <div class="flow-content">
        <div class="flow-header">
          <span class="flow-label">{{ node.label }}</span>
          <span
            v-if="node.status === 'active'"
            class="flow-badge flow-badge--active"
          >
            <span class="flow-badge-dot" />
            进行中
          </span>
          <span
            v-else-if="node.status === 'completed'"
            class="flow-badge flow-badge--done"
          >
            ✓ 完成
          </span>
          <span
            v-else-if="node.status === 'error'"
            class="flow-badge flow-badge--error"
          >
            ✗ 失败
          </span>
          <span v-else class="flow-badge flow-badge--pending">等待中</span>
        </div>
        <div v-if="node.detail" class="flow-detail">{{ node.detail }}</div>

        <!-- 工具调用展开详情 -->
        <div
          v-if="
            (node.id === 'tool_executing' || node.id === 'execute_tools') &&
            toolTraces.length
          "
          class="flow-tools"
        >
          <div class="flow-tools-toggle" @click="toggleTools">
            <IconifyIconOffline
              :icon="ChevronRight"
              :class="['flow-chevron', { 'is-expanded': expandedTools }]"
            />
            工具调用详情
          </div>
          <div v-if="expandedTools" class="flow-tools-body">
            <div
              v-for="(trace, ti) in toolTraces"
              :key="ti"
              :class="['flow-tool-item', { 'is-error': trace.hasError }]"
            >
              <div class="flow-tool-header">
                <span class="flow-tool-name">{{ trace.toolName }}</span>
                <span
                  :class="[
                    'flow-tool-status',
                    trace.hasError ? 'is-error' : 'is-success'
                  ]"
                >
                  {{ trace.hasError ? "失败" : "成功" }}
                </span>
              </div>
              <div v-if="trace.arguments" class="flow-tool-detail">
                <span class="flow-tool-label">参数:</span>
                <div class="flow-tool-kv-list">
                  <div
                    v-for="item in formatToolData(trace.arguments)"
                    :key="item.label"
                    class="flow-tool-kv"
                  >
                    <span class="flow-tool-kv-key">{{ item.label }}</span>
                    <span class="flow-tool-kv-value">{{ item.value }}</span>
                  </div>
                </div>
              </div>
              <div v-if="trace.result" class="flow-tool-detail">
                <span class="flow-tool-label">结果:</span>
                <div class="flow-tool-kv-list">
                  <div
                    v-for="item in formatToolResult(trace)"
                    :key="item.label"
                    class="flow-tool-kv"
                  >
                    <span class="flow-tool-kv-key">{{ item.label }}</span>
                    <span class="flow-tool-kv-value">{{ item.value }}</span>
                  </div>
                </div>
              </div>
            </div>
          </div>
        </div>

        <!-- 文件输出详情 -->
        <div
          v-if="node.id === 'file_output' && fileOutputs.length"
          class="flow-files"
        >
          <div
            v-for="(file, fi) in fileOutputs"
            :key="fi"
            class="flow-file-item"
          >
            <div class="flow-file-info">
              <IconifyIconOffline :icon="FileDone" class="flow-file-icon" />
              <div>
                <div class="flow-file-name">{{ file.filename }}</div>
                <div class="flow-file-meta">
                  ID: {{ file.docId }} · {{ (file.size / 1024).toFixed(1) }} KB
                </div>
              </div>
            </div>
            <el-button
              v-if="file.docId"
              type="primary"
              size="small"
              plain
              @click="handleDownload(file.docId, file.filename)"
            >
              下载文件
            </el-button>
          </div>
        </div>
      </div>
    </div>
  </div>
</template>

<style lang="scss" scoped>
.flow-node-diagram {
  position: relative;
  padding: 4px 0;
}

.flow-node {
  position: relative;
  display: flex;
  gap: 12px;
  padding-bottom: 16px;

  &.is-last {
    padding-bottom: 0;
  }
}

/* 竖线连接 */
.flow-line {
  position: absolute;
  top: 28px;
  left: 13px;
  width: 2px;
  height: calc(100% - 12px);
  background: var(--el-border-color);

  .flow-node--active ~ .flow-node & {
    background: var(--el-border-color);
  }

  .flow-node--completed & {
    background: var(--el-color-success);
  }

  .flow-node--active & {
    background: linear-gradient(
      to bottom,
      var(--el-color-success) 0%,
      var(--el-color-primary) 100%
    );
  }

  .flow-node--error & {
    background: var(--el-color-danger);
  }
}

/* 节点圆点 */
.flow-dot {
  position: relative;
  z-index: 1;
  display: flex;
  flex-shrink: 0;
  align-items: center;
  justify-content: center;
  width: 28px;
  height: 28px;
  border-radius: 50%;
  transition: all 300ms ease;

  &--pending {
    color: var(--el-text-color-placeholder);
    background: var(--el-fill-color-light);
    border: 2px solid var(--el-border-color);
  }

  &--active {
    color: #fff;
    background: var(--el-color-primary);
    border: 2px solid var(--el-color-primary);
    box-shadow: 0 0 0 4px var(--el-color-primary-light-8);
    animation: flow-pulse 2s ease-in-out infinite;
  }

  &--completed {
    color: #fff;
    background: var(--el-color-success);
    border: 2px solid var(--el-color-success);
  }

  &--error {
    color: #fff;
    background: var(--el-color-danger);
    border: 2px solid var(--el-color-danger);
  }
}

.flow-dot-icon {
  font-size: 14px;
}

/* 节点内容 */
.flow-content {
  flex: 1;
  min-width: 0;
}

.flow-header {
  display: flex;
  flex-wrap: wrap;
  gap: 8px;
  align-items: center;
}

.flow-label {
  font-size: 13px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.flow-badge {
  display: inline-flex;
  gap: 4px;
  align-items: center;
  padding: 1px 8px;
  font-size: 11px;
  font-weight: 500;
  border-radius: 10px;

  &--pending {
    color: var(--el-text-color-placeholder);
    background: var(--el-fill-color-light);
  }

  &--active {
    color: var(--el-color-primary);
    background: var(--el-color-primary-light-9);
  }

  &--done {
    color: var(--el-color-success);
    background: var(--el-color-success-light-9);
  }

  &--error {
    color: var(--el-color-danger);
    background: var(--el-color-danger-light-9);
  }
}

.flow-badge-dot {
  display: inline-block;
  width: 6px;
  height: 6px;
  background: currentcolor;
  border-radius: 50%;
  animation: flow-pulse 1.5s ease-in-out infinite;
}

.flow-detail {
  margin-top: 4px;
  font-size: 12px;
  line-height: 1.6;
  color: var(--el-text-color-secondary);
}

/* 工具调用区域 */
.flow-tools {
  padding: 10px 12px;
  margin-top: 8px;
  background: var(--el-fill-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 8px;
}

.flow-tools-toggle {
  display: flex;
  gap: 4px;
  align-items: center;
  font-size: 12px;
  font-weight: 500;
  color: var(--el-color-primary);
  cursor: pointer;
  user-select: none;
}

.flow-chevron {
  font-size: 14px;
  transition: transform 200ms ease;

  &.is-expanded {
    transform: rotate(90deg);
  }
}

.flow-tools-body {
  margin-top: 8px;
}

.flow-tool-item {
  padding: 8px;
  margin-top: 6px;
  background: var(--el-bg-color);
  border: 1px solid var(--el-border-color-lighter);
  border-radius: 6px;

  &.is-error {
    background: var(--el-color-danger-light-9);
    border-color: var(--el-color-danger-light-7);
  }
}

.flow-tool-header {
  display: flex;
  gap: 8px;
  align-items: center;
  justify-content: space-between;
}

.flow-tool-name {
  font-family: "SF Mono", "Cascadia Code", monospace;
  font-size: 12px;
  font-weight: 600;
  color: var(--el-text-color-primary);
}

.flow-tool-status {
  font-size: 11px;
  font-weight: 500;

  &.is-success {
    color: var(--el-color-success);
  }

  &.is-error {
    color: var(--el-color-danger);
  }
}

.flow-tool-detail {
  margin-top: 6px;

  .flow-tool-label {
    display: block;
    margin-bottom: 3px;
    font-size: 11px;
    font-weight: 500;
    color: var(--el-text-color-secondary);
  }

  .flow-tool-kv-list {
    display: flex;
    flex-direction: column;
    gap: 4px;
  }

  .flow-tool-kv {
    display: grid;
    grid-template-columns: minmax(72px, 120px) 1fr;
    gap: 8px;
    padding: 6px 8px;
    font-size: 12px;
    line-height: 1.5;
    background: var(--el-fill-color);
    border-radius: 5px;
  }

  .flow-tool-kv-key {
    color: var(--el-text-color-secondary);
  }

  .flow-tool-kv-value {
    color: var(--el-text-color-primary);
    overflow-wrap: anywhere;
    white-space: pre-wrap;
  }
}

/* 文件输出区域 */
.flow-files {
  margin-top: 8px;
}

.flow-file-item {
  display: flex;
  gap: 8px;
  align-items: center;
  justify-content: space-between;
  padding: 8px 10px;
  margin-top: 6px;
  background: var(--el-color-success-light-9);
  border: 1px solid var(--el-color-success-light-7);
  border-radius: 8px;
}

.flow-file-info {
  display: flex;
  gap: 8px;
  align-items: center;
  min-width: 0;
}

.flow-file-icon {
  flex-shrink: 0;
  font-size: 18px;
  color: var(--el-color-success);
}

.flow-file-name {
  font-size: 13px;
  font-weight: 500;
  color: var(--el-text-color-primary);
  word-break: break-all;
}

.flow-file-meta {
  margin-top: 2px;
  font-family: "SF Mono", "Cascadia Code", monospace;
  font-size: 11px;
  color: var(--el-text-color-secondary);
}

@keyframes flow-pulse {
  0%,
  100% {
    box-shadow: 0 0 0 4px var(--el-color-primary-light-8);
  }

  50% {
    box-shadow: 0 0 0 8px var(--el-color-primary-light-9);
  }
}
</style>
