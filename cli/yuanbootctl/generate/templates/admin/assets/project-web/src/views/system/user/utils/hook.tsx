import "./reset.css";
import dayjs from "dayjs";
import editForm from "../form/index.vue";
import { message } from "@/utils/message";
import userAvatar from "@/assets/user.jpg";
import { usePublicHooks } from "../../hooks";
import { addDialog } from "@/components/ReDialog";
import type { PaginationProps } from "@pureadmin/table";
import {
  getUserList,
  getUserDetail,
  getAllRoleList,
  getDeptList,
  saveUser,
  delUser,
  resetPwd
} from "@/api/system";
import {
  ElMessageBox,
  ElForm,
  ElInput,
  ElFormItem,
  ElProgress,
  ElRadioGroup,
  ElRadio,
  ElButton,
  ElNotification
} from "element-plus";
import { zxcvbn } from "@zxcvbn-ts/core";
import { hideTextAtIndex } from "@pureadmin/utils";
import type { FormItemProps } from "../utils/types";
import {
  type Ref,
  h,
  ref,
  toRaw,
  watch,
  computed,
  reactive,
  onMounted
} from "vue";

export function useUser(tableRef: Ref, treeRef: Ref) {
  const form = reactive({
    deptId: "",
    username: "",
    phone: "",
    status: ""
  });
  const ruleFormRef = ref();
  const dataList = ref([]);
  const loading = ref(true);
  const switchLoadMap = ref({});
  const { switchStyle } = usePublicHooks();
  const higherDeptOptions = ref();
  const treeData = ref([]);
  const treeLoading = ref(true);
  const selectedNum = ref(0);
  const pagination = reactive<PaginationProps>({
    total: 0,
    pageSize: 20,
    currentPage: 1,
    background: true
  });
  const columns: TableColumnList = [
    {
      label: "勾选列",
      type: "selection",
      fixed: "left",
      reserveSelection: true
    },
    {
      label: "用户编号",
      prop: "id",
      width: 90
    },
    {
      label: "用户头像",
      prop: "avatar",
      cellRenderer: ({ row }) => {
        const av = row.avatar || "";
        // base64 自定义头像
        if (av.startsWith("data:")) {
          return (
            <el-image
              fit="cover"
              preview-teleported={true}
              src={av}
              preview-src-list={Array.of(av)}
              class="size-6 rounded-full align-middle"
            />
          );
        }
        const src = av || userAvatar;
        return (
          <el-image
            fit="cover"
            preview-teleported={true}
            src={src}
            preview-src-list={Array.of(src)}
            class="size-6 rounded-full align-middle"
          />
        );
      },
      width: 90
    },
    {
      label: "用户名称",
      prop: "username",
      minWidth: 130
    },
    {
      label: "用户昵称",
      prop: "nickname",
      minWidth: 130
    },
    {
      label: "部门",
      prop: "deptName",
      minWidth: 120
    },
    {
      label: "手机号码",
      prop: "phone",
      minWidth: 120,
      formatter: ({ phone }) =>
        phone ? hideTextAtIndex(phone, { start: 3, end: 6 }) : "-"
    },
    {
      label: "状态",
      prop: "status",
      minWidth: 100,
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
      label: "创建时间",
      minWidth: 180,
      prop: "createTime",
      formatter: ({ createTime }) =>
        createTime ? dayjs(createTime).format("YYYY-MM-DD HH:mm:ss") : "-"
    },
    {
      label: "操作",
      fixed: "right",
      width: 260,
      slot: "operation"
    }
  ];

  const buttonClass = computed(() => [
    "h-5!",
    "reset-margin",
    "text-gray-500!",
    "dark:text-white!",
    "dark:hover:text-primary!"
  ]);

  const pwdForm = reactive({
    type: "auto",
    newPwd: ""
  });
  const pwdProgress = [
    { color: "#e74242", text: "非常弱" },
    { color: "#EFBD47", text: "弱" },
    { color: "#ffa500", text: "一般" },
    { color: "#1bbf1b", text: "强" },
    { color: "#008000", text: "非常强" }
  ];
  const curScore = ref();
  const roleOptions = ref([]);
  const defaultPasswordPolicy =
    /^(?=.*[a-z])(?=.*[A-Z])(?=.*\d)(?=.*[~!@#$%^&*()_+\-=[\]{};':"\\|,.<>/?]).{8,20}$/;

  function validatePasswordByPolicy(_rule, value: string, callback) {
    if (pwdForm.type !== "custom") {
      callback();
      return;
    }
    if (!value) {
      callback(new Error("请输入自定义密码"));
      return;
    }
    if (!defaultPasswordPolicy.test(value)) {
      callback(
        new Error("密码需为 8-20 个字符，包含大小写字母、数字和特殊字符")
      );
      return;
    }
    callback();
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

  function showResetPasswordResult(username: string, password: string) {
    ElNotification({
      title: "密码重置成功",
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
                  username
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
                  `用户名：${username}\n密码：${password}`
                );
                message("用户名和密码已复制", { type: "success" });
              }
            },
            () => "复制用户名和密码"
          )
        ])
    });
  }

  function onChange({ row, index }) {
    ElMessageBox.confirm(
      `确认要<strong>${
        row.status === 0 ? "停用" : "启用"
      }</strong><strong style='color:var(--el-color-primary)'>${
        row.username
      }</strong>用户吗?`,
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
          await saveUser({
            id: row.id,
            username: row.username,
            nickname: row.nickname,
            phone: row.phone,
            email: row.email,
            deptId: row.deptId,
            status: row.status,
            roleIds: row.roleIds,
            avatar: row.avatar
          });
          message("已成功修改用户状态", { type: "success" });
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
      `是否确认删除用户名称为<strong>${row.username}</strong>的数据?`,
      "系统提示",
      {
        confirmButtonText: "确定",
        cancelButtonText: "取消",
        type: "warning",
        dangerouslyUseHTMLString: true
      }
    )
      .then(async () => {
        await delUser({ id: row.id });
        message("删除成功", { type: "success" });
        onSearch();
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
    selectedNum.value = val.length;
    tableRef.value.setAdaptive();
  }

  function onSelectionCancel() {
    selectedNum.value = 0;
    tableRef.value.getTableRef().clearSelection();
  }

  function onbatchDel() {
    const curSelected = tableRef.value.getTableRef().getSelectionRows();
    message(`已删除用户编号为 ${curSelected.map(v => v.id).join(",")} 的数据`, {
      type: "success"
    });
    tableRef.value.getTableRef().clearSelection();
    onSearch();
  }

  async function onSearch() {
    loading.value = true;
    try {
      const { total, list } = await getUserList({
        ...toRaw(form),
        pageNum: pagination.currentPage,
        pageSize: pagination.pageSize
      });
      dataList.value = list ? list : [];
      pagination.total = total ? total : 0;
    } finally {
      setTimeout(() => {
        loading.value = false;
      }, 300);
    }
  }

  const resetForm = formEl => {
    if (!formEl) return;
    formEl.resetFields();
    form.deptId = "";
    treeRef.value?.setCurrentKey?.(null);
    treeRef.value?.onTreeReset?.();
    pagination.currentPage = 1;
    onSearch();
  };

  function onTreeSelect({ id, selected }) {
    form.deptId = selected === false ? "" : "" + id + "";
    pagination.currentPage = 1;
    onSearch();
  }

  function formatHigherDeptOptions(treeList) {
    if (!treeList || !treeList.length) return;
    const newTreeList = [];
    for (let i = 0; i < treeList.length; i++) {
      treeList[i].disabled = treeList[i].status === 0 ? true : false;
      formatHigherDeptOptions(treeList[i].children);
      newTreeList.push(treeList[i]);
    }
    return newTreeList;
  }

  const formRefInner = ref();

  function openDialog(title = "新增", row?: any) {
    addDialog({
      title: `${title}用户`,
      props: {
        formInline: {
          title,
          higherDeptOptions: formatHigherDeptOptions(higherDeptOptions.value),
          parentId: row?.deptId ?? 0,
          nickname: row?.nickname ?? "",
          username: row?.username ?? "",
          password: "",
          phone: row?.phone ?? "",
          email: row?.email ?? "",
          deptId: row?.deptId ?? 0,
          status: row?.status ?? 1,
          roleIds: row?.roleIds ?? [],
          avatar: row?.avatar ?? "",
          remark: row?.remark ?? ""
        }
      },
      width: "46%",
      draggable: true,
      closeOnClickModal: false,
      contentRenderer: ({ options }) =>
        h(editForm, {
          ref: formRefInner,
          formInline: options.props.formInline
        }),
      beforeSure: async (done, { options }) => {
        const FormRef = formRefInner.value?.getRef?.();
        if (FormRef) {
          const valid = await FormRef.validate().catch(() => false);
          if (!valid) return;
        }

        const curData =
          formRefInner.value?.getFormData?.() ??
          (options.props.formInline as FormItemProps);
        const password =
          title === "新增" ? generatePassword() : curData.password;
        try {
          await saveUser({
            id: row?.id ?? 0,
            username: curData.username,
            nickname: curData.nickname,
            password,
            phone: curData.phone,
            email: curData.email,
            deptId: curData.deptId,
            status: curData.status,
            roleIds: curData.roleIds,
            avatar: curData.avatar
          });
          message(`用户${title}成功`, { type: "success" });
          done();
          onSearch();
          if (title === "新增") {
            showResetPasswordResult(curData.username, password);
          }
        } catch (error) {
          console.error(error);
        }
      }
    });
  }

  function handleUpload(row) {
    openDialog("修改", row);
  }

  watch(
    () => pwdForm.newPwd,
    newPwd => (curScore.value = !newPwd ? -1 : zxcvbn(newPwd).score)
  );

  function handleReset(row) {
    pwdForm.type = "auto";
    pwdForm.newPwd = "";
    addDialog({
      title: `重置 ${row.username} 用户的密码`,
      width: "36%",
      draggable: true,
      closeOnClickModal: false,
      contentRenderer: () => (
        <>
          <ElForm ref={ruleFormRef} model={pwdForm} label-width="100px">
            <ElFormItem label="用户名">
              <ElInput modelValue={row.username} disabled />
            </ElFormItem>
            <ElFormItem label="密码类型" prop="type">
              <ElRadioGroup v-model={pwdForm.type}>
                <ElRadio value="auto">自动生成密码</ElRadio>
                <ElRadio value="custom">自定义密码</ElRadio>
              </ElRadioGroup>
            </ElFormItem>
            {pwdForm.type === "custom" ? (
              <>
                <ElFormItem
                  label="新密码"
                  prop="newPwd"
                  rules={[
                    { validator: validatePasswordByPolicy, trigger: "blur" }
                  ]}
                >
                  <ElInput
                    clearable
                    show-password
                    type="password"
                    v-model={pwdForm.newPwd}
                    placeholder="请输入新密码"
                  />
                </ElFormItem>
                <div class="my-4 flex">
                  {pwdProgress.map(({ color, text }, idx) => (
                    <div
                      class="w-[19vw]"
                      style={{ marginLeft: idx !== 0 ? "4px" : 0 }}
                    >
                      <ElProgress
                        striped
                        striped-flow
                        duration={curScore.value === idx ? 6 : 0}
                        percentage={curScore.value >= idx ? 100 : 0}
                        color={color}
                        stroke-width={10}
                        show-text={false}
                      />
                      <p
                        class="text-center"
                        style={{ color: curScore.value === idx ? color : "" }}
                      >
                        {text}
                      </p>
                    </div>
                  ))}
                </div>
                <div class="mb-2 ml-25 text-xs text-gray-500">
                  密码需为 8-20 个字符，包含大小写字母、数字和特殊字符
                </div>
              </>
            ) : null}
          </ElForm>
        </>
      ),
      closeCallBack: () => {
        pwdForm.type = "auto";
        pwdForm.newPwd = "";
      },
      beforeSure: async done => {
        const password =
          pwdForm.type === "auto" ? generatePassword() : pwdForm.newPwd;
        if (pwdForm.type === "custom") {
          const valid = await ruleFormRef.value
            ?.validate?.()
            .catch(() => false);
          if (!valid) return;
        }
        try {
          await resetPwd({
            id: row.id,
            password
          });
          message(`已成功重置 ${row.username} 用户的密码`, { type: "success" });
          done();
          showResetPasswordResult(row.username, password);
        } catch (error) {
          console.error(error);
        }
      }
    });
  }

  async function handleRole(row) {
    try {
      const userDetail = (await getUserDetail({ id: row.id })) as any;
      const user = userDetail?.data ?? userDetail ?? row;
      const ids = ref(user?.roleIds ?? row?.roleIds ?? []);

      addDialog({
        title: `分配 ${row.username} 用户的角色`,
        props: {
          formInline: {
            username: user?.username ?? row?.username ?? "",
            nickname: user?.nickname ?? row?.nickname ?? "",
            roleOptions: roleOptions.value ?? []
          }
        },
        width: "400px",
        draggable: true,
        closeOnClickModal: false,
        contentRenderer: () => (
          <el-select
            v-model={ids.value}
            multiple
            placeholder="请选择角色"
            class="w-full"
          >
            {roleOptions.value.map(role => (
              <el-option key={role.id} label={role.name} value={role.id} />
            ))}
          </el-select>
        ),
        beforeSure: async done => {
          await saveUser({
            id: row.id,
            username: user?.username ?? row.username,
            nickname: user?.nickname ?? row.nickname,
            phone: user?.phone ?? row.phone,
            email: user?.email ?? row.email,
            deptId: user?.deptId ?? row.deptId,
            status: user?.status ?? row.status,
            roleIds: ids.value,
            avatar: user?.avatar ?? row.avatar
          });
          row.roleIds = ids.value;
          message("角色分配成功", { type: "success" });
          done();
          onSearch();
        }
      });
    } catch (error) {
      console.error(error);
    }
  }

  onMounted(async () => {
    treeLoading.value = true;
    onSearch();

    try {
      // 归属部门
      const deptRes = (await getDeptList()) as any;
      const deptList = Array.isArray(deptRes) ? deptRes : deptRes?.data || [];
      higherDeptOptions.value = deptList;
      treeData.value = deptList;

      // 角色列表
      const roleRes = (await getAllRoleList()) as any;
      roleOptions.value = Array.isArray(roleRes)
        ? roleRes
        : roleRes?.data || [];
    } finally {
      treeLoading.value = false;
    }
  });

  return {
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
  };
}
