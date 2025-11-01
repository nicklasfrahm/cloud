import { NextResponse } from "next/server";
import { fsClient } from "../../../../lib/filesystemClient";

export async function POST(req: Request) {
  const body = await req.json();
  const { release, src, dst, title } = body;
  if (!release || !src || !dst)
    return NextResponse.json({ error: "missing" }, { status: 400 });
  const res = await fsClient.promote(release, src, dst, title || "promotion");
  if (!res.ok) return NextResponse.json({ error: res.error }, { status: 500 });
  return NextResponse.json({ prUrl: res.prUrl });
}
