import type { Request, Response } from "express";
import { BadRequestError } from "./errors.js";
export async function handlerValidateChirp(req: Request, res: Response, next: Function) {
    try {
        const { body } = req.body;
        const resBody = isValidChirpText(body);
        if ("error" in resBody) {
            if (resBody.error === "Chirp is too long") {
                throw new BadRequestError("");
            }
            res.status(400).send(JSON.stringify(resBody));
        } else {
            res.status(200).send(JSON.stringify(resBody));
        }
    } catch (e) {
        next(e)
    }
}

function isValidChirpText(text: string): Object {
    if (typeof text !== "string") {
        throw new BadRequestError("Chirp text must be a string");
    }
    if (text.length === 0){
        throw new BadRequestError("Chirp text cannot be empty");
    }
    if (text.length > 140) {
        throw new BadRequestError("Chirp is too long. Max length is 140");
    }
    return { cleanedBody: cleanBody(text) };
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