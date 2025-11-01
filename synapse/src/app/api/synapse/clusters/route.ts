import { NextResponse } from "next/server";
import { fsClient } from "@synapse/lib/filesystemClient";

export async function GET() {
  const clusters = await fsClient.listClusters();
  return NextResponse.json({ clusters });
}
