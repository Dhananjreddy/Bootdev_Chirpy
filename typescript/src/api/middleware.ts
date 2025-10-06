import type { Request, Response, NextFunction } from "express";
import { config } from "../config.js";
import { NotFoundError, UnauthorizedError, BadRequestError, ForbiddenError } from "./errors.js";

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
  config.api.fileServerHits++;
  next();
}

export function errorHandler(err: any, _req: Request, res: Response, _next: NextFunction) {
  let statusCode = 500;
  let message = "Something went wrong on our end";

  if (err instanceof BadRequestError) {
    statusCode = 400;
    message = err.message;
  } else if (err instanceof UnauthorizedError) {
    statusCode = 401;
    message = err.message;
  } else if (err instanceof ForbiddenError) {
    statusCode = 403;
    message = err.message;
  } else if (err instanceof NotFoundError) {
    statusCode = 404;
    message = err.message;
  }

  res.status(statusCode).json({ error: message });
  return;
}

