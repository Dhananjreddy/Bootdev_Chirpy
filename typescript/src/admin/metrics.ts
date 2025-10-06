import type { Request, Response } from "express";
import { config } from "../config.js";

export async function handlerCountHits(_: Request, res: Response) {
  const html = `
<html>
  <body>
    <h1>Welcome, Chirpy Admin</h1>
    <p>Chirpy has been visited ${config.api.fileServerHits} times!</p>
  </body>
</html>
`;

  res.set("Content-Type", "text/html; charset=utf-8");
  res.send(html);
}