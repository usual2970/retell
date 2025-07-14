import { Response as expressRes } from "express";

export type Response<T> = {
  code: number;
  message: string;
  data: T;
};

export const succ = <T>(res: expressRes, data: T) => {
  res.json({
    code: 0,
    message: "success",
    data,
  } as Response<T>);
};

export const fail = (res: expressRes, message: string, code: number = 100) => {
  res.json({
    code,
    message,
    data: null,
  } as Response<null>);
};
