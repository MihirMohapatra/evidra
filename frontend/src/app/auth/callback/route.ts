import { NextResponse } from "next/server";

export async function GET(request: Request) {
  const { searchParams } = new URL(request.url);
  const token = searchParams.get("token");
  const refreshToken = searchParams.get("refresh_token");

  if (!token) {
    return NextResponse.redirect(new URL("/login?error=no_token", request.url));
  }

  const response = NextResponse.redirect(new URL("/dashboard", request.url));
  response.cookies.set("token", token, { httpOnly: true, secure: true, path: "/" });
  if (refreshToken) {
    response.cookies.set("refresh_token", refreshToken, { httpOnly: true, secure: true, path: "/" });
  }

  return response;
}
