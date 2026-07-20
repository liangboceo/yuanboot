import { http } from "@/utils/http";

type Result = Array<any>;

export const getAsyncRoutes = () => {
  return http.request<Result>("post", "sys/GetRoutes");
};
