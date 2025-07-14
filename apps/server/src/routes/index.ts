import { Express } from "express";
import newTestHandler from "../infrastructure/rest/api/test";
import { TestService } from "../test/service";
import { Repository as EssayRepository } from "../infrastructure/repository/essay";
import { Service as EssayService } from "../essay/service";
import newEssayHandler from "../infrastructure/rest/api/essay";

export const register = (app: Express) => {
  const testSvc = new TestService();
  const testHandler = newTestHandler(testSvc);

  const essayRepo = new EssayRepository();
  const essaySvc = new EssayService(essayRepo);
  const essayHandler = newEssayHandler(essaySvc);

  app.use("/api/v1", testHandler, essayHandler);
};
