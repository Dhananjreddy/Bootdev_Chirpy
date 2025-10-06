import { Request, Response } from "express";
import { getUserByEmail } from "../db/queries/users.js";
import { verifyPassword } from "../auth.js";
import { BadRequestError, NotFoundError, UnauthorizedError } from "./errors.js";
import { SafeUser } from "./createUser.js";
import { config } from "../config.js";
import { makeJWT, verifyJWT, makeRefreshToken, getBearerToken } from "../auth.js";
import { revokeRefreshToken, saveRefreshToken, userForRefreshToken,} from "../db/queries/refresh.js";

type LoginResponse = SafeUser & { token: string, refreshToken: string };

export async function handlerLogin(req: Request, res: Response, next: Function) {
    try {
        const { email, password } = req.body;
        if (!email || typeof email !== "string") {
            throw new BadRequestError("Email is required and must be a string");
        }
        if (!password || typeof password !== "string") {
            throw new BadRequestError("Password is required and must be a string");
        }
        const users = await getUserByEmail(email);
        if (!users.length) {
            throw new NotFoundError("Invalid email or password");
        }
        const user = users[0];
        const isPasswordValid = await verifyPassword(user.hashedPassword, password);
        if (!isPasswordValid) {
            throw new UnauthorizedError("Incorrect email or password");
        }

        let duration = config.jwt.defaultDuration;

        const accessToken =  config.jwt.secret ?  makeJWT(user.id, duration, config.jwt.secret) : "";
        const refreshToken = makeRefreshToken();

        const saved = await saveRefreshToken(user.id, refreshToken);
        if (!saved) {
            throw new UnauthorizedError("could not save refresh token");
        }
        const safeUser: LoginResponse = { email: user.email, id: user.id, createdAt: user.createdAt, updatedAt: user.updatedAt, isChirpyRed: user.isChirpyRed, token: accessToken, refreshToken: refreshToken };
        res.status(200).json(safeUser);
        return;
    } catch (e) {
        next(e)
    }
}
export async function handlerRefresh(req: Request, res: Response) {
  let refreshToken = getBearerToken(req);

  const result = await userForRefreshToken(refreshToken);
  if (!result) {
    throw new UnauthorizedError("invalid refresh token");
  }

  const user = result.user;
  const accessToken = makeJWT(
    user.id,
    config.jwt.defaultDuration,
    config.jwt.secret,
  );

  res.status(200).json({ token: accessToken });
}

export async function handlerRevoke(req: Request, res: Response) {
  const refreshToken = getBearerToken(req);
  await revokeRefreshToken(refreshToken);
  res.status(204).send();
}