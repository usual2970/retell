import {
  EssayCreateReq,
  essayCreateSchema,
  EssayInfoResp,
} from "@/domain/essay";
import { Request, Response, Router } from "express";
import { fail, succ } from "@/pkg/resp";

import { BaseListReq, baseListReqSchema, BaseListResp } from "@/domain/base";

export interface EssayService {
  create(essayData: EssayCreateReq): Promise<EssayInfoResp>;
  list(req: BaseListReq): Promise<BaseListResp<EssayInfoResp>>;
}

class EssayHandler {
  constructor(private readonly essayService: EssayService) {}

  createEssay = async (req: Request, res: Response): Promise<void> => {
    try {
      const result = essayCreateSchema.safeParse(req.body);
      if (!result.success) {
        return fail(res, "Invalid request data", 400);
      }

      const essay = await this.essayService.create(result.data);
      succ(res, essay);
    } catch (error) {
      console.error("Error creating essay:", error);
      fail(res, "Failed to create essay", 100);
    }
  };

  list = async (req: Request, res: Response): Promise<void> => {
    try {
      const result = baseListReqSchema.safeParse(req.query);
      if (!result.success) {
        return fail(res, "Invalid request data", 400);
      }

      const essays = await this.essayService.list(result.data);
      succ(res, essays);
    } catch (error) {
      console.error("Error listing essays:", error);
      fail(res, "Failed to list essays", 100);
    }
  };
}

export default function newEssayHandler(essayService: EssayService): Router {
  const handler = new EssayHandler(essayService);

  const router = Router();
  router.post("/essay", handler.createEssay);
  router.get("/essay", handler.list);

  return router;
}
