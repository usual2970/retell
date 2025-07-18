import z from "zod";

export type ModelBase = {
  createdAt?: Date;
  updatedAt?: Date;
  deletedAt?: Date | null;
};

export const baseListReqSchema = z.object({
  page: z.number().optional(),
  pageSize: z.number().optional(),
  sortBy: z.string().optional(),
  sortOrder: z.enum(["asc", "desc"]).optional(),
  filters: z.record(z.any()).optional(),
});

export type BaseListReq = z.infer<typeof baseListReqSchema>;

export const getPage = (req: BaseListReq): number => {
  return req.page ? Math.max(req.page, 1) : 1;
};

export const getPageSize = (req: BaseListReq): number => {
  return req.pageSize ? Math.max(req.pageSize, 1) : 10;
};

export const getFilter = (req: BaseListReq, key: string): any => {
  return req.filters ? req.filters[key] : undefined;
};

export type BaseListResp<T> = {
  total: number;
  page: number;
  pageSize: number;
  items: T[];
};
