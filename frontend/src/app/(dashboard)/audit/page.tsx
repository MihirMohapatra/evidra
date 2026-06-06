"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api-client";

export default function AuditPage() {
  const [events, setEvents] = useState<any[]>([]);
  const [loading, setLoading] = useState(true);

  useEffect(() => {
    api
      .listAuditEvents()
      .then((res) => setEvents(res.events || []))
      .catch(() => setEvents([]))
      .finally(() => setLoading(false));
  }, []);

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Audit Log</h1>

      {loading ? (
        <div className="text-center py-12 text-gray-500">Loading...</div>
      ) : events.length === 0 ? (
        <div className="text-center py-12 text-gray-500">No audit events recorded yet.</div>
      ) : (
        <div className="bg-white rounded-xl border border-gray-200 overflow-hidden">
          <table className="w-full text-sm">
            <thead>
              <tr className="bg-gray-50 border-b border-gray-200">
                <th className="text-left px-4 py-3 font-medium text-gray-600">Action</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Target</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Actor</th>
                <th className="text-left px-4 py-3 font-medium text-gray-600">Timestamp</th>
              </tr>
            </thead>
            <tbody>
              {events.map((evt: any) => (
                <tr key={evt.id} className="border-b border-gray-100 hover:bg-gray-50">
                  <td className="px-4 py-3">
                    <span className="font-medium text-gray-900">{evt.action}</span>
                  </td>
                  <td className="px-4 py-3 text-gray-600">{evt.target_id || "-"}</td>
                  <td className="px-4 py-3 text-gray-600">
                    {evt.actor_id ? evt.actor_id.substring(0, 8) + "..." : "-"}
                  </td>
                  <td className="px-4 py-3 text-gray-500">
                    {evt.timestamp ? new Date(evt.timestamp).toLocaleString() : "-"}
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
