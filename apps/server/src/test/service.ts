import { ITestService } from "@/infrastructure/rest/api/test";

export class TestService implements ITestService {
  performTest = async (): Promise<string> => {
    // Simulate a test operation
    return "Test completed successfully";
  };
}
