import express from "express";
import helloRouter from "./routes/hello";

const app = express();
const port = process.env.PORT || 3000;

app.use(express.json());
app.use("/hello", helloRouter);

app.listen(port, () => {
  console.log(`🚀 Server is running at http://localhost:${port}`);
});
