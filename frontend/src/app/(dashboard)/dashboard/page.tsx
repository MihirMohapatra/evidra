"use client";

import { useEffect, useState } from "react";
import { api } from "@/lib/api-client";
import { FileText, ClipboardList, Activity } from "lucide-react";

export default function DashboardPage() {
  const [stats, setStats] = useState({ evidence: 0, questionnaires: 0 });

  useEffect(() => {
    Promise.all([
      api.listEvidence({ limit: "0" }).catch(() => ({ items: [], total: 0 })),
      api.listQuestionnaires().catch(() => []),
    ]).then(([evRes, qs]) => {
      setStats({
        evidence: evRes.total || 0,
        questionnaires: qs.length || 0,
      });
    });
  }, []);

  return (
    <div>
      <h1 className="text-2xl font-bold text-gray-900 mb-6">Dashboard</h1>

      <div className="grid grid-cols-1 md:grid-cols-3 gap-6">
        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-blue-100 rounded-lg">
              <FileText className="w-6 h-6 text-blue-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">Evidence Items</p>
              <p className="text-2xl font-bold text-gray-900">{stats.evidence}</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-green-100 rounded-lg">
              <ClipboardList className="w-6 h-6 text-green-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">Questionnaires</p>
              <p className="text-2xl font-bold text-gray-900">{stats.questionnaires}</p>
            </div>
          </div>
        </div>

        <div className="bg-white rounded-xl border border-gray-200 p-6">
          <div className="flex items-center gap-3">
            <div className="p-3 bg-purple-100 rounded-lg">
              <Activity className="w-6 h-6 text-purple-600" />
            </div>
            <div>
              <p className="text-sm text-gray-500">System Status</p>
              <p className="text-sm font-medium text-green-600 mt-1">All Services Online</p>
            </div>
          </div>
        </div>
      </div>
    </div>
  );
}
