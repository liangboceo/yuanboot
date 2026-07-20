import dayjs from "dayjs";
import { message } from "@/utils/message";
import { getDeptList, getDeptTree, saveDept, delDept } from "@/api/system";
import { usePublicHooks } from "../../hooks";
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
import { ElMessageBox } from "element-plus";

const ROOT_DEPT_ID = 0;

function appendRootDept(tree = []) {
  return [
    {
      id: ROOT_DEPT_ID,
      label: "主部门",
      name: "主部门",
      children: tree
    }
  ];
}

function removeSelfDept(tree = [], id) {
  return tree
    .filter(item => item.id !== id)
    .map(item => ({
      ...item,
      children: item.children?.length ? removeSelfDept(item.children, id) : []
    }));
}

function findDeptName(tree = [], id) {
  if (!id) return "主部门";
  for (const item of tree) {
    if (item.id === id) return item.label ?? item.name;
    const name = findDeptName(item.children || [], id);
    if (name) return name;
  }
  return "主部门";
}

export function useDept() {
  const form = reactive({
    deptName: "",
    status: null as null | number
  });

  const dataList = ref([]);
  const loading = ref(true);
  const treeData = ref([]);
  const { tagStyle } = usePublicHooks();

  const columns: TableColumnList = [
    {
      label: "部门名称",
      prop: "name",
      minWidth: 200,
      align: "left"
    },
    {
      label: "排序",
      prop: "sort",
      width: 100
    },
    {
      label: "负责人",
      prop: "leader",
      width: 120
    },
    {
      label: "联系电话",
      prop: "phone",
      width: 130
    },
    {
      label: "状态",
      prop: "status",
      width: 100,
      cellRenderer: ({ row, props }) => (
        <el-tag size={props.size} style={tagStyle.value(row.status)}>
          {row.status === 1 ? "正常" : "禁用"}
        </el-tag>
      )
    },
    {
      label: "创建时间",
      minWidth: 180,
      prop: "createTime",
      formatter: ({ createTime }) =>
        createTime ? dayjs(createTime).format("YYYY-MM-DD HH:mm:ss") : "-"
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
      const result = (await getDeptList({
        deptName: form.deptName,
        status: form.status
      })) as any;
      dataList.value = Array.isArray(result) ? result : result?.data || [];
    } finally {
      setTimeout(() => {
        loading.value = false;
      }, 300);
    }
  }

  async function loadTreeData() {
    try {
      const result = (await getDeptTree()) as any;
      treeData.value = Array.isArray(result) ? result : result?.data || [];
    } catch (error) {
      console.error(error);
    }
  }

  const innerFormRef = ref();
  function openDialog(title = "新增", row?: any) {
    const deptOptions = appendRootDept(removeSelfDept(treeData.value, row?.id));
    addDialog({
      title: `${title}部门`,
      props: {
        formInline: {
          title,
          id: row?.id ?? 0,
          parentId: row?.parentId ?? 0,
          parentName: findDeptName(treeData.value, row?.parentId ?? 0),
          deptOptions,
          name: row?.name ?? "",
          sort: row?.sort ?? 1,
          leader: row?.leader ?? "",
          phone: row?.phone ?? "",
          email: row?.email ?? "",
          status: row?.status ?? 1
        }
      },
      width: "650px",
      draggable: true,
      closeOnClickModal: false,
      contentRenderer: () => h(DeptForm, { ref: innerFormRef }),
      beforeSure: async (done, { options }) => {
        const FormRef = innerFormRef.value?.getRef?.();
        if (FormRef) {
          const valid = await FormRef.validate().catch(() => false);
          if (!valid) return;
        }

        const curData =
          innerFormRef.value?.getFormData?.() ?? options.props.formInline;
        try {
          await saveDept({
            id: curData.id || 0,
            parentId: curData.parentId || 0,
            name: curData.name,
            sort: curData.sort,
            leader: curData.leader,
            phone: curData.phone,
            email: curData.email,
            status: curData.status
          });
          message(`部门${title}成功`, { type: "success" });
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
      `是否确认删除部门名称为<strong>${row.name}</strong>的数据？${
        row.children?.length > 0
          ? "<br/><span style='color: #f56c6c'>注意：该部门下存在子部门，子部门也会一并删除</span>"
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
          await delDept({ id: row.id });
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

// 部门表单组件
const DeptForm = defineComponent({
  name: "DeptForm",
  props: {
    formInline: Object
  },
  setup(props, { expose }) {
    const form = ref({
      name: "",
      parentId: 0,
      parentName: "主部门",
      deptOptions: [],
      sort: 1,
      leader: "",
      phone: "",
      email: "",
      status: 1
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
        form.value.parentName = findDeptName(form.value.deptOptions, value);
      }
    );

    const getRef = () => formRef.value;
    const getFormData = () => form.value;

    expose({ getRef, getFormData });

    return () => (
      <el-form ref={formRef} model={form.value} label-width="100px">
        <el-form-item
          label="部门名称"
          prop="name"
          rules={[{ required: true, message: "请输入部门名称" }]}
        >
          <el-input v-model={form.value.name} placeholder="请输入部门名称" />
        </el-form-item>

        <el-form-item label="上级部门" prop="parentId">
          <el-tree-select
            v-model={form.value.parentId}
            data={form.value.deptOptions}
            check-strictly
            default-expand-all
            filterable
            node-key="id"
            props={{ label: "label", children: "children" }}
            placeholder="请选择上级部门"
            style={{ width: "100%" }}
          />
        </el-form-item>

        <el-form-item label="当前上级">
          <el-input modelValue={form.value.parentName} disabled />
        </el-form-item>

        <el-form-item label="显示顺序" prop="sort">
          <el-input-number v-model={form.value.sort} min={1} max={999} />
        </el-form-item>

        <el-form-item label="负责人" prop="leader">
          <el-input v-model={form.value.leader} placeholder="请输入负责人" />
        </el-form-item>

        <el-form-item label="联系电话" prop="phone">
          <el-input v-model={form.value.phone} placeholder="请输入联系电话" />
        </el-form-item>

        <el-form-item label="邮箱" prop="email">
          <el-input v-model={form.value.email} placeholder="请输入邮箱" />
        </el-form-item>

        <el-form-item label="状态" prop="status">
          <el-radio-group v-model={form.value.status}>
            <el-radio label={1}>正常</el-radio>
            <el-radio label={0}>禁用</el-radio>
          </el-radio-group>
        </el-form-item>
      </el-form>
    );
  }
});
