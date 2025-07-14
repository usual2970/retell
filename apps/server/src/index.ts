import express from "express";
import { register as routerRegister } from "./routes";

const app = express();
const port = process.env.PORT || "3000";

app.use(express.json());

// 添加路由
routerRegister(app);

app.listen(parseInt(port), "::", () => {
  console.log(`🚀 Server is running at http://localhost:${port}`);
});
