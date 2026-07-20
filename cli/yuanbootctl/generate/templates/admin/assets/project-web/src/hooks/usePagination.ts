import { reactive } from "vue";

export function usePagination() {
  const paginationData = reactive({
    currentPage: 1,
    pageSize: 20,
    total: 0,
    pageSizes: [10, 20, 50, 100],
    layout: "total, sizes, prev, pager, next, jumper"
  });

  function handleCurrentChange(value: number) {
    paginationData.currentPage = value;
  }

  function handleSizeChange(value: number) {
    paginationData.pageSize = value;
    paginationData.currentPage = 1;
  }

  return { paginationData, handleCurrentChange, handleSizeChange };
}
