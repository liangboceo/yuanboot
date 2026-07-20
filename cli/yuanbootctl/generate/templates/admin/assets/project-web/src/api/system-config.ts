import { http } from "@/utils/http";

type Result = {
  data: any;
  code: number;
  status: number;
  message: string;
};

// ==================== 系统配置 ====================

export const getSystemConfig = (data?: object) => {
  return http.request<Result>("post", "sys/GetSystemConfig", { data });
};

export const saveSystemConfig = (data?: object) => {
  return http.request<Result>("post", "sys/SaveSystemConfig", { data });
};
