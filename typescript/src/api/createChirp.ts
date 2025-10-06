import { Request, Response } from "express";
import { BadRequestError, ForbiddenError, NotFoundError, UnauthorizedError } from "./errors.js";
import { NewChirp } from "../db/schema.js";
import { createChirp, getChirpsByChirpId } from "../db/queries/chirp.js";
import { getBearerToken, verifyJWT } from "../auth.js";
import { config } from "../config.js";
import { deleteChirpById } from "../db/queries/chirp.js";
import { get } from "http";

export async function handlerCreateChirp(req: Request, res: Response, next: Function) {
    try {
        if (req.header("Content-Type") !== "application/json") {
            throw new BadRequestError("Content-Type must be application/json");
        }
        const chirpBody = req.body.body;
        if (!chirpBody || typeof chirpBody !== "string") {
            throw new BadRequestError("Body is required and must be a string");
        }

        const token = getBearerToken(req);
        const userId = verifyJWT(token, config.jwt.secret);
        const cleanedBody = isValidChirpText(chirpBody);
        const chirp: NewChirp = {
            body: cleanedBody,
            userId: userId, 
        };
        const result = await createChirp(chirp);

        res.status(201).json(result);

    } catch (e) {
        next(e)
    }
} 

export async function handlerDeleteChirp(req: Request, res: Response, next: Function) {
    const token = getBearerToken(req);
    const subject = verifyJWT(token, config.jwt.secret);
    const [chirp] = await getChirpsByChirpId(req.params.chirpID);
    if (chirp.userId !== subject) {
        throw new ForbiddenError("You are not authorized to delete this chirp");
    }

    try {
        const chirpId = req.params.chirpID;
        if (!chirpId || typeof chirpId !== "string") {
            throw new BadRequestError("Chirp ID is required and must be a string");
        }
        const deletedChirp = await deleteChirpById(chirpId);
        if (deletedChirp.length === 0) {
            throw new NotFoundError("Chirp not found or already deleted");
        }
        res.status(204).send();
    } catch (e) {
        next(e)
    }
}

function isValidChirpText(text: string): string {
    if (typeof text !== "string") {
        throw new BadRequestError("Chirp text must be a string");
    }
    if (text.length === 0){
        throw new BadRequestError("Chirp text cannot be empty");
    }
    if (text.length > 140) {
        throw new BadRequestError("Chirp is too long. Max length is 140");
    }
    return cleanBody(text);
}
function cleanBody(body: string): string {
    let words = body.split(/\s+/);
    words = words.map(word => {
        if (word.toLowerCase() == "kerfuffle" || word.toLowerCase() == "sharbert" || word.toLowerCase() == "fornax") {
            return "****";
        }
        return word;})
    body = words.join(" ");
    return body;
}