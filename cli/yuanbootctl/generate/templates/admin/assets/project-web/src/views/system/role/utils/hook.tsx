import dayjs from "dayjs";
import { message } from "@/utils/message";
import { ElMessageBox } from "element-plus";
import { usePublicHooks } from "../../hooks";
import { addDialog } from "@/components/ReDialog";
import { getRoleList, getRoleMenu, saveRole, delRole } from "@/api/system";
import {
  reactive,
  ref,
  onMounted,
  h,
  toRaw,
  watch,
  defineComponent,
  nextTick
} from "vue";

export function useRole() {
  const form = reactive({
    name: "",
    code: "",
    status: ""
  });
  const curRow = ref();
  const isShow = ref(false);
  const isLinkage = ref(true);
  const isExpandAll = ref(true);
  const isSelectAll = ref(false);
  const treeSearchValue = ref("");
  const dataList = ref([]);
  const treeData = ref([]);
  const loading = ref(true);
  const switchLoadMap = ref({});
  const { switchStyle } = usePublicHooks();
  const pagination = reactive({
    total: 0,
    pageSize: 20,
    currentPage: 1,
    background: true
  });
  const treeProps = {
    value: "id",
    label: "name",
    children: "children"
  };
  const rowStyle = ({ row }) => ({
    cursor: "pointer",
    background: curRow.value?.id === row.id ? "var(--el-fill-color-light)" : ""
  });

  const columns: TableColumnList = [
    {
      label: "角色编号",
      prop: "id",
      width: 100
    },
    {
      label: "角色名称",
      prop: "name",
      minWidth: 150
    },
    {
      label: "角色编码",
      prop: "code",
      minWidth: 120
    },
    {
      label: "状态",
      prop: "status",
      width: 120,
      cellRenderer: scope => (
        <el-switch
          size={scope.props.size === "small" ? "small" : "default"}
          loading={switchLoadMap.value[scope.index]?.loading}
          v-model={scope.row.status}
          active-value={1}
          inactive-value={0}
          active-text="正常"
          inactive-text="禁用"
          inline-prompt
          style={switchStyle.value}
          onChange={() => onChange(scope as any)}
        />
      )
    },
    {
      label: "备注",
      prop: "remark",
      minWidth: 160
    },
    {
      label: "创建时间",
      prop: "createTime",
      minWidth: 180,
      formatter: ({ createTime }) =>
        createTime ? dayjs(createTime).format("YYYY-MM-DD HH:mm:ss") : "-"
    },
    {
      label: "操作",
      fixed: "right",
      width: 240,
      slot: "operation"
    }
  ];

  function onChange({ row, index }) {
    ElMessageBox.confirm(
      `确认要<strong>${row.status === 0 ? "停用" : "启用"}</strong><strong style='color:var(--el-color-primary)'>${row.name}</strong>吗?`,
      "系统提示",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
        dangerouslyUseHTMLString: true,
        draggable: true
      }
    )
      .then(async () => {
        switchLoadMap.value[index] = Object.assign(
          {},
          switchLoadMap.value[index],
          { loading: true }
        );
        try {
          const menuData = (await getRoleMenu({ roleId: row.id })) as any;
          await saveRole({
            id: row.id,
            name: row.name,
            code: row.code,
            sort: row.sort,
            status: row.status,
            remark: row.remark,
            dataScope: row.dataScope,
            menuCheckStrictly: row.menuCheckStrictly,
            deptCheckStrictly: row.deptCheckStrictly,
            menuIds: menuData?.checkedKeys ?? []
          });
          message("状态修改成功", { type: "success" });
        } finally {
          switchLoadMap.value[index] = Object.assign(
            {},
            switchLoadMap.value[index],
            { loading: false }
          );
        }
      })
      .catch(() => {
        row.status === 0 ? (row.status = 1) : (row.status = 0);
      });
  }

  function handleDelete(row) {
    ElMessageBox.confirm(
      `是否确认删除角色名称为<strong>${row.name}</strong>的数据?`,
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
          await delRole({ id: row.id });
          message("删除成功", { type: "success" });
          onSearch();
        } catch (error) {
          console.error(error);
        }
      })
      .catch(() => {});
  }

  function handleSizeChange(val: number) {
    pagination.pageSize = val;
    onSearch();
  }

  function handleCurrentChange(val: number) {
    pagination.currentPage = val;
    onSearch();
  }

  function handleSelectionChange(val) {
    console.log("handleSelectionChange", val);
  }

  async function onSearch() {
    loading.value = true;
    try {
      const result = (await getRoleList({
        ...toRaw(form),
        pageNum: pagination.currentPage,
        pageSize: pagination.pageSize
      })) as any;
      dataList.value = result?.list ?? [];
      pagination.total = result?.total ?? 0;
    } finally {
      setTimeout(() => {
        loading.value = false;
      }, 300);
    }
  }

  const resetForm = formEl => {
    if (!formEl) return;
    formEl.resetFields();
    onSearch();
  };

  const innerFormRef = ref();

  function openDialog(title = "新增", row?: any) {
    addDialog({
      title: `${title}角色`,
      props: {
        formInline: {
          title,
          id: row?.id ?? 0,
          name: row?.name ?? "",
          code: row?.code ?? "",
          sort: row?.sort ?? 1,
          status: row?.status ?? 1,
          remark: row?.remark ?? "",
          dataScope: row?.dataScope ?? 1,
          menuCheckStrictly: row?.menuCheckStrictly ?? true,
          deptCheckStrictly: row?.deptCheckStrictly ?? true
        }
      },
      width: "550px",
      draggable: true,
      closeOnClickModal: false,
      contentRenderer: () => h(RoleForm, { ref: innerFormRef }),
      beforeSure: async (done, { options }) => {
        const FormRef = innerFormRef.value?.getRef?.();
        if (FormRef) {
          const valid = await FormRef.validate().catch(() => false);
          if (!valid) return;
        }

        const curData =
          innerFormRef.value?.getFormData?.() ?? options.props.formInline;
        try {
          await saveRole({
            id: curData.id || 0,
            name: curData.name,
            code: curData.code,
            sort: curData.sort,
            status: curData.status,
            remark: curData.remark,
            dataScope: curData.dataScope,
            menuCheckStrictly: curData.menuCheckStrictly,
            deptCheckStrictly: curData.deptCheckStrictly
          });
          message(`角色${title}成功`, { type: "success" });
          done();
          onSearch();
        } catch (error) {
          console.error(error);
        }
      }
    });
  }

  /** 菜单权限 */
  const menuFormRef = ref();
  const selectedMenuIds = ref<number[]>([]);

  async function handleMenu(row?: any) {
    try {
      const data = (await getRoleMenu({ roleId: row.id })) as any;
      selectedMenuIds.value = data?.checkedKeys ?? [];

      addDialog({
        title: `分配 ${row.name} 角色的菜单权限`,
        props: {
          formInline: {
            roleId: row.id,
            roleName: row.name
          }
        },
        width: "400px",
        draggable: true,
        closeOnClickModal: false,
        contentRenderer: () =>
          h(MenuPermForm, {
            ref: menuFormRef,
            menuData: data?.menuList || [],
            checkedKeys: selectedMenuIds.value
          }),
        beforeSure: async done => {
          try {
            await saveRole({
              id: row.id,
              name: row.name,
              code: row.code,
              sort: row.sort,
              status: row.status,
              remark: row.remark,
              dataScope: row.dataScope,
              menuCheckStrictly: row.menuCheckStrictly,
              deptCheckStrictly: row.deptCheckStrictly,
              menuIds: menuFormRef.value?.getCheckedKeys?.() || []
            });
            message("菜单权限分配成功", { type: "success" });
            done();
          } catch (error) {
            console.error(error);
          }
        }
      });
    } catch (error) {
      console.error(error);
    }
  }

  async function handleSave() {
    if (!curRow.value?.id) return;
    await saveRole({
      id: curRow.value.id,
      name: curRow.value.name,
      code: curRow.value.code,
      sort: curRow.value.sort,
      status: curRow.value.status,
      remark: curRow.value.remark,
      dataScope: curRow.value.dataScope,
      menuCheckStrictly: curRow.value.menuCheckStrictly,
      deptCheckStrictly: curRow.value.deptCheckStrictly,
      menuIds: menuFormRef.value?.getCheckedKeys?.() || selectedMenuIds.value
    });
    message("菜单权限保存成功", { type: "success" });
  }

  function filterMethod(query, node) {
    return !query || node.name?.includes(query) || node.label?.includes(query);
  }

  function transformI18n(value) {
    return value;
  }

  function onQueryChanged(query) {
    menuFormRef.value?.filter?.(query);
  }

  onMounted(async () => {
    onSearch();
  });

  return {
    form,
    isShow,
    curRow,
    loading,
    columns,
    rowStyle,
    dataList,
    treeData,
    treeProps,
    isLinkage,
    pagination,
    isExpandAll,
    isSelectAll,
    treeSearchValue,
    onSearch,
    resetForm,
    openDialog,
    handleMenu,
    handleSave,
    handleDelete,
    handleSizeChange,
    filterMethod,
    transformI18n,
    onQueryChanged,
    handleCurrentChange,
    handleSelectionChange
  };
}

// 角色表单组件
const RoleForm = defineComponent({
  name: "RoleForm",
  props: {
    formInline: Object
  },
  setup(props, { expose }) {
    const ruleFormRef = ref();
    const form = ref({
      id: 0,
      name: "",
      code: "",
      sort: 1,
      status: 1,
      remark: "",
      dataScope: 1,
      menuCheckStrictly: true,
      deptCheckStrictly: true
    });

    watch(
      () => props.formInline,
      val => {
        if (val) {
          Object.assign(form.value, val);
        }
      },
      { immediate: true, deep: true }
    );

    const getRef = () => ruleFormRef.value;
    const getFormData = () => form.value;

    expose({ getRef, getFormData });

    return () => (
      <el-form ref={ruleFormRef} model={form.value} label-width="100px">
        <el-form-item
          label="角色名称"
          prop="name"
          rules={[{ required: true, message: "请输入角色名称" }]}
        >
          <el-input v-model={form.value.name} placeholder="请输入角色名称" />
        </el-form-item>
        <el-form-item
          label="角色编码"
          prop="code"
          rules={[{ required: true, message: "请输入角色编码" }]}
        >
          <el-input v-model={form.value.code} placeholder="请输入角色编码" />
        </el-form-item>
        <el-form-item label="显示顺序" prop="sort">
          <el-input-number v-model={form.value.sort} min={1} max={999} />
        </el-form-item>
        <el-form-item label="状态" prop="status">
          <el-radio-group v-model={form.value.status}>
            <el-radio label={1}>正常</el-radio>
            <el-radio label={0}>禁用</el-radio>
          </el-radio-group>
        </el-form-item>
        <el-form-item label="备注" prop="remark">
          <el-input
            v-model={form.value.remark}
            type="textarea"
            rows={3}
            placeholder="请输入备注"
          />
        </el-form-item>
      </el-form>
    );
  }
});

// 菜单权限表单组件
const MenuPermForm = defineComponent({
  name: "MenuPermForm",
  props: {
    menuData: Array,
    checkedKeys: Array
  },
  setup(props, { expose }) {
    const treeRef = ref();
    const checkedKeys = ref<number[]>([]);

    watch(
      () => props.checkedKeys,
      val => {
        checkedKeys.value = (val || []) as number[];
        nextTick(() => {
          treeRef.value?.setCheckedKeys?.(checkedKeys.value);
        });
      },
      { immediate: true, deep: true }
    );

    const getCheckedKeys = () => treeRef.value?.getCheckedKeys?.() || [];
    const filter = query => treeRef.value?.filter?.(query);

    expose({ getCheckedKeys, filter });

    return () => (
      <div>
        <el-tree
          ref={treeRef}
          data={props.menuData || []}
          show-checkbox
          node-key="id"
          props={{ label: "name", children: "children" }}
          default-expand-all
          default-checked-keys={checkedKeys.value}
        />
      </div>
    );
  }
});
