"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api-client";
import { EvidenceCard } from "@/components/evidence/evidence-card";
import { Plus } from "lucide-react";
import Link from "next/link";

export default function EvidencePage() {
  const [items, setItems] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api
      .listEvidence()
      .then((res) => setItems(res.items || []))
      .catch(() => setItems([]))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Evidence</h1>
        <Link
          href="/evidence/new"
          className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors"
        >
          <Plus className="w-4 h-4" />
          New Evidence
        </Link>
      </div>

      {loading ? (
        <div className="text-center py-12 text-gray-500">Loading...</div>
      ) : items.length === 0 ? (
        <div className="text-center py-12 text-gray-500">
          <p className="text-lg">No evidence items yet</p>
          <p className="text-sm mt-1">Create your first evidence item to get started.</p>
        </div>
      ) : (
        <div className="grid grid-cols-1 md:grid-cols-2 lg:grid-cols-3 gap-4">
          {items.map((item) => (
            <EvidenceCard key={item.id} item={item} />
          ))}
        </div>
      )}
    </div>
  );
}
