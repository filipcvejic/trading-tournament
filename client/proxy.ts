import { NextRequest, NextResponse } from "next/server";

const publicRoutes = ["/login", "/register"];

export default function proxy(req: NextRequest) {
  const { pathname } = req.nextUrl;

  const token = req.cookies.get("access_token")?.value;
  const isLoggedIn = Boolean(token);

  const isPublic = publicRoutes.includes(pathname);

  // ✅ Competition paths should NEVER be redirected by "auth redirect" logic
  const isCompetitionPath =
    pathname === "/competition" || pathname.startsWith("/competition/");

  // 1) NOT logged in → everything except public goes to /login
  if (!isLoggedIn && !isPublic) {
    return NextResponse.redirect(new URL("/login", req.url));
  }

  // 2) logged in → prevent visiting login/register/root/home
  // BUT do not interfere with /competition and /competition/*
  if (
    isLoggedIn &&
    !isCompetitionPath &&
    (pathname === "/" ||
      pathname === "/home" ||
      pathname === "/login" ||
      pathname === "/register")
  ) {
    return NextResponse.redirect(new URL("/competition", req.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: [
    "/((?!api|_next/static|_next/image|favicon.ico|.*\\.(?:png|jpg|jpeg|svg|webp|ico)$).*)",
  ],
};
