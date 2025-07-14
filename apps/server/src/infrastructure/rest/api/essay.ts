import { EssayCreateReq, EssayInfoResp } from "../../../domain/essay";
import { Request, Response, Router } from "express";
import { fail, succ } from "../../../pkg/resp";
import z from "zod";

export interface EssayService {
  create(essayData: EssayCreateReq): Promise<EssayInfoResp>;
}

const essayCreateSchema = z.object({
  title: z.string().min(1, "Title is required"),
  content: z.string().min(1, "Content is required"),
});

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
}

export default function newEssayHandler(essayService: EssayService): Router {
  const handler = new EssayHandler(essayService);

  const router = Router();
  router.post("/essay", handler.createEssay);

  return router;
}
