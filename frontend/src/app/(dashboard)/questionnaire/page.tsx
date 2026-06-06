"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api-client";
import { ClipboardList, Upload } from "lucide-react";

export default function QuestionnairePage() {
  const [items, setItems] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api
      .listQuestionnaires()
      .then(setItems)
      .catch(() => setItems([]))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div>
      <div className="flex items-center justify-between mb-6">
        <h1 className="text-2xl font-bold text-gray-900">Questionnaires</h1>
        <button className="flex items-center gap-2 px-4 py-2 bg-blue-600 text-white rounded-lg text-sm font-medium hover:bg-blue-700 transition-colors">
          <Upload className="w-4 h-4" />
          Upload
        </button>
      </div>

      {loading ? (
        <div className="text-center py-12 text-gray-500">Loading...</div>
      ) : items.length === 0 ? (
        <div className="text-center py-12 text-gray-500">
          <ClipboardList className="w-12 h-12 mx-auto mb-3 text-gray-300" />
          <p className="text-lg">No questionnaires uploaded</p>
          <p className="text-sm mt-1">Upload a questionnaire to get started.</p>
        </div>
      ) : (
        <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="bg-gray-50 border-b border-gray-200">
                <th className="text-left px-4 py-3 font-medium text-gray-600">Title</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Status</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Questions</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Uploaded</th>
              </tr>
            </thead>
            <tbody>
              {items.map((q: any) => (
                <tr key={q.id} className="border-b border-gray-100 hover:bg-gray-50">
                  <td className="px-4 py-3 font-medium text-gray-900">{q.title}</td>
                  <td className="px-4 py-3">
                    <span className="text-xs font-medium px-2 py-0.5 rounded-full bg-gray-100">
                      {q.status}
                    </span>
                  </td>
                  <td className="px-4 py-3 text-gray-600">{q.question_count || "-"}</td>
                  <td className="px-4 py-3 text-gray-500">
                    {q.created_at ? new Date(q.created_at).toLocaleDateString() : "-"}
                  </td>
                </tr>
              ))}
            </tbody>
          </table>
        </div>
      )}
    </div>
  );
}
