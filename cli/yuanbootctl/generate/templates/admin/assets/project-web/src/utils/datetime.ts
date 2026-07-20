import dayjs from "dayjs";

export function formatDateTime(value?: string | number | Date) {
  return value ? dayjs(value).format("YYYY-MM-DD HH:mm:ss") : "-";
}
