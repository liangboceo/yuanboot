import { message } from "@/utils/message";
import { getMenuList, getMenuTree, saveMenu, delMenu } from "@/api/system";
import { addDialog } from "@/components/ReDialog";
import {
  reactive,
  ref,
  onMounted,
  h,
  defineComponent,
  watch,
  nextTick
} from "vue";
import { IconSelect } from "@/components/ReIcon";
import { useRenderIcon } from "@/components/ReIcon/src/hooks";
import {
  ElMessageBox,
  ElForm,
  ElFormItem,
  ElInput,
  ElInputNumber
} from "element-plus";

const ROOT_MENU_ID = 0;

function appendRootMenu(tree = []) {
  return [
    {
      id: ROOT_MENU_ID,
      label: "主目录",
      name: "主目录",
      children: tree
    }
  ];
}

function removeSelfMenu(tree = [], id) {
  return tree
    .filter(item => item.id !== id)
    .map(item => ({
      ...item,
      children: item.children?.length ? removeSelfMenu(item.children, id) : []
    }));
}

function findMenuName(tree = [], id) {
  if (!id) return "主目录";
  for (const item of tree) {
    if (item.id === id) return item.label ?? item.name;
    const name = findMenuName(item.children || [], id);
    if (name) return name;
  }
  return "主目录";
}

export function useMenu() {
  const form = reactive({
    menuName: ""
  });

  const formRef = ref();
  const dataList = ref([]);
  const loading = ref(true);
  const treeData = ref([]);

  const getMenuType = (type, text = false) => {
    switch (type) {
      case 0:
        return text ? "目录" : "primary";
      case 1:
        return text ? "菜单" : "warning";
      case 2:
        return text ? "按钮" : "info";
      case 3:
        return text ? "内页" : "success";
      default:
        return text ? "未知" : "info";
    }
  };

  const columns: TableColumnList = [
    {
      label: "菜单名称",
      prop: "name",
      align: "left",
      cellRenderer: ({ row }) => (
        <>
          <span class="inline-block mr-1">
            {row.icon
              ? h(useRenderIcon(row.icon), {
                  style: { paddingTop: "1px" }
                })
              : null}
          </span>
          <span>{row.name}</span>
        </>
      )
    },
    {
      label: "菜单类型",
      prop: "menuType",
      width: 100,
      cellRenderer: ({ row, props }) => (
        <el-tag size={props.size} type={getMenuType(row.menuType) as any}>
          {getMenuType(row.menuType, true)}
        </el-tag>
      )
    },
    {
      label: "路由路径",
      prop: "path"
    },
    {
      label: "组件路径",
      prop: "component"
    },
    {
      label: "权限标识",
      prop: "permission"
    },
    {
      label: "排序",
      prop: "orderNum",
      width: 80
    },
    {
      label: "状态",
      prop: "status",
      width: 80,
      cellRenderer: ({ row, props }) => (
        <el-tag
          size={props.size}
          type={row.status === 1 ? "success" : "danger"}
        >
          {row.status === 1 ? "正常" : "禁用"}
        </el-tag>
      )
    },
    {
      label: "操作",
      fixed: "right",
      width: 210,
      slot: "operation"
    }
  ];

  function handleSelectionChange(val) {
    console.log("handleSelectionChange", val);
  }

  function resetForm(formEl) {
    if (!formEl) return;
    formEl.resetFields();
    onSearch();
  }

  async function onSearch() {
    loading.value = true;
    try {
      const result = (await getMenuList({ menuName: form.menuName })) as any;
      dataList.value = Array.isArray(result) ? result : result?.data || [];
    } finally {
      setTimeout(() => {
        loading.value = false;
      }, 300);
    }
  }

  async function loadTreeData() {
    const result = (await getMenuTree()) as any;
    treeData.value = Array.isArray(result) ? result : result?.data || [];
  }

  function openDialog(title = "新增", row?: any) {
    const menuOptions = appendRootMenu(removeSelfMenu(treeData.value, row?.id));
    addDialog({
      title: `${title}菜单`,
      props: {
        formInline: {
          title,
          id: row?.id ?? 0,
          parentId: row?.parentId ?? 0,
          parentName: findMenuName(treeData.value, row?.parentId ?? 0),
          menuOptions,
          name: row?.name ?? "",
          path: row?.path ?? "",
          component: row?.component ?? "",
          menuType: row?.menuType ?? 1,
          orderNum: row?.orderNum ?? 1,
          icon: row?.icon ?? "",
          visible: row?.visible ?? 0,
          status: row?.status ?? 1,
          permission: row?.permission ?? "",
          perms: row?.perms ?? "",
          isFrame: row?.isFrame ?? 0,
          isCache: row?.isCache ?? 0,
          query: row?.query ?? ""
        }
      },
      width: "650px",
      draggable: true,
      closeOnClickModal: false,
      contentRenderer: () => h(MenuForm, { ref: formRef }),
      beforeSure: async (done, { options }) => {
        const FormRef = formRef.value?.getRef?.();
        if (FormRef) {
          const valid = await FormRef.validate().catch(() => false);
          if (!valid) return;
        }

        const curData =
          formRef.value?.getFormData?.() ?? options.props.formInline;
        try {
          await saveMenu({
            id: curData.id || 0,
            name: curData.name,
            parentId: curData.parentId || 0,
            path: curData.path,
            component: curData.component,
            menuType: curData.menuType,
            orderNum: curData.orderNum,
            icon: curData.icon,
            visible: curData.visible,
            status: curData.status,
            permission: curData.permission,
            perms: curData.perms,
            isFrame: curData.isFrame,
            isCache: curData.isCache,
            query: curData.query
          });
          message(`菜单${title}成功`, { type: "success" });
          done();
          onSearch();
          loadTreeData();
        } catch (error) {
          console.error(error);
        }
      }
    });
  }

  function handleDelete(row) {
    ElMessageBox.confirm(
      `是否确认删除菜单名称为<strong>${row.name}</strong>的数据？${
        row.children?.length > 0
          ? "<br/><span style='color: #f56c6c'>注意：该菜单下存在子菜单，子菜单也会一并删除</span>"
          : ""
      }`,
      "系统提示",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
        dangerouslyUseHTMLString: true
      }
    )
      .then(async () => {
        try {
          await delMenu({ id: row.id });
          message("删除成功", { type: "success" });
          onSearch();
          loadTreeData();
        } catch (error) {
          console.error(error);
        }
      })
      .catch(() => {});
  }

  onMounted(() => {
    onSearch();
    loadTreeData();
  });

  return {
    form,
    loading,
    columns,
    dataList,
    treeData,
    onSearch,
    resetForm,
    openDialog,
    handleDelete,
    handleSelectionChange
  };
}

const MenuForm = defineComponent({
  name: "MenuForm",
  props: {
    formInline: Object
  },
  setup(props, { expose }) {
    const form = ref({
      name: "",
      parentId: 0,
      parentName: "主目录",
      menuOptions: [],
      path: "",
      component: "",
      menuType: 1,
      orderNum: 1,
      icon: "",
      visible: 0,
      status: 1,
      permission: "",
      perms: "",
      isFrame: 0,
      isCache: 0,
      query: ""
    });

    const formRef = ref();

    watch(
      () => props.formInline,
      val => {
        if (val) {
          nextTick(() => {
            Object.assign(form.value, val);
          });
        }
      },
      { immediate: true, deep: true }
    );

    watch(
      () => form.value.parentId,
      value => {
        form.value.parentName = findMenuName(form.value.menuOptions, value);
      }
    );

    const getRef = () => formRef.value;
    const getFormData = () => form.value;

    expose({ getRef, getFormData });

    return () => (
      <ElForm ref={formRef} model={form.value} label-width="100px">
        <ElFormItem label="菜单类型" prop="menuType">
          <el-radio-group v-model={form.value.menuType}>
            <el-radio label={0}>目录</el-radio>
            <el-radio label={1}>菜单</el-radio>
            <el-radio label={2}>按钮</el-radio>
            <el-radio label={3}>内页</el-radio>
          </el-radio-group>
        </ElFormItem>

        <ElFormItem label="菜单名称" prop="name">
          <ElInput v-model={form.value.name} placeholder="请输入菜单名称" />
        </ElFormItem>

        <ElFormItem label="上级菜单" prop="parentId">
          <el-tree-select
            v-model={form.value.parentId}
            data={form.value.menuOptions}
            check-strictly
            default-expand-all
            filterable
            node-key="id"
            props={{ label: "label", children: "children" }}
            placeholder="请选择上级菜单"
            style={{ width: "100%" }}
          />
        </ElFormItem>

        <ElFormItem label="当前上级">
          <ElInput modelValue={form.value.parentName} disabled />
        </ElFormItem>

        <ElFormItem label="路由地址" prop="path">
          <ElInput v-model={form.value.path} placeholder="请输入路由地址" />
        </ElFormItem>

        {(form.value.menuType == 1 || form.value.menuType == 3) && (
          <ElFormItem label="组件路径" prop="component">
            <ElInput
              v-model={form.value.component}
              placeholder="请输入组件路径"
            />
          </ElFormItem>
        )}

        <ElFormItem label="显示顺序" prop="orderNum">
          <ElInputNumber v-model={form.value.orderNum} min={1} max={999} />
        </ElFormItem>

        <ElFormItem label="菜单图标" prop="icon">
          <IconSelect v-model={form.value.icon} />
        </ElFormItem>

        <ElFormItem label="权限标识" prop="permission">
          <ElInput
            v-model={form.value.permission}
            placeholder="请输入权限标识"
          />
        </ElFormItem>

        <ElFormItem label="路由参数" prop="query">
          <ElInput v-model={form.value.query} placeholder="请输入路由参数" />
        </ElFormItem>

        <ElFormItem label="显示状态" prop="visible">
          <el-radio-group v-model={form.value.visible}>
            <el-radio label={0}>显示</el-radio>
            <el-radio label={1}>隐藏</el-radio>
          </el-radio-group>
        </ElFormItem>

        <ElFormItem label="状态" prop="status">
          <el-radio-group v-model={form.value.status}>
            <el-radio label={1}>正常</el-radio>
            <el-radio label={0}>禁用</el-radio>
          </el-radio-group>
        </ElFormItem>
      </ElForm>
    );
  },
  methods: {
    getRef() {
      return this.$refs.formRef;
    }
  }
});
