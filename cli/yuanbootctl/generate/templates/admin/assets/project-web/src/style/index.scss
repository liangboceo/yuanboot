@use "theme";
@use "transition";
@use "element-plus";
@use "sidebar";
@use "dark";

/* 自定义全局 CssVar */
:root {
  /* 左侧菜单展开、收起动画时长 */
  --pure-transition-duration: 0.3s;

  /* 常用border-color 需要时可取用 */
  --pure-border-color: rgb(5 5 5 / 6%);

  /* switch关闭状态下的color 需要时可取用 */
  --pure-switch-off-color: #a6a6a6;

  /** 主题色 */
  --pure-theme-sub-menu-active-text: initial;
  --pure-theme-menu-bg: none;
  --pure-theme-menu-hover: none;
  --pure-theme-sub-menu-bg: transparent;
  --pure-theme-menu-text: initial;
  --pure-theme-sidebar-logo: none;
  --pure-theme-menu-title-hover: initial;
  --pure-theme-menu-active-before: transparent;
}

/* 灰色模式 */
.html-grey {
  filter: grayscale(100%);
}

/* 色弱模式 */
.html-weakness {
  filter: invert(80%);
}

/* 轻微缩小，再恢复原状的平滑缩放动画 */
@keyframes pure-scale-bounce {
  0% {
    transform: scale(1);
  }

  50% {
    transform: scale(0.9);
  }

  100% {
    transform: scale(1);
  }
}

/* 全局自适应兼容 */
* {
  box-sizing: border-box;
}

.main-content {
  min-width: 0;
}

.main,
.skill-plaza,
.resource-page,
.knowledge-page,
.session-page,
.skill-debug-page {
  min-width: 0;
}

.main > .flex {
  min-width: 0;
}

.main > .flex > * {
  min-width: 0;
}

.search-form {
  border-radius: 12px;
}

.el-table {
  width: 100%;
}

.el-table__body-wrapper,
.el-table__header-wrapper {
  overflow-x: auto;
}

@media screen and (width <= 1200px) {
  .resource-page,
  .knowledge-page {
    flex-direction: column;
  }

  .resource-page .category-panel,
  .knowledge-page .category-panel {
    width: 100%;
    min-width: 0;
  }

  .resource-page .category-scroll,
  .knowledge-page .category-scroll {
    height: auto;
    max-height: 280px;
  }

  .session-page {
    grid-template-columns: 280px minmax(0, 1fr);
  }

  .skill-debug-page.is-workspace {
    grid-template-columns: 220px minmax(0, 1fr);
    height: auto;
    min-height: calc(100vh - 116px);
  }

  .skill-debug-page.is-workspace .skill-chat-panel {
    grid-column: 1 / -1;
    min-height: 520px;
  }
}

@media screen and (width <= 992px) {
  .main-content {
    margin: 16px !important;
  }

  .main,
  .session-page,
  .skill-debug-page.is-workspace {
    height: auto !important;
    min-height: calc(100vh - 100px);
  }

  .session-page,
  .skill-debug-page.is-workspace {
    display: flex !important;
    flex-direction: column;
  }

  .session-page .session-sidebar,
  .session-page .chat-panel,
  .skill-debug-page .skill-file-panel,
  .skill-debug-page .skill-editor-panel,
  .skill-debug-page .skill-chat-panel {
    width: 100%;
  }

  .session-page .session-list {
    height: auto;
    max-height: 320px;
  }

  .session-page .chat-header,
  .skill-debug-page .chat-header,
  .skill-debug-page .editor-toolbar,
  .skill-plaza .hero-card,
  .skill-plaza .toolbar-card,
  .resource-page .hero-card,
  .resource-page .toolbar-card,
  .knowledge-page .hero-card,
  .knowledge-page .toolbar-card {
    flex-wrap: wrap;
    align-items: flex-start;
  }

  .session-page .chat-header > div:last-child,
  .skill-plaza .action-group,
  .resource-page .hero-actions,
  .knowledge-page .hero-actions,
  .resource-page .view-toggle,
  .knowledge-page .view-toggle {
    flex-wrap: wrap;
    width: 100%;
    margin-left: 0;
  }

  .session-page .chat-header .el-select,
  .skill-debug-page .model-select,
  .skill-plaza .search-input,
  .resource-page .toolbar-card .el-input,
  .knowledge-page .toolbar-card .el-input {
    width: 100% !important;
    max-width: none !important;
  }

  .search-form {
    display: flex;
    flex-wrap: wrap;
    gap: 8px 12px;
    padding: 12px !important;
  }

  .search-form .el-form-item {
    margin-right: 0;
    margin-bottom: 8px;
  }

  .main > .flex {
    flex-direction: column;
  }

  .main > .flex > .w-1\/5,
  .main > .flex > .min-w-300px {
    width: 100% !important;
    min-width: 0 !important;
  }

  .main > .flex > .flex-1 {
    margin-top: 12px;
    margin-left: 0 !important;
  }
}

/* 下载进度弹窗：统一右上角非模态展示 */
.el-overlay-dialog:has(.download-progress-dialog) {
  pointer-events: none;
}

.download-progress-dialog {
  margin: 72px 24px 0 auto !important;
  pointer-events: auto;
  border-radius: 12px;
  box-shadow: 0 18px 45px rgb(15 23 42 / 16%);
}

@media screen and (width <= 768px) {
  .download-progress-dialog {
    width: calc(100vw - 32px) !important;
    margin: 64px 16px 0 auto !important;
  }
}

@media screen and (width <= 768px) {
  .main-content {
    margin: 12px !important;
  }

  .resource-page,
  .knowledge-page,
  .skill-plaza {
    padding: 12px !important;
  }

  .resource-page .hero-card,
  .resource-page .toolbar-card,
  .resource-page .document-area,
  .resource-page .category-panel,
  .knowledge-page .hero-card,
  .knowledge-page .toolbar-card,
  .knowledge-page .document-area,
  .knowledge-page .category-panel,
  .session-page .session-sidebar,
  .session-page .chat-panel,
  .skill-debug-page .skill-file-panel,
  .skill-debug-page .skill-editor-panel,
  .skill-debug-page .skill-chat-panel,
  .skill-plaza .hero-card,
  .skill-plaza .toolbar-card {
    padding: 14px !important;
    border-radius: 14px !important;
  }

  .resource-page .doc-grid,
  .knowledge-page .doc-grid,
  .skill-plaza .skill-grid {
    grid-template-columns: 1fr !important;
  }

  .session-page .message-row,
  .skill-debug-page .message-row {
    gap: 8px;
  }

  .session-page .message-avatar,
  .skill-debug-page .message-avatar {
    display: none;
  }

  .session-page .message-bubble,
  .skill-debug-page .message-bubble {
    width: 100%;
    max-width: 100% !important;
  }

  .session-page .composer,
  .skill-debug-page .composer {
    flex-direction: column;
    align-items: stretch;
  }

  .session-page .composer-actions,
  .skill-debug-page .composer > .el-button {
    width: 100%;
  }

  .session-page .composer-actions .el-button,
  .skill-debug-page .composer > .el-button {
    flex: 1;
  }

  .el-pagination {
    flex-wrap: wrap;
    row-gap: 8px;
    justify-content: center;
  }
}

@media screen and (width <= 576px) {
  .main-content {
    margin: 8px !important;
  }

  .el-button + .el-button {
    margin-left: 0;
  }

  .resource-page .toolbar-card,
  .knowledge-page .toolbar-card,
  .skill-plaza .toolbar-card,
  .session-page .composer-actions,
  .search-form {
    flex-direction: column;
    align-items: stretch;
  }

  .resource-page .toolbar-card .el-button,
  .knowledge-page .toolbar-card .el-button,
  .skill-plaza .toolbar-card .el-button,
  .session-page .composer-actions .el-button,
  .search-form .el-button {
    width: 100%;
  }

  .search-form .el-form-item,
  .search-form .el-form-item__content {
    width: 100%;
  }

  .skill-plaza .hero-stats {
    width: 100%;
    overflow-x: auto;
  }

  .welcome-content h3 {
    flex-direction: column;
    font-size: 18px !important;
  }

  .welcome-question {
    margin-left: 0 !important;
  }
}
