import { NextResponse } from "next/server";
import { fsClient } from "../../../../lib/filesystemClient";

export async function POST(req: Request) {
  const body = await req.json();
  const { release, src, dst } = body;
  if (!release || !src || !dst) return NextResponse.json({ error: "missing" }, { status: 400 });
  const diff = await fsClient.getDiff(release, src, dst);
  return NextResponse.json({ diff });
}
