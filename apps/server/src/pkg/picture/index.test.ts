import { search } from "./index";

test("picture search", async () => {
  const resp = await search("cats");

  console.log(resp);
  expect(resp).toBeDefined();
});
