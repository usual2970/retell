import { search } from "./index";

test("picture search", async () => {
  const resp = await search("hello");

  console.log(resp);
  expect(resp).toBeDefined();
});
