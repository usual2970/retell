import { Router } from "express";

const router: Router = Router();

router.get("/", (req, res) => {
  res.json({ message: "Hello from Express + TypeScript!" });
});

export default router;
