import { Request, Response, NextFunction } from "express";
import { getAllChirps, getChirpsByAuthorId } from "../db/queries/chirp.js";
import { getChirpsByChirpId } from "../db/queries/chirp.js";
import { NotFoundError } from "./errors.js";

export async function handlerGetAllChirps(req: Request, res: Response, next: NextFunction) {
    try {
        // Handle optional authorId filter
        const authorIdQuery = req.query.authorId;
        if (typeof authorIdQuery === "string" && authorIdQuery.trim() !== "") {
            const chirpsByAuthor = await getChirpsByAuthorId(authorIdQuery);
            res.status(200).json(chirpsByAuthor);
            return;
        }

        // Handle optional sort query parameter
        let sort: "asc" | "desc" = "asc"; // default
        if (req.query.sort === "desc") sort = "desc";

        const allChirps = await getAllChirps();
        const sortedChirps = allChirps.sort((a, b) =>
            sort === "asc"
            ? a.createdAt.getTime() - b.createdAt.getTime()
            : b.createdAt.getTime() - a.createdAt.getTime()
        );
        res.status(200).json(sortedChirps);
    } catch (err) {
        next(err);
    }
}

export async function handlerGetChirpsByChirpId(req: Request, res: Response, next: Function) {
    try {
        const { chirpID } = req.params;
        const result = await getChirpsByChirpId(chirpID);
        const chirp = result[0];
        if (!result.length) {
            throw new NotFoundError("Chirp not found");
        }
        res.status(200).json(chirp);
    } catch (e) {
        next(e)
    }
}