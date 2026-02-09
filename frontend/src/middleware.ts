import { NextResponse } from "next/server";
import type { NextRequest } from "next/server";

export function middleware(request: NextRequest) {
  const { pathname } = request.nextUrl;
  const hasToken = request.cookies.has("access_token");

  const isProtectedRoute =
    pathname.startsWith("/home") || pathname.startsWith("/timeline");

  if (isProtectedRoute && !hasToken) {
    return NextResponse.redirect(new URL("/", request.url));
  }

  if (pathname === "/" && hasToken) {
    return NextResponse.redirect(new URL("/home", request.url));
  }

  return NextResponse.next();
}

export const config = {
  matcher: ["/", "/home/:path*", "/timeline/:path*"],
};
