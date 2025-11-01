import { NextResponse } from "next/server";
import { fsClient } from "../../../../../lib/filesystemClient";

export async function GET(
  _req: Request,
  { params }: { params: { release: string } }
) {
  const { release } = await params;
  const svc = await fsClient.getService(release);
  if (!svc) return NextResponse.json({ error: "not found" }, { status: 404 });
  return NextResponse.json({ service: svc });
}
