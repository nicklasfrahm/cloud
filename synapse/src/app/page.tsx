"use client";

import { useEffect, useMemo, useState } from "react";
import Link from "next/link";

interface Service {
  release: string;
  repository?: string;
  chart?: string;
  tag?: string;
}

export default function Home() {
  const [clusters, setClusters] = useState<string[]>([]);
  const [services, setServices] = useState<Service[]>([]);
  const [q, setQ] = useState("");

  useEffect(() => {
    fetch("/api/synapse/clusters")
      .then((r) => r.json())
      .then((d) => setClusters(d.clusters || []));
    fetch("/api/synapse/services")
      .then((r) => r.json())
      .then((d) => setServices(d.services || []));
  }, []);

  const filtered = useMemo(() => {
    if (!q) return services;
    return services.filter((s) =>
      s.release.toLowerCase().includes(q.toLowerCase())
    );
  }, [services, q]);
  return (
    <div className="min-h-screen bg-zinc-50 dark:bg-black p-8">
      <main className="mx-auto max-w-5xl">
        {/* Overview card */}
        <div className="bg-white dark:bg-black rounded-2xl p-8 shadow-sm mb-6">
          <div className="flex items-center justify-between mb-6">
            <h1 className="text-2xl font-semibold">Synapse</h1>
          </div>
          <div className="flex justify-between text-center">
            <div className="flex-1">
              <div className="text-sm text-zinc-500">Clusters</div>
              <div className="text-3xl font-bold">{clusters.length}</div>
            </div>
            <div className="flex-1">
              <div className="text-sm text-zinc-500">Tenants</div>
              <div className="text-3xl font-bold">{clusters.length}</div>
            </div>
            <div className="flex-1">
              <div className="text-sm text-zinc-500">Services</div>
              <div className="text-3xl font-bold">{services.length}</div>
            </div>
          </div>
        </div>

        {/* Search bar */}
        <div className="mb-6">
          <input
            value={q}
            onChange={(e) => setQ(e.target.value)}
            placeholder="Search infrastructure"
            className="w-full rounded border px-4 py-4 text-lg"
          />
        </div>

        {/* Services list */}
        <div className="bg-white dark:bg-black rounded-lg p-6">
          <h2 className="text-lg font-medium mb-4">Services</h2>
          <ul className="space-y-2">
            {filtered.map((s) => (
              <li
                key={s.release}
                className="p-3 border rounded flex items-center justify-between"
              >
                <div>
                  <div className="font-medium">{s.release}</div>
                  <div className="text-sm text-zinc-500">
                    {s.chart ?? s.repository}
                  </div>
                </div>
                <Link
                  href={`/services/${encodeURIComponent(s.release)}`}
                  className="text-sm text-blue-600"
                >
                  Open
                </Link>
              </li>
            ))}
          </ul>
        </div>
      </main>
    </div>
  );
}
