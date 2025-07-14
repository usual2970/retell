import { fail, succ } from "@/pkg/resp";
import { Router, Request, Response } from "express";

export interface ITestService {
  performTest(): Promise<string>;
}

class TestHandler {
  constructor(private readonly testService: ITestService) {}

  handleTestRequest = async (req: Request, res: Response): Promise<void> => {
    try {
      const result = await this.testService.performTest();
      succ(res, result);
    } catch (error) {
      fail(res, "Test failed", 500);
    }
  };
}

export default function newTestHandler(testService: ITestService): Router {
  const handler = new TestHandler(testService);

  const router = Router();
  router.get("/test", handler.handleTestRequest);

  return router;
}
