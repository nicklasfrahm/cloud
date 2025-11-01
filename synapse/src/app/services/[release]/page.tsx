"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { useSearchParams } from "next/navigation";

interface Service {
  release: string;
  overlays?: string[];
}

export default function ServicePage() {
  const { release } = useParams();
  const [service, setService] = useState<Service | null>(null);
  const [diff, setDiff] = useState<string | null>(null);
  const [modalOpen, setModalOpen] = useState(false);
  const [promoteStatus, setPromoteStatus] = useState<string | null>(null);

  useEffect(() => {
    fetch(`/api/synapse/services/${encodeURIComponent(release as string)}`)
      .then((r) => r.json())
      .then((d) => setService(d.service ?? null));
  }, [release]);

  async function showDiff(src: string, dst: string) {
    const res = await fetch(`/api/synapse/diff`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({ release, src, dst }),
    });
    const j = await res.json();
    setDiff(j.diff);
    setModalOpen(true);
  }

  async function promote(src: string, dst: string) {
    setPromoteStatus("running");
    const res = await fetch(`/api/synapse/promote`, {
      method: "POST",
      headers: { "Content-Type": "application/json" },
      body: JSON.stringify({
        release,
        src,
        dst,
        title: `Promote ${src} → ${dst}`,
      }),
    });
    const j = await res.json();
    if (j.prUrl) setPromoteStatus(`ok: ${j.prUrl}`);
    else setPromoteStatus(`error: ${j.error}`);
  }

  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-black p-8">
      <main className="mx-auto max-w-4xl bg-white dark:bg-black rounded-lg p-8">
        <h1 className="text-2xl font-semibold mb-4">Service: {release}</h1>
        {!service && <div>Loading...</div>}
        {service && (
          <div>
            <h2 className="font-medium mb-2">Overlays</h2>
            <div className="space-y-2">
              {service.overlays?.map((o, idx) => (
                <div key={o} className="p-3 border rounded">
                  <div className="flex items-center justify-between">
                    <div>{o}</div>
                    <div className="flex gap-2">
                      {idx > 0 && (
                        <>
                          <button
                            className="px-3 py-1 border rounded"
                            onClick={() =>
                              showDiff(service.overlays![idx - 1], o)
                            }
                          >
                            Show diff
                          </button>
                          <button
                            className="px-3 py-1 bg-blue-600 text-white rounded"
                            onClick={() =>
                              promote(service.overlays![idx - 1], o)
                            }
                          >
                            Promote
                          </button>
                        </>
                      )}
                    </div>
                  </div>
                </div>
              ))}
            </div>
          </div>
        )}

        {modalOpen && (
          <div className="fixed inset-0 bg-black/50 flex items-center justify-center">
            <div className="bg-white dark:bg-black p-6 rounded max-w-3xl w-full">
              <h3 className="font-semibold mb-2">Diff</h3>
              <pre className="max-h-96 overflow-auto p-2 bg-zinc-100 dark:bg-zinc-900 rounded">
                {diff}
              </pre>
              <div className="mt-4 flex justify-end">
                <button
                  className="px-4 py-2 border rounded"
                  onClick={() => {
                    setModalOpen(false);
                    setDiff(null);
                  }}
                >
                  Close
                </button>
              </div>
            </div>
          </div>
        )}

        {promoteStatus && <div className="mt-4">{promoteStatus}</div>}
      </main>
    </div>
  );
}
