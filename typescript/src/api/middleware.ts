import type { Request, Response } from "express";
import { config } from "../config.js";

export function middlewareLogResponses(
  req: Request,
  res: Response,
  next: Function){
    res.on("finish", () => {
        const code = res.statusCode;
        if (code < 200 || code >= 300) {
            console.log(`[NON-OK] ${req.method} ${req.url} - Status: ${code}`);
        }
    });
    next();
  }

export function middlewareMetricsInc(req: Request, res: Response, next: Function) {
  config.fileserverHits ++;
  next();
}