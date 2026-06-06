"use client";

import { useEffect, useState } from "react";
import { useParams, useRouter } from "next/navigation";
import { api } from "@/lib/api-client";
import { ArrowLeft, CheckCircle, XCircle } from "lucide-react";
import Link from "next/link";

export default function EvidenceDetailPage() {
  const { id } = useParams<{ id: string }>();
  const router = useRouter();
  const [item, setItem] = useState<any>(null);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api
      .getEvidence(id)
      .then(setItem)
      .catch(() => router.push("/evidence"))
      .finally(() => setLoading(false));
  }, [id, router]);

  async function handleApprove() {
    try {
      const updated = await api.approveEvidence(id, "", "Approved");
      setItem(updated);
    } catch (err) {
      alert("Failed to approve");
    }
  }

  async function handleReject() {
    const comment = prompt("Rejection reason:");
    if (!comment) return;
    try {
      const updated = await api.rejectEvidence(id, "", comment);
      setItem(updated);
    } catch (err) {
      alert("Failed to reject");
    }
  }

  if (loading) {
    return <div className="text-center py-12 text-gray-500">Loading...</div>;
  }

  if (!item) return null;

  const canReview = item.status === "needs_review";

  return (
    <div>
      <Link
        href="/evidence"
        className="inline-flex items-center gap-1 text-sm text-gray-500 hover:text-gray-700 mb-4"
      >
        <ArrowLeft className="w-4 h-4" />
        Back to Evidence
      </Link>

      <div className="bg-white rounded-xl border border-gray-200 p-6">
        <div className="flex items-start justify-between mb-4">
          <h1 className="text-xl font-bold text-gray-900">{item.title}</h1>
          <span className="text-sm font-medium px-3 py-1 rounded-full bg-gray-100">
            {item.status?.replace("_", " ")}
          </span>
        </div>

        <p className="text-gray-700 mb-4">{item.content}</p>

        <div className="flex flex-wrap gap-2 mb-4">
          <span className="text-xs bg-blue-100 text-blue-700 px-2 py-0.5 rounded-full">
            {item.category}
          </span>
          {item.tags?.map((tag: string) => (
            <span key={tag} className="text-xs bg-gray-100 text-gray-600 px-2 py-0.5 rounded-full">
              {tag}
            </span>
          ))}
        </div>

        {item.source_url && (
          <p className="text-sm text-gray-500 mb-2">
            Source:{" "}
            <a href={item.source_url} className="text-blue-600 hover:underline" target="_blank">
              {item.source_url}
            </a>
          </p>
        )}

        <div className="text-sm text-gray-400 space-y-1">
          <p>Version: {item.version}</p>
          <p>Created: {new Date(item.created_at).toLocaleString()}</p>
          <p>Expires: {new Date(item.expires_at).toLocaleDateString()}</p>
        </div>

        {canReview && (
          <div className="flex gap-3 mt-6 pt-4 border-t border-gray-100">
            <button
              onClick={handleApprove}
              className="flex items-center gap-2 px-4 py-2 bg-green-600 text-white rounded-lg text-sm font-medium hover:bg-green-700"
            >
              <CheckCircle className="w-4 h-4" />
              Approve
            </button>
            <button
              onClick={handleReject}
              className="flex items-center gap-2 px-4 py-2 bg-red-600 text-white rounded-lg text-sm font-medium hover:bg-red-700"
            >
              <XCircle className="w-4 h-4" />
              Reject
            </button>
          </div>
        )}
      </div>
    </div>
  );
}
