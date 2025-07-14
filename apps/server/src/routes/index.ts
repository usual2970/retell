import { Express } from "express";
import newTestHandler from "../infrastructure/rest/api/test";
import { TestService } from "../test/service";

export const register = (app: Express) => {
  const testSvc = new TestService();
  const testHandler = newTestHandler(testSvc);
  app.use("/api/v1", testHandler);
};

