import { Request, Response } from "express";
import { BadRequestError, ForbiddenError, NotFoundError } from "./errors.js";
import { redChirpyUser } from "../db/queries/users.js";
import { getAPIKey } from "../auth.js";
import { config } from "../config.js";

export async function handlerRedChirpy(req: Request, res: Response, next: Function) {
  try {

    const apiKey = getAPIKey(req);
    if (config.api.polkaKey !== apiKey) {
      throw new ForbiddenError("Invalid API key");
    }

    const { event, data } = req.body || {};
    if (event !== "user.upgraded") return res.status(204).send();

    const userId = data?.userId;
    if (typeof userId !== "string") {
      throw new BadRequestError("User ID is required and must be a string");
    }

    const result = await redChirpyUser(userId);
    if (!result) throw new NotFoundError("User not found");

    return res.status(204).send();
  } catch (e) {
    next(e);
  }
}