export interface Service {
  release: string;
  repository?: string;
  chart?: string;
  tag?: string;
  overlays?: string[]; // list of overlay file paths relative to deploy/services/<release>
}

export interface SynapseClient {
  listClusters(): Promise<string[]>;
  listServices(): Promise<Service[]>;
  getService(release: string): Promise<Service | null>;
  getDiff(release: string, src: string, dst: string): Promise<string>;
  promote(
    release: string,
    src: string,
    dst: string,
    title: string
  ): Promise<{ ok: boolean; prUrl?: string; error?: string }>;
}
