import { NextResponse } from "next/server";
import { fsClient } from "../../../../lib/filesystemClient";

export async function GET() {
  const services = await fsClient.listServices();
  return NextResponse.json({ services });
}
