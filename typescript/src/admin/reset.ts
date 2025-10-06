import type { Request, Response } from "express";
import { config } from "../config.js";
import { ForbiddenError } from "../api/errors.js";
import { dropUsers } from "../db/queries/dropUsers.js";

export async function handlerReset(_: Request, res: Response) {
    if (config.api.platform !== "dev") {
        throw new ForbiddenError("Reset is only allowed in dev environment");
    }
    await dropUsers();
    config.api.fileServerHits = 0;
    res.status(200).json({ message: "System reset successful" });
    return;
}

