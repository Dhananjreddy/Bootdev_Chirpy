import { Request, Response } from "express";
import { BadRequestError } from "./errors.js";
import type { NewUser } from "../db/schema.js";
import { createUser } from "../db/queries/users.js";
import { getBearerToken, hashPassword, verifyJWT } from "../auth.js";
import { config } from "../config.js";
import { updateUser } from "../db/queries/users.js";

export type SafeUser = Omit<NewUser, "hashedPassword">;

export async function handlerCreateUser(req: Request, res: Response) {
    if (req.header("Content-Type") !== "application/json") {
        throw new BadRequestError("Content-Type must be application/json");
    }
    const { email, password } = req.body;
    if (!password || typeof password !== "string") {
        throw new BadRequestError("Password is required and must be at least 8 characters long");
    }
    if (!email || typeof email !== "string") {
        throw new BadRequestError("Email is required and must be a string");
    }

    const hashedPassword = await hashPassword(password);
    const user: NewUser = { email, hashedPassword: hashedPassword };
    const body = await createUser(user);
    if (!body) {
        throw new Error("User creation failed");
    }
    const safeUser: SafeUser = { email: body.email, id: body.id, createdAt: body.createdAt, updatedAt: body.updatedAt, isChirpyRed: body.isChirpyRed };
    res.status(201).json(safeUser);
    return;
}

export async function handlerUpdateUser(req: Request, res: Response, next: Function){
    const token = getBearerToken(req);
    const subject = verifyJWT(token, config.jwt.secret);

    const { email, password } = req.body;
    if (!password || typeof password !== "string") {
        throw new BadRequestError("Password is required.");
    }
    if (!email || typeof email !== "string") {
        throw new BadRequestError("Email is required and must be a string");
    }

    const hashedPassword = await hashPassword(password);
    const user = await updateUser(subject, email, hashedPassword);

    res.status(200).json({ email: user?.email, id: user?.id, createdAt: user?.createdAt, updatedAt: user?.updatedAt, isChirpyRed: user?.isChirpyRed });
    return;
}