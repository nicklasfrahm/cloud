import fs from "fs";
import path from "path";
import { SynapseClient, Service } from "./types";

const ROOT = path.resolve(process.cwd(), "../"); // repo root from synapse/ folder

function readDirSafe(p: string): string[] {
  try {
    return fs.readdirSync(p);
  } catch (e) {
    return [];
  }
}

export class FilesystemClient implements SynapseClient {
  async listClusters(): Promise<string[]> {
    const clustersDir = path.join(ROOT, "deploy/clusters");
    return readDirSafe(clustersDir).filter((f) =>
      fs.statSync(path.join(clustersDir, f)).isDirectory()
    );
  }

  async listServices(): Promise<Service[]> {
    const servicesDir = path.join(ROOT, "deploy/services");
    const releases = readDirSafe(servicesDir).filter((f) =>
      fs.statSync(path.join(servicesDir, f)).isDirectory()
    );
    const services: Service[] = [];
    for (const r of releases) {
      const manifestDir = path.join(servicesDir, r);
      const overlays = readDirSafe(manifestDir).filter(
        (f) => f.endsWith(".yaml") || f.endsWith(".yml")
      );
      // try to read 00-base.yaml for metadata
      const basePath = path.join(manifestDir, "00-base.yaml");
      let repository: string | undefined;
      let chart: string | undefined;
      let tag: string | undefined;
      if (fs.existsSync(basePath)) {
        const content = fs.readFileSync(basePath, "utf8");
        // naive parse of small key: value pairs
        const repoMatch = content.match(/repository:\s*(.+)/);
        const chartMatch = content.match(/chart:\s*(.+)/);
        const tagMatch = content.match(/tag:\s*(.+)/);
        if (repoMatch) repository = repoMatch[1].trim();
        if (chartMatch) chart = chartMatch[1].trim();
        if (tagMatch) tag = tagMatch[1].trim();
      }
      services.push({ release: r, repository, chart, tag, overlays });
    }
    return services;
  }

  async getService(release: string): Promise<Service | null> {
    const services = await this.listServices();
    return services.find((s) => s.release === release) ?? null;
  }

  async getDiff(release: string, src: string, dst: string): Promise<string> {
    const servicesDir = path.join(ROOT, "deploy/services", release);
    const srcPath = path.join(servicesDir, src);
    const dstPath = path.join(servicesDir, dst);
    const srcContent = fs.existsSync(srcPath)
      ? fs.readFileSync(srcPath, "utf8")
      : "";
    const dstContent = fs.existsSync(dstPath)
      ? fs.readFileSync(dstPath, "utf8")
      : "";
    // simple line-by-line diff
    const srcLines = srcContent.split(/\r?\n/);
    const dstLines = dstContent.split(/\r?\n/);
    const max = Math.max(srcLines.length, dstLines.length);
    const lines: string[] = [];
    for (let i = 0; i < max; i++) {
      const a = srcLines[i] ?? "";
      const b = dstLines[i] ?? "";
      if (a === b) {
        lines.push(" " + a);
      } else {
        if (a) lines.push("-" + a);
        if (b) lines.push("+" + b);
      }
    }
    return lines.join("\n");
  }

  async promote(
    release: string,
    src: string,
    dst: string,
    title: string
  ): Promise<{ ok: boolean; prUrl?: string; error?: string }> {
    // local filesystem promote: copy src to dst with a header comment
    const servicesDir = path.join(ROOT, "deploy/services", release);
    const srcPath = path.join(servicesDir, src);
    const dstPath = path.join(servicesDir, dst);
    try {
      if (!fs.existsSync(srcPath)) return { ok: false, error: "src not found" };
      const content = fs.readFileSync(srcPath, "utf8");
      const header = `# promoted: ${new Date().toISOString()}\n# title: ${title}\n`;
      fs.writeFileSync(dstPath, header + content, "utf8");
      return { ok: true, prUrl: "file://" + dstPath };
    } catch (e: any) {
      return { ok: false, error: String(e) };
    }
  }
}

export const fsClient = new FilesystemClient();
